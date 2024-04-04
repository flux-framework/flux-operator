# ZeroMQ in Go Examples

> Pair to Pair with DEALER to ROUTER

Note that we are going to try to use the [DEALER to ROUTER](https://zguide.zeromq.org/docs/chapter3/#The-DEALER-to-ROUTER-Combination) design here.
I implemented this in two ways, the current approach in [main.go.txt](main.go.txt) that builds into the container, and an (opposite) design
in [main.go.v1](main.go.v1).

Create the kind cluster.

```bash
kind create cluster --config ./kind-config.yaml
```

Install the flux operator

```bash
kubectl apply -f ../../dist/flux-operator.yaml
```

## Local Test

You can automate the entire thing. Note that this first example has `--raw` added to the entrypoint to print the raw times.

```bash
./build.sh
```

And then get logs:

```console
Defaulted container "flux-sample" out of: flux-sample, flux-view (init)
Hello I'm host flux-sample-0
Hello I'm host flux-sample-3
Hello I'm host flux-sample-2
Hello I'm host flux-sample-1
  ⭐️ Times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: [2.022144ms 85.646µs 71.651µs 67.255µs 56.497µs 69.762µs 64.737µs 64.719µs 64.475µs 50.268µs]
  ⭐️ Times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: [187.388415ms 98.255µs 77.949µs 79.8µs 77.135µs 79.047µs 75.378µs 65.81µs 72.987µs 75.554µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: [1.910568ms 97.721µs 72.368µs 70.992µs 122.607µs 67.052µs 69.333µs 72.959µs 55.936µs 54.03µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: [2.766664ms 104.953µs 71.581µs 68.432µs 69.385µs 50.691µs 72.501µs 68.229µs 65.689µs 127.542µs]
  ⭐️ Times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: [312.753192ms 101.437µs 82.906µs 78.605µs 77.598µs 78.379µs 79.073µs 79.202µs 82.471µs 79.849µs]
  ⭐️ Times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: [1.627328ms 180.601µs 97.693µs 64.498µs 63.371µs 76.179µs 56.385µs 58.304µs 60.835µs 74.159µs]
  ⭐️ Times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: [1.748117ms 79.71µs 89.488µs 95.987µs 107.764µs 88.422µs 120.027µs 69.187µs 151.28µs 194.752µs]
  ⭐️ Times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: [2.39201ms 199.99µs 182.462µs 76.48µs 55.165µs 53.682µs 56.092µs 51.916µs 45.052µs 51.173µs]
  ⭐️ Times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: [532.644951ms 55.401µs 60.862µs 52.268µs 44.262µs 47.538µs 46.047µs 43.696µs 46.358µs 43.545µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: [127.942516ms 238.537µs 248.172µs 241.167µs 220.343µs 183.445µs 203.348µs 227.755µs 178.029µs 206.332µs]
  ⭐️ Times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: [451.017582ms 130.585µs 128.069µs 126.802µs 137.073µs 116.022µs 113.394µs 109.567µs 90.688µs 118.078µs]
  ⭐️ Times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: [661.313903ms 236.399µs 188.52µs 167.761µs 252.685µs 187.253µs 173.82µs 146.113µs 154.547µs 158.902µ
```

Note that we have 12 groups of 10 times that represent a matrix minus the diagonal, which would be a node to itself (which we don't record). So the above is 12 groups, which is 4x4 == 16 minus the diagonal of 4.
If you omit `--raw`, you'll get a matrix of mean times (over N=10 measurements).

```console
Defaulted container "flux-sample" out of: flux-sample, flux-view (init)
Hello I'm host flux-sample-0
Hello I'm host flux-sample-2
Hello I'm host flux-sample-3
Hello I'm host flux-sample-1
  ⭐️ Mean times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: 214.183µs
  ⭐️ Mean times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: 25.595208ms
  ⭐️ Mean times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: 324.54µs
  ⭐️ Mean times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: 324.835µs
  ⭐️ Mean times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: 32.615343ms
  ⭐️ Mean times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: 12.847121ms
  ⭐️ Mean times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: 252.283µs
  ⭐️ Mean times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: 239.2µs
  ⭐️ Mean times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: 279.983µs
  ⭐️ Mean times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: 12.813611ms
  ⭐️ Mean times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: 49.268451ms
  ⭐️ Mean times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: 33.067423ms
```

You can look at [build.sh](build.sh) for the build steps, and [entrypoint.sh](entrypoint.sh) for the start command,
and [main.go](main.go.txt) for the defaults and logic (we have to rename to .txt so the flux operator doesn't include
in its build).

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
Defaulted container "flux-sample" out of: flux-sample, flux-view (init)
Hello I'm host flux-sample-7
Hello I'm host flux-sample-3
Hello I'm host flux-sample-4
Hello I'm host flux-sample-6
Hello I'm host flux-sample-5
Hello I'm host flux-sample-2
Hello I'm host flux-sample-1
Hello I'm host flux-sample-0
  ⭐️ Times for 10 messages flux-sample-0.flux-service.default.svc.cluster.local:5555 to flux-sample-2.flux-service.default.svc.cluster.local:5555: [2.27578ms 124.99µs 101.27µs 114.43µs 99.71µs 97.76µs 112.5µs 99.08µs 100.39µs 106.33µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-0.flux-service.default.svc.cluster.local:5555: [509.558548ms 109.25µs 104.63µs 116.06µs 119.81µs 121.51µs 118.59µs 118.88µs 116.48µs 127.89µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-1.flux-service.default.svc.cluster.local:5555: [1.088216243s 126.44µs 121.869µs 113.149µs 114.86µs 102.871µs 95.871µs 94.16µs 94.6µs 94.609µs]
  ⭐️ Times for 10 messages flux-sample-1.flux-service.default.svc.cluster.local:5555 to flux-sample-7.flux-service.default.svc.cluster.local:5555: [291.263625ms 116.05µs 107.22µs 112.98µs 107.62µs 127.05µs 105.86µs 104.71µs 104.38µs 105.14µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-3.flux-service.default.svc.cluster.local:5555: [1.37243942s 119.05µs 104.24µs 100.7µs 98.49µs 92.93µs 103.75µs 95.28µs 101.65µs 97.46µs]
  ⭐️ Times for 10 messages flux-sample-3.flux-service.default.svc.cluster.local:5555 to flux-sample-6.flux-service.default.svc.cluster.local:5555: [188.336356ms 113.54µs 131.38µs 108.73µs 109.1µs 96.92µs 93.06µs 111.22µs 92.21µs 95.73µs]
  ⭐️ Times for 10 messages flux-sample-6.flux-service.default.svc.cluster.local:5555 to flux-sample-5.flux-service.default.svc.cluster.local:5555: [617.644567ms 122.01µs 131.17µs 108.52µs 105.43µs 111.13µs 111.14µs 108.96µs 107.06µs 109.07µs]
  ⭐️ Times for 10 messages flux-sample-2.flux-service.default.svc.cluster.local:5555 to flux-sample-4.flux-service.default.svc.cluster.local:5555: [2.335518399s 125.58µs 109.42µs 104.43µs 110.44µs 107.31µs 105.26µs 114.19µs 113.07µs 106.36µs]
...
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