# ZeroMQ in Go Examples

Note that we are going to try to use the [DEALER to ROUTER](https://zguide.zeromq.org/docs/chapter3/#The-DEALER-to-ROUTER-Combination) design here.

Create the kind cluster.

```bash
kind create cluster --config ../../kind-config.yaml
```

Install the flux operator

```bash
kubectl apply -f ../../dist/flux-operator.yaml
```

## Local Test

You can automate the entire thing:

```bash
./build.sh
```

And then get logs:

```console
Hello I'm host flux-sample-1
Hello I'm host flux-sample-0
  ⭐️ Times for flux-sample-0 to flux-sample-1: [5.279µs 15.12µs 66.039µs 27.815µs 11.285µs 15.554µs 8.282µs 7.801µs 3.895µs 4.585µs]
  ⭐️ Times for flux-sample-1 to flux-sample-0: [10.757µs 18.31µs 7.088µs 12.65µs 7.853µs 4.582µs 3.825µs 16.614µs 48.301µs 9.307µs]
```

You can look at [build.sh](build.sh) for the build steps, and [entrypoint.sh](entrypoint.sh) for the start command,
and [main.go](main.go) for the defaults and logic.

## Google Cloud

The times are probably fast since we are running on the same machine. Let's test on actual physical nodes on Google Cloud.

```bash
GOOGLE_PROJECT=myproject
gcloud container clusters create test-cluster \
    --threads-per-core=1 \
    --placement-type=COMPACT \
    --num-nodes=8 \
    --region=us-central1-a \
    --project=${GOOGLE_PROJECT} \
    --machine-type=c2d-standard-8
```

Install the flux operator.

```bash
kubectl apply -f https://raw.githubusercontent.com/flux-framework/flux-operator/main/examples/dist/flux-operator.yaml
```

We will build a container that the minicluster can actually pull.

```bash
docker build -t vanessa/gozmq:0 .
docker push vanessa/gozmq:0
```

And then apply the GKE minicluster.

```bash
kubectl apply -f minicluster-gke.yaml
kubectl logs flux-sample-0-2prsv -f
```
```console
Hello I'm host flux-sample-0
Hello I'm host flux-sample-5
Hello I'm host flux-sample-6
Hello I'm host flux-sample-3
Hello I'm host flux-sample-2
Hello I'm host flux-sample-1
Hello I'm host flux-sample-4
Hello I'm host flux-sample-7
  ⭐️ Times for flux-sample-0 to flux-sample-1: [3.01µs 2.8µs 3.82µs 2.72µs 3.21µs 2.6µs 2.6µs 3.3µs 3.37µs 3µs]
  ⭐️ Times for flux-sample-0 to flux-sample-2: [3.03µs 3.01µs 3.04µs 2.711µs 3.02µs 3.331µs 2.98µs 2.5µs 2.33µs 3.56µs]
  ⭐️ Times for flux-sample-0 to flux-sample-3: [2.37µs 3.37µs 2.44µs 2.4µs 3.43µs 2.65µs 2.6µs 2.94µs 2.3µs 3.35µs]
  ⭐️ Times for flux-sample-0 to flux-sample-4: [3.16µs 2.55µs 8.07µs 2.78µs 2.55µs 2.95µs 2.54µs 2.48µs 2.84µs 3.13µs]
  ⭐️ Times for flux-sample-0 to flux-sample-5: [3.74µs 2.42µs 2.92µs 3.33µs 3.36µs 4.41µs 3.41µs 5.29µs 3.09µs 2.95µs]
  ⭐️ Times for flux-sample-0 to flux-sample-6: [2.67µs 3.07µs 3.1µs 2.66µs 3.78µs 2.011µs 3.21µs 2.72µs 6.44µs 3.9µs]
  ⭐️ Times for flux-sample-0 to flux-sample-7: [4.46µs 9.2µs 3.11µs 8.29µs 3.54µs 3.15µs 3.43µs 2.48µs 2.37µs 2.75µs]
  ⭐️ Times for flux-sample-3 to flux-sample-0: [3.42µs 3.89µs 3.111µs 3.6µs 3.03µs 3.36µs 3.02µs 3.09µs 4.31µs 4.42µs]
  ⭐️ Times for flux-sample-1 to flux-sample-0: [5.649µs 3.15µs 3.18µs 3.86µs 4.18µs 3.19µs 4.19µs 3.31µs 4.159µs 3.74µs]
  ⭐️ Times for flux-sample-2 to flux-sample-0: [4.05µs 3.35µs 5.8µs 4.5µs 2.89µs 3.22µs 3.92µs 3.1µs 4.26µs 4.46µs]
  ⭐️ Times for flux-sample-6 to flux-sample-0: [3.02µs 2.91µs 3.66µs 3.19µs 2.77µs 2.89µs 4.02µs 4.969µs 3.08µs 3.04µs]
  ⭐️ Times for flux-sample-4 to flux-sample-0: [3.3µs 4.12µs 3.429µs 3.78µs 3.24µs 5.89µs 3.52µs 3.26µs 6.6µs 4.88µs]
  ⭐️ Times for flux-sample-1 to flux-sample-2: [5.96µs 9.9µs 3.82µs 3.011µs 3µs 3.68µs 3.18µs 2.77µs 5.73µs 4.24µs]
  ⭐️ Times for flux-sample-5 to flux-sample-0: [3.44µs 2.84µs 3.43µs 3.88µs 4.14µs 3.73µs 4.34µs 5.431µs 4.329µs 3.02µs]
  ⭐️ Times for flux-sample-7 to flux-sample-0: [2.95µs 3.269µs 2.88µs 3.14µs 2.491µs 3.84µs 2.86µs 3.529µs 2.9µs 2.93µs]
  ⭐️ Times for flux-sample-1 to flux-sample-3: [2.94µs 3.04µs 4.31µs 3.81µs 2.78µs 3.33µs 3.28µs 2.99µs 2.78µs 10.02µs]
  ⭐️ Times for flux-sample-1 to flux-sample-4: [2.66µs 14.87µs 3.04µs 3.14µs 2.98µs 2.7µs 2.7µs 2.86µs 2.55µs 3.61µs]
  ⭐️ Times for flux-sample-1 to flux-sample-5: [3.21µs 20.709µs 3.52µs 3.45µs 3.431µs 3.12µs 2.57µs 3.31µs 3.15µs 2.55µs]
  ⭐️ Times for flux-sample-3 to flux-sample-1: [9.88µs 5.01µs 4µs 4.65µs 3.471µs 3.22µs 3.29µs 3.25µs 3.11µs 4.04µs]
  ⭐️ Times for flux-sample-2 to flux-sample-1: [3.99µs 2.78µs 4.02µs 3.62µs 3.39µs 3.42µs 3.67µs 4.26µs 3.08µs 2.91µs]
  ⭐️ Times for flux-sample-1 to flux-sample-6: [2.87µs 3.13µs 3.33µs 2.89µs 3.15µs 2.53µs 3.06µs 3.1µs 11.16µs 3.05µs]
  ⭐️ Times for flux-sample-5 to flux-sample-1: [3.64µs 7.3µs 3.81µs 16.24µs 3.86µs 3.42µs 2.96µs 3.19µs 3.85µs 3.03µs]
  ⭐️ Times for flux-sample-2 to flux-sample-3: [3.78µs 2.87µs 3.54µs 3.43µs 2.9µs 9.64µs 2.38µs 2.9µs 2.58µs 2.58µs]
  ⭐️ Times for flux-sample-4 to flux-sample-1: [3.29µs 6.68µs 3.85µs 4.15µs 6.09µs 3.8µs 3.14µs 3.02µs 4.05µs 3.57µs]
  ⭐️ Times for flux-sample-1 to flux-sample-7: [3.13µs 3.04µs 4.04µs 3.19µs 2.99µs 2.611µs 2.94µs 4.34µs 3.82µs 9.65µs]
  ⭐️ Times for flux-sample-7 to flux-sample-1: [3.009µs 3.4µs 3.26µs 3.91µs 2.85µs 2.84µs 3.151µs 3.16µs 3.58µs 3.21µs]
  ⭐️ Times for flux-sample-6 to flux-sample-1: [3.43µs 4.65µs 3.38µs 3.02µs 2.95µs 3.03µs 3.08µs 3.071µs 3µs 2.71µs]
  ⭐️ Times for flux-sample-3 to flux-sample-2: [3.109µs 3.68µs 4.38µs 3.309µs 3.27µs 3.04µs 2.54µs 3µs 3.16µs 4.08µs]
  ⭐️ Times for flux-sample-2 to flux-sample-4: [2.75µs 2.99µs 3.25µs 2.25µs 2.81µs 3.24µs 5.2µs 4.08µs 4.14µs 2.94µs]
  ⭐️ Times for flux-sample-4 to flux-sample-2: [4.169µs 3.32µs 3.09µs 3.26µs 3.28µs 3.8µs 13.37µs 2.41µs 3.6µs 2.94µs]
  ⭐️ Times for flux-sample-2 to flux-sample-5: [3.28µs 3.16µs 2.46µs 2.88µs 2.9µs 3.3µs 4µs 3.24µs 3.49µs 2.87µs]
  ⭐️ Times for flux-sample-5 to flux-sample-2: [3.66µs 2.689µs 2.96µs 4.02µs 3.02µs 21.68µs 2.66µs 3.27µs 3.209µs 3.16µs]
  ⭐️ Times for flux-sample-3 to flux-sample-4: [3.77µs 2.67µs 3µs 2.74µs 3.75µs 2.76µs 3.77µs 3.09µs 2.92µs 5.22µs]
  ⭐️ Times for flux-sample-2 to flux-sample-6: [3.01µs 3.19µs 2.88µs 2.73µs 3.1µs 3.829µs 2.57µs 4.54µs 4.25µs 3.78µs]
  ⭐️ Times for flux-sample-4 to flux-sample-3: [3.87µs 3.66µs 3.2µs 10.11µs 3.54µs 2.95µs 2.31µs 2.93µs 2.5µs 3.76µs]
  ⭐️ Times for flux-sample-5 to flux-sample-3: [3µs 3.18µs 2.44µs 3.03µs 3.02µs 6.33µs 3.389µs 3.18µs 2.82µs 3.2µs]
  ⭐️ Times for flux-sample-6 to flux-sample-2: [3.06µs 3.53µs 3.47µs 2.85µs 2.62µs 3.07µs 3.07µs 2.93µs 2.689µs 3.18µs]
  ⭐️ Times for flux-sample-3 to flux-sample-5: [3.12µs 4.04µs 3.509µs 3.15µs 3.1µs 2.98µs 2.77µs 2.71µs 3.11µs 2.869µs]
  ⭐️ Times for flux-sample-7 to flux-sample-2: [2.63µs 3.47µs 3.08µs 12.26µs 3.15µs 2.82µs 2.891µs 5.41µs 3.26µs 2.98µs]
  ⭐️ Times for flux-sample-2 to flux-sample-7: [4.45µs 3.42µs 2.56µs 4.74µs 3.75µs 3.14µs 1.78µs 3.39µs 3.37µs 3.24µs]
  ⭐️ Times for flux-sample-5 to flux-sample-4: [3.7µs 3.02µs 3.23µs 4.611µs 3.22µs 4.38µs 3.25µs 2.47µs 7.02µs 3.29µs]
  ⭐️ Times for flux-sample-3 to flux-sample-6: [2.99µs 3.1µs 3.08µs 3.07µs 2.81µs 3.92µs 2.97µs 3.8µs 2.47µs 2.77µs]
  ⭐️ Times for flux-sample-4 to flux-sample-5: [2.99µs 6.609µs 3.59µs 3.591µs 3.43µs 3.31µs 3.16µs 3.611µs 3.6µs 2.94µs]
  ⭐️ Times for flux-sample-6 to flux-sample-3: [2.6µs 2.93µs 3.08µs 3.02µs 2.78µs 6.36µs 3.02µs 2.79µs 3.089µs 3.26µs]
  ⭐️ Times for flux-sample-3 to flux-sample-7: [2.86µs 3.529µs 3.41µs 2.8µs 2.91µs 2.78µs 2.43µs 3.08µs 5.46µs 2.84µs]
  ⭐️ Times for flux-sample-5 to flux-sample-6: [2.41µs 2.28µs 4.59µs 3.27µs 13.72µs 3.54µs 2.79µs 3.73µs 4.37µs 2.651µs]
  ⭐️ Times for flux-sample-7 to flux-sample-3: [3.8µs 2.98µs 2.73µs 5.53µs 6.36µs 3.5µs 3.22µs 2.62µs 3.04µs 3.07µs]
  ⭐️ Times for flux-sample-4 to flux-sample-6: [2.83µs 3.13µs 2.95µs 3.29µs 2.66µs 2.65µs 3.7µs 3.07µs 2.71µs 2.31µs]
  ⭐️ Times for flux-sample-6 to flux-sample-4: [7.01µs 3µs 3.02µs 2.81µs 2.99µs 2.95µs 3.43µs 3.17µs 2.991µs 2.81µs]
  ⭐️ Times for flux-sample-4 to flux-sample-7: [3.83µs 2.65µs 3.76µs 3.09µs 3.6µs 3.29µs 2.38µs 3.66µs 3.35µs 2.97µs]
  ⭐️ Times for flux-sample-5 to flux-sample-7: [3.58µs 4.73µs 4.58µs 3.16µs 3.21µs 4.66µs 4.55µs 2.42µs 3.16µs 3.4µs]
  ⭐️ Times for flux-sample-6 to flux-sample-5: [2.85µs 3.44µs 2.48µs 2.4µs 2.88µs 2.6µs 2.98µs 9.98µs 2.96µs 3.08µs]
  ⭐️ Times for flux-sample-7 to flux-sample-4: [2.87µs 2.57µs 2.531µs 3.53µs 3.14µs 2.56µs 2.58µs 2.74µs 2.39µs 3.16µs]
  ⭐️ Times for flux-sample-6 to flux-sample-7: [2.96µs 2.86µs 3.01µs 2.89µs 3.02µs 2.4µs 2.911µs 2.81µs 2.86µs 2.96µs]
  ⭐️ Times for flux-sample-7 to flux-sample-5: [2.88µs 2.6µs 2.88µs 2.68µs 2.48µs 3.81µs 3.27µs 14.84µs 3.66µs 3.13µs]
  ⭐️ Times for flux-sample-7 to flux-sample-6: [3.41µs 2.64µs 3.11µs 3.23µs 3.09µs 3.1µs 2.65µs 3.689µs 3.13µs 3.83µs]
```

Yes, they are running on different physical nodes:

```bash
$ kubectl get pods -o wide
```
```console
NAME                  READY   STATUS      RESTARTS   AGE    IP           NODE                                       
flux-sample-0-tzc6n   0/1     Completed   0          114s   10.64.4.4    gke-test-cluster-default-pool-1bf80ee1-x73j
flux-sample-1-74hg7   0/1     Completed   0          114s   10.64.6.4    gke-test-cluster-default-pool-1bf80ee1-6tv6
flux-sample-2-4mxf8   0/1     Completed   0          114s   10.64.2.4    gke-test-cluster-default-pool-1bf80ee1-9xst
flux-sample-3-m4ks9   0/1     Completed   0          113s   10.64.5.9    gke-test-cluster-default-pool-1bf80ee1-676j
flux-sample-4-grstb   0/1     Completed   0          113s   10.64.1.5    gke-test-cluster-default-pool-1bf80ee1-qwrl
flux-sample-5-8djxs   0/1     Completed   0          113s   10.64.3.4    gke-test-cluster-default-pool-1bf80ee1-6jng
flux-sample-6-67fdr   0/1     Completed   0          113s   10.64.7.5    gke-test-cluster-default-pool-1bf80ee1-h9sl
flux-sample-7-7w2ds   0/1     Completed   0          113s   10.64.0.16   gke-test-cluster-default-pool-1bf80ee1-n0pp
```

That should be a matrix size of times minus the diagonal (the process to itself, which we don't measure).
When you are done, clean up:

```bash
gcloud container clusters delete test-cluster --region us-central1-a
```