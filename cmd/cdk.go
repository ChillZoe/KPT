package main

import (
	"github.com/KPT/pkg/cli"
	_ "github.com/KPT/pkg/exploit" // register all exploits
	_ "github.com/KPT/pkg/task"    // register all task
)

func main() {
	cli.ParseKPTMain()
}
