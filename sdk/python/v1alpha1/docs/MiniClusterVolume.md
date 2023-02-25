# MiniClusterVolume

Mini Cluster local volumes available to mount (these are on the host)

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict(str, str)** | Annotations for persistent volume claim | [optional] 
**capacity** | **str** | Capacity (string) for PVC (storage request) to create PV | [optional] [default to '']
**_class** | **str** |  | [optional] [default to '']
**labels** | **dict(str, str)** |  | [optional] 
**path** | **str** |  | [default to '']
**secret** | **str** | Secret reference in Kubernetes with service account role | [optional] [default to '']
**secret_namespace** | **str** | Secret namespace | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


