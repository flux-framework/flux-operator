# ApiV1alpha1MiniClusterVolume

Mini Cluster local volumes available to mount (these are on the host)

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**path** | **str** |  | defaults to ""
**annotations** | **{str: (str,)}** | Annotations for persistent volume claim | [optional] 
**capacity** | **str** | Capacity (string) for PVC (storage request) to create PV | [optional]  if omitted the server will use the default value of ""
**_class** | **str** |  | [optional]  if omitted the server will use the default value of ""
**labels** | **{str: (str,)}** |  | [optional] 
**secret** | **str** | Secret reference in Kubernetes with service account role | [optional]  if omitted the server will use the default value of ""
**secret_namespace** | **str** | Secret namespace | [optional]  if omitted the server will use the default value of ""

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


