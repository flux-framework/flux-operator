# PodSpec

PodSpec controlls variables for the cluster pod

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **Dict[str, str]** | Annotations for each pod | [optional] 
**automount_service_account_token** | **bool** | Automatically mount the service account name | [optional] 
**host_ipc** | **bool** | Security | [optional] 
**host_network** | **bool** |  | [optional] 
**labels** | **Dict[str, str]** | Labels for each pod | [optional] 
**node_selector** | **Dict[str, str]** | NodeSelectors for a pod | [optional] 
**resources** | [**Dict[str, IntOrString]**](IntOrString.md) | Resources include limits and requests | [optional] 
**restart_policy** | **str** | Restart Policy | [optional] 
**runtime_class_name** | **str** | RuntimeClassName for the pod | [optional] 
**scheduler_name** | **str** | Scheduler name for the pod | [optional] 
**service_account_name** | **str** | Service account name for the pod | [optional] 

## Example

```python
from fluxoperator.models.pod_spec import PodSpec

# TODO update the JSON string below
json = "{}"
# create an instance of PodSpec from a JSON string
pod_spec_instance = PodSpec.from_json(json)
# print the JSON string representation of the object
print(PodSpec.to_json())

# convert the object into a dict
pod_spec_dict = pod_spec_instance.to_dict()
# create an instance of PodSpec from a dict
pod_spec_from_dict = PodSpec.from_dict(pod_spec_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


