# PodSpec

PodSpec controlls variables for the cluster pod

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **Dict[str, str]** | Annotations for each pod | [optional] 
**automount_service_account_token** | **bool** | Automatically mount the service account name | [optional] 
**dns_policy** | **str** | Pod DNS policy (defaults to ClusterFirst) | [optional] 
**host_ipc** | **bool** | Use Host IPC | [optional] 
**host_pid** | **bool** | Use Host PID | [optional] 
**labels** | **Dict[str, str]** | Labels for each pod | [optional] 
**node_affinity** | **Dict[str, List[str]]** | NodeAffinity is for a list of values assoicated with a label | [optional] 
**node_selector** | **Dict[str, str]** | NodeSelectors for a pod | [optional] 
**resources** | [**Dict[str, K8sIoApimachineryPkgApiResourceQuantity]**](K8sIoApimachineryPkgApiResourceQuantity.md) | Resources include limits and requests | [optional] 
**restart_policy** | **str** | Restart Policy | [optional] 
**runtime_class_name** | **str** | RuntimeClassName for the pod | [optional] 
**scheduler_name** | **str** | Scheduler name for the pod | [optional] 
**security_context** | [**K8sIoApiCoreV1PodSecurityContext**](K8sIoApiCoreV1PodSecurityContext.md) |  | [optional] 
**service_account_name** | **str** | Service account name for the pod | [optional] 
**tolerations** | [**List[Toleration]**](Toleration.md) | Tolerations for a pod | [optional] 

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


