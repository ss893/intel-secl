/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package verifier

//
// Implements 'Verifier' interface.
//

import (
	"github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/verifier/rules"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	"github.com/pkg/errors"
)

type verifierImpl struct {
	verifierCertificates VerifierCertificates
}

func (v *verifierImpl) Verify(hostManifest *types.HostManifest, signedFlavor *hvs.SignedFlavor, skipSignedFlavorVerification bool) (*hvs.TrustReport, error) {

	var err error

	if hostManifest == nil {
		return nil, errors.New("The host manifest cannot be nil")
	}

	if signedFlavor == nil {
		return nil, errors.New("The signed flavor cannot be nil")
	}

	ruleFactory := NewRuleFactory(v.verifierCertificates, hostManifest, signedFlavor, skipSignedFlavorVerification)
	verificationRules, policyName, err := ruleFactory.GetVerificationRules()
	if err != nil {
		return nil, err
	}

	results, overallTrust, err := v.applyRules(verificationRules, hostManifest, signedFlavor)
	if err != nil {
		return nil, err
	}

	trustReport := hvs.TrustReport{
		PolicyName:   policyName,
		Results:      results,
		Trusted:      overallTrust,
		HostManifest: *hostManifest,
	}

	return &trustReport, nil
}

func (v *verifierImpl) applyRules(rulesToApply []rules.Rule, hostManifest *types.HostManifest, signedFlavor *hvs.SignedFlavor) ([]hvs.RuleResult, bool, error) {

	var results []hvs.RuleResult

	// default overall trust to true, change to false during rule evaluation
	overallTrust := true

	for _, rule := range rulesToApply {

		log.Debugf("Applying verifier rule %T", rule)
		result, err := rule.Apply(hostManifest)
		if err != nil {
			return nil, overallTrust, errors.Wrapf(err, "Error ocrurred applying rule type '%T'", rule)
		}

		// if 'Apply' returned a result with any faults, then the
		// rule is not trusted
		if len(result.Faults) > 0 {
			result.Trusted = false
			overallTrust = false
		}

		// assign the flavor id to all rules
		fId := signedFlavor.Flavor.Meta.ID
		result.FlavorId = &fId

		results = append(results, *result)
	}

	return results, overallTrust, nil
}

func (v *verifierImpl) GetVerifierCerts() VerifierCertificates {
	return v.verifierCertificates
}
