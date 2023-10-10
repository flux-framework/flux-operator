# FluxContainer

A FluxContainer is equivalent to a MiniCluster container but has a different default image

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**cores** | **int** | Cores flux should use | [optional] [default to 0]
**image** | **str** |  | [optional] [default to 'ghcr.io/converged-computing/flux-view-rocky:tag-9']
**image_pull_secret** | **str** | Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec. | [optional] [default to '']
**mount_path** | **str** | Mount path for flux to be at (will be added to path) | [optional] [default to '/mnt/flux']
**name** | **str** | Container name is only required for non flux runners | [optional] [default to 'flux-view']
**pull_always** | **bool** | Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always | [optional] [default to False]
**python_path** | **str** | Customize python path for flux | [optional] [default to '']
**working_dir** | **str** | Working directory to run command from | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


