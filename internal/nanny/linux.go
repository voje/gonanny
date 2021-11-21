package nanny

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func (n *Nanny) systemMessage() {
	command := "notify-send -t 3000 -i face-smile 'The First Notification' 'Hello <b>World</b>'"
	cmd := exec.Command(command)
	err := cmd.Run()
	if err != nil {
		log.Infof("err: %s Command failed: %s", err, command)
	}
}

func (n *Nanny) systemShutdown() {
	command := "shutdown -h +1"
	cmd := exec.Command(command)
	err := cmd.Run()
	if err != nil {
		log.Infof("err: %s Command failed: %s", err, command)
	}
}
