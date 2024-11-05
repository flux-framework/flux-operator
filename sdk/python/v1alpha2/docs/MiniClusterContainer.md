# MiniClusterContainer


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**batch** | **bool** | Indicate that the command is a batch job that will be written to a file to submit | [optional] [default to False]
**batch_raw** | **bool** | Don&#39;t wrap batch commands in flux submit (provide custom logic myself) | [optional] [default to False]
**command** | **str** | Single user executable to provide to flux start | [optional] [default to '']
**commands** | [**Commands**](Commands.md) |  | [optional] 
**environment** | **Dict[str, str]** | Key/value pairs for the environment | [optional] 
**image** | **str** | Container image must contain flux and flux-sched install | [optional] [default to 'ghcr.io/rse-ops/accounting:app-latest']
**image_pull_secret** | **str** | Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec. | [optional] [default to '']
**launcher** | **bool** | Indicate that the command is a launcher that will ask for its own jobs (and provided directly to flux start) | [optional] [default to False]
**life_cycle** | [**LifeCycle**](LifeCycle.md) |  | [optional] 
**logs** | **str** | Log output directory | [optional] [default to '']
**name** | **str** | Container name is only required for non flux runners | [optional] [default to '']
**no_wrap_entrypoint** | **bool** | Do not wrap the entrypoint to wait for flux, add to path, etc? | [optional] [default to False]
**ports** | **List[int]** | Ports to be exposed to other containers in the cluster We take a single list of integers and map to the same | [optional] 
**pull_always** | **bool** | Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always | [optional] [default to False]
**resources** | [**ContainerResources**](ContainerResources.md) |  | [optional] 
**run_flux** | **bool** | Application container intended to run flux (broker) | [optional] [default to False]
**secrets** | [**Dict[str, Secret]**](Secret.md) | Secrets that will be added to the environment The user is expected to create their own secrets for the operator to find | [optional] 
**security_context** | [**SecurityContext**](SecurityContext.md) |  | [optional] 
**volumes** | [**Dict[str, ContainerVolume]**](ContainerVolume.md) | Existing volumes that can be mounted | [optional] 
**working_dir** | **str** | Working directory to run command from | [optional] [default to '']

## Example

```python
from fluxoperator.models.mini_cluster_container import MiniClusterContainer

# TODO update the JSON string below
json = "{}"
# create an instance of MiniClusterContainer from a JSON string
mini_cluster_container_instance = MiniClusterContainer.from_json(json)
# print the JSON string representation of the object
print(MiniClusterContainer.to_json())

# convert the object into a dict
mini_cluster_container_dict = mini_cluster_container_instance.to_dict()
# create an instance of MiniClusterContainer from a dict
mini_cluster_container_from_dict = MiniClusterContainer.from_dict(mini_cluster_container_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


