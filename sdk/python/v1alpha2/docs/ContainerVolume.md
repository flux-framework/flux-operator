# ContainerVolume


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**claim_name** | **str** | Claim name if the existing volume is a PVC | [optional] 
**config_map_name** | **str** | Config map name if the existing volume is a config map You should also define items if you are using this | [optional] 
**empty_dir** | **bool** |  | [optional] [default to False]
**empty_dir_medium** | **str** | Add an empty directory custom type | [optional] 
**empty_dir_size_limit** | **str** | Add an empty directory sizeLimit | [optional] 
**host_path** | **str** | An existing hostPath to bind to path | [optional] 
**items** | **Dict[str, str]** | Items (key and paths) for the config map | [optional] 
**path** | **str** | Path and claim name are always required if a secret isn&#39;t defined | [optional] 
**read_only** | **bool** |  | [optional] [default to False]
**secret_name** | **str** | An existing secret | [optional] 

## Example

```python
from fluxoperator.models.container_volume import ContainerVolume

# TODO update the JSON string below
json = "{}"
# create an instance of ContainerVolume from a JSON string
container_volume_instance = ContainerVolume.from_json(json)
# print the JSON string representation of the object
print(ContainerVolume.to_json())

# convert the object into a dict
container_volume_dict = container_volume_instance.to_dict()
# create an instance of ContainerVolume from a dict
container_volume_from_dict = ContainerVolume.from_dict(container_volume_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


