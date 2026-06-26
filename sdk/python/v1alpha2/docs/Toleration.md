# Toleration


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**effect** | **str** | The effect to have | [optional] 
**key** | **str** | The label key to tolerate | [optional] 
**operator** | **str** | The operator to use (e.g., Equal) | [optional] 
**value** | **str** | E.g., NoSchedule | [optional] 

## Example

```python
from fluxoperator.models.toleration import Toleration

# TODO update the JSON string below
json = "{}"
# create an instance of Toleration from a JSON string
toleration_instance = Toleration.from_json(json)
# print the JSON string representation of the object
print(Toleration.to_json())

# convert the object into a dict
toleration_dict = toleration_instance.to_dict()
# create an instance of Toleration from a dict
toleration_from_dict = Toleration.from_dict(toleration_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


