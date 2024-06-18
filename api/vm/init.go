package vm

import (
	log "github.com/wonderivan/logger"
	"kylin-lab/models"
)

func Setup() {
	var data models.LabVirtualMachine
	reseult, err := data.CheckAllVirtualMachineStatusAndUpdate("0")
	if err != nil {
		log.Error(err)
	}
	if len(reseult) > 0 {
		for _, v := range reseult {
			_, err = PostDeleteInstances(v.UUID)
			if err != nil {
				log.Error(err)
			}
			_, err = PostDeleteRecycleInstances(v.UUID)
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		log.Info("There are no outdated VMS to clean up")
	}

}
