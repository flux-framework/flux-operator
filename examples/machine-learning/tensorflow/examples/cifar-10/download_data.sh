#! /bin/bash

if [[ ! -f "cifar-10-python.tar.gz" ]]; then
    wget https://www.cs.toronto.edu/~kriz/cifar-10-python.tar.gz .
    tar -xvzf cifar-10-python.tar.gz
fi
