/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package tasks

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

type CreateEnvelopeKey struct {
	EnvelopePrivatekeyLocation string
	EnvelopePublickeyLocation  string
	KeyAlgorithmLength         int
}

// ValidateCreateKey method is used to check if the envelope keys exists on disk
func (ek CreateEnvelopeKey) Validate() error {
	log.Trace("tasks/create_envelope_key.go:Validate() Entering")
	defer log.Trace("pkg/setup/create_envelope_key.go:Validate() Leaving")

	log.Info("tasks/create_envelope_key.go:Validate() Validating envelope key creation")

	_, err := os.Stat(ek.EnvelopePrivatekeyLocation)
	if os.IsNotExist(err) {
		return errors.Wrap(err, "tasks/create_envelope_key.go:Validate() Private key does not exist")
	}

	_, err = os.Stat(ek.EnvelopePublickeyLocation)
	if os.IsNotExist(err) {
		return errors.Wrap(err, "tasks/create_envelope_key.go:Validate() Public key does not exist")
	}
	return nil
}

func (ek CreateEnvelopeKey) Run() error {
	log.Trace("tasks/create_envelope_key.go:Run() Entering")
	defer log.Trace("tasks/create_envelope_key.go:Run() Leaving")

	if ek.Validate() != nil {
		log.Info("tasks/create_envelope_key.go:Run() Creating envelope key")

		bitSize := ek.KeyAlgorithmLength
		keyPair, err := rsa.GenerateKey(rand.Reader, bitSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while generating new RSA key pair")
			return errors.Wrap(err, "tasks/create_envelope_key.go:Run() Error while generating a new RSA key pair")
		}

		// save private key
		privateKey := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(keyPair),
		}

		privateKeyFile, err := os.OpenFile(ek.EnvelopePrivatekeyLocation, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I/O error while saving private key file")
			return errors.Wrap(err, "tasks/create_envelope_key.go:Run() I/O error while saving private key file")
		}
		defer func() {
			derr := privateKeyFile.Close()
			if derr != nil {
				fmt.Fprintf(os.Stderr, "Error while closing file"+derr.Error())
			}
		}()
		err = pem.Encode(privateKeyFile, privateKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I/O error while encoding private key file")
			return errors.Wrap(err, "tasks/create_envelope_key.go:Run() Error while encoding the private key.")
		}

		// save public key
		publicKey := &keyPair.PublicKey

		pubkeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I/O error while encoding private key file")
			return errors.Wrap(err, "tasks/create_envelope_key.go:Run() Error while marshalling the public key.")
		}
		var publicKeyInPem = &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubkeyBytes,
		}

		publicKeyFile, err := os.OpenFile(ek.EnvelopePublickeyLocation, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I/O error while encoding public envelope key file")
			return errors.Wrap(err, "tasks/create_envelope_key.go:Run() Error while creating a new file. ")
		}
		defer func() {
			derr := publicKeyFile.Close()
			if derr != nil {
				fmt.Fprintf(os.Stderr, "Error while closing file"+derr.Error())
			}
		}()

		err = pem.Encode(publicKeyFile, publicKeyInPem)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while encoding the public envelope key")
			return errors.Wrap(err, "tasks/create_envelope_key.go:Run() Error while encoding the public key.")
		}
	}
	return nil
}

func (ek CreateEnvelopeKey) PrintHelp(w io.Writer) {}

func (ek CreateEnvelopeKey) SetName(n, e string) {}
