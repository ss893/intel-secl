/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v3/pkg/wpm"
	"os"
)

func openLogFiles() (logFile *os.File, secLogFile *os.File, err error) {

	logFile, err = os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	err = os.Chmod(LogFile, 0644)
	if err != nil {
		return nil, nil, err
	}

	secLogFile, err = os.OpenFile(SecurityLogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	err = os.Chmod(SecurityLogFile, 0644)
	if err != nil {
		return nil, nil, err
	}

	return
}

func main() {
	logFile, secLogFile, err := openLogFiles()
	var app *wpm.App
	if err != nil {
		app = &wpm.App{
			LogWriter: os.Stdout,
		}
	} else {
		defer func() {
			err = logFile.Close()
			if err != nil {
				fmt.Println("Failed close log file:", err.Error())
			}
		}()
		defer func() {
			err = secLogFile.Close()
			if err != nil {
				fmt.Println("Failed close log file:", err.Error())
			}
		}()
		app = &wpm.App{
			LogWriter:    logFile,
			SecLogWriter: secLogFile,
		}
	}

	err = app.Run(os.Args)
	if err != nil {
		fmt.Println("Application returned with error : ", err.Error())
		os.Exit(1)
	}
}
