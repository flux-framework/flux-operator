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


class LifeCycle(object):
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
        'post_start_exec': 'str'
    }

    attribute_map = {
        'post_start_exec': 'postStartExec'
    }

    def __init__(self, post_start_exec='', local_vars_configuration=None):  # noqa: E501
        """LifeCycle - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._post_start_exec = None
        self.discriminator = None

        if post_start_exec is not None:
            self.post_start_exec = post_start_exec

    @property
    def post_start_exec(self):
        """Gets the post_start_exec of this LifeCycle.  # noqa: E501


        :return: The post_start_exec of this LifeCycle.  # noqa: E501
        :rtype: str
        """
        return self._post_start_exec

    @post_start_exec.setter
    def post_start_exec(self, post_start_exec):
        """Sets the post_start_exec of this LifeCycle.


        :param post_start_exec: The post_start_exec of this LifeCycle.  # noqa: E501
        :type post_start_exec: str
        """

        self._post_start_exec = post_start_exec

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
        if not isinstance(other, LifeCycle):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, LifeCycle):
            return True

        return self.to_dict() != other.to_dict()
