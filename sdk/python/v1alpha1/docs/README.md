# flux-framework fluxoperator

> Python SDK for Flux-Operator

## Requirements.

Python >= 3.6

## Installation & Usage

### pip install

You can install directly from the git subfolder like:

```bash
pip install git+https://github.com/flux-framework/flux-operator.git#egg=subdir&subdirectory=sdk/python/v1alpha1
```

or install from pip:

```bash
$ pip install fluxoperator
```

Then import the package:

```python
import fluxoperator
```

### Setuptools

Install via [Setuptools](http://pypi.python.org/pypi/setuptools).

```sh
python setup.py install --user
```

(or `sudo python setup.py install` to install the package for all users)

Then import the package:

```python
import fluxoperator
```

## Getting Started

Please follow the [installation procedure](#installation--usage) and then run the following:

```python
import time
import fluxoperator
from pprint import pprint
```

## Documentation For Models

 - [Commands](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/Commands.md)
 - [ContainerResources](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/ContainerResources.md)
 - [ContainerVolume](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/ContainerVolume.md)
 - [FluxRestful](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/FluxRestful.md)
 - [FluxUser](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/FluxUser.md)
 - [LifeCycle](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/LifeCycle.md)
 - [LoggingSpec](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/LoggingSpec.md)
 - [MiniCluster](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniCluster.md)
 - [MiniClusterContainer](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniClusterContainer.md)
 - [MiniClusterList](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniClusterList.md)
 - [MiniClusterSpec](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniClusterSpec.md)
 - [MiniClusterStatus](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniClusterStatus.md)
 - [MiniClusterUser](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniClusterUser.md)
 - [MiniClusterVolume](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/MiniClusterVolume.md)
 - [PodSpec](https://github.com/flux-framework/flux-operator/tree/main/sdk/python/v1alpha1/docs/PodSpec.md)


## Documentation For Authorization

 All endpoints do not require authorization (but they do require you have permission via your kubernetes config)

## Author

- [@vsoch](https://github.com/vsoch)
