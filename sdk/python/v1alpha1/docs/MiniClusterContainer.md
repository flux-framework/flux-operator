# MiniClusterContainer


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**batch** | **bool** | Indicate that the command is a batch job that will be written to a file to submit | [optional] [default to False]
**batch_raw** | **bool** | Don&#39;t wrap batch commands in flux submit (provide custom logic myself) | [optional] [default to False]
**command** | **str** | Single user executable to provide to flux start | [optional] [default to '']
**commands** | [**Commands**](Commands.md) |  | [optional] 
**cores** | **int** | Cores the container should use | [optional] [default to 0]
**diagnostics** | **bool** | Run flux diagnostics on start instead of command | [optional] [default to False]
**environment** | **dict(str, str)** | Key/value pairs for the environment | [optional] 
**existing_volumes** | [**dict(str, MiniClusterExistingVolume)**](MiniClusterExistingVolume.md) | Existing Volumes to add to the containers | [optional] 
**flux_user** | [**FluxUser**](FluxUser.md) |  | [optional] 
**image** | **str** | Container image must contain flux and flux-sched install | [optional] [default to 'ghcr.io/rse-ops/accounting:app-latest']
**image_pull_secret** | **str** | Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec. | [optional] [default to '']
**launcher** | **bool** | Indicate that the command is a launcher that will ask for its own jobs (and provided directly to flux start) | [optional] [default to False]
**life_cycle** | [**LifeCycle**](LifeCycle.md) |  | [optional] 
**logs** | **str** | Log output directory | [optional] [default to '']
**name** | **str** | Container name is only required for non flux runners | [optional] [default to '']
**ports** | **list[int]** | Ports to be exposed to other containers in the cluster We take a single list of integers and map to the same | [optional] 
**pull_always** | **bool** | Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always | [optional] [default to False]
**resources** | [**ContainerResources**](ContainerResources.md) |  | [optional] 
**run_flux** | **bool** | Main container to run flux (only should be one) | [optional] [default to False]
**secrets** | [**dict(str, Secret)**](Secret.md) | Secrets that will be added to the environment The user is expected to create their own secrets for the operator to find | [optional] 
**security_context** | [**SecurityContext**](SecurityContext.md) |  | [optional] 
**volumes** | [**dict(str, ContainerVolume)**](ContainerVolume.md) | Volumes that can be mounted (must be defined in volumes) | [optional] 
**working_dir** | **str** | Working directory to run command from | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


