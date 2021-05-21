/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package flavor

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/model"
	"github.com/pkg/errors"
	"io/ioutil"

	"github.com/google/uuid"
)

/**
 *
 * @author purvades
 */

// ImageFlavor is a flavor for an image with the encryption requirement information
// and key details of an encrypted image.
type ImageFlavor struct {
	Image model.Image `json:"flavor"`
}

// GetImageFlavor is used to create a new image flavor with the specified label, encryption policy,
// key url, and digest of the encrypted image
func GetImageFlavor(label string, encryptionRequired bool, keyURL string, digest string) (*ImageFlavor, error) {
	log.Trace("flavor/image_flavor:GetImageFlavor() Entering")
	defer log.Trace("flavor/image_flavor:GetImageFlavor() Leaving")
	var encryption *model.Encryption

	description := model.Description{
		Label:      label,
		FlavorPart: "IMAGE",
	}

	meta := model.Meta{
		Description: description,
	}
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new UUID")
	}
	meta.ID = newUuid

	if encryptionRequired {
		encryption = &model.Encryption{
			KeyURL: keyURL,
			Digest: digest,
		}
	}

	imageflavor := model.Image{
		Meta:               meta,
		EncryptionRequired: encryptionRequired,
		Encryption:         encryption,
	}

	flavor := ImageFlavor{
		Image: imageflavor,
	}
	return &flavor, nil
}

// GetContainerImageFlavor is used to create a new container image flavor with the specified label, encryption policy,
// Key url of the encrypted image also integrity policy and notary url for docker image signature verification
func GetContainerImageFlavor(label string, encryptionRequired bool, keyURL string, integrityEnforced bool, notaryURL string) (*ImageFlavor, error) {
	log.Trace("flavor/image_flavor:GetContainerImageFlavor() Entering")
	defer log.Trace("flavor/image_flavor:GetContainerImageFlavor() Leaving")
	var encryption *model.Encryption
	var integrity *model.Integrity

	if label == "" {
		return nil, errors.Errorf("label cannot be empty")
	}

	description := model.Description{
		Label:      label,
		FlavorPart: "CONTAINER_IMAGE",
	}

	meta := model.Meta{
		Description: description,
	}
	newUuid, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new UUID")
	}
	meta.ID = newUuid

	encryption = &model.Encryption{
		KeyURL: keyURL,
	}

	integrity = &model.Integrity{
		NotaryURL: notaryURL,
	}

	containerImageFlavor := model.Image{
		Meta:               meta,
		EncryptionRequired: encryptionRequired,
		Encryption:         encryption,
		IntegrityEnforced:  integrityEnforced,
		Integrity:          integrity,
	}

	flavor := ImageFlavor{
		Image: containerImageFlavor,
	}
	return &flavor, nil
}

//GetSignedImageFlavor is used to sign image flavor
func GetSignedImageFlavor(flavorString string, rsaPrivateKeyLocation string) (string, error) {
	log.Trace("flavor/image_flavor:GetSignedImageFlavor() Entering")
	defer log.Trace("flavor/image_flavor:GetSignedImageFlavor() Leaving")
	var privateKey *rsa.PrivateKey
	var flavorInterface ImageFlavor
	if rsaPrivateKeyLocation == "" {
		log.Error("No RSA Key file path provided")
		return "", errors.New("No RSA Key file path provided")
	}

	priv, err := ioutil.ReadFile(rsaPrivateKeyLocation)
	if err != nil {
		log.Error("No RSA private key found")
		return "", err
	}

	privPem, _ := pem.Decode(priv)
	parsedKey, err := x509.ParsePKCS8PrivateKey(privPem.Bytes)
	if err != nil {
		log.Error("Cannot parse RSA private key from file")
		return "", err
	}

	privateKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		log.Error("Unable to parse RSA private key")
		return "", err
	}
	hashEntity := sha512.New384()
	hashEntity.Write([]byte(flavorString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA384, hashEntity.Sum(nil))
	signatureString := base64.StdEncoding.EncodeToString(signature)

	json.Unmarshal([]byte(flavorString), &flavorInterface)

	signedFlavor := SignedImageFlavor{
		ImageFlavor: flavorInterface.Image,
		Signature:   signatureString,
	}

	signedFlavorJSON, err := json.Marshal(signedFlavor)
	if err != nil {
		return "", errors.New("Error while marshalling signed image flavor: " + err.Error())
	}

	return string(signedFlavorJSON), nil
}
