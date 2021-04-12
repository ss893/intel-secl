/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/common/utils"

	"github.com/intel-secl/intel-secl/v3/pkg/ihub"
	"github.com/intel-secl/intel-secl/v3/pkg/ihub/constants"

	"os"
	"os/user"
	"strconv"
)

func openLogFiles(logDir string) (logFile *os.File, secLogFile *os.File, err error) {
	logFilePath := logDir + LogFile
	securityLogFilePath := logDir + SecurityLogFile
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open/create %s", LogFile)
	}
	err = os.Chmod(logFilePath, 0664)
	if err != nil {
		return nil, nil, fmt.Errorf("error in setting file permission for file : %s", LogFile)
	}

	secLogFile, err = os.OpenFile(securityLogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return nil, nil, fmt.Errorf("could not open/create %s", SecurityLogFile)
	}
	err = os.Chmod(securityLogFilePath, 0664)
	if err != nil {
		return nil, nil, fmt.Errorf("error in setting file permission for file : %s", SecurityLogFile)
	}

	// Containers are always run as non root users, does not require changing ownership of log directories
	if utils.IsContainerEnv() {
		return logFile, secLogFile, nil
	}

	ihubUser, err := user.Lookup(ServiceUserName)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find user '%s'", ServiceUserName)
	}

	uid, err := strconv.Atoi(ihubUser.Uid)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse ihub user id '%s'", ihubUser.Uid)
	}

	gid, err := strconv.Atoi(ihubUser.Gid)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse ihub group id '%s'", ihubUser.Gid)
	}

	err = os.Chown(securityLogFilePath, uid, gid)
	if err != nil {
		return nil, nil, fmt.Errorf("could not change file ownership for file: '%s'", SecurityLogFile)
	}
	err = os.Chown(logFilePath, uid, gid)
	if err != nil {
		return nil, nil, fmt.Errorf("could not change file ownership for file: '%s'", LogFile)
	}

	return
}

func main() {
	var app *ihub.App
	instanceName, configDir, logDir := getAppSubConfig(os.Args)
	logFile, secLogFile, err := openLogFiles(logDir)
	if err != nil {
		fmt.Println("Error in setting up Log files :", err.Error())
		app = &ihub.App{
			LogWriter:    os.Stdout,
			ConfigDir:    configDir,
			LogDir:       logDir,
			InstanceName: instanceName,
		}
	} else {
		defer func() {
			closeLogFiles(logFile, secLogFile)
		}()
		app = &ihub.App{
			LogWriter:    logFile,
			SecLogWriter: secLogFile,
			ConfigDir:    configDir,
			LogDir:       logDir,
			InstanceName: instanceName,
		}
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("Application returned with error : ", err.Error())
		closeLogFiles(logFile, secLogFile)
		os.Exit(1)
	}
}

func getAppSubConfig(args []string) (string, string, string) {
	instanceName := ""
	configDir := ""
	logDir := ""
	for i, flag := range args {
		if flag == "-i" || flag == "--instance" {
			if i+1 < len(args) {
				instanceName = args[i+1]
				configDir = constants.SysConfigDir + instanceName + "/"
				logDir = constants.SysLogDir + instanceName + "/"
				break
			}
		}
	}
	if instanceName == "" {
		instanceName = constants.ServiceName
	}
	if configDir == "" {
		configDir = constants.ConfigDir
	}
	if logDir == "" {
		logDir = constants.LogDir
	}
	return instanceName, configDir, logDir
}

func closeLogFiles(logFile, secLogFile *os.File) {
	var err error
	err = logFile.Close()
	if err != nil {
		fmt.Println("Failed to close default log file:", err.Error())
	}
	err = secLogFile.Close()
	if err != nil {
		fmt.Println("Failed to close security log file:", err.Error())
	}
}
