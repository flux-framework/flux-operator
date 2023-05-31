# coding: utf-8

"""
    fluxoperator

    Python SDK for Flux-Operator  # noqa: E501

    The version of the OpenAPI document: v1alpha1
    Generated by: https://openapi-generator.tech
"""


import inspect
import pprint
import re  # noqa: F401
import six

from fluxoperator.configuration import Configuration


class FluxSpec(object):
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
        'connect_timeout': 'str',
        'connection': 'str',
        'connection_size': 'int',
        'install_root': 'str',
        'log_level': 'int',
        'option_flags': 'str',
        'wrap': 'str'
    }

    attribute_map = {
        'connect_timeout': 'connectTimeout',
        'connection': 'connection',
        'connection_size': 'connectionSize',
        'install_root': 'installRoot',
        'log_level': 'logLevel',
        'option_flags': 'optionFlags',
        'wrap': 'wrap'
    }

    def __init__(self, connect_timeout='5s', connection=None, connection_size=None, install_root='/usr', log_level=6, option_flags='', wrap=None, local_vars_configuration=None):  # noqa: E501
        """FluxSpec - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._connect_timeout = None
        self._connection = None
        self._connection_size = None
        self._install_root = None
        self._log_level = None
        self._option_flags = None
        self._wrap = None
        self.discriminator = None

        if connect_timeout is not None:
            self.connect_timeout = connect_timeout
        if connection is not None:
            self.connection = connection
        if connection_size is not None:
            self.connection_size = connection_size
        if install_root is not None:
            self.install_root = install_root
        if log_level is not None:
            self.log_level = log_level
        if option_flags is not None:
            self.option_flags = option_flags
        if wrap is not None:
            self.wrap = wrap

    @property
    def connect_timeout(self):
        """Gets the connect_timeout of this FluxSpec.  # noqa: E501

        Single user executable to provide to flux start  # noqa: E501

        :return: The connect_timeout of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._connect_timeout

    @connect_timeout.setter
    def connect_timeout(self, connect_timeout):
        """Sets the connect_timeout of this FluxSpec.

        Single user executable to provide to flux start  # noqa: E501

        :param connect_timeout: The connect_timeout of this FluxSpec.  # noqa: E501
        :type connect_timeout: str
        """

        self._connect_timeout = connect_timeout

    @property
    def connection(self):
        """Gets the connection of this FluxSpec.  # noqa: E501

        Connect to this job in the same namespace (akin to BootServer but within cluster)  # noqa: E501

        :return: The connection of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._connection

    @connection.setter
    def connection(self, connection):
        """Sets the connection of this FluxSpec.

        Connect to this job in the same namespace (akin to BootServer but within cluster)  # noqa: E501

        :param connection: The connection of this FluxSpec.  # noqa: E501
        :type connection: str
        """

        self._connection = connection

    @property
    def connection_size(self):
        """Gets the connection_size of this FluxSpec.  # noqa: E501

        Additional number of nodes to allow from external boot-server This currently only allows local MiniCluster but could be extended to any general URI  # noqa: E501

        :return: The connection_size of this FluxSpec.  # noqa: E501
        :rtype: int
        """
        return self._connection_size

    @connection_size.setter
    def connection_size(self, connection_size):
        """Sets the connection_size of this FluxSpec.

        Additional number of nodes to allow from external boot-server This currently only allows local MiniCluster but could be extended to any general URI  # noqa: E501

        :param connection_size: The connection_size of this FluxSpec.  # noqa: E501
        :type connection_size: int
        """

        self._connection_size = connection_size

    @property
    def install_root(self):
        """Gets the install_root of this FluxSpec.  # noqa: E501

        Install root location  # noqa: E501

        :return: The install_root of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._install_root

    @install_root.setter
    def install_root(self, install_root):
        """Sets the install_root of this FluxSpec.

        Install root location  # noqa: E501

        :param install_root: The install_root of this FluxSpec.  # noqa: E501
        :type install_root: str
        """

        self._install_root = install_root

    @property
    def log_level(self):
        """Gets the log_level of this FluxSpec.  # noqa: E501

        Log level to use for flux logging (only in non TestMode)  # noqa: E501

        :return: The log_level of this FluxSpec.  # noqa: E501
        :rtype: int
        """
        return self._log_level

    @log_level.setter
    def log_level(self, log_level):
        """Sets the log_level of this FluxSpec.

        Log level to use for flux logging (only in non TestMode)  # noqa: E501

        :param log_level: The log_level of this FluxSpec.  # noqa: E501
        :type log_level: int
        """

        self._log_level = log_level

    @property
    def option_flags(self):
        """Gets the option_flags of this FluxSpec.  # noqa: E501

        Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \"runFlux\" true  # noqa: E501

        :return: The option_flags of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._option_flags

    @option_flags.setter
    def option_flags(self, option_flags):
        """Sets the option_flags of this FluxSpec.

        Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \"runFlux\" true  # noqa: E501

        :param option_flags: The option_flags of this FluxSpec.  # noqa: E501
        :type option_flags: str
        """

        self._option_flags = option_flags

    @property
    def wrap(self):
        """Gets the wrap of this FluxSpec.  # noqa: E501

        Commands for flux start --wrap  # noqa: E501

        :return: The wrap of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._wrap

    @wrap.setter
    def wrap(self, wrap):
        """Sets the wrap of this FluxSpec.

        Commands for flux start --wrap  # noqa: E501

        :param wrap: The wrap of this FluxSpec.  # noqa: E501
        :type wrap: str
        """

        self._wrap = wrap

    def to_dict(self, serialize=False):
        """Returns the model properties as a dict"""
        result = {}

        def convert(x):
            if hasattr(x, "to_dict"):
                args = inspect.getargspec(x.to_dict).args
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
        if not isinstance(other, FluxSpec):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, FluxSpec):
            return True

        return self.to_dict() != other.to_dict()
