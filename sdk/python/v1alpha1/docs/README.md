# flux-framework fluxoperator

> Python SDK for Flux-Operator

## Requirements.

Python >= 3.6

## Installation & Usage

### pip install

You can install directly from the git subfolder like:

```bash
pip install git+https://github.com/flux-framework/flux-operator.git#egg=subdir&subdirectory=python-sdk/v1alpha1
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

 - [ApiV1alpha1Commands](docs/ApiV1alpha1Commands.md)
 - [ApiV1alpha1ContainerResources](docs/ApiV1alpha1ContainerResources.md)
 - [ApiV1alpha1ContainerVolume](docs/ApiV1alpha1ContainerVolume.md)
 - [ApiV1alpha1FluxRestful](docs/ApiV1alpha1FluxRestful.md)
 - [ApiV1alpha1FluxUser](docs/ApiV1alpha1FluxUser.md)
 - [ApiV1alpha1LifeCycle](docs/ApiV1alpha1LifeCycle.md)
 - [ApiV1alpha1LoggingSpec](docs/ApiV1alpha1LoggingSpec.md)
 - [ApiV1alpha1MiniCluster](docs/ApiV1alpha1MiniCluster.md)
 - [ApiV1alpha1MiniClusterContainer](docs/ApiV1alpha1MiniClusterContainer.md)
 - [ApiV1alpha1MiniClusterList](docs/ApiV1alpha1MiniClusterList.md)
 - [ApiV1alpha1MiniClusterSpec](docs/ApiV1alpha1MiniClusterSpec.md)
 - [ApiV1alpha1MiniClusterStatus](docs/ApiV1alpha1MiniClusterStatus.md)
 - [ApiV1alpha1MiniClusterUser](docs/ApiV1alpha1MiniClusterUser.md)
 - [ApiV1alpha1MiniClusterVolume](docs/ApiV1alpha1MiniClusterVolume.md)
 - [ApiV1alpha1PodSpec](docs/ApiV1alpha1PodSpec.md)


## Documentation For Authorization

 All endpoints do not require authorization (but they do require you have permission via your kubernetes config)

## Author

- [@vsoch](https://github.com/vsoch)
