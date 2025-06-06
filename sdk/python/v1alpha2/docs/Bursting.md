# Bursting

Bursting Config For simplicity, we internally handle the name of the job (hostnames)

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**clusters** | [**List[BurstedCluster]**](BurstedCluster.md) | External clusters to burst to. Each external cluster must share the same listing to align ranks | [optional] 
**hostlist** | **str** | Hostlist is a custom hostlist for the broker.toml that includes the local plus bursted cluster. This is typically used for bursting to another resource type, where we can predict the hostnames but they don&#39;t follow the same convention as the Flux Operator | [optional] [default to '']
**lead_broker** | [**FluxBroker**](FluxBroker.md) |  | [optional] 

## Example

```python
from fluxoperator.models.bursting import Bursting

# TODO update the JSON string below
json = "{}"
# create an instance of Bursting from a JSON string
bursting_instance = Bursting.from_json(json)
# print the JSON string representation of the object
print(Bursting.to_json())

# convert the object into a dict
bursting_dict = bursting_instance.to_dict()
# create an instance of Bursting from a dict
bursting_from_dict = Bursting.from_dict(bursting_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


