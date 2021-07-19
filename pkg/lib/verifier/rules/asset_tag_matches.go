/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package rules

//
// Rule that validates that the host manifests matches what was supplied
// in the flavor.
//

import (
	"bytes"
	"encoding/base64"

	constants "github.com/intel-secl/intel-secl/v4/pkg/hvs/constants/verifier-rules-and-faults"
	asset_tag "github.com/intel-secl/intel-secl/v4/pkg/lib/asset-tag"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/flavor/common"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v4/pkg/model/hvs"
	"github.com/pkg/errors"
)

func NewAssetTagMatches(expectedAssetTagDigest []byte, tags []asset_tag.TagKvAttribute) (Rule, error) {

	assetTagMatches := assetTagMatches{
		expectedAssetTagDigest: expectedAssetTagDigest,
		tags:                   tags,
	}

	return &assetTagMatches, nil
}

type assetTagMatches struct {
	expectedAssetTagDigest []byte
	tags                   []asset_tag.TagKvAttribute
}

func (rule *assetTagMatches) Apply(hostManifest *types.HostManifest) (*hvs.RuleResult, error) {
	var fault *hvs.Fault
	result := hvs.RuleResult{}
	result.Trusted = true
	result.Rule.Name = constants.RuleAssetTagMatches
	result.Rule.ExpectedTag = rule.expectedAssetTagDigest
	result.Rule.Markers = append(result.Rule.Markers, common.FlavorPartAssetTag)
	tags := map[string]string{}
	for _, kvAttr := range rule.tags {
		tags[kvAttr.Key] = kvAttr.Value
	}
	result.Rule.Tags = tags

	if len(hostManifest.AssetTagDigest) == 0 {
		fault = &hvs.Fault{
			Name:        constants.FaultAssetTagMissing,
			Description: "AssetTag Reported is null",
		}
	} else if rule.expectedAssetTagDigest == nil {
		fault = &hvs.Fault{
			Name:        constants.FaultAssetTagNotProvisioned,
			Description: "AssetTag is not in provisioned by the management",
		}
	} else {
		actualAssetTagDigest, err := base64.StdEncoding.DecodeString(hostManifest.AssetTagDigest)
		if err != nil {
			return nil, errors.Wrap(err, "Could not decode TagCertDigest")
		}

		if bytes.Compare(actualAssetTagDigest, rule.expectedAssetTagDigest) != 0 {
			fault = &hvs.Fault{
				Name:        constants.FaultAssetTagMismatch,
				Description: "Asset tag provisioned does not match asset tag reported",
			}
		}
	}

	if fault != nil {
		result.Faults = append(result.Faults, *fault)
	}

	return &result, nil
}
