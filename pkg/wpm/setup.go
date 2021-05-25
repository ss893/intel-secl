/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package wpm

import (
	"crypto/x509/pkix"
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/setup"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/tasks"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
)

// input string slice should start with setup
func (a *App) setup(args []string) error {
	if len(args) < 2 {
		return errors.New("Invalid usage of setup")
	}
	// look for cli flags
	var ansFile string
	var force bool
	for i, s := range args {
		if s == "-f" || s == "--file" {
			if i+1 < len(args) {
				ansFile = args[i+1]
				break
			} else {
				return errors.New("Invalid answer file name")
			}
		}
		if s == "--force" {
			force = true
		}
	}
	// dump answer file to env
	if ansFile != "" {
		err := setup.ReadAnswerFileToEnv(ansFile)
		if err != nil {
			return errors.Wrap(err, "Failed to read answer file")
		}
	}
	runner, err := a.setupTaskRunner()
	if err != nil {
		return err
	}
	defer a.Config.Save(constants.DefaultConfigFilePath)
	cmd := args[1]
	// print help and return if applicable
	if len(args) > 2 && args[2] == "--help" {
		if cmd == "all" {
			err = runner.PrintAllHelp()
			if err != nil {
				return errors.Wrap(err, "Failed to write to console")
			}
		} else {
			err = runner.PrintHelp(cmd)
			if err != nil {
				return errors.Wrap(err, "Failed to write to console")
			}
		}
		return nil
	}
	if cmd == "all" {
		if err = runner.RunAll(force); err != nil {
			errCmds := runner.FailedCommands()
			fmt.Fprintln(a.errorWriter(), "Error(s) encountered when running all setup commands:")
			for errCmd, failErr := range errCmds {
				fmt.Fprintln(a.errorWriter(), errCmd+": "+failErr.Error())
				err = runner.PrintHelp(errCmd)
				if err != nil {
					return errors.Wrap(err, "Failed to write to console")
				}
			}
			return errors.New("Failed to run all tasks")
		}
		fmt.Fprintln(a.consoleWriter(), "All setup tasks succeeded")
	} else {
		if err = runner.Run(cmd, force); err != nil {
			fmt.Fprintln(a.errorWriter(), cmd+": "+err.Error())
			err = runner.PrintHelp(cmd)
			if err != nil {
				return errors.Wrap(err, "Failed to write to console")
			}
			return errors.New("Failed to run setup task " + cmd)
		}
	}
	return nil
}

// a helper function for setting up the task runner
func (a *App) setupTaskRunner() (*setup.Runner, error) {

	loadAlias()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if a.configuration() == nil {
		a.Config = defaultConfig()
	}
	runner := setup.NewRunner()
	runner.ConsoleWriter = a.consoleWriter()
	runner.ErrorWriter = a.errorWriter()

	runner.AddTask("download-ca-cert", "", &setup.DownloadCMSCert{
		CaCertDirPath: constants.TrustedCaCertsDir,
		ConsoleWriter: a.consoleWriter(),
		CmsBaseURL:    viper.GetString("cms-base-url"),
		TlsCertDigest: viper.GetString("cms-tls-cert-sha384"),
	})
	runner.AddTask("download-cert-flavor-signing", "flavor-signing", a.downloadCertTask("flavor-signing"))
	runner.AddTask("create-envelope-key", "", &tasks.CreateEnvelopeKey{
		EnvelopePrivatekeyLocation: constants.EnvelopePrivatekeyLocation,
		EnvelopePublickeyLocation:  constants.EnvelopePublickeyLocation,
		KeyAlgorithmLength:         constants.DefaultKeyAlgorithmLength,
	})
	return runner, nil
}

func (a *App) downloadCertTask(certType string) setup.Task {
	certTypeReq := certType
	var updateConfig = &a.configuration().FlavorSigning

	if updateConfig != nil {
		updateConfig.KeyFile = viper.GetString(certType + "-key-file")
		updateConfig.CertFile = viper.GetString(certType + "-cert-file")
		updateConfig.CommonName = viper.GetString(certType + "-common-name")
	}
	return &setup.DownloadCert{
		KeyFile:      viper.GetString(certType + "-key-file"),
		CertFile:     viper.GetString(certType + "-cert-file"),
		KeyAlgorithm: constants.DefaultKeyAlgorithm,
		KeyLength:    constants.DefaultKeyAlgorithmLength,
		Subject: pkix.Name{
			CommonName: viper.GetString(certType + "-common-name"),
		},
		CertType: certTypeReq,

		CaCertDirPath: constants.TrustedCaCertsDir,
		SanList:       constants.DefaultWpmSan,
		ConsoleWriter: a.consoleWriter(),
		CmsBaseURL:    viper.GetString("cms-base-url"),
		BearerToken:   viper.GetString("bearer-token"),
	}
}
