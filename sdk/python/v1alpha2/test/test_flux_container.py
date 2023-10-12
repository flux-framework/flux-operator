# coding: utf-8

"""
    fluxoperator

    Python SDK for Flux-Operator  # noqa: E501

    The version of the OpenAPI document: v1alpha2
    Generated by: https://openapi-generator.tech
"""


from __future__ import absolute_import

import unittest
import datetime

import fluxoperator
from fluxoperator.models.flux_container import FluxContainer  # noqa: E501
from fluxoperator.rest import ApiException

class TestFluxContainer(unittest.TestCase):
    """FluxContainer unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test FluxContainer
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # model = fluxoperator.models.flux_container.FluxContainer()  # noqa: E501
        if include_optional :
            return FluxContainer(
                cores = 56, 
                image = 'ghcr.io/converged-computing/flux-view-rocky:tag-9', 
                image_pull_secret = '', 
                mount_path = '/mnt/flux', 
                name = 'flux-view', 
                pull_always = True, 
                python_path = '', 
                working_dir = ''
            )
        else :
            return FluxContainer(
        )

    def testFluxContainer(self):
        """Test FluxContainer"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()