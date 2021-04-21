/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package keymanager

import (
	"strings"

	"github.com/intel-secl/intel-secl/v4/pkg/kbs/config"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/domain/models"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/kmipclient"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/intel-secl/intel-secl/v4/pkg/model/kbs"
	"github.com/pkg/errors"
)

var defaultLog = log.GetDefaultLogger()

func NewKeyManager(cfg *config.Configuration) (KeyManager, error) {
	defaultLog.Trace("keymanager/key_manager:NewKeyManager() Entering")
	defer defaultLog.Trace("keymanager/key_manager:NewKeyManager() Leaving")

	if strings.ToLower(cfg.KeyManager) == constants.KmipKeyManager {
		kmipClient := kmipclient.NewKmipClient()
		err := kmipClient.InitializeClient(cfg.Kmip.Version, cfg.Kmip.ServerIP, cfg.Kmip.ServerPort, cfg.Kmip.Hostname, cfg.Kmip.Username, cfg.Kmip.Password, cfg.Kmip.ClientKeyFilePath, cfg.Kmip.ClientCertificateFilePath, cfg.Kmip.RootCertificateFilePath)
		if err != nil {
			defaultLog.WithError(err).Error("keymanager/key_manager:NewKeyManager() Failed to initialize client")
			return nil, errors.New("Failed to initialize KeyManager")
		}
		return NewKmipManager(kmipClient), nil
	} else {
		defaultLog.Errorf("keymanager/key_manager:NewKeyManager() No Key Manager supported for provider: %s", cfg.KeyManager)
		return nil, errors.Errorf("No Key Manager supported for provider: %s", cfg.KeyManager)
	}
}

type KeyManager interface {
	CreateKey(*kbs.KeyRequest) (*models.KeyAttributes, error)
	DeleteKey(*models.KeyAttributes) error
	RegisterKey(*kbs.KeyRequest) (*models.KeyAttributes, error)
	TransferKey(*models.KeyAttributes) ([]byte, error)
}
