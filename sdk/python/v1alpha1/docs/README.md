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

 - [BurstedCluster](BurstedCluster.md)
 - [Bursting](Bursting.md)
 - [Commands](Commands.md)
 - [ContainerResources](dContainerResources.md)
 - [ContainerVolume](ContainerVolume.md)
 - [FluxBroker](FluxBrokerl.md)
 - [FluxRestful](FluxRestful.md)
 - [FluxSpec](FluxSpec.md)
 - [FluxUser](FluxUser.md)
 - [LifeCycle](LifeCycle.md)
 - [LoggingSpec](LoggingSpec.md)
 - [MiniCluster](MiniCluster.md)
 - [MiniClusterArchive](MiniClusterArchive.md)
 - [MiniClusterContainer](MiniClusterContainer.md)
 - [MiniClusterList](MiniClusterList.md)
 - [MiniClusterSpec](MiniClusterSpec.md)
 - [MiniClusterStatus](MiniClusterStatus.md)
 - [MiniClusterUser](MiniClusterUser.md)
 - [MiniClusterExistingVolume](MiniClusterExistingVolume.md)
 - [MiniClusterVolume](MiniClusterVolume.md)
 - [Network](Network.md)
 - [PodSpec](PodSpec.md)
 - [SecurityContext](SecurityContext.md)

## Documentation For Authorization

 All endpoints do not require authorization (but they do require you have permission via your kubernetes config)

## Author

- [@vsoch](https://github.com/vsoch)
