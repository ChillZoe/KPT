package cli

import (
	"fmt"

	"github.com/KPT/conf"
	"github.com/KPT/pkg/evaluate"
	"github.com/KPT/pkg/plugin"

	// "github.com/KPT/pkg/tool/dockerd_api"
	// "github.com/KPT/pkg/tool/kubectl"

	"os"

	// "github.com/KPT/pkg/tool/netcat"
	// "github.com/KPT/pkg/tool/network"
	// "github.com/KPT/pkg/tool/probe"
	// "github.com/KPT/pkg/tool/ps"
	// "github.com/KPT/pkg/tool/vi"
	"github.com/docopt/docopt-go"
)

func PassInnerArgs() {
	os.Args = os.Args[1:]
}

func ParseKPTMain() {

	if len(os.Args) == 1 {
		docopt.PrintHelpAndExit(nil, BannerContainer)
	}

	// // nc needs -v and -h , parse it outside
	// if os.Args[1] == "nc" {
	// 	// https://github.com/jiguangin/netcat
	// 	PassInnerArgs()
	// 	netcat.RunVendorNetcat()
	// 	return
	// }

	// docopt argparse start
	parseDocopt()
	// auto-escape TODO : Many batch
	// if Args["auto-escape"].(bool) {
	// 	plugin.RunSingleTask("auto-escape")
	// 	return
	// }
	// evaluate

	// support for cdk eva(Evangelion) and cdk evaluate
	// fmt.Print(Args)
	fok := Args["evaluate"]
	ok := Args["eva"]

	// docopt let fok = true, so we need to check it
	// fix #37 https://github.com/KPT/issues/37
	if ok.(bool) || fok.(bool) {

		// fmt.Printf("\n[Information Gathering - System Info]\n")
		// evaluate.BasicSysInfo()

		// fmt.Printf("\n[Information Gathering - Services]\n")
		// evaluate.SearchSensitiveEnv()
		// evaluate.SearchSensitiveService()

		// fmt.Printf("\n[Information Gathering - Commands and Capabilities]\n")
		// evaluate.SearchAvailableCommands()
		// evaluate.GetProcCapabilities()

		// fmt.Printf("\n[Information Gathering - Mounts]\n")
		// evaluate.MountEscape()

		// fmt.Printf("\n[Information Gathering - Net Namespace]\n")
		// evaluate.CheckNetNamespace()

		// fmt.Printf("\n[Information Gathering - Sysctl Variables]\n")
		// evaluate.CheckRouteLocalNetworkValue()

		fmt.Printf("\n[Discovery - K8s API Server]\n")
		evaluate.CheckK8sAnonymousLogin()

		fmt.Printf("\n[Discovery - K8s Service Account]\n")
		evaluate.CheckPrivilegedK8sServiceAccount(conf.K8sSATokenDefaultPath)

		// fmt.Printf("\n[Discovery - Cloud Provider Metadata API]\n")
		// evaluate.CheckCloudMetadataAPI()

		// if Args["--full"].(bool) {

		// 	fmt.Printf("\n[Information Gathering - Sensitive Files]\n")
		// 	evaluate.SearchLocalFilePath()

		// 	fmt.Printf("\n[Information Gathering - ASLR]\n")
		// 	evaluate.ASLR()

		// 	fmt.Printf("\n[Information Gathering - Cgroups]\n")
		// 	evaluate.DumpCgroup()

		// }
		return
	}
	// exploit
	if Args["run"].(bool) {
		if Args["--list"].(bool) {
			plugin.ListAllExploit()
			os.Exit(0)
		}
		name := Args["<exploit>"].(string)
		if plugin.Exploits[name] == nil {
			fmt.Printf("\nInvalid script name: %s , available scripts:\n", name)
			plugin.ListAllExploit()
			return
		}
		plugin.RunSingleExploit(name)
		return
	}
	// tools
	// if Args["<tool>"] != nil {
	// 	args := Args["<args>"].([]string)

	// 	switch Args["<tool>"] {
	// 	case "vi":
	// 		PassInnerArgs()
	// 		vi.RunVendorVi()
	// 	case "kcurl":
	// 		kubectl.KubectlToolApi(args)
	// 	case "ucurl":
	// 		dockerd_api.UcurlToolApi(args)
	// 	case "dcurl":
	// 		dockerd_api.DcurlToolApi(args)
	// 	case "ifconfig":
	// 		network.GetLocalAddresses()
	// 	case "ps":
	// 		ps.RunPs()
	// 	case "probe":
	// 		if len(args) != 4 {
	// 			log.Println("Invalid input args.")
	// 			log.Println("usage: cdk probe <ip> <port> <parallels> <timeout-ms>")
	// 			log.Fatal("example: cdk probe 192.168.1.0-255 22,80,100-110 50 1000")
	// 		}
	// 		parallel, err := strconv.ParseInt(args[2], 10, 64)
	// 		if err != nil {
	// 			log.Println("err found when parse input arg <parallel>")
	// 			log.Fatal(err)
	// 		}
	// 		timeout, err := strconv.Atoi(args[3])
	// 		if err != nil {
	// 			log.Println("err found when parse input arg <timeout-ms>")
	// 			log.Fatal(err)
	// 		}
	// 		probe.TCPScanToolAPI(args[0], args[1], parallel, timeout)
	// 	default:
	// 		docopt.PrintHelpAndExit(nil, BannerContainer)
	// 	}
	// }
}
