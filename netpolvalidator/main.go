package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Used to gather connection information from kube-burner proxy pod
type connection struct {
	Addresses []string `json:"addresses"`
	Ports     []int    `json:"ports"`
	Netpol    string   `json:"netpol"`
}

var connections []connection

// Local copy of connection information and also stores the result i.e succesful connection timestamp
type connTest struct {
	Address    string    `json:"address"`
	Port       int       `json:"port"`
	IngressIdx int       `json:"connectionidx"`
	NpName     string    `json:"npname"`
	Timestamp  time.Time `json:"timestamp"`
}

var parallelConnections = 10

var (
	resultsLock    sync.Mutex
	connTestLock   sync.Mutex
	results        = make([]connTest, 0)
	wg             sync.WaitGroup
	failedConnChan = make(chan connTest, 1000)
	netpolTimeout  = time.Second * 10
	httpClient     = http.Client{Timeout: 1500 * time.Millisecond}
	gotConnectins  = make(chan bool)
)

var allConnTests []connTest

func sendRequest(address string, port int) (bool, time.Time) {
	url := fmt.Sprintf("http://%s:%d", address, port)
	log.Printf("Sending request to address %s", address)
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Printf("Error connecting to address %v %v", address, err)
		return false, time.Now().UTC()
	}
	defer resp.Body.Close()
	log.Printf("Got %v response to address %s", resp.StatusCode, address)
	return resp.StatusCode == http.StatusOK, time.Now().UTC()
}

// Send requests to provided addresses as part of connections.
// Return succesful and failed connections
func testConnections(connTests []connTest) ([]connTest, []connTest) {
	var successConn []connTest
	var failedConn []connTest
	for _, ct := range connTests {
		for attempt := 1; attempt <= 3; attempt++ {
			success, timestamp := sendRequest(ct.Address, ct.Port)
			if success {
				ct.Timestamp = timestamp
				successConn = append(successConn, ct)
				break
			} else if attempt == 3 {
				failedConn = append(failedConn, ct)
				log.Printf("failed request to %v after 3 attempts", ct.Address)
			}
		}
	}
	if len(successConn) > 0 {
		resultsLock.Lock()
		for _, conn := range successConn {
			results = append(results, conn)
		}
		resultsLock.Unlock()
	}
	return successConn, failedConn
}

// This pod has to wait till kube-burner starts creating network policies.
// To implement this wait, we test connections using first 3 network policies.
// Once kube-burner starts running the job and creating network policies,
// checking the 3 network policies will be succesful, then we exit the wait
// and start processing all connections
func waitForJobStarted(connTests []connTest) {
	var jobStartConn map[string]connTest = make(map[string]connTest)
	addConn := 1
	success := false
	for _, conn := range connTests {
		if _, ok := jobStartConn[conn.NpName]; !ok {
			jobStartConn[conn.NpName] = conn
			if addConn == 3 {
				break
			}
			addConn += 1
		}
	}
	for {
		for _, conn := range jobStartConn {
			success, _ = sendRequest(conn.Address, conn.Port)
			if success {
				break
			}
		}
		if success {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// Dedicated thread to test failed connections
func processFailed() {
	batchSize := 10
	semaphore := make(chan struct{}, parallelConnections)
	for {
		batch := make([]connTest, 0, batchSize)
		for i := 0; i < batchSize; i++ {
			select {
			case val := <-failedConnChan:
				batch = append(batch, val)
			default:
				// No more data available, process what we have
				break
			}
		}
		if len(batch) > 0 {
			semaphore <- struct{}{}
			go func(batch []connTest, semaphore chan struct{}) {
				defer func() { <-semaphore }()
				_, failedConn := testConnections(batch)
				if len(failedConn) > 0 {
					for _, c := range failedConn {
						failedConnChan <- c
					}
				}
			}(batch, semaphore)
		}
		time.Sleep(100 * time.Millisecond) // Avoid tight loop
	}
}

// Test sending requests to the provided addresses.
// Move failed connections to a dedicated thread.
func processIngress(connTests []connTest, semaphore chan struct{}, tc int) {
	defer wg.Done()
	defer func() { <-semaphore }()
	done := make(chan bool, 1)
	success := make(chan bool, 1)
	testChan := make(chan []connTest, 1)
	ticker := time.NewTicker(time.Second)
	go func(threadCount int) {
		var failedConn []connTest
		for {
			select {
			case <-done:
				testChan <- failedConn
				return
			case <-ticker.C:
				_, failedConn = testConnections(connTests)
				if len(failedConn) != len(connTests) {
					testChan <- failedConn
					success <- true
					return
				}
			}
		}
	}(tc)
	select {
	case <-success:
	case <-time.After(netpolTimeout):
		log.Println("Timeout reached processing network policy.")
	}
	done <- true
	ticker.Stop()
	// retry failed conn in 3 attempts
	failedConn := <-testChan
	for attempt := 1; attempt <= 3; attempt++ {
		_, failedConn = testConnections(failedConn)
		if len(failedConn) == 0 {
			break
		}
	}
	// if still they fail move them to another thread
	if len(failedConn) > 0 {
		for _, c := range failedConn {
			failedConnChan <- c
		}
	}
}

// Return the results to kube-burner proxy pod
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	resultsLock.Lock()
	defer resultsLock.Unlock()
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get connections from kube-burner proxy pod and create a local copy of
// this information.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Check Request received, processing...")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &connections)
	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusBadRequest)
		return
	}
	log.Println("Start sending Connections info ")
	for connectionIdx, connection := range connections {
		for _, address := range connection.Addresses {
			for _, port := range connection.Ports {
				allConnTests = append(allConnTests, connTest{Address: address, Port: port, IngressIdx: connectionIdx, NpName: connection.Netpol})
			}
		}
	}
	r.Body.Close()
	log.Println("Finished sending Connections info")
	gotConnectins <- true
}

// function to process connections
func sendRequests() {
	// Wait till we get connections from kube-burner proxy pod
	<-gotConnectins
	// Wait till kube-burner creates job's objects
	log.Println("Start waiting for Network policy object creation ", allConnTests)
	waitForJobStarted(allConnTests)
	log.Println("Finished waiting for Network policy object creation ")
	// Start a dedicated thread for processing failed connections
	go processFailed()
	// Use 20 parallel threads for connection testing.
	// Each thread will process 20 requests
	semaphore := make(chan struct{}, parallelConnections)
	for i := 0; i < len(allConnTests); i += parallelConnections {
		end := i + parallelConnections
		if end > len(allConnTests) {
			end = len(allConnTests)
		}
		semaphore <- struct{}{}
		wg.Add(1)
		go processIngress(allConnTests[i:end], semaphore, i)

	}
	wg.Wait()
}

func processEnvVars() {
	var err error
	parallelConnectionsStr := os.Getenv("PARALLEL_CONNECTIONS")
	if parallelConnectionsStr != "" {
		parallelConnections, err = strconv.Atoi(parallelConnectionsStr)
		if err != nil {
			panic(fmt.Sprintf("failed to parse env PARALLEL_CONNECTIONS: %v", err))
		}
	}
}

func main() {
	processEnvVars()
	go sendRequests()
	http.HandleFunc("/check", handleRequest)
	http.HandleFunc("/results", resultsHandler)
	log.Println("Server started on 127.0.0.1:9001")
	go func() {
		log.Fatal(http.ListenAndServe(":9001", nil))
	}()

	select {} // Keep the server running
}
