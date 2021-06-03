/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hostinfo

import (
	"encoding/binary"
	"os"

	model "github.com/intel-secl/intel-secl/v4/pkg/model/ta"
	"github.com/pkg/errors"
)

type secureBootParser struct{}

func (secureBootParser *secureBootParser) Init() error {
	return nil
}

func (secureBootParser *secureBootParser) Parse(hostInfo *model.HostInfo) error {

	// if the secure-boot file does not exists (ex. on older Purley systems) then assume
	// secure boot is disabled
	if _, err := os.Stat(secureBootFile); os.IsNotExist(err) {
		hostInfo.HardwareFeatures.UEFI.Meta.SecureBootEnabled = false
		return nil
	}

	// otherwise, read a 32bit int from the file -- any number other than zero indicates
	// secure-boot is enabled
	var results uint32
	file, err := os.Open(secureBootFile)
	if err != nil {
		return errors.Errorf("Failed to open secure-boot file %q", secureBootFile)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Errorf("Failed close secure-boot file %q: %s", secureBootFile, err.Error())
		}
	}()

	err = binary.Read(file, binary.LittleEndian, &results)
	if err != nil {
		log.Warnf("The secure-boot file %q is too small.  SecureBoot will be considered disabled", secureBootFile)
		return nil
	}

	if results == 0 {
		hostInfo.HardwareFeatures.UEFI.Meta.SecureBootEnabled = false
	} else {
		hostInfo.HardwareFeatures.UEFI.Meta.SecureBootEnabled = true
	}

	return nil
}
