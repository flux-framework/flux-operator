# MiniClusterArchive


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**path** | **str** | Save or load from this directory path | [optional] 

## Example

```python
from fluxoperator.models.mini_cluster_archive import MiniClusterArchive

# TODO update the JSON string below
json = "{}"
# create an instance of MiniClusterArchive from a JSON string
mini_cluster_archive_instance = MiniClusterArchive.from_json(json)
# print the JSON string representation of the object
print(MiniClusterArchive.to_json())

# convert the object into a dict
mini_cluster_archive_dict = mini_cluster_archive_instance.to_dict()
# create an instance of MiniClusterArchive from a dict
mini_cluster_archive_from_dict = MiniClusterArchive.from_dict(mini_cluster_archive_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


