/*
 *  Copyright (C) 2021 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package config

import (
	"fmt"

	"github.com/intel-secl/intel-secl/v4/pkg/hvs/config"
	"github.com/spf13/viper"
)

//LoadConfig fetches the configuration details from config.yml file
func LoadConfig(path string) (config *config.Configuration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Error in fetching Configuration details")
	}

	return config, nil
}
