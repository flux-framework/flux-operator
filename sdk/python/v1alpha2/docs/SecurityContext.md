# SecurityContext


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**add_capabilities** | **List[str]** | Capabilities to add | [optional] 
**privileged** | **bool** | Privileged container | [optional] 

## Example

```python
from fluxoperator.models.security_context import SecurityContext

# TODO update the JSON string below
json = "{}"
# create an instance of SecurityContext from a JSON string
security_context_instance = SecurityContext.from_json(json)
# print the JSON string representation of the object
print(SecurityContext.to_json())

# convert the object into a dict
security_context_dict = security_context_instance.to_dict()
# create an instance of SecurityContext from a dict
security_context_from_dict = SecurityContext.from_dict(security_context_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


