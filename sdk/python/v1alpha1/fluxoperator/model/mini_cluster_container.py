"""
    fluxoperator

    Python SDK for Flux-Operator  # noqa: E501

    The version of the OpenAPI document: v1alpha1
    Generated by: https://openapi-generator.tech
"""


import re  # noqa: F401
import sys  # noqa: F401

from fluxoperator.model_utils import (  # noqa: F401
    ApiTypeError,
    ModelComposed,
    ModelNormal,
    ModelSimple,
    cached_property,
    change_keys_js_to_python,
    convert_js_args_to_python_args,
    date,
    datetime,
    file_type,
    none_type,
    validate_get_composed_info,
)

def lazy_import():
    from fluxoperator.model.commands import Commands
    from fluxoperator.model.container_resources import ContainerResources
    from fluxoperator.model.container_volume import ContainerVolume
    from fluxoperator.model.flux_user import FluxUser
    from fluxoperator.model.life_cycle import LifeCycle
    globals()['Commands'] = Commands
    globals()['ContainerResources'] = ContainerResources
    globals()['ContainerVolume'] = ContainerVolume
    globals()['FluxUser'] = FluxUser
    globals()['LifeCycle'] = LifeCycle


class MiniClusterContainer(ModelNormal):
    """NOTE: This class is auto generated by OpenAPI Generator.
    Ref: https://openapi-generator.tech

    Do not edit the class manually.

    Attributes:
      allowed_values (dict): The key is the tuple path to the attribute
          and the for var_name this is (var_name,). The value is a dict
          with a capitalized key describing the allowed value and an allowed
          value. These dicts store the allowed enum values.
      attribute_map (dict): The key is attribute name
          and the value is json key in definition.
      discriminator_value_class_map (dict): A dict to go from the discriminator
          variable value to the discriminator class name.
      validations (dict): The key is the tuple path to the attribute
          and the for var_name this is (var_name,). The value is a dict
          that stores validations for max_length, min_length, max_items,
          min_items, exclusive_maximum, inclusive_maximum, exclusive_minimum,
          inclusive_minimum, and regex.
      additional_properties_type (tuple): A tuple of classes accepted
          as additional properties values.
    """

    allowed_values = {
    }

    validations = {
    }

    additional_properties_type = None

    _nullable = False

    @cached_property
    def openapi_types():
        """
        This must be a method because a model may have properties that are
        of type self, this must run after the class is loaded

        Returns
            openapi_types (dict): The key is attribute name
                and the value is attribute type.
        """
        lazy_import()
        return {
            'image': (str,),  # noqa: E501
            'command': (str,),  # noqa: E501
            'commands': (Commands,),  # noqa: E501
            'cores': (int,),  # noqa: E501
            'diagnostics': (bool,),  # noqa: E501
            'environment': ({str: (str,)},),  # noqa: E501
            'flux_log_level': (int,),  # noqa: E501
            'flux_option_flags': (str,),  # noqa: E501
            'flux_user': (FluxUser,),  # noqa: E501
            'image_pull_secret': (str,),  # noqa: E501
            'life_cycle': (LifeCycle,),  # noqa: E501
            'name': (str,),  # noqa: E501
            'ports': ([int],),  # noqa: E501
            'pre_command': (str,),  # noqa: E501
            'pull_always': (bool,),  # noqa: E501
            'resources': (ContainerResources,),  # noqa: E501
            'run_flux': (bool,),  # noqa: E501
            'volumes': ({str: (ContainerVolume,)},),  # noqa: E501
            'working_dir': (str,),  # noqa: E501
        }

    @cached_property
    def discriminator():
        return None


    attribute_map = {
        'image': 'image',  # noqa: E501
        'command': 'command',  # noqa: E501
        'commands': 'commands',  # noqa: E501
        'cores': 'cores',  # noqa: E501
        'diagnostics': 'diagnostics',  # noqa: E501
        'environment': 'environment',  # noqa: E501
        'flux_log_level': 'fluxLogLevel',  # noqa: E501
        'flux_option_flags': 'fluxOptionFlags',  # noqa: E501
        'flux_user': 'fluxUser',  # noqa: E501
        'image_pull_secret': 'imagePullSecret',  # noqa: E501
        'life_cycle': 'lifeCycle',  # noqa: E501
        'name': 'name',  # noqa: E501
        'ports': 'ports',  # noqa: E501
        'pre_command': 'preCommand',  # noqa: E501
        'pull_always': 'pullAlways',  # noqa: E501
        'resources': 'resources',  # noqa: E501
        'run_flux': 'runFlux',  # noqa: E501
        'volumes': 'volumes',  # noqa: E501
        'working_dir': 'workingDir',  # noqa: E501
    }

    _composed_schemas = {}

    required_properties = set([
        '_data_store',
        '_check_type',
        '_spec_property_naming',
        '_path_to_item',
        '_configuration',
        '_visited_composed_classes',
    ])

    @convert_js_args_to_python_args
    def __init__(self, *args, **kwargs):  # noqa: E501
        """MiniClusterContainer - a model defined in OpenAPI

        Args:

        Keyword Args:
            image (str): Container image must contain flux and flux-sched install. defaults to ""  # noqa: E501
            _check_type (bool): if True, values for parameters in openapi_types
                                will be type checked and a TypeError will be
                                raised if the wrong type is input.
                                Defaults to True
            _path_to_item (tuple/list): This is a list of keys or values to
                                drill down to the model in received_data
                                when deserializing a response
            _spec_property_naming (bool): True if the variable names in the input data
                                are serialized names, as specified in the OpenAPI document.
                                False if the variable names in the input data
                                are pythonic names, e.g. snake case (default)
            _configuration (Configuration): the instance to use when
                                deserializing a file_type parameter.
                                If passed, type conversion is attempted
                                If omitted no type conversion is done.
            _visited_composed_classes (tuple): This stores a tuple of
                                classes that we have traveled through so that
                                if we see that class again we will not use its
                                discriminator again.
                                When traveling through a discriminator, the
                                composed schema that is
                                is traveled through is added to this set.
                                For example if Animal has a discriminator
                                petType and we pass in "Dog", and the class Dog
                                allOf includes Animal, we move through Animal
                                once using the discriminator, and pick Dog.
                                Then in Dog, we will make an instance of the
                                Animal class but this time we won't travel
                                through its discriminator because we passed in
                                _visited_composed_classes = (Animal,)
            command (str): Single user executable to provide to flux start IMPORTANT: This is left here, but not used in favor of exposing Flux via a Restful API. We Can remove this when that is finalized.. [optional] if omitted the server will use the default value of ""  # noqa: E501
            commands (Commands): [optional]  # noqa: E501
            cores (int): Cores the container should use. [optional] if omitted the server will use the default value of 0  # noqa: E501
            diagnostics (bool): Run flux diagnostics on start instead of command. [optional] if omitted the server will use the default value of False  # noqa: E501
            environment ({str: (str,)}): Key/value pairs for the environment. [optional]  # noqa: E501
            flux_log_level (int): Log level to use for flux logging (only in non TestMode). [optional] if omitted the server will use the default value of 0  # noqa: E501
            flux_option_flags (str): Flux option flags, usually provided with -o optional - if needed, default option flags for the server These can also be set in the user interface to override here. This is only valid for a FluxRunner \"runFlux\" true. [optional] if omitted the server will use the default value of ""  # noqa: E501
            flux_user (FluxUser): [optional]  # noqa: E501
            image_pull_secret (str): Allow the user to pull authenticated images By default no secret is selected. Setting this with the name of an already existing imagePullSecret will specify that secret in the pod spec.. [optional] if omitted the server will use the default value of ""  # noqa: E501
            life_cycle (LifeCycle): [optional]  # noqa: E501
            name (str): Container name is only required for non flux runners. [optional] if omitted the server will use the default value of ""  # noqa: E501
            ports ([int]): Ports to be exposed to other containers in the cluster We take a single list of integers and map to the same. [optional]  # noqa: E501
            pre_command (str): Special command to run at beginning of script, directly after asFlux is defined as sudo -u flux -E (so you can change that if desired.) This is only valid if FluxRunner is set (that writes a wait.sh script) This is for the indexed job pods and the certificate generation container.. [optional] if omitted the server will use the default value of ""  # noqa: E501
            pull_always (bool): Allow the user to dictate pulling By default we pull if not present. Setting this to true will indicate to pull always. [optional] if omitted the server will use the default value of False  # noqa: E501
            resources (ContainerResources): [optional]  # noqa: E501
            run_flux (bool): Main container to run flux (only should be one). [optional] if omitted the server will use the default value of False  # noqa: E501
            volumes ({str: (ContainerVolume,)}): Volumes that can be mounted (must be defined in volumes). [optional]  # noqa: E501
            working_dir (str): Working directory to run command from. [optional] if omitted the server will use the default value of ""  # noqa: E501
        """

        image = kwargs.get('image', "")
        _check_type = kwargs.pop('_check_type', True)
        _spec_property_naming = kwargs.pop('_spec_property_naming', False)
        _path_to_item = kwargs.pop('_path_to_item', ())
        _configuration = kwargs.pop('_configuration', None)
        _visited_composed_classes = kwargs.pop('_visited_composed_classes', ())

        if args:
            raise ApiTypeError(
                "Invalid positional arguments=%s passed to %s. Remove those invalid positional arguments." % (
                    args,
                    self.__class__.__name__,
                ),
                path_to_item=_path_to_item,
                valid_classes=(self.__class__,),
            )

        self._data_store = {}
        self._check_type = _check_type
        self._spec_property_naming = _spec_property_naming
        self._path_to_item = _path_to_item
        self._configuration = _configuration
        self._visited_composed_classes = _visited_composed_classes + (self.__class__,)

        self.image = image
        for var_name, var_value in kwargs.items():
            if var_name not in self.attribute_map and \
                        self._configuration is not None and \
                        self._configuration.discard_unknown_keys and \
                        self.additional_properties_type is None:
                # discard variable.
                continue
            setattr(self, var_name, var_value)
