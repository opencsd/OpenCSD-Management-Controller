package controller

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	manager "opencsd-controller/src/manager"
)

type CreateInstanceInfo struct {
	InstanceName  string `json:"instanceName"`
	InstanceType  string `json:"instanceType"`
	DbRootPwd     string `json:"dbRootPwd"`
	DbName        string `json:"dbName"`
	OperationNode string `json:"operationNode"`
	StorageNode   string `json:"storageNode"`
	VolumeName    string `json:"volumeName"`
	Validator     bool   `json:"validator"`
}

func getRandomNodeport() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(32768-30000) + 30000
}

func CreateNamespace(namespace string) error {
	_, err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		return nil
	}

	newNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err = manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().Namespaces().Create(context.TODO(), newNamespace, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create namespace '%s': %v", namespace, err)
	}
	return nil
}

func CreateQueryEngineDeployment(instanceInfo CreateInstanceInfo) error {
	nodePort := getRandomNodeport()

	err := createNodePortService(instanceInfo, nodePort, "query-engine")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	err = createQueryEngineDeployment(instanceInfo)
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Deployment created successfully")

	return nil
}

func CreateStorageEngineDeployment(instanceInfo CreateInstanceInfo) error {
	err := createMysqlConfigMap(instanceInfo, "storage-engine-instance")
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Configmap created successfully")

	err = createClusterIpService(instanceInfo, "storage-engine")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	err = createClusterIpService(instanceInfo, "storage-engine-dbms")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	nodePort := getRandomNodeport()
	err = createNodePortService(instanceInfo, nodePort, "storage-engine")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	err = createStorageEngineDeployment(instanceInfo)
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Deployment created successfully")

	return nil
}

func CreateValidatorDeployment(instanceInfo CreateInstanceInfo) error {
	err := createClusterIpService(instanceInfo, "validator")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	err = createValidatorDeployment(instanceInfo)
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Deployment created successfully")

	return nil
}

func CreateMySQLDeployment(instanceInfo CreateInstanceInfo) error {
	err := createMysqlConfigMap(instanceInfo, "mysql")
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Configmap created successfully")

	err = createClusterIpService(instanceInfo, "mysql")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	err = createMysqlDeployment(instanceInfo)
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Deployment created successfully")

	return nil
}

func CreateGraphDBDeployment(instanceInfo CreateInstanceInfo) error {
	err := createClusterIpService(instanceInfo, "graphdb")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		return err
	}
	fmt.Println("Service created successfully")

	err = createGraphdbDeployment(instanceInfo)
	if err != nil {
		fmt.Printf("Error creating deployment: %v\n", err)
		return err
	}
	fmt.Println("Deployment created successfully")

	return nil
}

func DeleteInstance(instanceName string) error {
	err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().Namespaces().Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
	return err
}

func DeleteQueryEngineDeployment(instanceName string) error {
	err := deleteService(instanceName, "query-engine-instance-svc")
	if err != nil {
		fmt.Printf("Error deleting Query Engine Service: %v\n", err)
		return err
	}

	err = deleteDeployment(instanceName, "query-engine-instance")
	if err != nil {
		fmt.Printf("Error deleting Query Engine Deployment: %v\n", err)
		return err
	}
	return nil
}

func DeleteStorageEngineDeployment(instanceName string) error {
	err := deleteService(instanceName, "storage-engine-instance-svc")
	if err != nil {
		fmt.Printf("Error deleting Storage Engine Service: %v\n", err)
		return err
	}

	err = deleteService(instanceName, "storage-engine-dbms-svc")
	if err != nil {
		fmt.Printf("Error deleting Storage Engine DBMS Service: %v\n", err)
		return err
	}

	err = deleteConfigmap(instanceName, "mysql-config")
	if err != nil {
		fmt.Printf("Error deleting Storage Engine Configmap: %v\n", err)
		return err
	}

	err = deleteDeployment(instanceName, "storage-engine-instance")
	if err != nil {
		fmt.Printf("Error deleting Storage Engine Deployment: %v\n", err)
		return err
	}

	return nil
}

func DeleteValidatorDeployment(instanceName string) error {
	err := deleteService(instanceName, "validator-svc")
	if err != nil {
		fmt.Printf("Error deletingValidator Service: %v\n", err)
		return err
	}

	err = deleteDeployment(instanceName, "validator")
	if err != nil {
		fmt.Printf("Error deleting Validator Deployment: %v\n", err)
		return err
	}
	return nil
}

func DeleteMySQLDeployment(instanceName string) error {
	err := deleteService(instanceName, "mysql-dbms-svc")
	if err != nil {
		fmt.Printf("Error deleting MySQL Service: %v\n", err)
		return err
	}

	err = deleteConfigmap(instanceName, "mysql-config")
	if err != nil {
		fmt.Printf("Error deleting MySQL Configmap: %v\n", err)
		return err
	}

	err = deleteStatefulset(instanceName, "mysql")
	if err != nil {
		fmt.Printf("Error deleting MySQL Statefulset: %v\n", err)
		return err
	}

	return nil
}

func DeleteGraphDBDeployment(instanceName string) error {
	err := deleteService(instanceName, "graphdb-dbms-svc")
	if err != nil {
		fmt.Printf("Error deleting GraphDB Service: %v\n", err)
		return err
	}

	err = deleteStatefulset(instanceName, "graphdb")
	if err != nil {
		fmt.Printf("Error deleting GraphDB Statefulset: %v\n", err)
		return err
	}

	return nil
}

func createQueryEngineDeployment(instanceInfo CreateInstanceInfo) error {
	namespace := instanceInfo.InstanceName
	storageEngineDns := "storage-engine-instance-svc." + namespace + ".svc.cluster.local"
	nodeName := instanceInfo.OperationNode
	instanceType := instanceInfo.InstanceType

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "query-engine-instance",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "query-engine-instance",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "query-engine-instance",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "query-engine-instance",
						"tier": "opencsd",
					},
				},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					NodeSelector: map[string]string{
						"layer": "operation",
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: "regcred"},
					},
					Containers: []corev1.Container{
						{
							Name:            "query-engine",
							Image:           "ketidevit2/query-engine-instance:v3.0",
							ImagePullPolicy: corev1.PullAlways,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tz-config",
									MountPath: "/etc/localtime",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("128Mi"),
									corev1.ResourceCPU:    resourceQuantity("250m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("1Gi"),
									corev1.ResourceCPU:    resourceQuantity("1"),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "LOG_LEVEL",
									Value: "INFO",
								},
								{
									Name:  "INSTANCE_TYPE",
									Value: instanceType,
								},
								{
									Name:  "STORAGE_ENGINE_DNS",
									Value: storageEngineDns,
								},
								{
									Name: "INSTANCE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tz-config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/usr/share/zoneinfo/Asia/Seoul",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func createStorageEngineDeployment(instanceInfo CreateInstanceInfo) error {
	namespace := instanceInfo.InstanceName
	storageNode := instanceInfo.StorageNode
	storageNodeIp := manager.InstanceManager_.StorageLayer[storageNode].IP
	nodeName := instanceInfo.OperationNode
	volumeName := instanceInfo.VolumeName
	volumePath := manager.InstanceManager_.VolumeInfo[volumeName].VolumePath + "/" + namespace

	fmt.Println("createStorageEngineDeployment ", storageNodeIp, volumePath)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "storage-engine-instance",
			Namespace: namespace,
			Labels: map[string]string{
				"app": "storage-engine-instance",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "storage-engine-instance"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "storage-engine-instance",
						"tier": "opencsd",
					},
				},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					NodeSelector: map[string]string{
						"layer": "operation",
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: "regcred"},
					},
					Containers: []corev1.Container{
						{
							Name:  "interface-module",
							Image: "ketidevit2/storage-engine-interface:v3.0",
							Env: []corev1.EnvVar{
								{Name: "LOG_LEVEL", Value: "Info"},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("128Mi"),
									corev1.ResourceCPU:    resourceQuantity("250m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("1Gi"),
									corev1.ResourceCPU:    resourceQuantity("1"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tz-config",
									MountPath: "/etc/localtime",
								},
							},
						},
						{
							Name:  "offloading-module",
							Image: "ketidevit2/offloading-module:v3.0",
							Env: []corev1.EnvVar{
								{Name: "LOG_LEVEL", Value: "INFO"},
								{Name: "STORAGE_NODE_IP", Value: storageNodeIp},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("128Mi"),
									corev1.ResourceCPU:    resourceQuantity("250m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("1Gi"),
									corev1.ResourceCPU:    resourceQuantity("1"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tz-config",
									MountPath: "/etc/localtime",
								},
							},
						},
						{
							Name:  "merging-module",
							Image: "ketidevit2/merging-module:v3.0",
							Env: []corev1.EnvVar{
								{Name: "LOG_LEVEL", Value: "INFO"},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("1Gi"),
									corev1.ResourceCPU:    resourceQuantity("1"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("50Gi"),
									corev1.ResourceCPU:    resourceQuantity("3"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tz-config",
									MountPath: "/etc/localtime",
								},
							},
						},
						{
							Name:  "myrocks",
							Image: "percona:5.7.43",
							Env: []corev1.EnvVar{
								{Name: "MYSQL_ROOT_PASSWORD", Value: instanceInfo.DbRootPwd},
								{Name: "MYSQL_DATABASE", Value: instanceInfo.DbName},
								{Name: "INIT_ROCKSDB", Value: "yes"},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("1G"),
									corev1.ResourceCPU:    resourceQuantity("200m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("10G"),
									corev1.ResourceCPU:    resourceQuantity("500m"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/var/lib/mysql",
								},
								{
									Name:      "myrocks-cm",
									MountPath: "/etc/mysql",
								},
								{
									Name:      "tz-config",
									MountPath: "/etc/localtime",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tz-config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/usr/share/zoneinfo/Asia/Seoul",
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/root/.kube",
								},
							},
						},
						{
							Name: volumeName, // 가변
							VolumeSource: corev1.VolumeSource{
								NFS: &corev1.NFSVolumeSource{
									Server: storageNodeIp,
									Path:   volumePath,
								},
							},
						},
						{
							Name: "myrocks-cm",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mysql-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func createValidatorDeployment(instanceInfo CreateInstanceInfo) error {
	namespace := instanceInfo.InstanceName
	nodeName := instanceInfo.OperationNode

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "validator",
			Namespace: namespace,
			Labels: map[string]string{
				"app":  "validator",
				"tier": "opencsd",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "validator",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "validator",
					},
				},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					NodeSelector: map[string]string{
						"layer": "operation",
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{Name: "regcred"},
					},
					Containers: []corev1.Container{
						{
							Name:  "validator",
							Image: "ketidevit2/validator:v1.0",
							Env: []corev1.EnvVar{
								{
									Name:  "LOG_LEVEL",
									Value: "INFO",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("128Mi"),
									corev1.ResourceCPU:    resourceQuantity("250m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resourceQuantity("1Gi"),
									corev1.ResourceCPU:    resourceQuantity("1"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "tz-config",
									MountPath: "/etc/localtime",
								},
								{
									Name:      "config",
									MountPath: "/root/.kube/",
								},
							},
							ImagePullPolicy: corev1.PullAlways,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "tz-config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/usr/share/zoneinfo/Asia/Seoul",
								},
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/root/.kube",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	return err
}

func createMysqlDeployment(instanceInfo CreateInstanceInfo) error {
	namespace := instanceInfo.InstanceName
	storageNode := instanceInfo.StorageNode
	storageNodeIp := manager.InstanceManager_.StorageLayer[storageNode].IP
	nodeName := instanceInfo.OperationNode
	volumeName := instanceInfo.VolumeName
	volumePath := manager.InstanceManager_.VolumeInfo[volumeName].VolumePath + "/" + namespace

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql",
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    int32Ptr(1),
			ServiceName: "mysql-dbms-svc", // Headless Service 이름
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "mysql",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "mysql",
						"tier": "opencsd",
					},
				},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					NodeSelector: map[string]string{
						"key": "worker",
					},
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "percona:5.7.43",
							Env: []corev1.EnvVar{
								{Name: "MYSQL_DATABASE", Value: "keti_opencsd"},
								{Name: "MYSQL_PASSWORD", Value: "ketilinux"},
								{Name: "MYSQL_ROOT_PASSWORD", Value: "ketilinux"},
								{Name: "MYSQL_USER", Value: "keti"},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 3306},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/var/lib/mysql",
								},
								{
									Name:      "mysql-cm",
									MountPath: "/etc/mysql",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: volumeName,
							VolumeSource: corev1.VolumeSource{
								NFS: &corev1.NFSVolumeSource{
									Server: storageNodeIp,
									Path:   volumePath,
								},
							},
						},
						{
							Name: "mysql-cm",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mysql-config",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	return err
}

func createGraphdbDeployment(instanceInfo CreateInstanceInfo) error {
	namespace := instanceInfo.InstanceName
	storageNode := instanceInfo.StorageNode
	storageNodeIp := manager.InstanceManager_.StorageLayer[storageNode].IP
	nodeName := instanceInfo.OperationNode
	volumeName := instanceInfo.VolumeName
	volumePath := manager.InstanceManager_.VolumeInfo[volumeName].VolumePath + "/" + namespace

	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "graphdb",
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    int32Ptr(1),
			ServiceName: "graphdb-dbms-svc", // Headless Service 이름
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "graphdb",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  "graphdb",
						"tier": "opencsd",
					},
				},
				Spec: corev1.PodSpec{
					NodeName: nodeName,
					NodeSelector: map[string]string{
						"layer": "operation",
					},
					Containers: []corev1.Container{
						{
							Name:  "neo4j",
							Image: "neo4j:5.8",
							Ports: []corev1.ContainerPort{
								{ContainerPort: 7474}, // HTTP Browser
								{ContainerPort: 7687}, // Bolt Protocol
							},
							Env: []corev1.EnvVar{
								{
									Name:  "NEO4J_AUTH",
									Value: "neo4j/neo4jpass", // 사용자 인증 정보
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      volumeName,
									MountPath: "/data",
								},
								{
									Name:      "graphdb-logs",
									MountPath: "/logs",
								},
								{
									Name:      "graphdb-import",
									MountPath: "/import",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: volumeName,
							VolumeSource: corev1.VolumeSource{
								NFS: &corev1.NFSVolumeSource{
									Server: storageNodeIp,
									Path:   volumePath + "/data",
								},
							},
						},
						{
							Name: "graphdb-logs",
							VolumeSource: corev1.VolumeSource{
								NFS: &corev1.NFSVolumeSource{
									Server: storageNodeIp,
									Path:   volumePath + "/logs",
								},
							},
						},
						{
							Name: "graphdb-import",
							VolumeSource: corev1.VolumeSource{
								NFS: &corev1.NFSVolumeSource{
									Server: storageNodeIp,
									Path:   volumePath + "/import",
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().StatefulSets(namespace).Create(context.TODO(), statefulSet, metav1.CreateOptions{})
	return err
}

func createNodePortService(instanceInfo CreateInstanceInfo, nodePort int, kind string) error {
	namespace := instanceInfo.InstanceName
	label, name := "", ""
	var ports []corev1.ServicePort

	if kind == "query-engine" {
		label = "query-engine-instance"
		name = "query-engine-instance-svc"
		ports = []corev1.ServicePort{
			{
				Name:       "main",
				Port:       40100,
				TargetPort: intstr.FromInt(40100),
				NodePort:   int32(nodePort),
				Protocol:   corev1.ProtocolTCP,
			},
		}
	} else if kind == "storage-engine" {
		label = "storage-engine-instance"
		name = "buffer-manager-svc"
		ports = []corev1.ServicePort{
			{
				Port:       40204,
				TargetPort: intstr.FromInt(40204),
				NodePort:   int32(nodePort),
				Protocol:   corev1.ProtocolTCP,
			},
		}
	} else {
		return fmt.Errorf("undefined dbms type")
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": label,
			},
			Ports: ports,
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})

	return err
}

func createClusterIpService(instanceInfo CreateInstanceInfo, kind string) error {
	namespace := instanceInfo.InstanceName
	label, name := "", ""
	var ports []corev1.ServicePort

	if kind == "storage-engine" {
		label = "storage-engine-instance"
		name = "storage-engine-interface-svc"
		ports = []corev1.ServicePort{
			{
				Port:       40200,
				TargetPort: intstr.FromInt(40200),
				Protocol:   corev1.ProtocolTCP,
			},
		}
	} else if kind == "storage-engine-dbms" {
		label = "storage-engine-instance"
		name = "storage-engine-dbms-svc"
		ports = []corev1.ServicePort{
			{
				Name:       "main",
				Port:       3306,
				TargetPort: intstr.FromInt(3306),
				Protocol:   corev1.ProtocolTCP,
			},
		}
	} else if kind == "mysql" {
		label = "mysql"
		name = "mysql-dbms-svc"
		ports = []corev1.ServicePort{
			{
				Name:       "main",
				Port:       3306,
				TargetPort: intstr.FromInt(3306),
				Protocol:   corev1.ProtocolTCP,
			},
		}
	} else if kind == "graphdb" {
		label = "graphdb"
		name = "graphdb-dbms-svc"
		ports = []corev1.ServicePort{
			{
				Name:       "browser",
				Port:       7474,
				TargetPort: intstr.FromInt(7474),
			},
			{
				Name:       "bolt",
				Port:       7687,
				TargetPort: intstr.FromInt(7687),
			},
		}
	} else if kind == "validator" {
		label = "validator"
		name = "validator-svc"
		ports = []corev1.ServicePort{
			{
				Name:       "main",
				Port:       40000,
				TargetPort: intstr.FromInt(40000),
				Protocol:   corev1.ProtocolTCP,
			},
			{
				Name:       "qemu",
				Port:       40001,
				TargetPort: intstr.FromInt(40001),
				Protocol:   corev1.ProtocolTCP,
			},
		}
	} else {
		return fmt.Errorf("undefined dbms type")
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": label,
			},
			Ports: ports,
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}

func createMysqlConfigMap(instanceInfo CreateInstanceInfo, label string) error {
	namespace := instanceInfo.InstanceName

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-config",
			Namespace: namespace,
			Labels: map[string]string{
				"app": label,
			},
		},
		Data: map[string]string{
			"my.cnf": `[mysqld]
user=mysql
plugin-load-add = ha_rocksdb.so
default-storage-engine=rocksdb
`,
		},
	}

	_, err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMap, metav1.CreateOptions{})
	return err
}

func deleteDeployment(namespace string, deploymentName string) error {
	err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	return err
}

func deleteService(namespace string, serviceName string) error {
	err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	return err
}

func deleteStatefulset(namespace string, statefulSetName string) error {
	err := manager.InstanceManager_.ClusterConfig.Clientset.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	return err
}

func deleteConfigmap(namespace string, configMapName string) error {
	err := manager.InstanceManager_.ClusterConfig.Clientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configMapName, metav1.DeleteOptions{})
	return err
}

func resourceQuantity(val string) resource.Quantity {
	quantity, err := resource.ParseQuantity(val)
	if err != nil {
		panic(err)
	}
	return quantity
}

func int32Ptr(i int32) *int32 { return &i }
