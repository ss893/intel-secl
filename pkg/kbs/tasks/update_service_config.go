/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package tasks

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/config"
	commConfig "github.com/intel-secl/intel-secl/v4/pkg/lib/common/config"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/setup"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"strings"
)

type UpdateServiceConfig struct {
	ServerConfig  commConfig.ServerConfig
	ServiceConfig config.KBSConfig
	DefaultPort   int
	AASApiUrl     string
	AppConfig     **config.Configuration
	ConsoleWriter io.Writer
}

const envHelpPrompt = "Following environment variables are required for update-service-config setup:"

var allowedSKCChallengeTypes = map[string]bool{"sgx": true}
var allowedKeyManagers = map[string]bool{"kmip": true}

var envHelp = map[string]string{
	"SERVICE_USERNAME":           "The service username as configured in AAS",
	"SERVICE_PASSWORD":           "The service password as configured in AAS",
	"LOG_LEVEL":                  "Log level",
	"LOG_MAX_LENGTH":             "Max length of log statement",
	"LOG_ENABLE_STDOUT":          "Enable console log",
	"AAS_BASE_URL":               "AAS Base URL",
	"KMIP_SERVER_IP":             "IP of KMIP server",
	"KMIP_SERVER_PORT":           "PORT of KMIP server",
	"KMIP_HOSTNAME":              "HOSTNAME of KMIP server",
	"KMIP_USERNAME":              "USERNAME of KMIP server",
	"KMIP_PASSWORD":              "PASSWORD of KMIP server",
	"KMIP_CLIENT_CERT_PATH":      "KMIP Client certificate path",
	"KMIP_CLIENT_KEY_PATH":       "KMIP Client key path",
	"KMIP_ROOT_CERT_PATH":        "KMIP Root Certificate path",
	"SKC_CHALLENGE_TYPE":         "SKC challenge type",
	"SQVS_URL":                   "SQVS URL",
	"SESSION_EXPIRY_TIME":        "Session Expiry Time",
	"SERVER_PORT":                "The Port on which Server Listens to",
	"SERVER_READ_TIMEOUT":        "Request Read Timeout Duration in Seconds",
	"SERVER_READ_HEADER_TIMEOUT": "Request Read Header Timeout Duration in Seconds",
	"SERVER_WRITE_TIMEOUT":       "Request Write Timeout Duration in Seconds",
	"SERVER_IDLE_TIMEOUT":        "Request Idle Timeout in Seconds",
	"SERVER_MAX_HEADER_BYTES":    "Max Length Of Request Header in Bytes ",
}

func (uc UpdateServiceConfig) Run() error {
	log.Trace("tasks/update_config:Run() Entering")
	defer log.Trace("tasks/update_config:Run() Leaving")

	if uc.AASApiUrl == "" {
		return errors.New("KBS configuration not provided: AAS_BASE_URL is not set")
	}

	if uc.ServerConfig.Port < 1024 ||
		uc.ServerConfig.Port > 65535 {
		uc.ServerConfig.Port = uc.DefaultPort
	}
	(*uc.AppConfig).KBS = uc.ServiceConfig

	(*uc.AppConfig).Server = uc.ServerConfig
	(*uc.AppConfig).AASApiUrl = uc.AASApiUrl

	(*uc.AppConfig).Log = commConfig.LogConfig{
		MaxLength:    viper.GetInt("log-max-length"),
		EnableStdout: viper.GetBool("log-enable-stdout"),
		Level:        viper.GetString("log-level"),
	}
	(*uc.AppConfig).EndpointURL = viper.GetString("endpoint-url")
	(*uc.AppConfig).Kmip = config.KmipConfig{
		Version:                   viper.GetString("kmip-version"),
		ServerIP:                  viper.GetString("kmip-server-ip"),
		ServerPort:                viper.GetString("kmip-server-port"),
		Hostname:                  viper.GetString("kmip-hostname"),
		Username:                  viper.GetString("kmip-username"),
		Password:                  viper.GetString("kmip-password"),
		ClientKeyFilePath:         viper.GetString("kmip-client-key-path"),
		ClientCertificateFilePath: viper.GetString("kmip-client-cert-path"),
		RootCertificateFilePath:   viper.GetString("kmip-root-cert-path"),
	}
	(*uc.AppConfig).Skc = config.SKCConfig{
		StmLabel:          viper.GetString("skc-challenge-type"),
		SQVSUrl:           viper.GetString("sqvs-url"),
		SessionExpiryTime: viper.GetInt("session-expiry-time"),
	}
	(*uc.AppConfig).KeyManager = viper.GetString("key-manager")
	return nil
}

func (uc UpdateServiceConfig) Validate() error {
	if uc.AASApiUrl == "" {
		return errors.New("KBS configuration not provided: AAS_BASE_URL is not set")
	}
	if (*uc.AppConfig).Server.Port < 1024 ||
		(*uc.AppConfig).Server.Port > 65535 {
		return errors.New("Configured port is not valid")
	}
	if _, validInput := allowedKeyManagers[strings.ToLower((*uc.AppConfig).KeyManager)]; !validInput {
		return errors.New("Invalid value provided for KEY_MANAGER. Value should be kmip")
	}
	if (*uc.AppConfig).Skc.StmLabel != "" {
		if _, validInput := allowedSKCChallengeTypes[strings.ToLower((*uc.AppConfig).Skc.StmLabel)]; !validInput {
			return errors.New("Invalid value provided for SKC_CHALLENGE_TYPE. allowed value is SGX")
		}
	}
	return nil
}
func (uc UpdateServiceConfig) PrintHelp(w io.Writer) {
	setup.PrintEnvHelp(w, envHelpPrompt, "", envHelp)
	fmt.Fprintln(w, "")
}

func (uc UpdateServiceConfig) SetName(n, e string) {
}
