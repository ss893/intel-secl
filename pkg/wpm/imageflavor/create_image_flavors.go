/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package imageflavor

/*
 *
 * @author srege
 *
 */
import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	cLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/flavor"
	consts "github.com/intel-secl/intel-secl/v4/pkg/wpm/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	log = cLog.GetDefaultLogger()
)

//CreateImageFlavor is used to create flavor of an encrypted image
func CreateImageFlavor(flavorLabel, outputFlavorFilename, inputImageFilename,
	outputEncImageFilename, keyID string, integrityRequired bool) (string, error) {
	log.Trace("pkg/wpm/imageflavor/create_image_flavors.go:CreateImageFlavor() Entering")
	defer log.Trace("pkg/wpm/imageflavor/create_image_flavors.go:CreateImageFlavor() Leaving")

	var err error
	var wrappedKey []byte
	var keyUrlString string
	encRequired := true

	if filepath.IsAbs(inputImageFilename) {
		return "", errors.New("Image filename should not be an absolute path")
	}
	if filepath.IsAbs(outputFlavorFilename) {
		return "", errors.New("Image flavor filename should not be an absolute path")
	}
	if filepath.IsAbs(outputEncImageFilename) {
		return "", errors.New("Output image filename should not be an absolute path")
	}
	imageFilePath := filepath.Join(consts.VmImagesDir, inputImageFilename)
	inputImageFilePath := filepath.Clean(imageFilePath)
	outputEncImageFilePath := filepath.Join(consts.EncryptedVmImagesDir, outputEncImageFilename)
	outputFlavorFilePath := filepath.Join(consts.FlavorsDir, outputFlavorFilename)

	//Determine if encryption is required
	outputEncImageFilename = strings.TrimSpace(outputEncImageFilename)
	if len(outputEncImageFilename) <= 0 {
		encRequired = false
	}

	// set logger fields
	log = log.WithFields(logrus.Fields{
		"flavorLabel":            flavorLabel,
		"encryptionRequired":     encRequired,
		"integrityrequired":      integrityRequired,
		"inputImageFilePath":     inputImageFilePath,
		"outputFlavorFilePath":   outputFlavorFilePath,
		"outputEncImageFilePath": outputEncImageFilePath,
	})

	//Error if image specified doesn't exist
	_, err = os.Stat(inputImageFilePath)
	if os.IsNotExist(err) {
		return "", errors.Wrap(err, "I/O error reading image file: "+err.Error())
	}

	//Encrypt the image with the key
	if encRequired {
		// fetch the key to encrypt the image
		wrappedKey, keyUrlString, err = util.FetchKey(keyID, "")
		if err != nil {
			return "", errors.Wrap(err, "Fetch key failed: "+err.Error())
		}
		// encrypt the image with key retrieved from KBS
		err = util.Encrypt(inputImageFilePath, consts.EnvelopePrivatekeyLocation, outputEncImageFilePath, wrappedKey)
		if err != nil {
			return "", errors.Wrap(err, "Image encryption failed: "+err.Error())
		}
		imageFilePath = outputEncImageFilePath
	}

	//Check the encrypted image output file
	imageFile, err := ioutil.ReadFile(imageFilePath)
	if err != nil {
		return "", errors.Wrap(err, "I/O Error creating encrypted image file: "+err.Error())
	}

	//Take the digest of the encrypted image
	digest := sha512.Sum384(imageFile)

	//Create image flavor
	imageFlavor, err := flavor.GetImageFlavor(flavorLabel, encRequired, keyUrlString, base64.StdEncoding.EncodeToString(digest[:]))
	if err != nil {
		return "", errors.Wrap(err, "Error creating image flavor: "+err.Error())
	}

	//Marshall the image flavor to a JSON string
	imageFlavorJSON, err := json.Marshal(imageFlavor)
	if err != nil {
		return "", errors.Wrap(err, "Error while marshalling image flavor: "+err.Error())
	}

	signedFlavor, err := flavor.GetSignedImageFlavor(string(imageFlavorJSON), consts.FlavorSigningKeyFile)
	if err != nil {
		return "", errors.Wrap(err, "Error signing flavor for image: "+err.Error())
	}

	log.Info("pkg/wpm/imageflavor/create_image_flavors.go:CreateImageFlavor() Successfully created image flavor")
	log.Debugf("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Successfully created image flavor %s", signedFlavor)

	//If no output flavor file path was specified, return the marshalled image flavor
	if len(strings.TrimSpace(outputFlavorFilename)) <= 0 {
		return signedFlavor, nil
	}

	//Otherwise, write it to the specified file
	err = ioutil.WriteFile(outputFlavorFilePath, []byte(signedFlavor), 0600)
	if err != nil {
		return "", errors.Wrapf(err, "I/O Error writing image flavor to output file %s", outputFlavorFilePath)
	}

	log.Info("pkg/imageflavor/create_image_flavors.go:CreateImageFlavor() Successfully wrote image flavor to file")
	return "", nil
}
