#!/bin/bash

if [ "$1" == "a" ] ; then   
    while [ -z $PODNAME ]
    do
        PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-api-server`
        PODNAME="${PODNAME:4}"
    done
    kubectl logs $PODNAME -n management-controller -f 
elif [ "$1" == "e" ] ; then  
    while [ -z $PODNAME ]
    do
        PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-controller`
        PODNAME="${PODNAME:4}"
    done
    kubectl logs $PODNAME -n management-controller -f 
else 
    echo arg error
fi



