package observe

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/KPT/gadget/kubectl"
)

func CheckPrivilegedK8sServiceAccount(tokenPath string) bool {
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    "",
			Api:       "/apis",
			Method:    "get",
			PostData:  "",
			Anonymous: false,
		})
	if err != nil {
		log.Error(err)
		return false
	}
	if len(resp) > 0 && strings.Contains(resp, "APIGroupList") {
		log.WithFields(log.Fields{"SUCCESS": true}).Info("Default Service Account can access Kubernetes API Server")

		// check if the current service-account can list namespaces
		log.Info("Trying to get Pod.")
		resp, err := kubectl.ServerAccountRequest(
			kubectl.K8sRequestOption{
				TokenPath: "",
				Server:    "",
				Api:       "/api/v1/pods",
				Method:    "get",
				PostData:  "",
				Anonymous: false,
			})
		if err != nil {
			fmt.Println(err)
			return false
		}
		if len(resp) > 0 && strings.Contains(resp, "kube-system") {
			log.WithFields(log.Fields{"SUCCESS": true}).Info("Default Service Account has a high authority.")
			return true
		} else {
			log.Info("\tfailed")
			return false
		}
	} else {
		log.WithFields(log.Fields{"SUCCESS": false}).Info("Default Service Account is not available.")
		return false
	}
}
