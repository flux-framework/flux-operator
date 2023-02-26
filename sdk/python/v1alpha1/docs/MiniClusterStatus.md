# MiniClusterStatus

MiniClusterStatus defines the observed state of Flux

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**conditions** | [**list[V1Condition]**](V1Condition.md) | conditions hold the latest Flux Job and MiniCluster states | [optional]
**jobid** | **str** | The Jobid is set internally to associate to a miniCluster This isn&#39;t currently in use, we only have one! | [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
