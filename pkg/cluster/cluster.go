package cluster

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"os"
)

func getKubernetesClientOutsideCluster() *kubernetes.Clientset {
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

	log.Info("Successfully constructed k8s client outside cluster")

	return cl
}

func getMetricsClientOutsideCLuster() *metrics.Clientset {
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

	log.Info("Successfully constructed k8s metrics client outside cluster")

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
