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
from fluxoperator.models.mini_cluster import MiniCluster  # noqa: E501
from fluxoperator.rest import ApiException

class TestMiniCluster(unittest.TestCase):
    """MiniCluster unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def make_instance(self, include_optional):
        """Test MiniCluster
            include_option is a boolean, when False only required
            params are included, when True both required and
            optional params are included """
        # model = fluxoperator.models.mini_cluster.MiniCluster()  # noqa: E501
        if include_optional :
            return MiniCluster(
                api_version = '', 
                kind = '', 
                metadata = None, 
                spec = fluxoperator.models.mini_cluster_spec.MiniClusterSpec(
                    archive = fluxoperator.models.mini_cluster_archive.MiniClusterArchive(
                        path = '', ), 
                    cleanup = True, 
                    containers = [
                        fluxoperator.models.mini_cluster_container.MiniClusterContainer(
                            batch = True, 
                            batch_raw = True, 
                            command = '', 
                            commands = fluxoperator.models.commands.Commands(
                                broker_pre = '', 
                                init = '', 
                                post = '', 
                                pre = '', 
                                prefix = '', 
                                run_flux_as_root = True, 
                                worker_pre = '', ), 
                            cores = 56, 
                            diagnostics = True, 
                            environment = {
                                'key' : ''
                                }, 
                            existing_volumes = {
                                'key' : fluxoperator.models.mini_cluster_existing_volume.MiniClusterExistingVolume(
                                    claim_name = '', 
                                    config_map_name = '', 
                                    items = {
                                        'key' : ''
                                        }, 
                                    path = '', 
                                    read_only = True, 
                                    secret_name = '', )
                                }, 
                            flux_user = fluxoperator.models.flux_user.FluxUser(
                                name = 'flux', 
                                uid = 56, ), 
                            image = 'ghcr.io/rse-ops/accounting:app-latest', 
                            image_pull_secret = '', 
                            launcher = True, 
                            life_cycle = fluxoperator.models.life_cycle.LifeCycle(
                                post_start_exec = '', 
                                pre_stop_exec = '', ), 
                            logs = '', 
                            name = '', 
                            ports = [
                                56
                                ], 
                            pull_always = True, 
                            resources = fluxoperator.models.container_resources.ContainerResources(
                                limits = {
                                    'key' : None
                                    }, 
                                requests = {
                                    'key' : None
                                    }, ), 
                            run_flux = True, 
                            secrets = {
                                'key' : fluxoperator.models.secret.Secret(
                                    key = '', 
                                    name = '', )
                                }, 
                            security_context = fluxoperator.models.security_context.SecurityContext(
                                add_capabilities = [
                                    ''
                                    ], 
                                privileged = True, ), 
                            volumes = {
                                'key' : fluxoperator.models.container_volume.ContainerVolume(
                                    path = '', 
                                    read_only = True, )
                                }, 
                            working_dir = '', )
                        ], 
                    deadline_seconds = 56, 
                    flux = fluxoperator.models.flux_spec.FluxSpec(
                        broker_config = '', 
                        bursting = fluxoperator.models.bursting.Bursting(
                            clusters = [
                                fluxoperator.models.bursted_cluster.BurstedCluster(
                                    name = '', 
                                    size = 56, )
                                ], 
                            hostlist = '', 
                            lead_broker = fluxoperator.models.flux_broker.FluxBroker(
                                address = '', 
                                name = '', 
                                port = 56, 
                                size = 56, ), ), 
                        connect_timeout = '5s', 
                        curve_cert = '', 
                        curve_cert_secret = '', 
                        install_root = '/usr', 
                        log_level = 56, 
                        minimal_service = True, 
                        munge_secret = '', 
                        option_flags = '', 
                        wrap = '', ), 
                    flux_restful = fluxoperator.models.flux_restful.FluxRestful(
                        branch = 'main', 
                        port = 56, 
                        secret_key = '', 
                        token = '', 
                        username = '', ), 
                    interactive = True, 
                    job_labels = {
                        'key' : ''
                        }, 
                    logging = fluxoperator.models.logging_spec.LoggingSpec(
                        debug = True, 
                        quiet = True, 
                        strict = True, 
                        timed = True, 
                        zeromq = True, ), 
                    max_size = 56, 
                    network = fluxoperator.models.network.Network(
                        headless_name = 'flux-service', ), 
                    pod = fluxoperator.models.pod_spec.PodSpec(
                        annotations = {
                            'key' : ''
                            }, 
                        labels = {
                            'key' : ''
                            }, 
                        node_selector = {
                            'key' : ''
                            }, 
                        service_account_name = '', ), 
                    services = [
                        fluxoperator.models.mini_cluster_container.MiniClusterContainer(
                            batch = True, 
                            batch_raw = True, 
                            command = '', 
                            cores = 56, 
                            diagnostics = True, 
                            image = 'ghcr.io/rse-ops/accounting:app-latest', 
                            image_pull_secret = '', 
                            launcher = True, 
                            logs = '', 
                            name = '', 
                            pull_always = True, 
                            run_flux = True, 
                            working_dir = '', )
                        ], 
                    share_process_namespace = True, 
                    size = 56, 
                    tasks = 56, 
                    users = [
                        fluxoperator.models.mini_cluster_user.MiniClusterUser(
                            name = '', 
                            password = '', )
                        ], 
                    volumes = {
                        'key' : fluxoperator.models.mini_cluster_volume.MiniClusterVolume(
                            attributes = {
                                'key' : ''
                                }, 
                            capacity = '5Gi', 
                            claim_annotations = {
                                'key' : ''
                                }, 
                            delete = True, 
                            driver = '', 
                            path = '', 
                            secret = '', 
                            secret_namespace = 'default', 
                            storage_class = 'hostpath', 
                            volume_handle = '', )
                        }, ), 
                status = fluxoperator.models.mini_cluster_status.MiniClusterStatus(
                    conditions = [
                        None
                        ], 
                    jobid = '', 
                    maximum_size = 56, 
                    selector = '', 
                    size = 56, )
            )
        else :
            return MiniCluster(
        )

    def testMiniCluster(self):
        """Test MiniCluster"""
        inst_req_only = self.make_instance(include_optional=False)
        inst_req_and_optional = self.make_instance(include_optional=True)

if __name__ == '__main__':
    unittest.main()
