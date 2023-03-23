# Commands


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**init** | **str** | init command is run before anything | [optional] [default to '']
**post** | **str** | post command is run in the entrypoint when the broker exits / finishes | [optional] [default to '']
**pre** | **str** | pre command is run after global PreCommand, after asFlux is set (can override) | [optional] [default to '']
**prefix** | **str** | Prefix to flux start / submit / broker Typically used for a wrapper command to mount, etc. | [optional] [default to '']
**run_flux_as_root** | **bool** | Run flux start as root - required for some storage binds | [optional] [default to False]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


