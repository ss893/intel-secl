/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package wpm

import (
	"os"

	commConfig "github.com/intel-secl/intel-secl/v3/pkg/lib/common/config"
	"github.com/intel-secl/intel-secl/v3/pkg/wpm/config"
	"github.com/intel-secl/intel-secl/v3/pkg/wpm/constants"
	"github.com/spf13/viper"
)

// this func sets the default values for viper keys
func init() {
	// set default values for log
	viper.SetDefault("log-max-length", constants.DefaultLogMaxlength)
	viper.SetDefault("log-enable-stdout", false)
	viper.SetDefault("log-level", constants.DefaultLogLevel)
	viper.SetDefault("flavor-signing-cert-file", constants.FlavorSigningCertFile)
	viper.SetDefault("flavor-signing-key-file", constants.FlavorSigningKeyFile)
	viper.SetDefault("flavor-signing-common-name", constants.DefaultWpmFlavorSigningCn)

}

func defaultConfig() *config.Configuration {
	loadAlias()
	return &config.Configuration{
		AASApiUrl:        viper.GetString("aas-base-url"),
		CMSBaseURL:       viper.GetString("cms-base-url"),
		CmsTlsCertDigest: viper.GetString("cms-tls-cert-sha384"),
		KBSApiUrl:        viper.GetString("kbs-base-url"),
		WPM: commConfig.ServiceConfig{
			Username: viper.GetString("wpm-service-username"),
			Password: viper.GetString("wpm-service-password"),
		},
		Log: commConfig.LogConfig{
			MaxLength:    viper.GetInt("log-max-length"),
			Level:        viper.GetString("log-level"),
			EnableStdout: viper.GetBool("log-enable-stdout"),
		},
		FlavorSigning: commConfig.SigningCertConfig{
			CertFile:   viper.GetString("flavor-signing-cert-file"),
			KeyFile:    viper.GetString("flavor-signing-key-file"),
			CommonName: viper.GetString("flavor-signing-common-name"),
		},
	}
}

func loadAlias() {
	alias := map[string]string{
		"aas-base-url": "AAS_API_URL",
		"kbs-base-url": "KMS_API_URL",
	}
	for k, v := range alias {
		if env := os.Getenv(v); env != "" {
			viper.Set(k, env)
		}
	}
}
