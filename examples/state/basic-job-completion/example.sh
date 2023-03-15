#!/bin/bash

# Usage: /bin/bash script/test.sh $name 30
HERE=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

set -eEu -o pipefail

# make sure we clean up old pods / data
make clean > /dev/null 2>&1
minikube ssh -- rm -rf /tmp/data/* || true

# If there is a pre-run script
printf "🌀️ Creating first MiniCluster...\n"
kubectl apply -f ${HERE}/minicluster.yaml 

printf "\n🥱️ Sleeping 20 seconds to wait for cluster..."
sleep 20

broker_pod=$(kubectl get pods -n flux-operator --no-headers -o custom-columns=":metadata.name" | grep flux-sample-0)
printf "Broker pod is ${broker_pod}\n"

printf "\n🤓️ Contents of /tmp/data in MiniKube\n"
minikube ssh ls /tmp/data

# Then submit jobs
printf "\n✨️ Submitting jobs\n"
for i in {1..5}
do
   kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux submit sleep ${i}
   kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux submit whoami
done

printf "\n🥱️ Waiting for jobs...\n"
sleep 10

printf "Jobs finished...\n"
kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux jobs -a

kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux queue stop
kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux queue idle
kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux dump /state/archive.tar.gz

printf "\n🥱️ Wait a minute to be sure we have saved...\n"
sleep 30

printf "\n🧊️ Current state directory at /var/lib/flux...\n"
kubectl exec -it -n flux-operator ${broker_pod} -- ls -l /var/lib/flux

printf "\n🧊️ Current archive directory at /state... should be empty, not saved yet\n"
kubectl exec -it -n flux-operator ${broker_pod} -- ls -l /state

printf "Cleaning up...\n"
kubectl delete -f ${HERE}/minicluster.yaml 

sleep 10
minikube ssh -- ls -l /tmp/data
sleep 10

# If there is a pre-run script
printf "\n🌀️ Creating second MiniCluster\n"
kubectl apply -f ${HERE}/minicluster.yaml 

printf "\n🥱️ Sleeping a minute to wait for cluster...\n"
sleep 60

broker_pod=$(kubectl get pods -n flux-operator --no-headers -o custom-columns=":metadata.name" | grep flux-sample-0)
printf "Broker pod is ${broker_pod}\n"

printf "\n🤓️ Contents of /tmp/data in MiniKube - should be populated with archive from first\n"
minikube ssh -- ls -l /tmp/data
sleep 5

printf "\n🤓️ Inspecting state directory in new cluster...\n"
kubectl exec -it -n flux-operator ${broker_pod} -- ls -l /var/lib/flux

printf "\n😎️ Looking to see if old job history exists...\n"
kubectl exec -it -n flux-operator ${broker_pod} -- sudo -u flux flux proxy local:///var/run/flux/local flux jobs -a

sleep 5
printf "Cleaning up..\n"
kubectl delete -f ${HERE}/minicluster.yaml
minikube ssh -- rm -rf /tmp/data/* || true
