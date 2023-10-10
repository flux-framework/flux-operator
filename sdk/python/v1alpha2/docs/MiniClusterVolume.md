# MiniClusterVolume

Mini Cluster local volumes available to mount (these are on the host)

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict[str, str]** | Annotations for the volume | [optional] 
**attributes** | **dict[str, str]** | Optional volume attributes | [optional] 
**capacity** | **str** | Capacity (string) for PVC (storage request) to create PV | [optional] [default to '5Gi']
**claim_annotations** | **dict[str, str]** | Annotations for the persistent volume claim | [optional] 
**delete** | **bool** | Delete the persistent volume on cleanup | [optional] [default to True]
**driver** | **str** | Storage driver, e.g., gcs.csi.ofek.dev Only needed if not using hostpath | [optional] [default to '']
**labels** | **dict[str, str]** |  | [optional] 
**path** | **str** |  | [default to '']
**secret** | **str** | Secret reference in Kubernetes with service account role | [optional] [default to '']
**secret_namespace** | **str** | Secret namespace | [optional] [default to 'default']
**storage_class** | **str** |  | [optional] [default to 'hostpath']
**volume_handle** | **str** | Volume handle, falls back to storage class name if not defined | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


