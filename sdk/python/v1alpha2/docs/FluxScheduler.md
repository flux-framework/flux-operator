# FluxScheduler

FluxScheduler attributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**queue_policy** | **str** | Scheduler queue policy, defaults to \&quot;fcfs\&quot; can also be \&quot;easy\&quot; | [optional] [default to '']
**simple** | **bool** | Use sched-simple (no support for GPU) | [optional] [default to False]

## Example

```python
from fluxoperator.models.flux_scheduler import FluxScheduler

# TODO update the JSON string below
json = "{}"
# create an instance of FluxScheduler from a JSON string
flux_scheduler_instance = FluxScheduler.from_json(json)
# print the JSON string representation of the object
print(FluxScheduler.to_json())

# convert the object into a dict
flux_scheduler_dict = flux_scheduler_instance.to_dict()
# create an instance of FluxScheduler from a dict
flux_scheduler_from_dict = FluxScheduler.from_dict(flux_scheduler_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


