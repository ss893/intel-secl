/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package wlsclient

import (
	"bytes"
	"crypto/x509"
	"github.com/intel-secl/intel-secl/v4/pkg/clients/util"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
)

type ReportsClient interface {
	PostVMReport([]byte) error
}

type reportsClientImpl struct {
	caCerts []x509.Certificate
	cfg     *wlsClientConfig
}

//PostVMReport method is used to upload the VM trust report to workload service
func (client reportsClientImpl) PostVMReport(report []byte) error {
	log.Trace("wlsclient/reports_client:PostVMReport() Entering")
	defer log.Trace("wlsclient/reports_client:PostVMReport() Leaving")
	var err error

	//Add client here
	requestURL, err := url.Parse(client.cfg.BaseURL)
	if err != nil {
		return errors.New("wlsclient/reports_client:PostVMReport() error retrieving WLS API URL")
	}

	requestURL.Path = path.Join(requestURL.Path, "reports")

	log.Debugf("wlsclient/reports_client:PostVMReport() WLS VM reports POST Request URL: %s", requestURL.String())

	httpRequest, err := http.NewRequest("POST", requestURL.String(), bytes.NewBuffer(report))
	if err != nil {
		return errors.Wrap(err, "wlsclient/reports_client:PostVMReport() Failed to create WLS POST API request for vm reports")
	}
	// set POST request Accept and Content-Type headers
	httpRequest.Header.Set("Accept", "application/json")
	httpRequest.Header.Set("Content-Type", "application/json")

	_, err = util.SendRequest(httpRequest, client.cfg.AasApiURL, client.cfg.Username, client.cfg.Password, client.caCerts)
	if err != nil {
		return errors.Wrap(err, "wlsclient/reports_client:PostVMReport() Error while getting response for Post WLS VM reports API")
	}

	return nil
}
