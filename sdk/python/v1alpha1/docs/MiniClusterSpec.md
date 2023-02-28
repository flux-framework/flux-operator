# MiniClusterSpec


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**cleanup** | **bool** | Cleanup the pods and storage when the index broker pod is complete | [optional] [default to False]
**containers** | [**list[MiniClusterContainer]**](MiniClusterContainer.md) | Containers is one or more containers to be created in a pod. There should only be one container to run flux with runFlux | 
**deadline_seconds** | **int** | Should the job be limited to a particular number of seconds? Approximately one year. This cannot be zero or job won&#39;t start | [optional] [default to 31500000]
**flux_restful** | [**FluxRestful**](FluxRestful.md) |  | [optional] 
**job_labels** | **dict(str, str)** | Labels for the job | [optional] 
**logging** | [**LoggingSpec**](LoggingSpec.md) |  | [optional] 
**pod** | [**PodSpec**](PodSpec.md) |  | [optional] 
**size** | **int** | Size (number of job pods to run, size of minicluster in pods) | [optional] [default to 1]
**tasks** | **int** | Total number of CPUs being run across entire cluster | [optional] [default to 1]
**users** | [**list[MiniClusterUser]**](MiniClusterUser.md) | Users of the MiniCluster | [optional] 
**volumes** | [**dict(str, MiniClusterVolume)**](MiniClusterVolume.md) | Volumes accessible to containers from a host | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


