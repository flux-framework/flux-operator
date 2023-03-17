# Commands


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**post** | **str** | post command is run by all pods on finish | [optional] [default to '']
**pre** | **str** | pre command is run after global PreCommand, before anything else | [optional] [default to '']
**prefix** | **str** | Prefix to flux start / submit / broker Typically used for a wrapper command to mount, etc. | [optional] [default to '']
**run_flux_as_root** | **bool** | Run flux start as root - required for some storage binds | [optional] [default to False]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


