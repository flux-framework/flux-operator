# MiniClusterSpec

MiniCluster is an HPC cluster in Kubernetes you can control Either to submit a single job (and go away) or for a persistent single- or multi- user cluster

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**archive** | [**MiniClusterArchive**](MiniClusterArchive.md) |  | [optional] 
**cleanup** | **bool** | Cleanup the pods and storage when the index broker pod is complete | [optional] [default to False]
**containers** | [**List[MiniClusterContainer]**](MiniClusterContainer.md) | Containers is one or more containers to be created in a pod. There should only be one container to run flux with runFlux | 
**deadline_seconds** | **int** | Should the job be limited to a particular number of seconds? Approximately one year. This cannot be zero or job won&#39;t start | [optional] [default to 31500000]
**flux** | [**FluxSpec**](FluxSpec.md) |  | [optional] 
**interactive** | **bool** | Run a single-user, interactive minicluster | [optional] [default to False]
**job_labels** | **Dict[str, str]** | Labels for the job | [optional] 
**logging** | [**LoggingSpec**](LoggingSpec.md) |  | [optional] 
**max_size** | **int** | MaxSize (maximum number of pods to allow scaling to) | [optional] 
**min_size** | **int** | MinSize (minimum number of pods that must be up for Flux) Note that this option does not edit the number of tasks, so a job could run with fewer (and then not start) | [optional] 
**network** | [**Network**](Network.md) |  | [optional] 
**pod** | [**PodSpec**](PodSpec.md) |  | [optional] 
**services** | [**List[MiniClusterContainer]**](MiniClusterContainer.md) | Services are one or more service containers to bring up alongside the MiniCluster. | [optional] 
**share_process_namespace** | **bool** | Share process namespace? | [optional] [default to False]
**size** | **int** | Size (number of job pods to run, size of minicluster in pods) This is also the minimum number required to start Flux | [optional] [default to 1]
**tasks** | **int** | Total number of CPUs being run across entire cluster | [optional] [default to 1]

## Example

```python
from fluxoperator.models.mini_cluster_spec import MiniClusterSpec

# TODO update the JSON string below
json = "{}"
# create an instance of MiniClusterSpec from a JSON string
mini_cluster_spec_instance = MiniClusterSpec.from_json(json)
# print the JSON string representation of the object
print(MiniClusterSpec.to_json())

# convert the object into a dict
mini_cluster_spec_dict = mini_cluster_spec_instance.to_dict()
# create an instance of MiniClusterSpec from a dict
mini_cluster_spec_from_dict = MiniClusterSpec.from_dict(mini_cluster_spec_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


