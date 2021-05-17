/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package constants

const (
	ServiceName         = "WPM"
	ExtendedServiceName = "Workload Policy Manager"
	ServiceDir          = "wpm/"
	ServiceUserName     = "wpm"

	HomeDir               = "/opt/" + ServiceDir
	ExecLinkPath          = "/usr/bin/" + ServiceUserName
	LogDir                = "/var/log/" + ServiceDir
	ConfigDir             = "/etc/" + ServiceDir
	ConfigFile            = "config"
	DefaultConfigFilePath = ConfigDir + "config.yml"
	FlavorsDir            = HomeDir + "flavors/"
	VmImagesDir           = HomeDir + "vm-images/"
	EncryptedVmImagesDir  = HomeDir + "encrypted-vm-images/"

	// certificates' path
	FlavorSigningCertDir = ConfigDir + "certs/flavorsign/"
	TrustedCaCertsDir    = ConfigDir + "certs/trustedca/"
	EnvelopekeyDir       = ConfigDir + "certs/kbs/"

	// flavor signing key and cert
	FlavorSigningCertFile     = FlavorSigningCertDir + "flavor-signing.pem"
	FlavorSigningKeyFile      = FlavorSigningCertDir + "flavor-signing.key"
	DefaultWpmFlavorSigningCn = "WPM Flavor Signing Certificate"
	DefaultWpmSan             = "127.0.0.1,localhost"

	EnvelopePublickeyLocation  = EnvelopekeyDir + "envelopePublicKey.pub"
	EnvelopePrivatekeyLocation = EnvelopekeyDir + "envelopePrivateKey.pem"

	//log config
	DefaultLogLevel     = "info"
	DefaultLogMaxlength = 1500

	// create key parameters
	KbsEncryptAlgo            = "AES"
	KbsKeyLength              = 256
	KbsCipherMode             = "GCM"
	DefaultKeyAlgorithm       = "rsa"
	DefaultKeyAlgorithmLength = 3072
	CertApproverGroupName     = "CertApprover"
	KBSKeyRetrievalGroupName  = "KeyCRUD"
	SampleUUID                = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
)
