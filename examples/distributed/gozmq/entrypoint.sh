#!/bin/bash

workers="${1:-2}"
rank=${FLUX_TASK_RANK}

# Get the host name
host=$(hostname)
echo "Hello I'm host $host"
go run /code/main.go run --size ${workers} --prefix flux-sample --suffix "flux-service.default.svc.cluster.local" --port 5555 --rank ${rank}
