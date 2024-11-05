# ContainerResources

ContainerResources include limits and requests

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**limits** | [**Dict[str, IntOrString]**](IntOrString.md) |  | [optional] 
**requests** | [**Dict[str, IntOrString]**](IntOrString.md) |  | [optional] 

## Example

```python
from fluxoperator.models.container_resources import ContainerResources

# TODO update the JSON string below
json = "{}"
# create an instance of ContainerResources from a JSON string
container_resources_instance = ContainerResources.from_json(json)
# print the JSON string representation of the object
print(ContainerResources.to_json())

# convert the object into a dict
container_resources_dict = container_resources_instance.to_dict()
# create an instance of ContainerResources from a dict
container_resources_from_dict = ContainerResources.from_dict(container_resources_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


