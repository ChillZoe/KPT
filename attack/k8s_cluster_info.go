package attack

import (
	"fmt"
	"log"

	"github.com/KPT/gadget/kubectl"
)

var configmapApi = "/api/v1/configmaps"

func GetNamespaces(serverAddr string) string {
	log.Println("requesting ", configmapApi)
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    serverAddr, // default
			Api:       "/api/v1/namespaces",
			Method:    "get",
			PostData:  "",
			Anonymous: true,
		})
	if err != nil {
		fmt.Println(err)
	}

	return resp
}

func GetNodes(serverAddr string) string {
	log.Println("requesting ", configmapApi)
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: "",
			Server:    serverAddr, // default
			Api:       "/api/v1/namespaces",
			Method:    "get",
			PostData:  "",
			Anonymous: true,
		})
	if err != nil {
		fmt.Println(err)
	}

	return resp
}
