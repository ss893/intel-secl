/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package rules

import (
	"testing"

	constants "github.com/intel-secl/intel-secl/v4/pkg/hvs/constants/verifier-rules-and-faults"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/host-connector/types"
	"github.com/stretchr/testify/assert"
)

func TestAssetTagMatchesNotProvisionedFault(t *testing.T) {

	hostManifest := types.HostManifest{
		AssetTagDigest: validAssetTagString, // valid tag in host
	}

	// provide a nil certificate value to the rule
	rule, err := NewAssetTagMatches(nil, assetTags)
	assert.NoError(t, err)

	// no faults should be returned...
	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultAssetTagNotProvisioned, result.Faults[0].Name)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestAssetTagMissingFromManifest(t *testing.T) {

	hostManifest := types.HostManifest{
		AssetTagDigest: "", // not in host manifest
	}

	// simulate adding valid asset tag bytes from the flavor...
	rule, err := NewAssetTagMatches(validAssetTagBytes, assetTags)
	assert.NoError(t, err)

	// we should get a "missing asset tag" fault...
	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultAssetTagMissing, result.Faults[0].Name)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestAssetTagMismatch(t *testing.T) {

	hostManifest := types.HostManifest{
		AssetTagDigest: invalidAssetTagString, // in valid from the host
	}

	// simulate adding valid asset tag bytes from the flavor...
	rule, err := NewAssetTagMatches(validAssetTagBytes, assetTags)
	assert.NoError(t, err)

	// we should get a "asset tag mismatch" fault...
	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Faults))
	assert.Equal(t, constants.FaultAssetTagMismatch, result.Faults[0].Name)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestAssetTagMatchesNoFault(t *testing.T) {

	hostManifest := types.HostManifest{
		AssetTagDigest: validAssetTagString, // valid tag in host
	}

	// simulate adding valid asset tag bytes from the flavor...
	rule, err := NewAssetTagMatches(validAssetTagBytes, assetTags)
	assert.NoError(t, err)

	// no faults should be returned...
	result, err := rule.Apply(&hostManifest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Faults))
}
