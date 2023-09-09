package sdk

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type Quit struct {
	QuitChan chan bool
}

func NewQuit() *Quit {
	return &Quit{QuitChan: make(chan bool)}
}

func (q *Quit) WatchOsSignal() {
	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGUSR2) //nolint:govet
		log.Info("wait a QuitChan signal")
		s := <-signals
		log.Info("receive a QuitChan signal", "signal", s)
		q.QuitChan <- true
	}()
}

func (q *Quit) IsQuit() bool {
	select {
	case <-q.QuitChan:
		log.Info("QuitChan now!")
		return true
	}
}
