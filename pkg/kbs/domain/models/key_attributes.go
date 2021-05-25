/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/intel-secl/intel-secl/v4/pkg/model/kbs"
)

// KeyAttributes - Contains all possible key attributes.
type KeyAttributes struct {
	ID               uuid.UUID `json:"id"`
	Algorithm        string    `json:"algorithm"`
	KeyLength        int       `json:"key_length,omitempty"`
	KeyData          string    `json:"key,omitempty"`
	CurveType        string    `json:"curve_type,omitempty"`
	PublicKey        string    `json:"public_key,omitempty"`
	PrivateKey       string    `json:"private_key,omitempty"`
	KmipKeyID        string    `json:"kmip_key_id,omitempty"`
	TransferPolicyId uuid.UUID `json:"transfer_policy_id,omitempty"`
	TransferLink     string    `json:"transfer_link,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	Label            string    `json:"label,omitempty"`
	Usage            string    `json:"usage,omitempty"`
}

func (ka *KeyAttributes) ToKeyResponse() *kbs.KeyResponse {

	keyInformation := kbs.KeyInformation{
		ID:        ka.ID,
		Algorithm: ka.Algorithm,
		KeyLength: ka.KeyLength,
		CurveType: ka.CurveType,
		KmipKeyID: ka.KmipKeyID,
	}

	keyResponse := kbs.KeyResponse{
		KeyInformation:   &keyInformation,
		TransferPolicyID: ka.TransferPolicyId,
		TransferLink:     ka.TransferLink,
		CreatedAt:        ka.CreatedAt,
		Label:            ka.Label,
		Usage:            ka.Usage,
	}

	return &keyResponse
}
