package kubectl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/KPT/conf"
	"github.com/KPT/pkg/errors"
	"github.com/KPT/pkg/util"
)

// MaybeSuccessfulStatuscodeList from https://www.w3.org/Protocols/HTTP/HTRESP.html
var MaybeSuccessfulStatuscodeList = []int{
	100, // RFC 7231, 6.2.1
	101, // RFC 7231, 6.2.2
	102, // RFC 2518, 10.1
	103, // RFC 8297

	200, // RFC 7231, 6.3.1
	201, // RFC 7231, 6.3.2
	202, // RFC 7231, 6.3.3
	203, // RFC 7231, 6.3.4
	204, // RFC 7231, 6.3.5
	205, // RFC 7231, 6.3.6
	206, // RFC 7233, 4.1
	207, // RFC 4918, 11.1
	208, // RFC 5842, 7.1
	226, // RFC 3229, 10.4.1
}

type K8sRequestOption struct {
	TokenPath string
	Server    string
	Api       string
	Method    string
	PostData  string
	Url       string
	Anonymous bool
}

func ApiServerAddr() (string, error) {
	protocol := ""
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		text := "err: cannot find kubernetes api host in ENV"
		return "", errors.New(text)
	}
	if port == "8080" || port == "8001" {
		protocol = "http://"
	} else {
		protocol = "https://"
	}
	return protocol + net.JoinHostPort(host, port), nil
}

func GetServiceAccountToken(tokenPath string) (string, error) {
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

/*
curl -s https://192.168.0.234:6443/api/v1/nodes?watch  --header "Authorization: Bearer xxx" --cacert ca.crt
*/
//https://github.com/kubernetes/client-go/blob/66db2540991da169fb60fce735064a55bfc52b71/rest/config.go#L483
func ServerAccountRequest(opts K8sRequestOption) (string, error) {

	// parse token
	var token string
	var tokenErr error
	if opts.Anonymous {
		token = ""
	} else if opts.TokenPath == "" {
		token, tokenErr = GetServiceAccountToken(conf.K8sSATokenDefaultPath)
	} else {
		token, tokenErr = GetServiceAccountToken(opts.TokenPath)
	}
	if tokenErr != nil {
		return "", &errors.CDKRuntimeError{Err: tokenErr, CustomMsg: "load K8s service account token error."}
	}

	// parse url if opts.Url is ""
	if len(opts.Url) == 0 {
		var server string
		var urlErr error
		if opts.Server == "" {
			server, urlErr = ApiServerAddr()
			opts.Url = server + opts.Api
		} else {
			opts.Url = opts.Server + opts.Api
			urlErr = nil
		}
		if urlErr != nil {
			return "", &errors.CDKRuntimeError{Err: urlErr, CustomMsg: "err found while searching local K8s apiserver addr."}
		}
	}

	// http client
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	var request *http.Request
	opts.Method = strings.ToUpper(opts.Method)

	request, err := http.NewRequest(opts.Method, opts.Url, bytes.NewBuffer([]byte(opts.PostData)))
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "err found while generate post request in net.http ."}
	}

	// set request header
	if opts.Method == "POST" {
		request.Header.Set("Content-Type", "application/json")
	}
	// auth token
	if len(token) > 0 {
		token = strings.TrimSpace(token)
		request.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "err found in post request."}
	}
	//defer resp.Body.Close()

	// Fix a bug reported by the author of crossc2 on whc2021.
	// When the DeployBackdoorDaemonset call fails and returns an error, it will still feedback true.
	if !util.IntContains(MaybeSuccessfulStatuscodeList, resp.StatusCode) {
		errMsg := fmt.Sprintf("err found in post request, error response code: %v.", resp.Status)
		return "", &errors.CDKRuntimeError{
			Err:       err,
			CustomMsg: errMsg,
		}
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", &errors.CDKRuntimeError{Err: err, CustomMsg: "err found in post request."}
	}

	return string(content), nil
}

func GetServerVersion(serverAddr string) (string, error) {
	opts := K8sRequestOption{
		TokenPath: "",
		Server:    serverAddr,
		Api:       "/version",
		Method:    "GET",
		PostData:  "",
		Anonymous: true,
	}
	resp, err := ServerAccountRequest(opts)
	if err != nil {
		return "", err
	}
	// use regexp to find gitVersion
	versionPattern := regexp.MustCompile(`"gitVersion":.*?"(.*?)"`)
	results := versionPattern.FindStringSubmatch(resp)
	if len(results) != 2 {
		return "", errors.New("field gitVersion not found in response")
	}
	return results[1], nil
}
