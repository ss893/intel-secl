/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/utils"
	"os"
	"os/user"
	"strconv"

	"github.com/intel-secl/intel-secl/v4/pkg/kbs"
)

func openLogFiles() (logFile *os.File, httpLogFile *os.File, secLogFile *os.File, err error) {

	logFile, err = os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, nil, nil, err
	}
	err = os.Chmod(LogFile, 0640)
	if err != nil {
		return nil, nil, nil, err
	}

	httpLogFile, err = os.OpenFile(HttpLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, nil, nil, err
	}
	err = os.Chmod(HttpLogFile, 0640)
	if err != nil {
		return nil, nil, nil, err
	}

	secLogFile, err = os.OpenFile(SecurityLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0640)
	if err != nil {
		return nil, nil, nil, err
	}
	err = os.Chmod(SecurityLogFile, 0640)
	if err != nil {
		return nil, nil, nil, err
	}

	// Containers are always run as non root users, does not require changing ownership of log directories
	if utils.IsContainerEnv() {
		return logFile, httpLogFile, secLogFile, nil
	}

	kbsUser, err := user.Lookup(ServiceUserName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not find user '%s'", ServiceUserName)
	}

	uid, err := strconv.Atoi(kbsUser.Uid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not parse kbs user id '%s'", kbsUser.Uid)
	}

	gid, err := strconv.Atoi(kbsUser.Gid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not parse kbs group id '%s'", kbsUser.Gid)
	}

	err = os.Chown(HttpLogFile, uid, gid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not change file ownership for file: '%s'", HttpLogFile)
	}
	err = os.Chown(SecurityLogFile, uid, gid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not change file ownership for file: '%s'", SecurityLogFile)
	}
	err = os.Chown(LogFile, uid, gid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not change file ownership for file: '%s'", LogFile)
	}
	return
}

func main() {
	logFile, httpLogFile, secLogFile, err := openLogFiles()
	var app *kbs.App
	if err != nil {
		app = &kbs.App{
			LogWriter: os.Stdout,
		}
	} else {
		defer func() {
			closeLogFiles(logFile, httpLogFile, secLogFile)
		}()
		app = &kbs.App{
			LogWriter:     logFile,
			HTTPLogWriter: httpLogFile,
			SecLogWriter:  secLogFile,
		}
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("Application returned with error : ", err.Error())
		closeLogFiles(logFile, httpLogFile, secLogFile)
		os.Exit(1)
	}
}

func closeLogFiles(logFile, httpLogFile, secLogFile *os.File) {
	var err error
	err = logFile.Close()
	if err != nil {
		fmt.Println("Failed to close default log file:", err.Error())
	}
	err = httpLogFile.Close()
	if err != nil {
		fmt.Println("Failed to close http log file:", err.Error())
	}
	err = secLogFile.Close()
	if err != nil {
		fmt.Println("Failed to close security log file:", err.Error())
	}
}
