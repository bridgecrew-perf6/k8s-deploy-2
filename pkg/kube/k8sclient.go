package kube

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func GetKubernetesClient(env string) (*kubernetes.Clientset, error) {
	config, err := getKubernetesConfig(env)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getKubernetesConfig(env string) (*rest.Config, error) {
	var configName string
	//config, err := rest.InClusterConfig()
	//if err == nil {
	//	return config, err
	//} else if err != rest.ErrNotInCluster {
	//	return nil, err
	//}
	if env == "pre" || env == "prod" {
		configName = "prod-k8s-config"
	} else if env == "dev" {
		configName = "dev-k8s-config"
	} else if env == "qa" {
		configName = "qa-k8s-config"
	} else {
		return nil, errors.New("env for k8s config is wrong")
	}
	
	configPath := filepath.Join("./config", ".kube", configName)
	if !isExist(configPath) {
		return nil, errors.New("k8s config is not exist")
	}
	return clientcmd.BuildConfigFromFlags("", configPath)
}

func isExist(path string) bool  {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}