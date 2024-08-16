# Network


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**disable_affinity** | **bool** | Disable affinity rules that guarantee one network address / node | [optional] 
**headless_name** | **str** | Name for cluster headless service | [optional] [default to 'flux-service']

## Example

```python
from fluxoperator.models.network import Network

# TODO update the JSON string below
json = "{}"
# create an instance of Network from a JSON string
network_instance = Network.from_json(json)
# print the JSON string representation of the object
print(Network.to_json())

# convert the object into a dict
network_dict = network_instance.to_dict()
# create an instance of Network from a dict
network_from_dict = Network.from_dict(network_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


