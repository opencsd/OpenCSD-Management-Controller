#!/bin/bash 

while [ -z $PODNAME ]
do
    PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-instance-manager`
    PODNAME="${PODNAME:4}"
done

kubectl logs $PODNAME -n management-controller -f 



