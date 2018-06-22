package paw

import (
	"log"
)

func Cron() {
	log.Print("paw.Cron start")
	// Can be on startup if we trust people not to remove them
	delete_old_aws_stack()
	// Refreshing state of cloudformation stack + redis
	log.Print("paw.Cron end")
	//log.Print(config.GlobalConfig)
	//log.Fatal("Failing hard")
}

// Loop over all aws account + region to create the s3 bucket for cloudformation
func delete_old_aws_stack() {
	log.Print("delete_old_aws_stack")
}
