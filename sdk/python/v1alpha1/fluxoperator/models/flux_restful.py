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


class FluxRestful(object):
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
        'branch': 'str',
        'port': 'int',
        'token': 'str',
        'username': 'str'
    }

    attribute_map = {
        'branch': 'branch',
        'port': 'port',
        'token': 'token',
        'username': 'username'
    }

    def __init__(self, branch='', port=0, token='', username='', local_vars_configuration=None):  # noqa: E501
        """FluxRestful - a model defined in OpenAPI"""  # noqa: E501
        if local_vars_configuration is None:
            local_vars_configuration = Configuration.get_default_copy()
        self.local_vars_configuration = local_vars_configuration

        self._branch = None
        self._port = None
        self._token = None
        self._username = None
        self.discriminator = None

        if branch is not None:
            self.branch = branch
        if port is not None:
            self.port = port
        if token is not None:
            self.token = token
        if username is not None:
            self.username = username

    @property
    def branch(self):
        """Gets the branch of this FluxRestful.  # noqa: E501

        Branch to clone Flux Restful API from  # noqa: E501

        :return: The branch of this FluxRestful.  # noqa: E501
        :rtype: str
        """
        return self._branch

    @branch.setter
    def branch(self, branch):
        """Sets the branch of this FluxRestful.

        Branch to clone Flux Restful API from  # noqa: E501

        :param branch: The branch of this FluxRestful.  # noqa: E501
        :type branch: str
        """

        self._branch = branch

    @property
    def port(self):
        """Gets the port of this FluxRestful.  # noqa: E501

        Port to run Flux Restful Server On  # noqa: E501

        :return: The port of this FluxRestful.  # noqa: E501
        :rtype: int
        """
        return self._port

    @port.setter
    def port(self, port):
        """Sets the port of this FluxRestful.

        Port to run Flux Restful Server On  # noqa: E501

        :param port: The port of this FluxRestful.  # noqa: E501
        :type port: int
        """

        self._port = port

    @property
    def token(self):
        """Gets the token of this FluxRestful.  # noqa: E501

        Token to use for RestFul API  # noqa: E501

        :return: The token of this FluxRestful.  # noqa: E501
        :rtype: str
        """
        return self._token

    @token.setter
    def token(self, token):
        """Sets the token of this FluxRestful.

        Token to use for RestFul API  # noqa: E501

        :param token: The token of this FluxRestful.  # noqa: E501
        :type token: str
        """

        self._token = token

    @property
    def username(self):
        """Gets the username of this FluxRestful.  # noqa: E501

        These two should not actually be set by a user, but rather generated by tools and provided Username to use for RestFul API  # noqa: E501

        :return: The username of this FluxRestful.  # noqa: E501
        :rtype: str
        """
        return self._username

    @username.setter
    def username(self, username):
        """Sets the username of this FluxRestful.

        These two should not actually be set by a user, but rather generated by tools and provided Username to use for RestFul API  # noqa: E501

        :param username: The username of this FluxRestful.  # noqa: E501
        :type username: str
        """

        self._username = username

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
        if not isinstance(other, FluxRestful):
            return False

        return self.to_dict() == other.to_dict()

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        if not isinstance(other, FluxRestful):
            return True

        return self.to_dict() != other.to_dict()
