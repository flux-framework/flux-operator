#!/usr/bin/env python3

# We aren't able to set defaults in the annotations, so this small script
# reads in the swagger file and sets them.

import json
import sys
import os

# Updates organized by parent and then definition and property
updates = {
    "LoggingSpec": {"strict": True},
    "MiniClusterContainer": {
        "cores": 1,
        "image": "ghcr.io/rse-ops/accounting:app-latest",
    },
    "MiniClusterSpec": {"deadlineSeconds": 31500000, "tasks": 1, "size": 1},
    "MiniClusterVolume": {"capacity": "5Gi", "secretNamespace": "default"},
}


def main():
    if len(sys.argv) < 2:
        sys.exit("Please provide the swagger.json as the only input")
    swagger_file = os.path.abspath(sys.argv[1])
    if not os.path.exists(swagger_file):
        sys.exit("Swagger file {swagger_file} does not exist.")

    # Read in the swagger to json
    with open(swagger_file, "r") as fd:
        data = json.loads(fd.read())

    for name, definition in data.get("definitions", {}).items():
        if name not in updates:
            continue
        print(f"Looking to update defaults for parent: {name}")
        parent = updates[name]
        for prop, meta in definition.get("properties", {}).items():
            if prop in parent:
                default = parent[prop]
                print(f'    "{prop}" has default {default}')
                meta["default"] = default

    # Save back to file
    with open(swagger_file, "w") as fd:
        fd.write(json.dumps(data, indent=4))


if __name__ == "__main__":
    main()
