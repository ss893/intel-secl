/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package rules

import (
	"testing"

	constants "github.com/intel-secl/intel-secl/v4/pkg/hvs/constants/verifier-rules-and-faults"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/flavor/common"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/host-connector/types"
	"github.com/stretchr/testify/assert"
)

func TestPcrMatchesConstantNoFault(t *testing.T) {
	expectedPcr := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: PCR_VALID_256,
	}

	hostManifest := types.HostManifest{
		PcrManifest: types.PcrManifest{
			Sha256Pcrs: []types.HostManifestPcrs{
				{
					Index:   0,
					Value:   PCR_VALID_256,
					PcrBank: types.SHA256,
				},
			},
		},
	}

	rule, err := NewPcrMatchesConstant(&expectedPcr, common.FlavorPartPlatform)
	assert.NoError(t, err)

	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 0)
	assert.True(t, result.Trusted)
	t.Logf("Pcr matches constant rule verified")
}

func TestPcrMatchesConstantNoMeasurementFault(t *testing.T) {
	expectedPcr := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
	}

	_, err := NewPcrMatchesConstant(&expectedPcr, common.FlavorPartPlatform)
	assert.Error(t, err)
}

func TestPcrMatchesConstantNoExpectedPcrFault(t *testing.T) {

	_, err := NewPcrMatchesConstant(nil, common.FlavorPartPlatform)
	assert.Error(t, err)
}

func TestPcrMatchesConstantPcrManifestMissingFault(t *testing.T) {
	expectedPcr := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: PCR_VALID_256,
	}

	rule, err := NewPcrMatchesConstant(&expectedPcr, common.FlavorPartPlatform)
	assert.NoError(t, err)

	// provide a manifest without a PcrManifest and expect FaultPcrManifestMissing
	result, err := rule.Apply(&types.HostManifest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, constants.FaultPcrManifestMissing)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestPcrMatchesConstantMismatchFault(t *testing.T) {
	expectedPcr := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: PCR_VALID_256,
	}

	// host manifest with 'invalid' value for pcr0
	hostManifest := types.HostManifest{
		PcrManifest: types.PcrManifest{
			Sha256Pcrs: []types.HostManifestPcrs{
				{
					Index:   0,
					Value:   PCR_INVALID_256,
					PcrBank: types.SHA256,
				},
			},
		},
	}

	rule, err := NewPcrMatchesConstant(&expectedPcr, common.FlavorPartPlatform)
	assert.NoError(t, err)

	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, constants.FaultPcrValueMismatchSHA256)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestPcrMatchesConstantMissingFault(t *testing.T) {
	// empty manifest will result in 'missing' fault
	hostManifest := types.HostManifest{
		PcrManifest: types.PcrManifest{
			Sha256Pcrs: []types.HostManifestPcrs{
				{
					Index:   1,
					Value:   PCR_VALID_256,
					PcrBank: types.SHA256,
				},
			},
		},
	}

	expectedPcr := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: PCR_VALID_256,
	}

	rule, err := NewPcrMatchesConstant(&expectedPcr, common.FlavorPartPlatform)
	assert.NoError(t, err)

	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, constants.FaultPcrValueMissing)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}
