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


class FluxScheduler(object):
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
        'queue_policy': 'str'
    }

    attribute_map = {
        'queue_policy': 'queuePolicy'
    }

    def __init__(self, queue_policy='', local_vars_configuration=None):  # noqa: E501
        """FluxScheduler - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._queue_policy = None
        self.discriminator = None

        if queue_policy is not None:
            self.queue_policy = queue_policy

    @property
    def queue_policy(self):
        """Gets the queue_policy of this FluxScheduler.  # noqa: E501

        Scheduler queue policy, defaults to \"fcfs\" can also be \"easy\"  # noqa: E501

        :return: The queue_policy of this FluxScheduler.  # noqa: E501
        :rtype: str
        """
        return self._queue_policy

    @queue_policy.setter
    def queue_policy(self, queue_policy):
        """Sets the queue_policy of this FluxScheduler.

        Scheduler queue policy, defaults to \"fcfs\" can also be \"easy\"  # noqa: E501

        :param queue_policy: The queue_policy of this FluxScheduler.  # noqa: E501
        :type queue_policy: str
        """

        self._queue_policy = queue_policy

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
        if not isinstance(other, FluxScheduler):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, FluxScheduler):
            return True

        return self.to_dict() != other.to_dict()