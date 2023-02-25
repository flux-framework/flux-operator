# ApiV1alpha1MiniClusterSpec


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**containers** | [**[ApiV1alpha1MiniClusterContainer]**](ApiV1alpha1MiniClusterContainer.md) | Containers is one or more containers to be created in a pod. There should only be one container to run flux with runFlux | 
**cleanup** | **bool** | Cleanup the pods and storage when the index broker pod is complete | [optional]  if omitted the server will use the default value of False
**deadline_seconds** | **int** | Should the job be limited to a particular number of seconds? Approximately one year. This cannot be zero or job won&#39;t start | [optional]  if omitted the server will use the default value of 0
**flux_restful** | [**ApiV1alpha1FluxRestful**](ApiV1alpha1FluxRestful.md) |  | [optional] 
**job_labels** | **{str: (str,)}** | Labels for the job | [optional] 
**logging** | [**ApiV1alpha1LoggingSpec**](ApiV1alpha1LoggingSpec.md) |  | [optional] 
**pod** | [**ApiV1alpha1PodSpec**](ApiV1alpha1PodSpec.md) |  | [optional] 
**size** | **int** | Size (number of job pods to run, size of minicluster in pods) | [optional]  if omitted the server will use the default value of 0
**tasks** | **int** | Total number of CPUs being run across entire cluster | [optional]  if omitted the server will use the default value of 0
**users** | [**[ApiV1alpha1MiniClusterUser]**](ApiV1alpha1MiniClusterUser.md) | Users of the MiniCluster | [optional] 
**volumes** | [**{str: (ApiV1alpha1MiniClusterVolume,)}**](ApiV1alpha1MiniClusterVolume.md) | Volumes accessible to containers from a host | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


