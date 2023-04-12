#!/usr/bin/env python3

from distutils.core import setup

setup(
    name="tensorflow_flux",
    version="0.1.0",
    description="Run distributed TensorFlow on the Flux Operator",
    author="Vanessasaurus",
    author_email="sochat1@llnl.gov",
    url="https://github.com/flux-framework/flux-operator/tree/main/examples/machine-learning/tensorflow",
    packages=["tensorflow_flux"],
)
