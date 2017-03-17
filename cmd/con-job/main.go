package main

import (
	"github.com/allen13/con-job/pkg/scheduler"
	"github.com/docopt/docopt-go"
	"log"
)

func main() {
	usage := `Con Job.

		Usage:
		  con-job start (scheduler|api|node) [--config=<config>]
		  con-job -h | --help
		  con-job --version

		Options:
		  -h --help     Show this screen.
		  --version     Show version.
		  --config=<config>  Config file [default: /etc/con-job/con-job.toml].`

	args, err := docopt.Parse(usage, nil, true, "Con Job 2.0", false)
	if err != nil {
		log.Fatal(err)
	}
	if args["scheduler"].(bool) {
		scheduler, err := scheduler.Build()
		if err != nil {
			log.Fatal(err)
		}
		scheduler.Start()
	}
}
