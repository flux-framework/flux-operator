# coding: utf-8

"""
    fluxoperator

    Python SDK for Flux-Operator

    The version of the OpenAPI document: v1alpha2
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


from __future__ import annotations
import pprint
import re  # noqa: F401
import json

from pydantic import BaseModel, ConfigDict, Field, StrictBool
from typing import Any, ClassVar, Dict, List, Optional
from typing import Optional, Set
from typing_extensions import Self

class LoggingSpec(BaseModel):
    """
    LoggingSpec
    """ # noqa: E501
    debug: Optional[StrictBool] = Field(default=False, description="Debug mode adds extra verbosity to Flux")
    quiet: Optional[StrictBool] = Field(default=False, description="Quiet mode silences all output so the job only shows the test running")
    strict: Optional[StrictBool] = Field(default=False, description="Strict mode ensures any failure will not continue in the job entrypoint")
    timed: Optional[StrictBool] = Field(default=False, description="Timed mode adds timing to Flux commands")
    zeromq: Optional[StrictBool] = Field(default=False, description="Enable Zeromq logging")
    __properties: ClassVar[List[str]] = ["debug", "quiet", "strict", "timed", "zeromq"]

    model_config = ConfigDict(
        populate_by_name=True,
        validate_assignment=True,
        protected_namespaces=(),
    )


    def to_str(self) -> str:
        """Returns the string representation of the model using alias"""
        return pprint.pformat(self.model_dump(by_alias=True))

    def to_json(self) -> str:
        """Returns the JSON representation of the model using alias"""
        # TODO: pydantic v2: use .model_dump_json(by_alias=True, exclude_unset=True) instead
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str: str) -> Optional[Self]:
        """Create an instance of LoggingSpec from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self) -> Dict[str, Any]:
        """Return the dictionary representation of the model using alias.

        This has the following differences from calling pydantic's
        `self.model_dump(by_alias=True)`:

        * `None` is only added to the output dict for nullable fields that
          were set at model initialization. Other fields with value `None`
          are ignored.
        """
        excluded_fields: Set[str] = set([
        ])

        _dict = self.model_dump(
            by_alias=True,
            exclude=excluded_fields,
            exclude_none=True,
        )
        return _dict

    @classmethod
    def from_dict(cls, obj: Optional[Dict[str, Any]]) -> Optional[Self]:
        """Create an instance of LoggingSpec from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return cls.model_validate(obj)

        _obj = cls.model_validate({
            "debug": obj.get("debug") if obj.get("debug") is not None else False,
            "quiet": obj.get("quiet") if obj.get("quiet") is not None else False,
            "strict": obj.get("strict") if obj.get("strict") is not None else False,
            "timed": obj.get("timed") if obj.get("timed") is not None else False,
            "zeromq": obj.get("zeromq") if obj.get("zeromq") is not None else False
        })
        return _obj


