/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package verifier

import (
	"github.com/intel-secl/intel-secl/v4/pkg/lib/flavor/model"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/verifier/rules"
)

type vendorTrustPolicyReader interface {
	Rules() []rules.Rule
}

type hostTrustPolicyManager struct {
}

func NewHostTrustPolicyManager(model.Flavor, *types.HostManifest) *hostTrustPolicyManager {
	return &hostTrustPolicyManager{}
}

func (htpm *hostTrustPolicyManager) GetVendorTrustPolicyReader() vendorTrustPolicyReader {
	return nil
}
