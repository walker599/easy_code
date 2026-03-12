package ui

import (
	"fmt"
	"time"
)

type Spinner struct {
	stop    chan bool
	stopped chan bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
}

func (s *Spinner) Start(message string) {
	go func() {
		chars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Print("\r\033[K") // Clear line
				s.stopped <- true
				return
			default:
				fmt.Printf("\r%s %s", chars[i], message)
				i = (i + 1) % len(chars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.stop <- true
	<-s.stopped
}
