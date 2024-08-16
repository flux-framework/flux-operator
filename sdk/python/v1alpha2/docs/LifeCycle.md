# LifeCycle


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**post_start_exec** | **str** |  | [optional] [default to '']
**pre_stop_exec** | **str** |  | [optional] [default to '']

## Example

```python
from fluxoperator.models.life_cycle import LifeCycle

# TODO update the JSON string below
json = "{}"
# create an instance of LifeCycle from a JSON string
life_cycle_instance = LifeCycle.from_json(json)
# print the JSON string representation of the object
print(LifeCycle.to_json())

# convert the object into a dict
life_cycle_dict = life_cycle_instance.to_dict()
# create an instance of LifeCycle from a dict
life_cycle_from_dict = LifeCycle.from_dict(life_cycle_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


