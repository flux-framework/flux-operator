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
        'broker_config': 'str',
        'bursting': 'Bursting',
        'connect_timeout': 'str',
        'curve_cert': 'str',
        'curve_cert_secret': 'str',
        'install_root': 'str',
        'log_level': 'int',
        'minimal_service': 'bool',
        'munge_secret': 'str',
        'option_flags': 'str',
        'wrap': 'str'
    }

    attribute_map = {
        'broker_config': 'brokerConfig',
        'bursting': 'bursting',
        'connect_timeout': 'connectTimeout',
        'curve_cert': 'curveCert',
        'curve_cert_secret': 'curveCertSecret',
        'install_root': 'installRoot',
        'log_level': 'logLevel',
        'minimal_service': 'minimalService',
        'munge_secret': 'mungeSecret',
        'option_flags': 'optionFlags',
        'wrap': 'wrap'
    }

    def __init__(self, broker_config='', bursting=None, connect_timeout='5s', curve_cert='', curve_cert_secret='', install_root='/usr', log_level=6, minimal_service=False, munge_secret='', option_flags='', wrap=None, local_vars_configuration=None):  # noqa: E501
        """FluxSpec - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._broker_config = None
        self._bursting = None
        self._connect_timeout = None
        self._curve_cert = None
        self._curve_cert_secret = None
        self._install_root = None
        self._log_level = None
        self._minimal_service = None
        self._munge_secret = None
        self._option_flags = None
        self._wrap = None
        self.discriminator = None

        if broker_config is not None:
            self.broker_config = broker_config
        if bursting is not None:
            self.bursting = bursting
        if connect_timeout is not None:
            self.connect_timeout = connect_timeout
        if curve_cert is not None:
            self.curve_cert = curve_cert
        self.curve_cert_secret = curve_cert_secret
        if install_root is not None:
            self.install_root = install_root
        if log_level is not None:
            self.log_level = log_level
        if minimal_service is not None:
            self.minimal_service = minimal_service
        if munge_secret is not None:
            self.munge_secret = munge_secret
        if option_flags is not None:
            self.option_flags = option_flags
        if wrap is not None:
            self.wrap = wrap

    @property
    def broker_config(self):
        """Gets the broker_config of this FluxSpec.  # noqa: E501

        Optionally provide a manually created broker config this is intended for bursting to remote clusters  # noqa: E501

        :return: The broker_config of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._broker_config

    @broker_config.setter
    def broker_config(self, broker_config):
        """Sets the broker_config of this FluxSpec.

        Optionally provide a manually created broker config this is intended for bursting to remote clusters  # noqa: E501

        :param broker_config: The broker_config of this FluxSpec.  # noqa: E501
        :type broker_config: str
        """

        self._broker_config = broker_config

    @property
    def bursting(self):
        """Gets the bursting of this FluxSpec.  # noqa: E501


        :return: The bursting of this FluxSpec.  # noqa: E501
        :rtype: Bursting
        """
        return self._bursting

    @bursting.setter
    def bursting(self, bursting):
        """Sets the bursting of this FluxSpec.


        :param bursting: The bursting of this FluxSpec.  # noqa: E501
        :type bursting: Bursting
        """

        self._bursting = bursting

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
    def curve_cert(self):
        """Gets the curve_cert of this FluxSpec.  # noqa: E501

        Optionally provide an already existing curve certificate this is intended for bursting to remote clusters  # noqa: E501

        :return: The curve_cert of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._curve_cert

    @curve_cert.setter
    def curve_cert(self, curve_cert):
        """Sets the curve_cert of this FluxSpec.

        Optionally provide an already existing curve certificate this is intended for bursting to remote clusters  # noqa: E501

        :param curve_cert: The curve_cert of this FluxSpec.  # noqa: E501
        :type curve_cert: str
        """

        self._curve_cert = curve_cert

    @property
    def curve_cert_secret(self):
        """Gets the curve_cert_secret of this FluxSpec.  # noqa: E501

        Expect a secret for a curve cert here. This is ideal over the curveCert (as a string) above.  # noqa: E501

        :return: The curve_cert_secret of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._curve_cert_secret

    @curve_cert_secret.setter
    def curve_cert_secret(self, curve_cert_secret):
        """Sets the curve_cert_secret of this FluxSpec.

        Expect a secret for a curve cert here. This is ideal over the curveCert (as a string) above.  # noqa: E501

        :param curve_cert_secret: The curve_cert_secret of this FluxSpec.  # noqa: E501
        :type curve_cert_secret: str
        """
        if self.local_vars_configuration.client_side_validation and curve_cert_secret is None:  # noqa: E501
            raise ValueError("Invalid value for `curve_cert_secret`, must not be `None`")  # noqa: E501

        self._curve_cert_secret = curve_cert_secret

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
    def minimal_service(self):
        """Gets the minimal_service of this FluxSpec.  # noqa: E501

        Only expose the broker service (to reduce load on DNS)  # noqa: E501

        :return: The minimal_service of this FluxSpec.  # noqa: E501
        :rtype: bool
        """
        return self._minimal_service

    @minimal_service.setter
    def minimal_service(self, minimal_service):
        """Sets the minimal_service of this FluxSpec.

        Only expose the broker service (to reduce load on DNS)  # noqa: E501

        :param minimal_service: The minimal_service of this FluxSpec.  # noqa: E501
        :type minimal_service: bool
        """

        self._minimal_service = minimal_service

    @property
    def munge_secret(self):
        """Gets the munge_secret of this FluxSpec.  # noqa: E501

        Expect a secret (named according to this string) for a munge key. This is intended for bursting. Assumed to be at /etc/munge/munge.key This is binary data.  # noqa: E501

        :return: The munge_secret of this FluxSpec.  # noqa: E501
        :rtype: str
        """
        return self._munge_secret

    @munge_secret.setter
    def munge_secret(self, munge_secret):
        """Sets the munge_secret of this FluxSpec.

        Expect a secret (named according to this string) for a munge key. This is intended for bursting. Assumed to be at /etc/munge/munge.key This is binary data.  # noqa: E501

        :param munge_secret: The munge_secret of this FluxSpec.  # noqa: E501
        :type munge_secret: str
        """

        self._munge_secret = munge_secret

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
