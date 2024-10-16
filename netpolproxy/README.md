# Kube Burner Network Policy Proxy Pod for Connection Testing and Latency Measurement
Kube-burner employs a proxy pod to interact with client pods, which helps streamline communication and avoid the need for direct routes or executing commands on each client pod. This is particularly beneficial during large-scale tests, where a significant number of client pods are created. The proxy pod facilitates both the delivery of connection information to client pods and the retrieval of results, reducing overhead and complexity.

The proxy pod is built using a specific image and listens on port 9002. This port is enabled on worker nodes by default via AWS security groups. The proxy pod is equipped with 5 handlers and operates across two primary flows:

- Sending connection information to client pods
- Retrieving connection results from client pods

### Workflow:
- Initialization and Connection Setup:
  + The proxy pod initially waits to receive connection information from Kube-burner.
  + Once Kube-burner sends the connection information via the `/initiate` endpoint, the proxy pod uses 20 parallel Goroutines to distribute this information to all the client pods efficiently.
  + Kube-burner waits for up to 30 minutes, periodically checking every 5 seconds using the `/checkConnectionsStatus` endpoint to confirm whether the proxy pod has successfully delivered the connection details to all client pods.
  
- Retrieving Test Results:
  + After the testing phase is complete, Kube-burner triggers the `/stop` endpoint on the proxy pod. This signals the proxy pod to begin retrieving results from all the client pods.
  + Similar to the connection phase, the proxy pod employs 20 parallel Goroutines to gather results from the client pods.
  + Kube-burner again waits for a maximum of 30 minutes, checking every 5 seconds via the `/checkStopStatus` endpoint to ensure that the proxy pod has retrieved results from all client pods.
  + Once all results are collected, Kube-burner retrieves the final data by querying the `/results` endpoint on the proxy pod.


Log from one of the client pods

```shell
$ oc logs -n network-policy-proxy network-policy-proxy -f
2024/10/01 11:18:02 Client server started on :9002
2024/10/01 11:18:02 Wait for Connections from kube-burner..
2024/10/01 11:18:25 Number of connections got from kube-burner 2
2024/10/01 11:18:25 Got connections from kube-burner, sending them to 2 pods
2024/10/01 11:18:25 Connections sent to http://10.128.2.50:9001/check successfully
2024/10/01 11:18:25 Connections sent to http://10.128.2.51:9001/check successfully
2024/10/01 11:18:47 Address: 10.128.2.50, Port: 8080, IngressIdx: 0, NpName: ingress-0-1 Timestamp: 2024-10-01 11:18:33.247061614 +0000 UTC
2024/10/01 11:18:47 Address: 10.128.2.51, Port: 8080, IngressIdx: 0, NpName: ingress-1-1 Timestamp: 2024-10-01 11:18:33.248359122 +0000 UTC
```
