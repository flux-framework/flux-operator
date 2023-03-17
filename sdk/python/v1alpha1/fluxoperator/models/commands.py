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


class Commands(object):
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
        'broker_post': 'str',
        'post': 'str',
        'pre': 'str',
        'prefix': 'str',
        'run_flux_as_root': 'bool'
    }

    attribute_map = {
        'broker_post': 'brokerPost',
        'post': 'post',
        'pre': 'pre',
        'prefix': 'prefix',
        'run_flux_as_root': 'runFluxAsRoot'
    }

    def __init__(self, broker_post='', post='', pre='', prefix='', run_flux_as_root=False, local_vars_configuration=None):  # noqa: E501
        """Commands - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._broker_post = None
        self._post = None
        self._pre = None
        self._prefix = None
        self._run_flux_as_root = None
        self.discriminator = None

        if broker_post is not None:
            self.broker_post = broker_post
        if post is not None:
            self.post = post
        if pre is not None:
            self.pre = pre
        if prefix is not None:
            self.prefix = prefix
        if run_flux_as_root is not None:
            self.run_flux_as_root = run_flux_as_root

    @property
    def broker_post(self):
        """Gets the broker_post of this Commands.  # noqa: E501

        post command is run by the broker on finish  # noqa: E501

        :return: The broker_post of this Commands.  # noqa: E501
        :rtype: str
        """
        return self._broker_post

    @broker_post.setter
    def broker_post(self, broker_post):
        """Sets the broker_post of this Commands.

        post command is run by the broker on finish  # noqa: E501

        :param broker_post: The broker_post of this Commands.  # noqa: E501
        :type broker_post: str
        """

        self._broker_post = broker_post

    @property
    def post(self):
        """Gets the post of this Commands.  # noqa: E501

        post command is run by all pods on finish  # noqa: E501

        :return: The post of this Commands.  # noqa: E501
        :rtype: str
        """
        return self._post

    @post.setter
    def post(self, post):
        """Sets the post of this Commands.

        post command is run by all pods on finish  # noqa: E501

        :param post: The post of this Commands.  # noqa: E501
        :type post: str
        """

        self._post = post

    @property
    def pre(self):
        """Gets the pre of this Commands.  # noqa: E501

        pre command is run after global PreCommand, before anything else  # noqa: E501

        :return: The pre of this Commands.  # noqa: E501
        :rtype: str
        """
        return self._pre

    @pre.setter
    def pre(self, pre):
        """Sets the pre of this Commands.

        pre command is run after global PreCommand, before anything else  # noqa: E501

        :param pre: The pre of this Commands.  # noqa: E501
        :type pre: str
        """

        self._pre = pre

    @property
    def prefix(self):
        """Gets the prefix of this Commands.  # noqa: E501

        Prefix to flux start / submit / broker Typically used for a wrapper command to mount, etc.  # noqa: E501

        :return: The prefix of this Commands.  # noqa: E501
        :rtype: str
        """
        return self._prefix

    @prefix.setter
    def prefix(self, prefix):
        """Sets the prefix of this Commands.

        Prefix to flux start / submit / broker Typically used for a wrapper command to mount, etc.  # noqa: E501

        :param prefix: The prefix of this Commands.  # noqa: E501
        :type prefix: str
        """

        self._prefix = prefix

    @property
    def run_flux_as_root(self):
        """Gets the run_flux_as_root of this Commands.  # noqa: E501

        Run flux start as root - required for some storage binds  # noqa: E501

        :return: The run_flux_as_root of this Commands.  # noqa: E501
        :rtype: bool
        """
        return self._run_flux_as_root

    @run_flux_as_root.setter
    def run_flux_as_root(self, run_flux_as_root):
        """Sets the run_flux_as_root of this Commands.

        Run flux start as root - required for some storage binds  # noqa: E501

        :param run_flux_as_root: The run_flux_as_root of this Commands.  # noqa: E501
        :type run_flux_as_root: bool
        """

        self._run_flux_as_root = run_flux_as_root

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
        if not isinstance(other, Commands):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, Commands):
            return True

        return self.to_dict() != other.to_dict()
