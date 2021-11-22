package nanny

import (
	"fmt"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

func (n *Nanny) systemMessage(msg string) {
	cmd := exec.Command(
		"notify-send",
		"-t",
		"3000",
		"-i",
		"face-smile",
		"'The First Notification'",
		fmt.Sprintf("'%s'", msg),
	)
	err := cmd.Run()
	if err != nil {
		log.Error(err)
	}
}

func (n *Nanny) systemShutdown() {
	nsec := 60
	log.Infof("Shutting down in %d seconds.", nsec)
	time.Sleep(time.Duration(nsec) * time.Second)
	cmd := exec.Command(
		"shutdown",
		"now",
	)
	err := cmd.Run()
	if err != nil {
		log.Error(err)
	}
}
