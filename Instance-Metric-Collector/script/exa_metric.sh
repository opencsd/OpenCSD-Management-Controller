#!/bin/bash

target="hpc-metric-collector"
namespace="keti-system"

kubectl get pods -n $namespace -l name=$target -o wide

# pod_list=$(kubectl get pods -o=name -n $namespace | sed "s/^.\{4\}//")
# for pod_name in $pod_list; do
#   echo "------------------------------------"
#   if [[ "$pod_name" == *$target* ]]; then
    
#     kubectl logs --tail=32 $pod_name -n $namespace
#     echo "------------------------------------"
#   fi
# done

node_names=$(kubectl get nodes -o jsonpath='{.items[*].metadata.name}')

for node in $node_names; do
    pod_name=$(kubectl get pods -n $namespace --field-selector=status.phase=Running -o jsonpath="{.items[0].metadata.name}" --selector=name=$target --node=$node)
    echo "Node: $node, HPC Metric Collector Pod: $pod_name"
done
