package observe

import (
	"strings"

	"github.com/KPT/gadget/kubectl"
	log "github.com/sirupsen/logrus"
)

func CheckK8sAnonymousLogin() bool {

	// check if api-server allows system:anonymous request
	log.Info("Checking if Kubernetes API Server allows system:anonymous request.")

	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    "", // default
			Api:       "/",
			Method:    "get",
			PostData:  "",
			Anonymous: true,
		})
	if err != nil {
		log.Error(err)
	}

	if strings.Contains(resp, "/api") {
		log.WithFields(log.Fields{"SUCCESS": true}).Info("API Server allows anonymous request !")
		log.Info("trying to list namespaces")

		// check if system:anonymous can list namespaces
		resp, err := kubectl.ServerAccountRequest(
			kubectl.K8sRequestOption{
				TokenPath: "",
				Server:    "", // default
				Api:       "/api/v1/namespaces",
				Method:    "get",
				PostData:  "",
				Anonymous: true,
			})
		if err != nil {
			log.Error(err)
		}
		if len(resp) > 0 && strings.Contains(resp, "kube-system") {
			log.WithFields(log.Fields{"SUCCESS": true}).Info("Anonymous role have a high authority !")
			// log.Info("\tsuccess, the system:anonymous role have a high authority.")
			// log.Info("\tnow you can make your own request to takeover the entire k8s cluster with `./kpt kcurl` command\n\tgood luck and have fun.")
			return true
		} else {
			log.WithFields(log.Fields{"SUCCESS": false}).Info("\tresponse:" + resp)
			return true
		}
	} else {
		log.WithFields(log.Fields{"SUCCESS": false}).Info("Kubernetes API Server forbids anonymous request.")
		// log.Info("\tresponse:" + resp)
		return false
	}
}
