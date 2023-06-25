# coding: utf-8

"""
    fluxoperator

    Python SDK for Flux-Operator  # noqa: E501

    The version of the OpenAPI document: v1alpha1
    Generated by: https://openapi-generator.tech
"""


from __future__ import absolute_import

import unittest
import datetime

import fluxoperator
from fluxoperator.models.mini_cluster_archive import MiniClusterArchive  # noqa: E501
from fluxoperator.rest import ApiException


class TestMiniClusterArchive(unittest.TestCase):
    """MiniClusterArchive unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test MiniClusterArchive
        include_option is a boolean, when False only required
        params are included, when True both required and
        optional params are included"""
        # model = fluxoperator.models.mini_cluster_archive.MiniClusterArchive()  # noqa: E501
        if include_optional:
            return MiniClusterArchive(path="")
        else:
            return MiniClusterArchive()

    def testMiniClusterArchive(self):
        """Test MiniClusterArchive"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)


if __name__ == "__main__":
    unittest.main()
