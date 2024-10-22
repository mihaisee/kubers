package client

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
)

func GetKubernetesClientOutsideCluster() *kubernetes.Clientset {
	cp := getK8sConfigPath()

	// create the config from the path
	cfg, err := clientcmd.BuildConfigFromFlags("", cp)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
		panic(err.Error())
	}

	// generate the client based off of the config
	cl, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
		panic(err.Error())
	}

	return cl
}

func GetMetricsClientOutsideCLuster() *metrics.Clientset {
	cp := getK8sConfigPath()

	// create the config from the path
	cfg, err := clientcmd.BuildConfigFromFlags("", cp)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
		panic(err.Error())
	}

	// generate the client based off of the config
	mc, err := metrics.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("getMetricsClientOutsideCLuster: %v", err)
		panic(err.Error())
	}

	return mc
}

func getHomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE") // windows
}

func getK8sConfigPath() string {
	// construct the path to resolve to `~/.kube/config`
	return getHomeDir() + "/.kube/config"
}
