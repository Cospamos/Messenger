package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

const (
	ngrokURL  = "https://zonia-interparliament-nonnormally.ngrok-free.dev"
	localPort = "8080"
)

var (
	ngrokCmd  *exec.Cmd
	ngrokOnce sync.Once
)

func StartNgrok() error {
	var startErr error

	ngrokOnce.Do(func() {
		ngrokCmd = exec.Command(
			"ngrok",
			"http",
			"--url="+ngrokURL[len("https://"):],
			localPort,
		)

		ngrokCmd.Stdout = nil
		ngrokCmd.Stderr = nil

		if err := ngrokCmd.Start(); err != nil {
			startErr = err
			return
		}
	})

	return startErr
}


func StopNgrok() error {
	if ngrokCmd != nil && ngrokCmd.Process != nil {
		return ngrokCmd.Process.Kill()
	}
	return nil
}

func GetPublicURL() string {
	return ngrokURL
}

func UseNgrok() error {
	if err := StartNgrok(); err != nil {
		return fmt.Errorf("failed to start ngrok: %w", err)
	}

	// graceful shutdown
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		_ = StopNgrok()
		os.Exit(0)
	}()

	return nil
}
