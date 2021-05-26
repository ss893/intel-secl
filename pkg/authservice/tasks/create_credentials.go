/*
 *  Copyright (C) 2021 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package tasks

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/common"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/config"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/setup"
	"github.com/nats-io/nkeys"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

type CreateCredentials struct {
	CreateCredentials bool
	NatsConfig        config.NatsConfig
	ConsoleWriter     io.Writer
}

const createCredentialsHelpPrompt = "Following environment variables are optional for create-credentials setup:"

var createCredentialsEnvHelp = map[string]string{
	"CREATE_CREDENTIALS": "Trigger to run create-credentials setup task when set to True. Default is False",
}

func (cc *CreateCredentials) Run() error {
	defaultLog.Trace("tasks/create_credentials:Run() Entering")
	defer defaultLog.Trace("tasks/create_credentials:Run() Leaving")

	operatorKeyPair, err := nkeys.CreateOperator()
	if err != nil {
		return errors.Wrap(err, "Error creating operator nkeys")
	}

	operatorSeed, err := operatorKeyPair.Seed()
	if err != nil {
		return errors.Wrap(err, "Error fetching operator seed")
	}

	operatorToken, err := common.CreateJWTToken(operatorKeyPair, operatorKeyPair, constants.Operator, "",
		cc.NatsConfig.OperatorName)
	if err != nil {
		return errors.Wrap(err, "Error creating operator JWT")
	}
	log.Debug("Operator token is : ", operatorToken)

	err = ioutil.WriteFile(constants.OperatorSeedFile, operatorSeed, 0600)
	if err != nil {
		return errors.Wrap(err, "Error writing operator seed and private key to file")
	}

	accountKeyPair, err := nkeys.CreateAccount()
	if err != nil {
		return errors.Wrap(err, "Error creating account nkeys")
	}

	accountPublicKey, err := accountKeyPair.PublicKey()
	if err != nil {
		return errors.Wrap(err, "Error fetching public key of account")
	}
	log.Debug("Account Public Key : ", accountPublicKey)

	accountSeed, err := accountKeyPair.Seed()
	if err != nil {
		return errors.Wrap(err, "Error fetching seed of account")
	}

	accountToken, err := common.CreateJWTToken(accountKeyPair, operatorKeyPair, constants.Account, "",
		cc.NatsConfig.AccountName)
	if err != nil {
		return errors.Wrap(err, "Error creating account JWT")
	}
	log.Debug("Account token is : ", accountToken)

	err = ioutil.WriteFile(constants.AccountSeedFile, accountSeed, 0600)
	if err != nil {
		return errors.Wrap(err, "Error writing private key and seed of account to file")
	}

	//Create fixed format for server configuration
	formattedConf := fmt.Sprintf("// Operator %s\noperator: %s\n\nresolver: MEMORY\n\nresolver_preload: {\n"+
		" // Account %s\n %s: %s\n}", cc.NatsConfig.OperatorName, operatorToken, cc.NatsConfig.AccountName,
		accountPublicKey, accountToken)

	err = ioutil.WriteFile(constants.AccountConfigurationFile, []byte(formattedConf), 0600)
	if err != nil {
		return errors.Wrap(err, "Error writing server configuration to file")
	}

	fmt.Fprintln(cc.ConsoleWriter, "\nPlease copy the NATS server configuration printed on console:")
	fmt.Fprintln(cc.ConsoleWriter, "\n"+formattedConf+"\n")

	return nil
}

func (cc *CreateCredentials) Validate() error {
	defaultLog.Trace("tasks/create_credentials:Validate() Entering")
	defer defaultLog.Trace("tasks/create_credentials:Validate() Leaving")

	if !cc.CreateCredentials {
		return nil
	}
	_, err := os.Stat(constants.OperatorSeedFile)
	if err != nil {
		return err
	}
	_, err = os.Stat(constants.AccountSeedFile)
	if err != nil {
		return err
	}
	_, err = os.Stat(constants.AccountConfigurationFile)
	if err != nil {
		return err
	}
	return nil
}

func (cc *CreateCredentials) PrintHelp(w io.Writer) {
	setup.PrintEnvHelp(w, createCredentialsHelpPrompt, "", createCredentialsEnvHelp)
	fmt.Fprintln(w, "")
}

func (cc *CreateCredentials) SetName(n, e string) {

}
