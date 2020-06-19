/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package verifier

import (
	"crypto/x509"
	"testing"
	"time"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/model"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/constants"
	"github.com/stretchr/testify/assert"
)


func TestTagCertificateTrustedNoFault(t *testing.T) {

	// create a CA certpool...
	caPemBytes, caPrivateKey, err := newCACertificate()
	assert.NoError(t, err)

	trustedAuthorityCerts := x509.NewCertPool()
	ok := trustedAuthorityCerts.AppendCertsFromPEM(caPemBytes)
	assert.True(t, ok)

	// create the attribute certificate...
	tagCertificate, err := newCertificateTemplate()
	assert.NoError(t, err)

	tagCertificateBytes, err := getCertificateBytes(tagCertificate, caPrivateKey)
	assert.NoError(t, err)

	attributeCertificate := model.X509AttributeCertificate {
		Encoded: tagCertificateBytes,
		NotBefore: time.Now().AddDate(-1, 0, 0).Format(constants.FlavorTimestampFormat),
		NotAfter: time.Now().AddDate(1, 0, 0).Format(constants.FlavorTimestampFormat),
	}

	// create the rule
	rule, err := newTagCertificateTrusted(trustedAuthorityCerts, &attributeCertificate)
	assert.NoError(t, err)

	// apply the rule, the hostManifest has no impact on TagCertificateTrusted rule
	result, err := rule.Apply(&types.HostManifest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 0)
}

func TestTagCertificateTrustedMissingFault(t *testing.T) {

	// create a CA certpool...
	caPemBytes, _, err := newCACertificate()
	assert.NoError(t, err)

	trustedAuthorityCerts := x509.NewCertPool()
	ok := trustedAuthorityCerts.AppendCertsFromPEM(caPemBytes)
	assert.True(t, ok)	

	// create the rule, not provding the attribute certificate to invoke
	// FaultTagCertificateMissing.
	rule, err := newTagCertificateTrusted(trustedAuthorityCerts, nil)
	assert.NoError(t, err)

	// apply the rule, the hostManifest has no impact on TagCertificateTrusted rule
	result, err := rule.Apply(&types.HostManifest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, FaultTagCertificateMissing)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestTagCertificateTrustedNotTrusted(t *testing.T) {

	// create an empty CA certpool to force a FaultTagCertificateNotTrusted fault
	_, caPrivateKey, err := newCACertificate()
	assert.NoError(t, err)

	trustedAuthorityCerts := x509.NewCertPool() // empty

	// create the attribute certificate...
	tagCertificate, err := newCertificateTemplate()
	assert.NoError(t, err)

	tagCertificateBytes, err := getCertificateBytes(tagCertificate, caPrivateKey)
	assert.NoError(t, err)

	attributeCertificate := model.X509AttributeCertificate {
		Encoded: tagCertificateBytes,
		NotBefore: time.Now().AddDate(-1, 0, 0).Format(constants.FlavorTimestampFormat),
		NotAfter: time.Now().AddDate(1, 0, 0).Format(constants.FlavorTimestampFormat),
	}

	// create the rule
	rule, err := newTagCertificateTrusted(trustedAuthorityCerts, &attributeCertificate)
	assert.NoError(t, err)

	// apply the rule, the hostManifest has no impact on TagCertificateTrusted rule
	result, err := rule.Apply(&types.HostManifest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, FaultTagCertificateNotTrusted)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestTagCertificateTrustedExpiredFault(t *testing.T) {

	// create a CA certpool...
	caPemBytes, caPrivateKey, err := newCACertificate()
	assert.NoError(t, err)

	trustedAuthorityCerts := x509.NewCertPool()
	ok := trustedAuthorityCerts.AppendCertsFromPEM(caPemBytes)
	assert.True(t, ok)

	// create the attribute certificate, providing a 'NotAfter' in the past
	// to invoke the FaultTagCertificateExpired fault.
	tagCertificate, err := newCertificateTemplate()
	assert.NoError(t, err)

	tagCertificateBytes, err := getCertificateBytes(tagCertificate, caPrivateKey)
	assert.NoError(t, err)

	attributeCertificate := model.X509AttributeCertificate {
		Encoded: tagCertificateBytes,
		NotBefore: time.Now().AddDate(-1, 0, 0).Format(constants.FlavorTimestampFormat),
		NotAfter: time.Now().AddDate(-11, 0, 0).Format(constants.FlavorTimestampFormat),
	}

	// create the rule
	rule, err := newTagCertificateTrusted(trustedAuthorityCerts, &attributeCertificate)
	assert.NoError(t, err)

	// apply the rule, the hostManifest has no impact on TagCertificateTrusted rule
	result, err := rule.Apply(&types.HostManifest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, FaultTagCertificateExpired)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}

func TestTagCertificateTrustedNotYetValidFault(t *testing.T) {

	// create a CA certpool...
	caPemBytes, caPrivateKey, err := newCACertificate()
	assert.NoError(t, err)

	trustedAuthorityCerts := x509.NewCertPool()
	ok := trustedAuthorityCerts.AppendCertsFromPEM(caPemBytes)
	assert.True(t, ok)

	// create the attribute certificate, providing a 'NotBefore' in the future
	// to invoke the FaultTagCertificateNotYetValid fault.
	tagCertificate, err := newCertificateTemplate()
	assert.NoError(t, err)

	tagCertificateBytes, err := getCertificateBytes(tagCertificate, caPrivateKey)
	assert.NoError(t, err)

	attributeCertificate := model.X509AttributeCertificate {
		Encoded: tagCertificateBytes,
		NotBefore: time.Now().AddDate(1, 0, 0).Format(constants.FlavorTimestampFormat),
		NotAfter: time.Now().AddDate(1, 0, 0).Format(constants.FlavorTimestampFormat),
	}

	// create the rule
	rule, err := newTagCertificateTrusted(trustedAuthorityCerts, &attributeCertificate)
	assert.NoError(t, err)

	// apply the rule, the hostManifest has no impact on TagCertificateTrusted rule
	result, err := rule.Apply(&types.HostManifest{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result.Faults), 1)
	assert.Equal(t, result.Faults[0].Name, FaultTagCertificateNotYetValid)
	t.Logf("Fault description: %s", result.Faults[0].Description)
}