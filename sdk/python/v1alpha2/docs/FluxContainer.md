# FluxContainer

A FluxContainer is equivalent to a MiniCluster container but has a different default image

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**disable** | **bool** | Disable the sidecar container, assuming that the main application container has flux | [optional] [default to False]
**image** | **str** |  | [optional] [default to 'ghcr.io/converged-computing/flux-view-rocky:tag-9']
**image_pull_secret** | **str** | Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec. | [optional] [default to '']
**mount_path** | **str** | Mount path for flux to be at (will be added to path) | [optional] [default to '/mnt/flux']
**name** | **str** | Container name is only required for non flux runners | [optional] [default to 'flux-view']
**pull_always** | **bool** | Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always | [optional] [default to False]
**python_path** | **str** | Customize python path for flux | [optional] [default to '']
**resources** | [**ContainerResources**](ContainerResources.md) |  | [optional] 
**working_dir** | **str** | Working directory to run command from | [optional] [default to '']

## Example

```python
from fluxoperator.models.flux_container import FluxContainer

# TODO update the JSON string below
json = "{}"
# create an instance of FluxContainer from a JSON string
flux_container_instance = FluxContainer.from_json(json)
# print the JSON string representation of the object
print(FluxContainer.to_json())

# convert the object into a dict
flux_container_dict = flux_container_instance.to_dict()
# create an instance of FluxContainer from a dict
flux_container_from_dict = FluxContainer.from_dict(flux_container_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


