/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package wpm

import (
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/google/uuid"
	commLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/log/message"
	commLogInt "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log/setup"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/setup"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/validation"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/config"
	consts "github.com/intel-secl/intel-secl/v4/pkg/wpm/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/containerimageflavor"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/imageflavor"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/url"
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

func (a *App) getContainerImageId(args []string) error {
	NameSpaceDNS := uuid.Must(uuid.Parse(consts.SampleUUID))
	imageUUID := uuid.NewHash(md5.New(), NameSpaceDNS, []byte(args[1]), 4)
	log.Infof("app:getContainerImageId() Successfully retrieved container image ID: %s\n", imageUUID)
	fmt.Println(imageUUID)
	return nil
}

func (a *App) unwrapKey(args []string) error {
	wrappedKeyFilePath := flag.String("i", "", "wrapped key file path")
	flag.StringVar(wrappedKeyFilePath, "in", "", "wrapped key file path")
	flag.Usage = func() { a.printUnwrapKeyUsage() }
	err := flag.CommandLine.Parse(args[2:])
	if err != nil {
		a.printUnwrapKeyUsage()
		return errors.Wrap(err, "Error parsing arguments")
	}

	// validate input strings
	inputArr := []string{*wrappedKeyFilePath}
	if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
		log.WithError(validationErr).Errorf("app:unwrapKey() %s : Error unwrapping key: %s\n", message.AppRuntimeErr, validationErr.Error())
		a.printUnwrapKeyUsage()
		return errors.Wrap(validationErr, "Error unwrapping key")
	}

	wrappedKey, err := ioutil.ReadFile(*wrappedKeyFilePath)
	if err != nil {
		log.WithError(err).Errorf("app:unwrapKey() %s : Error unwrapping key: Unable to read from wrapped key file %s: %s\n", message.AppRuntimeErr, *wrappedKeyFilePath, err.Error())
		a.printUnwrapKeyUsage()
		return errors.Wrap(err, "Unable to read from wrapped key file")
	}

	unwrappedKey, err := util.UnwrapKey(wrappedKey, consts.EnvelopePrivatekeyLocation)
	if err != nil {
		log.WithError(err).Errorf("app:unwrapKey() %s : Error unwrapping key: %s\n", message.AppRuntimeErr, err.Error())
		a.printUnwrapKeyUsage()
		return errors.Wrap(err, "Error unwrapping key")
	}
	log.Info("app:unwrapKey() Successfully unwrapped key")
	fmt.Println(base64.StdEncoding.EncodeToString(unwrappedKey))
	return nil
}

func (a *App) createContainerImageFlavor(args []string) error {
	imageName := flag.String("i", "", "docker image name")
	flag.StringVar(imageName, "img-name", "", "docker image name")
	tagName := flag.String("t", "latest", "docker image tag")
	flag.StringVar(tagName, "tag", "latest", "docker image tag")
	dockerFilePath := flag.String("f", "", "Dockerfile path")
	flag.StringVar(dockerFilePath, "docker-file", "", "Dockerfile path")
	buildDir := flag.String("d", "", "build directory path containing source to build the docker image")
	flag.StringVar(buildDir, "build-dir", "", "build directory path containing source to build the docker image")
	keyID := flag.String("k", "", "key ID of key used for encrypting the image")
	flag.StringVar(keyID, "key-id", "", "key ID of key used for encrypting the image")
	encryptionRequired := flag.Bool("e", false, "specifies if image needs to be encrypted")
	flag.BoolVar(encryptionRequired, "encryption-required", false, "specifies if image needs to be encrypted")
	integrityEnforced := flag.Bool("s", false, "specifies if container image should be signed")
	flag.BoolVar(integrityEnforced, "integrity-enforced", false, "specifies if container image needs to be signed")
	notaryServerURL := flag.String("n", "", "notary server url to pull signed images")
	flag.StringVar(notaryServerURL, "notary-server", "", "notary server url to pull signed images")
	outputFlavorFilename := flag.String("o", "", "output flavor file name")
	flag.StringVar(outputFlavorFilename, "out-file", "", "output flavor file name")
	flag.Usage = func() { a.printContainerFlavorUsage() }
	err := flag.CommandLine.Parse(args[2:])
	if err != nil {
		a.printContainerFlavorUsage()
		return errors.Wrap(err, "Error parsing arguments")
	}
	if len(strings.TrimSpace(*imageName)) <= 0 {
		a.printContainerFlavorUsage()
		return errors.New("Flavor label and image file path are required arguments")
	}

	// validate input strings
	inputArr := []string{*imageName, *tagName, *dockerFilePath, *buildDir, *outputFlavorFilename}
	if validationErr := validation.ValidateStrings(inputArr); validationErr != nil {
		log.WithError(validationErr).Errorf("app:createContainerImageFlavor() %s : Error Creating Container Flavor: Validation error for input args: %s\n", message.InvalidInputBadParam, inputArr)
		a.printContainerFlavorUsage()
		return errors.Wrap(validationErr, "Error Creating Container Flavor: Validation error for input args")
	}

	//If the key ID is specified, make sure it's a valid UUID
	if len(strings.TrimSpace(*keyID)) > 0 {
		if validatekeyIDErr := validation.ValidateUUIDv4(*keyID); validatekeyIDErr != nil {
			log.WithError(validatekeyIDErr).Errorf("app:createContainerImageFlavor() %s : Error Creating Container Flavor: %s\n", message.InvalidInputBadParam, validatekeyIDErr.Error())
			a.printContainerFlavorUsage()
			return errors.Wrap(validatekeyIDErr, "Error Creating Container Flavor")
		}
	}

	if *notaryServerURL != "" {
		notaryServerURIValue, _ := url.Parse(*notaryServerURL)
		protocol := make(map[string]byte)
		protocol["https"] = 0
		if validateURLErr := validation.ValidateURL(*notaryServerURL, protocol, notaryServerURIValue.RequestURI()); validateURLErr != nil {
			log.WithError(validateURLErr).Errorf("app:createContainerImageFlavor() %s : Error Creating Container Flavor: Invalid key URL format %s\n", message.InvalidInputBadParam, validateURLErr.Error())
			a.printContainerFlavorUsage()
			return errors.Wrap(validateURLErr, "Error Creating Container Flavor: Invalid key URL format")
		}
	}

	containerImageFlavor, err := containerimageflavor.CreateContainerImageFlavor(*imageName, *tagName, *dockerFilePath, *buildDir,
		*keyID, *encryptionRequired, *integrityEnforced, *notaryServerURL, *outputFlavorFilename)
	if err != nil {
		log.WithError(err).Errorf("app:createContainerImageFlavor() %s : Error Creating Container Flavor: %s\n", message.AppRuntimeErr, err.Error())
		a.printContainerFlavorUsage()
		return errors.Wrap(err, "Error Creating Container Flavor")
	}

	if len(containerImageFlavor) > 0 {
		log.Info("app:createContainerImageFlavor() Successfully created container image flavor")
		fmt.Println(containerImageFlavor)
	}
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
	case "get-container-image-id":
		configuration := a.configuration()
		if err := a.configureLogs(configuration.Log.EnableStdout, true); err != nil {
			return err
		}

		if len(os.Args[1:]) < 1 {
			a.printGetContainerImageIdUsage()
			return errors.New("Invalid number of parameters")
		}
		return a.getContainerImageId(os.Args[1:])
	case "unwrap-key":
		// logs cannot be output to stdout as it will be mixed with the
		// unwrapped key output
		_ = a.configuration()
		if err := a.configureLogs(false, true); err != nil {
			return err
		}
		return a.unwrapKey(os.Args[:])
	case "create-container-image-flavor":
		configuration := a.configuration()
		if err := a.configureLogs(configuration.Log.EnableStdout, true); err != nil {
			return err
		}
		return a.createContainerImageFlavor(os.Args[:])
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
