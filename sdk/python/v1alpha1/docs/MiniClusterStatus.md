# MiniClusterStatus

MiniClusterStatus defines the observed state of Flux

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**completed** | **bool** | Label to indicate that job is completed, comes from Job.Completed The user can also look at conditions -&gt; JobFinished | [default to False]
**conditions** | [**list[V1Condition]**](V1Condition.md) | conditions hold the latest Flux Job and MiniCluster states | [optional] 
**jobid** | **str** | The Jobid is set internally to associate to a miniCluster This isn&#39;t currently in use, we only have one! | [default to '']
**maximum_size** | **int** | We keep the original size of the MiniCluster request as this is the absolute maximum | [default to 0]
**selector** | **str** |  | [default to '']
**size** | **int** | These are for the sub-resource scale functionality | [default to 0]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


