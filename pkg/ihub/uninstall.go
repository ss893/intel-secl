/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package ihub

import (
	"fmt"
	"os"

	"github.com/intel-secl/intel-secl/v4/pkg/ihub/constants"
	commonExec "github.com/intel-secl/intel-secl/v4/pkg/lib/common/exec"
)

func (app *App) executablePath() string {
	if app.ExecutablePath != "" {
		return app.ExecutablePath
	}
	osExec, err := os.Executable()
	if err != nil {
		log.WithError(err).Error("Unable to find ihub executable")
		panic(err)
	}
	return osExec
}

func (app *App) homeDir() string {
	if app.HomeDir != "" {
		return app.HomeDir
	}
	return constants.HomeDir
}

func (app *App) configDir() string {
	if app.ConfigDir != "" {
		return app.ConfigDir
	}
	return constants.ConfigDir
}

func (app *App) logDir() string {
	if app.LogDir != "" {
		return app.LogDir
	}
	return constants.LogDir
}

func (app *App) execLinkPath() string {
	if app.ExecLinkPath != "" {
		return app.ExecLinkPath
	}
	return constants.ExecLinkPath
}

func (app *App) runDirPath() string {
	if app.RunDirPath != "" {
		return app.RunDirPath
	}
	return constants.RunDirPath
}

func (app *App) uninstall(purge, exec bool) {
	fmt.Println("Stopping Integration Hub - " + app.InstanceName)
	err := app.stop()
	if err != nil {
		log.WithError(err).Error("error stopping service")
	}

	fmt.Println("Uninstalling Integration Hub - " + app.InstanceName)
	removeService(app.InstanceName)

	if exec {
		fmt.Println("Removing : ", app.executablePath())
		err := os.Remove(app.executablePath())
		if err != nil {
			log.WithError(err).Error("error removing executable")
		}

		fmt.Println("Removing : ", app.execLinkPath())
		err = os.Remove(app.execLinkPath())
		if err != nil {
			log.WithError(err).Error("Error removing executable link")
		}
		fmt.Println("Removing : ", app.homeDir())
		err = os.RemoveAll(app.homeDir())
		if err != nil {
			log.WithError(err).Error("Error removing home dir")
		}
		// Remove\Disable multi instance service file
		removeService("")
	}
	if purge {
		fmt.Println("Removing : ", app.configDir())
		err = os.RemoveAll(app.configDir())
		if err != nil {
			log.WithError(err).Error("Error removing config dir")
		}
	}
	fmt.Println("Removing : ", app.logDir())
	err = os.RemoveAll(app.logDir())
	if err != nil {
		log.WithError(err).Error("Error removing log dir")
	}

	fmt.Fprintln(app.consoleWriter(), "Integration Hub uninstalled")
}

func removeService(instanceName string) {
	serviceName := constants.InstancePrefix + instanceName + ".service"
	_, _, err := commonExec.RunCommandWithTimeout(constants.ServiceRemoveCmd+serviceName, 5)
	if err != nil {
		fmt.Println("Could not remove Integration Hub Service")
		fmt.Println("Error : ", err)
	}
}
