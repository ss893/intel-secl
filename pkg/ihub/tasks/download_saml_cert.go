/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package tasks

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/intel-secl/intel-secl/v4/pkg/clients/vs"
	"github.com/intel-secl/intel-secl/v4/pkg/ihub/config"
	"github.com/pkg/errors"
)

// DownloadSamlCert task for downloading SAML Certificate
type DownloadSamlCert struct {
	AttestationConfig *config.AttestationConfig
	ConsoleWriter     io.Writer
	SamlCertPath      string
}

// Run Runs the setup Task
func (samlCert DownloadSamlCert) Run() error {

	attestationHVSURL := samlCert.AttestationConfig.HVSBaseURL
	if attestationHVSURL == "" {
		fmt.Fprintln(samlCert.ConsoleWriter, "Skipping Download SAML Cert Task for SGX Attestation Service")
		return nil
	}

	if !strings.HasSuffix(attestationHVSURL, "/") {
		attestationHVSURL = attestationHVSURL + "/"
	}

	baseURL, err := url.Parse(attestationHVSURL)
	if err != nil {
		return errors.Wrap(err, "tasks/download_saml_cert:Run() Error in parsing Host Verification Service URL")
	}

	vsClient := &vs.Client{
		BaseURL: baseURL,
	}

	caCerts, err := vsClient.GetCaCerts("saml")
	if err != nil {
		return errors.Wrap(err, "tasks/download_saml_cert:Run() Failed to get SAML ca-certificates from HVS")
	}

	// write the output to a file
	err = ioutil.WriteFile(samlCert.SamlCertPath, caCerts, 0640)
	if err != nil {
		return errors.Wrapf(err, "tasks/download_saml_cert:Run() Error while writing file:%s", samlCert.SamlCertPath)
	}
	err = os.Chmod(samlCert.SamlCertPath, 0640)
	if err != nil {
		return errors.Wrapf(err, "tasks/download_saml_cert:Run() Error while changing file permission for file :%s", samlCert.SamlCertPath)
	}

	return nil
}

// Validate validates the downloaded certificate
func (samlCert DownloadSamlCert) Validate() error {

	if samlCert.AttestationConfig.HVSBaseURL == "" {
		fmt.Fprintln(samlCert.ConsoleWriter, "Skipping Download SAML Cert Task for SGX Attestation Service")
		return nil
	}

	if _, err := os.Stat(samlCert.SamlCertPath); os.IsNotExist(err) {
		return errors.Wrap(err, "tasks/download_saml_cert:Validate() Saml certificate does not exist")
	}

	_, err := ioutil.ReadFile(samlCert.SamlCertPath)
	if err != nil {
		return errors.Wrap(err, "tasks/download_saml_cert:Validate() Error while reading Saml CA Certificate file")
	}

	return nil
}

func (samlCert DownloadSamlCert) PrintHelp(w io.Writer) {}

func (samlCert DownloadSamlCert) SetName(n, e string) {}
