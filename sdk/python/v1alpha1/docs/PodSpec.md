# PodSpec

PodSpec controlls variables for the cluster pod

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**annotations** | **dict(str, str)** | Annotations for each pod | [optional] 
**labels** | **dict(str, str)** | Labels for each pod | [optional] 
**node_selector** | **dict(str, str)** | NodeSelectors for a pod | [optional] 
**resources** | [**dict(str, IntOrString)**](IntOrString.md) | Resources include limits and requests | [optional] 
**service_account_name** | **str** | Service account name for the pod | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


