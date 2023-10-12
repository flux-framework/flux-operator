# coding: utf-8

"""
    fluxoperator

    Python SDK for Flux-Operator  # noqa: E501

    The version of the OpenAPI document: v1alpha2
    Generated by: https://openapi-generator.tech
"""


try:
    from inspect import getfullargspec
except ImportError:
    from inspect import getargspec as getfullargspec
import pprint
import re  # noqa: F401
import six

from fluxoperator.configuration import Configuration


class FluxContainer(object):
    """NOTE: This class is auto generated by OpenAPI Generator.
    Ref: https://openapi-generator.tech

    Do not edit the class manually.
    """

    """
    Attributes:
      openapi_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    openapi_types = {
        'cores': 'int',
        'image': 'str',
        'image_pull_secret': 'str',
        'mount_path': 'str',
        'name': 'str',
        'pull_always': 'bool',
        'python_path': 'str',
        'working_dir': 'str'
    }

    attribute_map = {
        'cores': 'cores',
        'image': 'image',
        'image_pull_secret': 'imagePullSecret',
        'mount_path': 'mountPath',
        'name': 'name',
        'pull_always': 'pullAlways',
        'python_path': 'pythonPath',
        'working_dir': 'workingDir'
    }

    def __init__(self, cores=0, image='ghcr.io/converged-computing/flux-view-rocky:tag-9', image_pull_secret='', mount_path='/mnt/flux', name='flux-view', pull_always=False, python_path='', working_dir='', local_vars_configuration=None):  # noqa: E501
        """FluxContainer - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._cores = None
        self._image = None
        self._image_pull_secret = None
        self._mount_path = None
        self._name = None
        self._pull_always = None
        self._python_path = None
        self._working_dir = None
        self.discriminator = None

        if cores is not None:
            self.cores = cores
        if image is not None:
            self.image = image
        if image_pull_secret is not None:
            self.image_pull_secret = image_pull_secret
        if mount_path is not None:
            self.mount_path = mount_path
        if name is not None:
            self.name = name
        if pull_always is not None:
            self.pull_always = pull_always
        if python_path is not None:
            self.python_path = python_path
        if working_dir is not None:
            self.working_dir = working_dir

    @property
    def cores(self):
        """Gets the cores of this FluxContainer.  # noqa: E501

        Cores flux should use  # noqa: E501

        :return: The cores of this FluxContainer.  # noqa: E501
        :rtype: int
        """
        return self._cores

    @cores.setter
    def cores(self, cores):
        """Sets the cores of this FluxContainer.

        Cores flux should use  # noqa: E501

        :param cores: The cores of this FluxContainer.  # noqa: E501
        :type cores: int
        """

        self._cores = cores

    @property
    def image(self):
        """Gets the image of this FluxContainer.  # noqa: E501


        :return: The image of this FluxContainer.  # noqa: E501
        :rtype: str
        """
        return self._image

    @image.setter
    def image(self, image):
        """Sets the image of this FluxContainer.


        :param image: The image of this FluxContainer.  # noqa: E501
        :type image: str
        """

        self._image = image

    @property
    def image_pull_secret(self):
        """Gets the image_pull_secret of this FluxContainer.  # noqa: E501

        Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec.  # noqa: E501

        :return: The image_pull_secret of this FluxContainer.  # noqa: E501
        :rtype: str
        """
        return self._image_pull_secret

    @image_pull_secret.setter
    def image_pull_secret(self, image_pull_secret):
        """Sets the image_pull_secret of this FluxContainer.

        Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec.  # noqa: E501

        :param image_pull_secret: The image_pull_secret of this FluxContainer.  # noqa: E501
        :type image_pull_secret: str
        """

        self._image_pull_secret = image_pull_secret

    @property
    def mount_path(self):
        """Gets the mount_path of this FluxContainer.  # noqa: E501

        Mount path for flux to be at (will be added to path)  # noqa: E501

        :return: The mount_path of this FluxContainer.  # noqa: E501
        :rtype: str
        """
        return self._mount_path

    @mount_path.setter
    def mount_path(self, mount_path):
        """Sets the mount_path of this FluxContainer.

        Mount path for flux to be at (will be added to path)  # noqa: E501

        :param mount_path: The mount_path of this FluxContainer.  # noqa: E501
        :type mount_path: str
        """

        self._mount_path = mount_path

    @property
    def name(self):
        """Gets the name of this FluxContainer.  # noqa: E501

        Container name is only required for non flux runners  # noqa: E501

        :return: The name of this FluxContainer.  # noqa: E501
        :rtype: str
        """
        return self._name

    @name.setter
    def name(self, name):
        """Sets the name of this FluxContainer.

        Container name is only required for non flux runners  # noqa: E501

        :param name: The name of this FluxContainer.  # noqa: E501
        :type name: str
        """

        self._name = name

    @property
    def pull_always(self):
        """Gets the pull_always of this FluxContainer.  # noqa: E501

        Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always  # noqa: E501

        :return: The pull_always of this FluxContainer.  # noqa: E501
        :rtype: bool
        """
        return self._pull_always

    @pull_always.setter
    def pull_always(self, pull_always):
        """Sets the pull_always of this FluxContainer.

        Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always  # noqa: E501

        :param pull_always: The pull_always of this FluxContainer.  # noqa: E501
        :type pull_always: bool
        """

        self._pull_always = pull_always

    @property
    def python_path(self):
        """Gets the python_path of this FluxContainer.  # noqa: E501

        Customize python path for flux  # noqa: E501

        :return: The python_path of this FluxContainer.  # noqa: E501
        :rtype: str
        """
        return self._python_path

    @python_path.setter
    def python_path(self, python_path):
        """Sets the python_path of this FluxContainer.

        Customize python path for flux  # noqa: E501

        :param python_path: The python_path of this FluxContainer.  # noqa: E501
        :type python_path: str
        """

        self._python_path = python_path

    @property
    def working_dir(self):
        """Gets the working_dir of this FluxContainer.  # noqa: E501

        Working directory to run command from  # noqa: E501

        :return: The working_dir of this FluxContainer.  # noqa: E501
        :rtype: str
        """
        return self._working_dir

    @working_dir.setter
    def working_dir(self, working_dir):
        """Sets the working_dir of this FluxContainer.

        Working directory to run command from  # noqa: E501

        :param working_dir: The working_dir of this FluxContainer.  # noqa: E501
        :type working_dir: str
        """

        self._working_dir = working_dir

    def to_dict(self, serialize=False):
        """Returns the model properties as a dict"""
        result = {}

        def convert(x):
            if hasattr(x, "to_dict"):
                args = getfullargspec(x.to_dict).args
                if len(args) == 1:
                    return x.to_dict()
                else:
                    return x.to_dict(serialize)
            else:
                return x

        for attr, _ in six.iteritems(self.openapi_types):
            value = getattr(self, attr)
            attr = self.attribute_map.get(attr, attr) if serialize else attr
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: convert(x),
                    value
                ))
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], convert(item[1])),
                    value.items()
                ))
            else:
                result[attr] = convert(value)

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, FluxContainer):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, FluxContainer):
            return True

        return self.to_dict() != other.to_dict()