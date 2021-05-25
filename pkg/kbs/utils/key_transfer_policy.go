/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/model/kbs"
	"github.com/pkg/errors"
)

func GetDefaultKeyTransferPolicyId() (uuid.UUID, error) {
	defaultLog.Trace("utils/key_transfer_policy:GetDefaultKeyTransferPolicyId() Entering")
	defer defaultLog.Trace("utils/key_transfer_policy:GetDefaultKeyTransferPolicyId() Leaving")

	var id uuid.UUID
	bytes, err := ioutil.ReadFile(constants.DefaultTransferPolicyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return id, errors.New("Default key transfer policy file was not created. Please run setup task again.")
		} else {
			return id, errors.Wrapf(err, "Unable to read default key transfer policy file : %s", constants.DefaultTransferPolicyFile)
		}
	}

	var policy kbs.KeyTransferPolicyAttributes
	err = json.Unmarshal(bytes, &policy)
	if err != nil {
		return id, errors.Wrap(err, "Failed to unmarshal default key transfer policy")
	}

	return policy.ID, nil
}
