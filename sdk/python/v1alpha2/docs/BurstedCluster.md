# BurstedCluster


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** | The hostnames for the bursted clusters If set, the user is responsible for ensuring uniqueness. The operator will set to burst-N | [optional] [default to '']
**size** | **int** | Size of bursted cluster. Defaults to same size as local minicluster if not set | [optional] 

## Example

```python
from fluxoperator.models.bursted_cluster import BurstedCluster

# TODO update the JSON string below
json = "{}"
# create an instance of BurstedCluster from a JSON string
bursted_cluster_instance = BurstedCluster.from_json(json)
# print the JSON string representation of the object
print(BurstedCluster.to_json())

# convert the object into a dict
bursted_cluster_dict = bursted_cluster_instance.to_dict()
# create an instance of BurstedCluster from a dict
bursted_cluster_from_dict = BurstedCluster.from_dict(bursted_cluster_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


