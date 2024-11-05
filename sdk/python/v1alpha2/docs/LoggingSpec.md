# LoggingSpec


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**debug** | **bool** | Debug mode adds extra verbosity to Flux | [optional] [default to False]
**quiet** | **bool** | Quiet mode silences all output so the job only shows the test running | [optional] [default to False]
**strict** | **bool** | Strict mode ensures any failure will not continue in the job entrypoint | [optional] [default to False]
**timed** | **bool** | Timed mode adds timing to Flux commands | [optional] [default to False]
**zeromq** | **bool** | Enable Zeromq logging | [optional] [default to False]

## Example

```python
from fluxoperator.models.logging_spec import LoggingSpec

# TODO update the JSON string below
json = "{}"
# create an instance of LoggingSpec from a JSON string
logging_spec_instance = LoggingSpec.from_json(json)
# print the JSON string representation of the object
print(LoggingSpec.to_json())

# convert the object into a dict
logging_spec_dict = logging_spec_instance.to_dict()
# create an instance of LoggingSpec from a dict
logging_spec_from_dict = LoggingSpec.from_dict(logging_spec_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


