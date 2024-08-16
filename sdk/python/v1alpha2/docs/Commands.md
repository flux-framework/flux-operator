# Commands


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**broker_pre** | **str** | A single command for only the broker to run | [optional] [default to '']
**init** | **str** | init command is run before anything | [optional] [default to '']
**post** | **str** | post command is run in the entrypoint when the broker exits / finishes | [optional] [default to '']
**pre** | **str** | pre command is run after global PreCommand, after asFlux is set (can override) | [optional] [default to '']
**prefix** | **str** | Prefix to flux start / submit / broker Typically used for a wrapper command to mount, etc. | [optional] [default to '']
**script** | **str** | Custom script for submit (e.g., multiple lines) | [optional] [default to '']
**service_pre** | **str** | A command only for service start.sh tor run | [optional] [default to '']
**worker_pre** | **str** | A command only for workers to run | [optional] [default to '']

## Example

```python
from fluxoperator.models.commands import Commands

# TODO update the JSON string below
json = "{}"
# create an instance of Commands from a JSON string
commands_instance = Commands.from_json(json)
# print the JSON string representation of the object
print(Commands.to_json())

# convert the object into a dict
commands_dict = commands_instance.to_dict()
# create an instance of Commands from a dict
commands_from_dict = Commands.from_dict(commands_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


