# MiniClusterContainer


## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**command** | **str** | Single user executable to provide to flux start IMPORTANT: This is left here, but not used in favor of exposing Flux via a Restful API. We Can remove this when that is finalized. | [optional] [default to '']
**commands** | [**Commands**](Commands.md) |  | [optional] 
**cores** | **int** | Cores the container should use | [optional] [default to 0]
**diagnostics** | **bool** | Run flux diagnostics on start instead of command | [optional] [default to False]
**environment** | **dict(str, str)** | Key/value pairs for the environment | [optional] 
**flux_log_level** | **int** | Log level to use for flux logging (only in non TestMode) | [optional] [default to 6]
**flux_option_flags** | **str** | Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \&quot;runFlux\&quot; true | [optional] [default to '']
**flux_user** | [**FluxUser**](FluxUser.md) |  | [optional] 
**image** | **str** | Container image must contain flux and flux-sched install | [optional] [default to 'ghcr.io/rse-ops/accounting:app-latest']
**image_pull_secret** | **str** | Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec. | [optional] [default to '']
**life_cycle** | [**LifeCycle**](LifeCycle.md) |  | [optional] 
**name** | **str** | Container name is only required for non flux runners | [optional] [default to '']
**ports** | **list[int]** | Ports to be exposed to other containers in the cluster We take a single list of integers and map to the same | [optional] 
**pre_command** | **str** | Special command to run at beginning of script, directly after asFlux is defined as sudo -u flux -E (so you can change that if desired.) This is only valid if FluxRunner is set (that writes a wait.sh script) This is for the indexed job pods and the certificate generation container. | [optional] [default to '']
**pull_always** | **bool** | Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always | [optional] [default to False]
**resources** | [**ContainerResources**](ContainerResources.md) |  | [optional] 
**run_flux** | **bool** | Main container to run flux (only should be one) | [optional] [default to False]
**volumes** | [**dict(str, ContainerVolume)**](ContainerVolume.md) | Volumes that can be mounted (must be defined in volumes) | [optional] 
**working_dir** | **str** | Working directory to run command from | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


