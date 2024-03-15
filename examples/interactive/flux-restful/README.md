# Flux Restful

This example demonstrates installing and running a basic Flux Restful Server. You likely want to customize this further,
see [the repository](https://github.com/flux-framework/flux-restful-api) for details.

## Usage

After installing the flux operator and creating the cluster, create the minicluster:

```bash
kubectl apply -f minicluster.yaml
```

Ensure that it is running:

```bash
kubectl logs flux-sample-0-xxx -f
```
```console
🍓 Require auth: False
🍓   Secret key ********************************
🍓    Flux user: ****
🍓   Flux token: ****
INFO:     Started server process [307]
INFO:     Waiting for application startup.
INFO:     Application startup complete.
INFO:     Uvicorn running on http://0.0.0.0:5000 (Press CTRL+C to quit)
```

Now we can expose a service.

```bash
$ kubectl port-forward flux-sample-0-xjcx7 5000:5000
```

In a different terminal try interacting with it:

```bash
$ curl -ks http://localhost:5000/v1/jobs | jq
{
  "jobs": []
}
```

You can see other endpoints (and the tool we provide in Python) in the documentation [here](https://flux-framework.org/flux-restful-api/getting_started/api.html).
Note that this API could be available within the cluster as well. After shelling in:

```bash
dnf install -y jq
curl -ks flux-sample-0.flux-service.default.svc.cluster.local:5000/v1/jobs | jq
{
  "jobs": []
}
```