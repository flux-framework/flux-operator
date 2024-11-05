# FluxBroker

A FluxBroker defines a broker for flux

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**address** | **str** | Lead broker address (ip or hostname) | [default to '']
**name** | **str** | We need the name of the lead job to assemble the hostnames | [default to '']
**port** | **int** | Lead broker port - should only be used for external cluster | [optional] [default to 8050]
**size** | **int** | Lead broker size | [default to 0]

## Example

```python
from fluxoperator.models.flux_broker import FluxBroker

# TODO update the JSON string below
json = "{}"
# create an instance of FluxBroker from a JSON string
flux_broker_instance = FluxBroker.from_json(json)
# print the JSON string representation of the object
print(FluxBroker.to_json())

# convert the object into a dict
flux_broker_dict = flux_broker_instance.to_dict()
# create an instance of FluxBroker from a dict
flux_broker_from_dict = FluxBroker.from_dict(flux_broker_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


