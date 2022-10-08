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
	"github.com/KPT/conf"
	"github.com/KPT/observe"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// observeCmd represents the auxiliary command
var observeCmd = &cobra.Command{
	Use:   "observe",
	Short: "Observe module of KPT",
	Long:  `Observe module is used to gather information within the container and in the cluster, which will be useful during exploit Kubernetes cluster.`,
}

var accessAPICmd = &cobra.Command{
	Use:   "checkServer",
	Short: "Try to access API Server anonymously or with Service Account.",
	Long:  `Try to access API Server anonymously or with Service Account(the token is default).`,
	// Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Check whether K8s API Server can be accessed anonymously")
		observe.CheckK8sAnonymousLogin()
		// TODO client-go enhangce
		log.Info("Check whether K8s API Server can be accessed by defult")
		observe.CheckPrivilegedK8sServiceAccount(conf.K8sSATokenDefaultPath)
	},
}

var checkCloudCmd = &cobra.Command{
	Use:   "checkCloud",
	Short: "Try to find whether Pod is on Cloud Platform",
	Long:  `Try to access API Server anonymously or with Service Account(the token is default).`,
	// Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Check whether K8s cluster is on Cloud Platform")
		observe.CheckCloudMetadataAPI()
	},
}

var checkPodCmd = &cobra.Command{
	Use:     "checkPod",
	Short:   "Try to access API Server anonymously or with Service Account.",
	Long:    `Try to access API Server anonymously or with Service Account(the token is default).`,
	Example: "privileged ps",
	// Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Check whether container is in K8s cluster")
		observe.IsInK8sPod()
	},
}

func init() {
	rootCmd.AddCommand(observeCmd)
	observeCmd.AddCommand(checkPodCmd, accessAPICmd, checkCloudCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// observe.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// observeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
