# FluxSpec


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**broker_config** | **str** | Optionally provide a manually created broker config this is intended for bursting to remote clusters | [optional] [default to '']
**connect_timeout** | **str** | Single user executable to provide to flux start | [optional] [default to '5s']
**curve_cert** | **str** | Optionally provide an already existing curve certificate this is intended for bursting to remote clusters | [optional] [default to '']
**install_root** | **str** | Install root location | [optional] [default to '/usr']
**lead_broker** | [**FluxBroker**](FluxBroker.md) |  | [optional] 
**log_level** | **int** | Log level to use for flux logging (only in non TestMode) | [optional] [default to 6]
**minimal_service** | **bool** | Only expose the broker service (to reduce load on DNS) | [optional] [default to False]
**munge_key** | **str** | Optionally provide an already existing munge key this is intended for bursting to remote clusters. Assumed to be at /etc/munge/munge.key | [optional] [default to '']
**option_flags** | **str** | Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \&quot;runFlux\&quot; true | [optional] [default to '']
**wrap** | **str** | Commands for flux start --wrap | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


