package utils

import (
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

func NginxReload() {
	shellCmd := exec.Command("/usr/bin/systemctl", "reload", "nginx")
	err := shellCmd.Start()
	if err != nil {
		log.Panicln(err)
	}
	log.Debugln(shellCmd.CombinedOutput())
}
