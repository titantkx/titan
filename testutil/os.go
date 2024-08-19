package testutil

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stretchr/testify/require"
)

func MustCreateTemp(t TestingT, name string) *os.File {
	file, err := os.CreateTemp("", name)
	require.NoError(t, err)
	require.NotNil(t, file)
	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})
	return file
}

func HandleOSInterrupt(f func()) {
	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("Received interrupt signal, cleaning up...")
		f()
		os.Exit(1)
	}()
}
