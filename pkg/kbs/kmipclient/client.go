/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kmipclient

import (
	kmip "github.com/gemalto/kmip-go"
	"github.com/gemalto/kmip-go/kmip14"
	"github.com/gemalto/kmip-go/ttlv"
)

type KmipClient interface {
	InitializeClient(string, string, string, string, string, string, string, string, string) error
	CreateSymmetricKey(int) (string, error)
	CreateAsymmetricKeyPair(string, string, int) (string, error)
	DeleteKey(string) error
	GetKey(string, string) ([]byte, error)
	SendRequest(interface{}, kmip14.Operation) (*kmip.ResponseBatchItem, *ttlv.Decoder, error)
}
