# coding: utf-8

"""
    fluxoperator

    Python SDK for Flux-Operator  # noqa: E501

    The version of the OpenAPI document: v1alpha1
    Generated by: https://openapi-generator.tech
"""


from __future__ import absolute_import

import datetime
import unittest

import fluxoperator
from fluxoperator.models.life_cycle import LifeCycle  # noqa: E501
from fluxoperator.rest import ApiException


class TestLifeCycle(unittest.TestCase):
    """LifeCycle unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test LifeCycle
        include_option is a boolean, when False only required
        params are included, when True both required and
        optional params are included"""
        # model = fluxoperator.models.life_cycle.LifeCycle()  # noqa: E501
        if include_optional:
            return LifeCycle(post_start_exec="")
        else:
            return LifeCycle()

    def testLifeCycle(self):
        """Test LifeCycle"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)


if __name__ == "__main__":
    unittest.main()
