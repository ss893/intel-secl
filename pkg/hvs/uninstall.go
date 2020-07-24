/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hvs

import (
	"fmt"
	"os"

	"github.com/intel-secl/intel-secl/v3/pkg/hvs/constants"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/postgres"
	"github.com/intel-secl/intel-secl/v3/pkg/hvs/tasks"
	e "github.com/intel-secl/intel-secl/v3/pkg/lib/common/exec"
	"github.com/pkg/errors"
)

func executablePath() string {
	exc, err := os.Executable()
	if err != nil {
		// if we can't find self-executable path, we're probably in a state that is panic() worthy
		panic(err)
	}
	return exc
}

func (a *App) executablePath() string {
	if a.ExecutablePath != "" {
		return a.ExecutablePath
	}
	exc, err := os.Executable()
	if err != nil {
		// if we can't find self-executable path, we're probably in a state that is panic() worthy
		panic(err)
	}
	return exc
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
		return a.ConfigDir
	}
	return constants.LogDir
}

func (a *App) execLinkPath() string {
	if a.ExecLinkPath != "" {
		return a.ExecLinkPath
	}
	return constants.ExecLinkPath
}

func (a *App) runDirPath() string {
	if a.RunDirPath != "" {
		return a.RunDirPath
	}
	return constants.RunDirPath
}

func (a *App) uninstall(purge bool) error {
	fmt.Println("Uninstalling HVS Service")
	// remove service
	_, _, err := e.RunCommandWithTimeout(constants.ServiceRemoveCmd, 5)
	if err != nil {
		fmt.Println("Could not disable HVS Service")
		fmt.Println("Error : ", err)
	}

	fmt.Println("removing : ", a.executablePath())
	err = os.Remove(a.executablePath())
	if err != nil {
		defaultLog.WithError(err).Error("error removing executable")
	}
	fmt.Println("removing : ", a.runDirPath())
	err = os.Remove(a.runDirPath())
	if err != nil {
		defaultLog.WithError(err).Error("error removing ", a.runDirPath())
	}
	fmt.Println("removing : ", a.execLinkPath())
	err = os.Remove(a.execLinkPath())
	if err != nil {
		defaultLog.WithError(err).Error("error removing ", a.execLinkPath())
	}
	// if purge is set
	if purge {
		fmt.Println("removing : ", a.configDir())
		err = os.RemoveAll(a.configDir())
		if err != nil {
			defaultLog.WithError(err).Error("error removing config dir")
		}
	}
	fmt.Println("removing : ", a.logDir())
	err = os.RemoveAll(a.logDir())
	if err != nil {
		defaultLog.WithError(err).Error("error removing log dir")
	}
	fmt.Println("removing : ", a.homeDir())
	err = os.RemoveAll(a.homeDir())
	if err != nil {
		defaultLog.WithError(err).Error("error removing home dir")
	}
	fmt.Fprintln(a.consoleWriter(), "HVS Service uninstalled")
	a.stop()
	return nil
}

var tablesToDrop = []string{
	"esxi_cluster",
	"esxi_cluster_host",
	"flavor_group",
	"host",
	"host_credential",
	"host_status",
	"report",
	"tag_certificate",
	"tpm_endorsement",
	// link_flavor_host,
	"flavorgroup_flavor",
	"host_flavorgroup",
	"queue",
	"flavor",
	"trust_cache",
	// tag_certificate_request,
	// host_tpm_password,
	"audit_log_entry",
}

func (a *App) eraseData() error {
	if a.configuration() == nil {
		return errors.New("Failed to load configuration file")
	}
	dbConf := a.configuration().DB
	// test connection and create schemas
	dataStore, err := postgres.NewDataStore(postgres.NewDatabaseConfig(constants.DBTypePostgres, &dbConf))
	if err != nil {
		return errors.Wrap(err, "Failed to connect database")
	}
	for _, t := range tablesToDrop {
		sqlCmd := "DROP TABLE IF EXISTS " + t + " CASCADE;"
		dataStore.ExecuteSql(&sqlCmd)
	}
	dataStore.Migrate()
	// create default flavor group
	t := tasks.CreateDefaultFlavor{
		DBConfig: dbConf,
	}
	if err := t.Run(); err != nil {
		return errors.Wrap(err, "Failed to run setup task CreateDefaultFlavor")
	}
	if err := t.Validate(); err != nil {
		return errors.Wrap(err, "Failed to validate setup task CreateDefaultFlavor")
	}
	return nil
}
