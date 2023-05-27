# FluxSpec


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**connect_timeout** | **str** | Single user executable to provide to flux start | [optional] [default to '5s']
**connection** | **str** | Connect to this job in the same namespace (akin to BootServer but within cluster) | [optional] 
**connection_size** | **int** | Additional number of nodes to allow from external boot-server This currently only allows local MiniCluster but could be extended to any general URI | [optional] 
**install_root** | **str** | Install root location | [optional] [default to '/usr']
**log_level** | **int** | Log level to use for flux logging (only in non TestMode) | [optional] [default to 6]
**option_flags** | **str** | Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \&quot;runFlux\&quot; true | [optional] [default to '']
**wrap** | **str** | Commands for flux start --wrap | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


