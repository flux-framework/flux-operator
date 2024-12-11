# PodSecurityContext


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sysctls** | **Dict[str, str]** | Sysctls | [optional] 

## Example

```python
from fluxoperator.models.pod_security_context import PodSecurityContext

# TODO update the JSON string below
json = "{}"
# create an instance of PodSecurityContext from a JSON string
pod_security_context_instance = PodSecurityContext.from_json(json)
# print the JSON string representation of the object
print(PodSecurityContext.to_json())

# convert the object into a dict
pod_security_context_dict = pod_security_context_instance.to_dict()
# create an instance of PodSecurityContext from a dict
pod_security_context_from_dict = PodSecurityContext.from_dict(pod_security_context_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


