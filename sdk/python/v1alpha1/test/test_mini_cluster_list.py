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
from fluxoperator.models.mini_cluster_list import MiniClusterList  # noqa: E501
from fluxoperator.rest import ApiException

class TestMiniClusterList(unittest.TestCase):
    """MiniClusterList unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test MiniClusterList
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # model = fluxoperator.models.mini_cluster_list.MiniClusterList()  # noqa: E501
        if include_optional :
            return MiniClusterList(
                api_version = '', 
                items = [
                    fluxoperator.models.mini_cluster.MiniCluster(
                        api_version = '', 
                        kind = '', 
                        metadata = None, 
                        spec = fluxoperator.models.mini_cluster_spec.MiniClusterSpec(
                            cleanup = True, 
                            containers = [
                                fluxoperator.models.mini_cluster_container.MiniClusterContainer(
                                    command = '', 
                                    commands = fluxoperator.models.commands.Commands(
                                        pre = '', 
                                        run_flux_as_root = True, ), 
                                    cores = 56, 
                                    diagnostics = True, 
                                    environment = {
                                        'key' : ''
                                        }, 
                                    flux_log_level = 56, 
                                    flux_option_flags = '', 
                                    flux_user = fluxoperator.models.flux_user.FluxUser(
                                        name = '', 
                                        uid = 56, ), 
                                    image = '', 
                                    image_pull_secret = '', 
                                    life_cycle = fluxoperator.models.life_cycle.LifeCycle(
                                        post_start_exec = '', ), 
                                    name = '', 
                                    ports = [
                                        56
                                        ], 
                                    pre_command = '', 
                                    pull_always = True, 
                                    resources = fluxoperator.models.container_resources.ContainerResources(
                                        limits = {
                                            'key' : None
                                            }, 
                                        requests = {
                                            'key' : None
                                            }, ), 
                                    run_flux = True, 
                                    volumes = {
                                        'key' : fluxoperator.models.container_volume.ContainerVolume(
                                            path = '', 
                                            read_only = True, )
                                        }, 
                                    working_dir = '', )
                                ], 
                            deadline_seconds = 56, 
                            flux_restful = fluxoperator.models.flux_restful.FluxRestful(
                                branch = '', 
                                port = 56, 
                                token = '', 
                                username = '', ), 
                            job_labels = {
                                'key' : ''
                                }, 
                            logging = fluxoperator.models.logging_spec.LoggingSpec(
                                debug = True, 
                                quiet = True, 
                                strict = True, 
                                timed = True, ), 
                            pod = fluxoperator.models.pod_spec.PodSpec(
                                annotations = {
                                    'key' : ''
                                    }, 
                                labels = {
                                    'key' : ''
                                    }, ), 
                            size = 56, 
                            tasks = 56, 
                            users = [
                                fluxoperator.models.mini_cluster_user.MiniClusterUser(
                                    name = '', 
                                    password = '', )
                                ], 
                            volumes = {
                                'key' : fluxoperator.models.mini_cluster_volume.MiniClusterVolume(
                                    capacity = '', 
                                    path = '', 
                                    secret = '', 
                                    secret_namespace = '', 
                                    storage_class = '', )
                                }, ), 
                        status = fluxoperator.models.mini_cluster_status.MiniClusterStatus(
                            conditions = [
                                None
                                ], 
                            jobid = '', ), )
                    ], 
                kind = '', 
                metadata = None
            )
        else :
            return MiniClusterList(
                items = [
                    fluxoperator.models.mini_cluster.MiniCluster(
                        api_version = '', 
                        kind = '', 
                        metadata = None, 
                        spec = fluxoperator.models.mini_cluster_spec.MiniClusterSpec(
                            cleanup = True, 
                            containers = [
                                fluxoperator.models.mini_cluster_container.MiniClusterContainer(
                                    command = '', 
                                    commands = fluxoperator.models.commands.Commands(
                                        pre = '', 
                                        run_flux_as_root = True, ), 
                                    cores = 56, 
                                    diagnostics = True, 
                                    environment = {
                                        'key' : ''
                                        }, 
                                    flux_log_level = 56, 
                                    flux_option_flags = '', 
                                    flux_user = fluxoperator.models.flux_user.FluxUser(
                                        name = '', 
                                        uid = 56, ), 
                                    image = '', 
                                    image_pull_secret = '', 
                                    life_cycle = fluxoperator.models.life_cycle.LifeCycle(
                                        post_start_exec = '', ), 
                                    name = '', 
                                    ports = [
                                        56
                                        ], 
                                    pre_command = '', 
                                    pull_always = True, 
                                    resources = fluxoperator.models.container_resources.ContainerResources(
                                        limits = {
                                            'key' : None
                                            }, 
                                        requests = {
                                            'key' : None
                                            }, ), 
                                    run_flux = True, 
                                    volumes = {
                                        'key' : fluxoperator.models.container_volume.ContainerVolume(
                                            path = '', 
                                            read_only = True, )
                                        }, 
                                    working_dir = '', )
                                ], 
                            deadline_seconds = 56, 
                            flux_restful = fluxoperator.models.flux_restful.FluxRestful(
                                branch = '', 
                                port = 56, 
                                token = '', 
                                username = '', ), 
                            job_labels = {
                                'key' : ''
                                }, 
                            logging = fluxoperator.models.logging_spec.LoggingSpec(
                                debug = True, 
                                quiet = True, 
                                strict = True, 
                                timed = True, ), 
                            pod = fluxoperator.models.pod_spec.PodSpec(
                                annotations = {
                                    'key' : ''
                                    }, 
                                labels = {
                                    'key' : ''
                                    }, ), 
                            size = 56, 
                            tasks = 56, 
                            users = [
                                fluxoperator.models.mini_cluster_user.MiniClusterUser(
                                    name = '', 
                                    password = '', )
                                ], 
                            volumes = {
                                'key' : fluxoperator.models.mini_cluster_volume.MiniClusterVolume(
                                    capacity = '', 
                                    path = '', 
                                    secret = '', 
                                    secret_namespace = '', 
                                    storage_class = '', )
                                }, ), 
                        status = fluxoperator.models.mini_cluster_status.MiniClusterStatus(
                            conditions = [
                                None
                                ], 
                            jobid = '', ), )
                    ],
        )

    def testMiniClusterList(self):
        """Test MiniClusterList"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
