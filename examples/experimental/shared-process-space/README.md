# Testing Communication between Containers

We are going to test running this application in the context of a [shared process namespace between containers in a pod](https://kubernetes.io/docs/tasks/configure-pod-container/share-process-namespace/).

## Go Experiment

Create a cluster, and install the Flux Operator

```bash
kind create cluster
kubectl apply -f ../../dist/flux-operator-dev.yaml
```

Create the interactive Minicluster. The [goshare](https://github.com/converged-computing/goshare) client and server will
be installed to two containers. The server has the application we want to run, and the client has flux.

```bash
$ kubectl create namespace flux-operator
$ kubectl apply -f minicluster.yaml
```

We will test this interactively for now. In the future we will want to:

- install the client/server depending on container
- find the correct PID for the running server based on matching some name or similar
- start the client with the common socket path.

Wait until your pods are all running:

```bash
$ kubectl get pods -n flux-operator 
```
```console
NAME                  READY   STATUS    RESTARTS   AGE
flux-sample-0-k5ccg   2/2     Running   0          7m36s
flux-sample-1-bb8ks   2/2     Running   0          7m36s
flux-sample-2-5cwk4   2/2     Running   0          7m36s
flux-sample-3-jggrg   2/2     Running   0          7m36s
```

You can then watch the logs of a server container to see the command being run.
```bash
$ kubectl logs -n flux-operator flux-sample-0-wpsnj -c server
```
```console
task: [build] GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/client cmd/client/client.go
🟦️ service: 2023/07/26 22:42:52 server.go:38: starting service at socket /dinosaur.sock
🟦️ service: 2023/07/26 22:42:52 server.go:50: creating a new service to listen at /dinosaur.sock
🟦️ service: 2023/07/26 22:43:57 command.go:26: start new stream request
🟦️ service: 2023/07/26 22:43:57 command.go:54: Received command echo hello world
🟦️ service: 2023/07/26 22:43:57 command.go:67: send new pid=3025
🟦️ service: 2023/07/26 22:43:57 command.go:70: Process started with PID: 3025
🟦️ service: 2023/07/26 22:43:57 command.go:76: send final output: hello world
```

Note that this experiment has it running twice - once outside of flux, and one with flux run.
The latter doesn't seem to work, at least I haven't figured out why it works outside flux but not within it.