package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/airbytehq/abctl/internal/build"
	"github.com/airbytehq/abctl/internal/cmd"
	"github.com/airbytehq/abctl/internal/status"
	"github.com/airbytehq/abctl/internal/update"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// check for update
	updateCtx, updateCancel := context.WithTimeout(ctx, 2*time.Second)
	defer updateCancel()

	updateChan := make(chan updateInfo)
	go func() {
		info := updateInfo{}
		info.version, info.err = update.Check(updateCtx, http.DefaultClient, build.Version)
		updateChan <- info
	}()

	// listen for shutdown signals
	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
		<-signalCh

		cancel()
	}()

	// ensure the pterm info width matches the other printers
	//pterm.Info.Prefix.Text = " INFO  "

	root := cmd.NewCmd()
	cmd.Execute(ctx, root)

	newRelease := <-updateChan
	if newRelease.err != nil {
		if errors.Is(newRelease.err, update.ErrDevVersion) {
			status.Debug("Release checking is disabled for dev builds")
		}
	} else if newRelease.version != "" {
		status.Empty()
		status.Info(fmt.Sprintf("A new release of abctl is available: %s -> %s\nUpdating to the latest version is highly recommended", build.Version, newRelease.version))
	}
}

type updateInfo struct {
	version string
	err     error
}
