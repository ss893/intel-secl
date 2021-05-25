/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package wpm

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/exec"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/constants"
	"os"
)

func (a *App) executablePath() string {
	if a.ExecutablePath != "" {
		return a.ExecutablePath
	}
	exec, err := os.Executable()
	if err != nil {
		log.WithError(err).Error("app:executablePath() Unable to find WPM executable")
		// if we can't find self-executable path, we're probably in a state that is panic() worthy
		panic(err)
	}
	return exec
}

func (a *App) homeDir() string {
	if a.HomeDir != "" {
		return a.HomeDir
	}
	return constants.HomeDir
}

func (a *App) configDir() string {
	if a.ConfigDir != "" {
		return a.ConfigDir
	}
	return constants.ConfigDir
}

func (a *App) logDir() string {
	if a.LogDir != "" {
		return a.LogDir
	}
	return constants.LogDir
}

func (a *App) execLinkPath() string {
	if a.ExecLinkPath != "" {
		return a.ExecLinkPath
	}
	return constants.ExecLinkPath
}

func removeSecureDockerDaemon() {
	fmt.Println("Uninstalling secure-docker-daemon")

	commandArgs := []string{constants.HomeDir + "secure-docker-daemon/uninstall-secure-docker-daemon.sh"}

	_, err := exec.ExecuteCommand("/bin/sh", commandArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while removing secure-docker-daemon %s:", err.Error())
	}
}

func (a *App) uninstall(purge bool) error {
	log.Trace("app:uninstall() Entering")
	defer log.Trace("app:uninstall() Leaving")

	fmt.Println("Uninstalling Workload Policy Manager")

	_, err := os.Stat(constants.HomeDir + "secure-docker-daemon")
	if err == nil {
		removeSecureDockerDaemon()
		// restart docker daemon
		if err == nil {
			commandArgs := []string{"start", "docker"}
			_, err = exec.ExecuteCommand("systemctl", commandArgs)
			if err != nil {
				fmt.Print("Error starting docker daemon post-uninstall. Refer dockerd logs for more information.")
			}
		}
	}

	fmt.Println("removing : ", a.executablePath())
	err = os.Remove(a.executablePath())
	if err != nil {
		log.WithError(err).Error("error removing executable")
	}

	fmt.Println("removing : ", a.execLinkPath())
	err = os.Remove(a.execLinkPath())
	if err != nil {
		log.WithError(err).Error("error removing ", a.execLinkPath())
	}

	if purge {
		fmt.Println("removing : ", a.configDir())
		err = os.RemoveAll(a.configDir())
		if err != nil {
			log.WithError(err).Error("error removing config dir")
		}
	}
	fmt.Println("removing : ", a.logDir())
	err = os.RemoveAll(a.logDir())
	if err != nil {
		log.WithError(err).Error("error removing log dir")
	}
	fmt.Println("removing : ", a.homeDir())
	err = os.RemoveAll(a.homeDir())
	if err != nil {
		log.WithError(err).Error("error removing home dir")
	}

	fmt.Fprintln(a.consoleWriter(), "Workload Policy Manager uninstalled")
	return nil
}
