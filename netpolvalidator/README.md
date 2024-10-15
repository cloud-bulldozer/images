# Connection Testing and Latency Measurement for Network Policies Validation

Kube-burner utilizes a client pod, built from this image, to send connections to designated server IPs. Kube-burner provides the client pod with multiple IP addresses, which the pod uses to send requests. Once a request is successfully processed, the client pod records the timestamp and sends this information back to Kube-burner. A proxy pod is employed by Kube-burner to facilitate communication between itself and the client pod.

Key Details:

- All client pods listen on port 9001 for incoming requests.
- The Kube-burner proxy pod communicates with the client pod via the `/check` endpoint. This triggers the client pod to continuously send HTTP requests to the first three remote addresses from the connection list. Once any of these requests succeed, it confirms that the network policies have been applied, and then proceeds to send requests to all remaining connections.
- To handle high volume and scale efficiently, the client pod uses 20 parallel Goroutines. Each Goroutine handles sending requests to 20 different IP addresses concurrently.
- Each HTTP request is retried up to 3 times, with a timeout of 1.5 seconds per request.
- If a request fails after 3 attempts, it is considered failed and added to a dedicated channel. A separate Goroutine monitors this channel and retries the failed requests.
- Upon successful completion of a request, the timestamp is recorded for future reference.
- Finally, the proxy pod gathers all the results from the client pod by querying the `/results` endpoint.

Log from one of the client pods

```shell
$ oc logs -n network-policy-perf-0 -c curlapp test-pod-1 -f
2024/10/01 10:39:17 Server started on 127.0.0.1:9001
2024/10/01 10:39:34 Start sending Connections info 
2024/10/01 10:39:34 Finished sending Connections info
2024/10/01 10:39:34 Start waiting for Network policy object creation  [{10.128.2.44 8080 0 ingress-1-1 0001-01-01 00:00:00 +0000 UTC}]
2024/10/01 10:39:34 Sending request to address 10.128.2.44
2024/10/01 10:39:35 Error connecting to address 10.128.2.44 Get "http://10.128.2.44:8080": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2024/10/01 10:39:36 Sending request to address 10.128.2.44
2024/10/01 10:39:37 Error connecting to address 10.128.2.44 Get "http://10.128.2.44:8080": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2024/10/01 10:39:37 Sending request to address 10.128.2.44
2024/10/01 10:39:39 Error connecting to address 10.128.2.44 Get "http://10.128.2.44:8080": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2024/10/01 10:39:39 Sending request to address 10.128.2.44
2024/10/01 10:39:40 Error connecting to address 10.128.2.44 Get "http://10.128.2.44:8080": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
2024/10/01 10:39:40 Sending request to address 10.128.2.44
2024/10/01 10:39:40 Got 200 response to address 10.128.2.44
2024/10/01 10:39:40 Finished waiting for Network policy object creation 
2024/10/01 10:39:41 Sending request to address 10.128.2.44
2024/10/01 10:39:41 Got 200 response to address 10.128.2.44

```
