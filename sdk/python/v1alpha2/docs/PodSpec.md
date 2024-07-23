# PodSpec

PodSpec controlls variables for the cluster pod

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict[str, str]** | Annotations for each pod | [optional] 
**automount_service_account_token** | **bool** | Automatically mount the service account name | [optional] 
**labels** | **dict[str, str]** | Labels for each pod | [optional] 
**node_selector** | **dict[str, str]** | NodeSelectors for a pod | [optional] 
**resources** | [**dict[str, IntOrString]**](IntOrString.md) | Resources include limits and requests | [optional] 
**restart_policy** | **str** | Restart Policy | [optional] 
**runtime_class_name** | **str** | RuntimeClassName for the pod | [optional] 
**scheduler_name** | **str** | Scheduler name for the pod | [optional] 
**service_account_name** | **str** | Service account name for the pod | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


