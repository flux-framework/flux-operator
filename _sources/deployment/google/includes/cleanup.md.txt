## Clean up

Whatever tutorial you choose, don't forget to clean up at the end!
You can optionally undeploy the operator (this is again at the root of the operator repository clone)

```bash
$ make undeploy
```

Or the file you used to deploy it:

```bash
$ kubectl delete -f examples/dist/flux-operator.yaml
$ kubectl delete -f examples/dist/flux-operator-dev.yaml
```

And then to delete the cluster with gcloud:

```bash
$ gcloud container clusters delete --zone us-central1-a flux-cluster
```

I like to check in the cloud console to ensure that it was actually deleted.

