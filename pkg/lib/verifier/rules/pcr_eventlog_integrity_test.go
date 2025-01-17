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
	"github.com/intel-secl/intel-secl/v4/pkg/lib/host-connector/util"
	"github.com/stretchr/testify/assert"
)

func TestPcrEventLogIntegrityNoFault(t *testing.T) {
	expectedCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	expectedPcrLog := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: expectedCumulativeHash,
	}

	expectedPcrLog1 := types.HostManifestPcrs{
		Index:   0,
		PcrBank: "SHA256",
		Value:   expectedCumulativeHash,
	}

	hostManifest := types.HostManifest{}
	hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs = append(hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs, testExpectedPcrEventLogEntry)
	hostManifest.PcrManifest.Sha256Pcrs = append(hostManifest.PcrManifest.Sha256Pcrs, expectedPcrLog1)

	rule, err := NewPcrEventLogIntegrity(&expectedPcrLog, common.FlavorPartPlatform)
	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Faults))
	t.Logf("Integrity rule verified")
}

func TestPcrEventLogIntegrityFault(t *testing.T) {

	_, err := NewPcrEventLogIntegrity(nil, common.FlavorPartPlatform)
	assert.Error(t, err)
}

// Provide the empty pcr manifest values in the host manifest and when applying PcrEventLogEquals rule, expecting
// 'PcrManifestMissing' fault.
func TestIntegrityPcrManifestMissingFault(t *testing.T) {
	expectedCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	expectedPcrLog := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: expectedCumulativeHash,
	}

	hostManifest := types.HostManifest{
		PcrManifest: types.PcrManifest{},
	}

	rule, err := NewPcrEventLogIntegrity(&expectedPcrLog, common.FlavorPartPlatform)
	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultPcrManifestMissing, result.Faults[0].Name)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

// Provide unsupported SHA algorithm
func TestPcrEventLogIntegrityUnsupportedSHAFault(t *testing.T) {
	expectedCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	expectedPcrLog := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA512",
		},
		Measurement: expectedCumulativeHash,
	}

	expectedPcrLog1 := types.HostManifestPcrs{
		Index:   0,
		PcrBank: "SHA256",
		Value:   expectedCumulativeHash,
	}

	hostManifest := types.HostManifest{}
	hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs = append(hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs, testExpectedPcrEventLogEntry)
	hostManifest.PcrManifest.Sha256Pcrs = append(hostManifest.PcrManifest.Sha256Pcrs, expectedPcrLog1)

	rule, err := NewPcrEventLogIntegrity(&expectedPcrLog, common.FlavorPartPlatform)
	result, err := rule.Apply(&hostManifest)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPcrEventLogIntegrityPcrValueMissingFault(t *testing.T) {
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

	expectedCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	expectedPcrLog := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: expectedCumulativeHash,
	}

	hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs = append(hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs, testExpectedPcrEventLogEntry)

	// if the pcr is no incuded, the PcrEventLogIntegrity rule should return
	// a PcrMissingFault
	// hostManifest.PcrManifest.Sha256Pcrs = ...not set
	rule, err := NewPcrEventLogIntegrity(&expectedPcrLog, common.FlavorPartPlatform)

	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultPcrValueMissing, result.Faults[0].Name)
	assert.NotNil(t, result.Faults[0].PcrIndex) // should report the missing pcr
	assert.Equal(t, types.PCR0, *result.Faults[0].PcrIndex)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestPcrEventLogIntegrityPcrEventLogMissingFault(t *testing.T) {
	expectedCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	expectedPcrLog1 := types.HostManifestPcrs{
		Index:   types.PCR0,
		PcrBank: types.SHA256,
		Value:   expectedCumulativeHash,
	}
	expectedPcrLog := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: expectedCumulativeHash,
	}

	hostManifest := types.HostManifest{}
	hostManifest.PcrManifest.Sha256Pcrs = append(hostManifest.PcrManifest.Sha256Pcrs, expectedPcrLog1)
	// omit the event log from the host manifest to invoke "PcrEventLogMissing" fault...
	//hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs = append(hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs, eventLogEntry)
	rule, err := NewPcrEventLogIntegrity(&expectedPcrLog, common.FlavorPartPlatform)

	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultPcrEventLogMissing, result.Faults[0].Name)
	assert.NotNil(t, result.Faults[0].PcrIndex) // should report the missing pcr
	assert.Equal(t, types.PCR0, *result.Faults[0].PcrIndex)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestPcrEventLogIntegrityPcrEventLogInvalidFault(t *testing.T) {
	expectedCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	expectedPcrLog := types.FlavorPcrs{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		Measurement: expectedCumulativeHash,
	}

	invalidPcrEventLogEntry := types.TpmEventLog{
		Pcr: types.Pcr{
			Index: 0,
			Bank:  "SHA256",
		},
		TpmEvent: []types.EventLog{
			{
				TypeName:    util.EVENT_LOG_DIGEST_SHA256,
				Measurement: zeros,
			},
		},
	}

	invalidCumulativeHash, err := testExpectedPcrEventLogEntry.Replay()
	assert.NoError(t, err)

	invalidPcrLog := types.HostManifestPcrs{
		Index:   types.PCR0,
		PcrBank: types.SHA256,
		Value:   invalidCumulativeHash,
	}

	hostManifest := types.HostManifest{}
	hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs = append(hostManifest.PcrManifest.PcrEventLogMap.Sha256EventLogs, invalidPcrEventLogEntry)
	hostManifest.PcrManifest.Sha256Pcrs = append(hostManifest.PcrManifest.Sha256Pcrs, invalidPcrLog)

	rule, err := NewPcrEventLogIntegrity(&expectedPcrLog, common.FlavorPartPlatform)

	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultPcrEventLogInvalid, result.Faults[0].Name)
	assert.NotNil(t, result.Faults[0].PcrIndex) // should report the missing pcr
	assert.Equal(t, types.PCR0, *result.Faults[0].PcrIndex)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}
