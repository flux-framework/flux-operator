# MiniCluster

MiniCluster is the Schema for a Flux job launcher on K8s

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**api_version** | **str** | APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources | [optional] 
**kind** | **str** | Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds | [optional] 
**metadata** | [**V1ObjectMeta**](V1ObjectMeta.md) |  | [optional] 
**spec** | [**MiniClusterSpec**](MiniClusterSpec.md) |  | [optional] 
**status** | [**MiniClusterStatus**](MiniClusterStatus.md) |  | [optional] 

## Example

```python
from fluxoperator.models.mini_cluster import MiniCluster

# TODO update the JSON string below
json = "{}"
# create an instance of MiniCluster from a JSON string
mini_cluster_instance = MiniCluster.from_json(json)
# print the JSON string representation of the object
print(MiniCluster.to_json())

# convert the object into a dict
mini_cluster_dict = mini_cluster_instance.to_dict()
# create an instance of MiniCluster from a dict
mini_cluster_from_dict = MiniCluster.from_dict(mini_cluster_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


