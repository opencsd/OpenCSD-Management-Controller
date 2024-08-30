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
        PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-engine-deployer`
        PODNAME="${PODNAME:4}"
    done
    kubectl logs $PODNAME -n management-controller -f 
elif [ "$1" == "i" ] ; then  
    while [ -z $PODNAME ]
    do
        PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-instance-manager`
        PODNAME="${PODNAME:4}"
    done
    kubectl logs $PODNAME -n management-controller -f
elif [ "$1" == "v" ] ; then  
    while [ -z $PODNAME ]
    do
        PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-volume-allocator`
        PODNAME="${PODNAME:4}"
    done
    kubectl logs $PODNAME -n management-controller -f
elif [ "$1" == "mc" ] ; then  
    while [ -z $PODNAME ]
    do
        PODNAME=`kubectl get po -o=name -A --field-selector=status.phase=Running | grep opencsd-metric-collector`
        PODNAME="${PODNAME:4}"
    done
    kubectl logs $PODNAME -n management-controller -f 
else 
    echo arg error
fi



