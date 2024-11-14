#$1 create/c or delete/d

if [ "$1" == "delete" ] || [ "$1" == "d" ]; then   
    echo kubectl delete -f deploy/
    kubectl delete -f deploy/
else
    echo kubectl create -f deploy/
    kubectl create -f deploy/
fi