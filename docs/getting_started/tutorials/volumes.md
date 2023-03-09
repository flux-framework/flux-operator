# Volumes

These short examples will describe advanced functionality for volumes. For examples,
see our [storage examples directory](https://github.com/flux-framework/flux-operator/tree/main/examples/storage).

## Existing Persistent Volume

It might be the case that you've already defined a persistent volume, and you simply want
to ask for claims (for your pods). And technically, this could work for the PVC too.
Either way, as long as you create your object with the matching:

1. namespace for metadata (for a PVC)
2. name for metadata (e.g., flux-sample)
3. name for the volume (e.g., data)

The operator should detect that it already exists, and not try to re-create it, but
still use the names. As an example, here is a complete PVC and PV created by the operator
for a simple volume (in the MiniCluster CRD) named "data." Here is an example
creating a PVC that is a driver type (Flex) that is deprecated that the operator does
not support:

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:

  # IMPORTANT: this needs to correspond with your volume name "e.g., data" in the minicluster
  name: data
  annotations:
    ibm.io/auto-create-bucket: "false"
    ibm.io/auto-delete-bucket: "false"
    ibm.io/bucket: flux-operator-storage
    ibm.io/chunk-size-mb: "40"
    ibm.io/cos-service: ""
    ibm.io/debug-level: warn
    # ibm.io/iam-endpoint: https://private.iam.cloud.ibm.com
    ibm.io/iam-endpoint: https://iam.cloud.ibm.com
    ibm.io/multireq-max: "20"
    ibm.io/parallel-count: "20"
    ibm.io/s3fs-fuse-retry-count: "5"
    ibm.io/secret-name: s3-secret
    ibm.io/secret-namespace: flux-operator
    ibm.io/stat-cache-size: "100000"
    pv.kubernetes.io/provisioned-by: ibm.io/ibmc-s3fs

spec:
  accessModes:
  - ReadWriteMany
  capacity:
    storage: 25Gi
  claimRef:
    apiVersion: v1
    kind: PersistentVolumeClaim

    # IMPORTANT: this must be the name of the volume plus "-claim"
    name: data-claim
    namespace: flux-operator

  flexVolume:
    driver: ibm/ibmc-s3fs
    options:
      access-mode: ReadWriteMany
      bucket: flux-operator-storage
      chunk-size-mb: "40"
      curl-debug: "false"
      debug-level: warn
      iam-endpoint: https://iam.cloud.ibm.com
      kernel-cache: "true"
      multireq-max: "20"
      object-store-storage-class: us-east-standard
      parallel-count: "20"
      s3fs-fuse-retry-count: "5"
      stat-cache-size: "100000"
      tls-cipher-suite: AESGCM
    secretRef:
      name: s3-secret
      namespace: flux-operator
  persistentVolumeReclaimPolicy: Delete
  storageClassName: ibmc-s3fs-standard-regional
  volumeMode: Filesystem
```

Note that we've given it the name from our minicluster volume "data" and the claim name is always the volume 
name + "-claim" (e.g., "data-claim") and the namespace and name for the objects match the one for our MiniCluster. 