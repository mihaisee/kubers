package cluster

import (
	"fmt"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RemovePoFromDeployment(deploymentName string, namespace string, oneByOne bool, wait bool, surge int) {
	cl := getKubernetesClientOutsideCluster()

	options := metaV1.ListOptions{}
	deployment, _ := cl.AppsV1().Deployments(namespace).List(options)

	fmt.Println(deployment)
}