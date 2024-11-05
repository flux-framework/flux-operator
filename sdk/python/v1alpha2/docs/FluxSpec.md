# FluxSpec


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**arch** | **str** | Change the arch string - determines the binaries that are downloaded to run the entrypoint | [optional] 
**broker_config** | **str** | Optionally provide a manually created broker config this is intended for bursting to remote clusters | [optional] [default to '']
**bursting** | [**Bursting**](Bursting.md) |  | [optional] 
**complete_workers** | **bool** | Complete workers when they fail This is ideal if you don&#39;t want them to restart | [optional] [default to False]
**connect_timeout** | **str** | Single user executable to provide to flux start | [optional] [default to '5s']
**container** | [**FluxContainer**](FluxContainer.md) |  | [optional] 
**curve_cert** | **str** | Optionally provide an already existing curve certificate This is not recommended in favor of providing the secret name as curveCertSecret, below | [optional] [default to '']
**disable_socket** | **bool** | Disable specifying the socket path | [optional] [default to False]
**log_level** | **int** | Log level to use for flux logging (only in non TestMode) | [optional] [default to 6]
**minimal_service** | **bool** | Only expose the broker service (to reduce load on DNS) | [optional] [default to False]
**munge_secret** | **str** | Expect a secret (named according to this string) for a munge key. This is intended for bursting. Assumed to be at /etc/munge/munge.key This is binary data. | [optional] [default to '']
**no_wait_socket** | **bool** | Do not wait for the socket | [optional] [default to False]
**option_flags** | **str** | Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \&quot;runFlux\&quot; true | [optional] [default to '']
**scheduler** | [**FluxScheduler**](FluxScheduler.md) |  | [optional] 
**submit_command** | **str** | Modify flux submit to be something else | [optional] 
**topology** | **str** | Specify a custom Topology | [optional] [default to '']
**wrap** | **str** | Commands for flux start --wrap | [optional] 

## Example

```python
from fluxoperator.models.flux_spec import FluxSpec

# TODO update the JSON string below
json = "{}"
# create an instance of FluxSpec from a JSON string
flux_spec_instance = FluxSpec.from_json(json)
# print the JSON string representation of the object
print(FluxSpec.to_json())

# convert the object into a dict
flux_spec_dict = flux_spec_instance.to_dict()
# create an instance of FluxSpec from a dict
flux_spec_from_dict = FluxSpec.from_dict(flux_spec_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


