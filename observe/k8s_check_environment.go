package observe

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// DotDockerEnvExists checks if /.dockerenv file exists
func DotDockerEnvExists() bool {
	_, err := os.Stat("/.dockerenv")
	return !os.IsNotExist(err)
}

// CheckCgroup checks if `docker` can be found in /proc/1/cgroup
func CheckCgroup() bool {
	content, _ := ioutil.ReadFile("/proc/1/cgroup")
	//log.Debug("Content of /proc/1/cgroup:\n", string(content))
	return bytes.Contains(content, []byte("docker"))
}

// CheckInitProcName checks if executable name of init process(process with PID 1) is systemd/init
func CheckInitProcName() bool {
	var (
		comm []byte
		err  error
	)
	// https://man7.org/linux/man-pages/man5/proc.5.html
	// /proc/[pid]/comm exposes the process's comm value—that is, the command name associated with the process.
	if comm, err = ioutil.ReadFile("/proc/1/comm"); err != nil {
		return false
	}

	name := strings.TrimSpace(string(comm))

	return !(name == "init" || name == "systemd")
}

// CheckKubeEnv checks if environment variables include `KUBERNETES_SERVICE_HOST`
// referring to https://kubernetes.io/docs/concepts/services-networking/connect-applications-service/#environment-variables, When a Pod runs on a Node, the kubelet adds a set of environment variables for each active Service.
func CheckKubeEnv() bool {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	return len(host) > 0
}

func IsInK8sPod() {
	if CheckKubeEnv() {
		log.Info("KUBERNETES_SERVICE_HOST found in env")
		log.WithFields(log.Fields{"SUCCESS": true}).Info("Container is in K8s Pod")
	} else {
		log.Warn("KUBERNETES_SERVICE_HOST not found in env")
		log.WithFields(log.Fields{"SUCCESS": false}).Error("Container is not in K8s Pod!")
	}
}
