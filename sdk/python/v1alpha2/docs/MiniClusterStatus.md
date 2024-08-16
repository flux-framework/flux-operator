# MiniClusterStatus

MiniClusterStatus defines the observed state of Flux

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**conditions** | [**List[V1Condition]**](V1Condition.md) | conditions hold the latest Flux Job and MiniCluster states | [optional] 
**jobid** | **str** | The Jobid is set internally to associate to a miniCluster This isn&#39;t currently in use, we only have one! | [default to '']
**maximum_size** | **int** | We keep the original size of the MiniCluster request as this is the absolute maximum | [default to 0]
**selector** | **str** |  | [default to '']
**size** | **int** | These are for the sub-resource scale functionality | [default to 0]

## Example

```python
from fluxoperator.models.mini_cluster_status import MiniClusterStatus

# TODO update the JSON string below
json = "{}"
# create an instance of MiniClusterStatus from a JSON string
mini_cluster_status_instance = MiniClusterStatus.from_json(json)
# print the JSON string representation of the object
print(MiniClusterStatus.to_json())

# convert the object into a dict
mini_cluster_status_dict = mini_cluster_status_instance.to_dict()
# create an instance of MiniClusterStatus from a dict
mini_cluster_status_from_dict = MiniClusterStatus.from_dict(mini_cluster_status_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


