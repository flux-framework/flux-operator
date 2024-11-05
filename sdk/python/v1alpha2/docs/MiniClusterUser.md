# MiniClusterUser


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **str** | If a user is defined, the username is required | [default to '']
**password** | **str** |  | [optional] [default to '']

## Example

```python
from fluxoperator.models.mini_cluster_user import MiniClusterUser

# TODO update the JSON string below
json = "{}"
# create an instance of MiniClusterUser from a JSON string
mini_cluster_user_instance = MiniClusterUser.from_json(json)
# print the JSON string representation of the object
print(MiniClusterUser.to_json())

# convert the object into a dict
mini_cluster_user_dict = mini_cluster_user_instance.to_dict()
# create an instance of MiniClusterUser from a dict
mini_cluster_user_from_dict = MiniClusterUser.from_dict(mini_cluster_user_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


