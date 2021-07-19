/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kbs

import (
	"crypto/x509"
	"net/url"

	commLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/intel-secl/intel-secl/v4/pkg/model/kbs"
)

var log = commLog.GetDefaultLogger()

type KBSClient interface {
	CreateKey(*kbs.KeyRequest) (*kbs.KeyResponse, error)
	TransferKey(string, string) (*kbs.KeyTransferAttributes, error)
	TransferKeyWithSaml(string, string) ([]byte, error)
}

func NewKBSClient(aasURL, kbsURL *url.URL, username, password string, certs []x509.Certificate) KBSClient {
	return &kbsClient{
		AasURL:   aasURL,
		BaseURL:  kbsURL,
		UserName: username,
		Password: password,
		CaCerts:  certs,
	}
}

type kbsClient struct {
	AasURL   *url.URL
	BaseURL  *url.URL
	UserName string
	Password string
	CaCerts  []x509.Certificate
}
