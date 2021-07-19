/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package wpm

import (
	"flag"
	"fmt"
	commLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/log/message"
	commLogInt "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log/setup"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/setup"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/validation"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/config"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/imageflavor"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"strings"
)

var errInvalidCmd = errors.New("Invalid input after command")
var log = commLog.GetDefaultLogger()
var secLog = commLog.GetSecurityLogger()

type App struct {
	HomeDir        string
	ConfigDir      string
	LogDir         string
	ExecutablePath string
	ExecLinkPath   string
	RunDirPath     string
	Config         *config.Configuration
	ConsoleWriter  io.Writer
	ErrorWriter    io.Writer
	LogWriter      io.Writer
	SecLogWriter   io.Writer
}

func (a *App) consoleWriter() io.Writer {
	if a.ConsoleWriter != nil {
		return a.ConsoleWriter
	}
	return os.Stdout
}
func (a *App) errorWriter() io.Writer {
	if a.ErrorWriter != nil {
		return a.ErrorWriter
	}
	return os.Stderr
}

func (a *App) secLogWriter() io.Writer {
	if a.SecLogWriter != nil {
		return a.SecLogWriter
	}
	return os.Stdout
}

func (a *App) logWriter() io.Writer {
	if a.LogWriter != nil {
		return a.LogWriter
	}
	return os.Stderr
}

func (a *App) configuration() *config.Configuration {
	if a.Config != nil {
		return a.Config
	}
	viper.AddConfigPath(a.configDir())
	c, err := config.LoadConfiguration()
	if err == nil {
		a.Config = c
		return a.Config
	}
	return nil
}

func (a *App) configureLogs(isStdOut, isFileOut bool) error {
	var ioWriterDefault io.Writer
	ioWriterDefault = a.LogWriter
	if isStdOut {
		if isFileOut {
			ioWriterDefault = io.MultiWriter(os.Stdout, a.logWriter())
		} else {
			ioWriterDefault = os.Stdout
		}
	}

	ioWriterSecurity := io.MultiWriter(ioWriterDefault, a.secLogWriter())
	logConfig := a.Config.Log
	lv, err := logrus.ParseLevel(logConfig.Level)
	if err != nil {
		return errors.Wrap(err, "Failed to initiate loggers. Invalid log level: "+logConfig.Level)
	}
	commLogInt.SetLogger(commLog.DefaultLoggerName, lv, &commLog.LogFormatter{MaxLength: logConfig.MaxLength}, ioWriterDefault, false)
	commLogInt.SetLogger(commLog.SecurityLoggerName, lv, &commLog.LogFormatter{MaxLength: logConfig.MaxLength}, ioWriterSecurity, false)

	secLog.Info(message.LogInit)
	log.Info(message.LogInit)
	return nil
}

func (a *App) fetchKey(args []string) error {
	keyID := flag.String("k", "", "existing key ID")
	flag.StringVar(keyID, "key", "", "existing key ID")
	assetTag := flag.String("t", "", "asset tags associated with the new key")
	flag.StringVar(assetTag, "asset-tag", "", "asset tags associated with the new key")
	flag.Usage = func() { a.printFetchKeyUsage() }
	err := flag.CommandLine.Parse(args[2:])
	if err != nil {
		a.printFetchKeyUsage()
		return errors.Wrap(err, "Error parsing arguments")
	}

	//If the key ID is specified, make sure it's a valid UUID
	if len(strings.TrimSpace(*keyID)) > 0 {
		if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
			log.WithError(validatekeyIDErr).Errorf("app:fetchKey() %s : Error fetching key: Invalid UUID - %s\n", message.InvalidInputBadParam, *keyID)
			a.printFetchKeyUsage()
			return errors.Wrap(validatekeyIDErr, "Error fetching key: Invalid UUID")
		}
	}

	keyInfo, err := util.FetchKeyForAssetTag(*keyID, *assetTag)
	if err != nil {
		log.WithError(err).Errorf("app:fetchKey() %s - Error fetching: %s\n", message.AppRuntimeErr, err.Error())
		a.printFetchKeyUsage()
		return errors.Wrap(err, "Error fetching key")
	}
	if len(keyInfo) > 0 {
		fmt.Println(string(keyInfo))
	}
	return nil
}

func (a *App) createImageFlavor(args []string) error {
	flavorLabel := flag.String("l", "", "flavor label")
	flag.StringVar(flavorLabel, "label", "", "flavor label")
	inputImageFilename := flag.String("i", "", "input image file name")
	flag.StringVar(inputImageFilename, "in", "", "input image file name")
	outputFlavorFilename := flag.String("o", "", "output flavor file name")
	flag.StringVar(outputFlavorFilename, "out", "", "output flavor file name")
	outputEncImageFilename := flag.String("e", "", "output encrypted image file name")
	flag.StringVar(outputEncImageFilename, "encout", "", "output encrypted image file name")
	keyID := flag.String("k", "", "existing key ID")
	flag.StringVar(keyID, "key", "", "existing key ID")
	flag.Usage = func() { a.printImageFlavorUsage() }
	err := flag.CommandLine.Parse(args[2:])
	if err != nil {
		a.printImageFlavorUsage()
		return errors.Wrap(err, "Error parsing arguments")
	}

	if len(strings.TrimSpace(*flavorLabel)) <= 0 || len(strings.TrimSpace(*inputImageFilename)) <= 0 {
		log.Errorf("app:createImageFlavor() %s : Error creating VM image flavor: Missing arguments Flavor label and image file path\n", message.InvalidInputBadParam)
		a.printImageFlavorUsage()
		return errors.New("Error creating VM image flavor: Missing arguments Flavor label and image file path")
	}

	// validate input strings
	inputArr := []string{*flavorLabel, *outputFlavorFilename, *inputImageFilename, *outputEncImageFilename}
	if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
		log.WithError(validationErr).Errorf("app:createImageFlavor() %s : Error creating VM image flavor. Parse error for input args: [ %s ] - %s\n", message.InvalidInputBadParam, inputArr, validationErr.Error())
		a.printImageFlavorUsage()
		return errors.Wrap(validationErr, "Error creating VM image flavor: Invalid input arguments format")
	}

	//If the key ID is specified, make sure it's a valid UUID
	if len(strings.TrimSpace(*keyID)) > 0 {
		if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
			log.WithError(validatekeyIDErr).Errorf("app:createImageFlavor() %s : Error creating VM image flavor: Invalid UUID - %s\n", message.InvalidInputBadParam, *keyID)
			a.printImageFlavorUsage()
			return errors.Wrap(validatekeyIDErr, "Error creating VM image flavor: Invalid key UUID")
		}
	}

	imageFlavor, err := imageflavor.CreateImageFlavor(*flavorLabel, *outputFlavorFilename, *inputImageFilename,
		*outputEncImageFilename, *keyID, false)
	if err != nil {
		log.WithError(err).Errorf("app:createImageFlavor() %s - Error creating VM image flavor: %s\n", message.AppRuntimeErr, err.Error())
		a.printImageFlavorUsage()
		return errors.Wrap(err, "Error creating VM image flavor")
	}
	if len(imageFlavor) > 0 {
		fmt.Println(imageFlavor)
	}
	return nil
}

func (a *App) Run(args []string) error {

	if len(args) < 2 {
		a.printUsage()
		return nil
	}
	cmd := args[1]
	switch cmd {
	default:
		err := errors.New("Invalid command: " + cmd)
		a.printUsageWithError(err)
		return err
	case "help", "-h", "--help":
		a.printUsage()
		return nil
	case "uninstall":
		// the only allowed flag is --purge
		purge := false
		if len(args) == 3 {
			if args[2] != "--purge" {
				return errors.New("Invalid flag: " + args[2])
			}
			purge = true
		} else if len(args) != 2 {
			return errInvalidCmd
		}
		return a.uninstall(purge)
	case "version", "--version", "-v":
		a.printVersion()
		return nil
	case "setup":
		if err := a.setup(args[1:]); err != nil {
			if errors.Cause(err) == setup.ErrTaskNotFound {
				a.printUsageWithError(err)
			} else {
				fmt.Fprintln(a.errorWriter(), err.Error())
			}
			return err
		}
	case "fetch-key":
		configuration := a.configuration()
		if err := a.configureLogs(configuration.Log.EnableStdout, true); err != nil {
			return err
		}
		return a.fetchKey(os.Args[:])
	case "create-image-flavor":
		configuration := a.configuration()
		if err := a.configureLogs(configuration.Log.EnableStdout, true); err != nil {
			return err
		}
		return a.createImageFlavor(os.Args[:])
	}
	return nil
}
