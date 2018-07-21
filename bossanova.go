package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vtfr/bossanova/server"
)

func main() {
	if err := server.Start(); err != nil {
		logrus.WithField("error", err).Fatalln("Failed running server")
	}
}
