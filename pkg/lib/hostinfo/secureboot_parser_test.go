/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package hostinfo

import (
	"reflect"
	"testing"

	model "github.com/intel-secl/intel-secl/v4/pkg/model/ta"
)

func testSecureBootParser(t *testing.T, expectedResults *model.HostInfo) {
	hostInfo := model.HostInfo{}

	secureBootParser := secureBootParser{}

	secureBootParser.Parse(&hostInfo)

	if !reflect.DeepEqual(hostInfo.HardwareFeatures.UEFI, expectedResults.HardwareFeatures.UEFI) {
		t.Errorf("The parsed UEFI data does not match the expected results.\nExpected: %+v\nActual: %+v\n", expectedResults.HardwareFeatures.UEFI, hostInfo.HardwareFeatures.UEFI)
	}
}

func TestSecureBootWhitley(t *testing.T) {

	secureBootFile = "test_data/whitley/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c"

	expectedResults := model.HostInfo{}
	expectedResults.HardwareFeatures.UEFI.Meta.SecureBootEnabled = true

	testSecureBootParser(t, &expectedResults)
}

func TestSecureBootPurley(t *testing.T) {

	// purley doesn't have efi var files -- provide a non-existent path
	secureBootFile = "test_data/purley/nosuchfile"

	expectedResults := model.HostInfo{}
	expectedResults.HardwareFeatures.UEFI.Meta.SecureBootEnabled = false

	testSecureBootParser(t, &expectedResults)
}

func TestSecureBootShortFile(t *testing.T) {

	// test a file that doesn't have enough data -- it not error and
	// show secure-boot as disabled
	secureBootFile = "test_data/misc/SecureBootShortFile"

	expectedResults := model.HostInfo{}
	expectedResults.HardwareFeatures.UEFI.Meta.SecureBootEnabled = false

	testSecureBootParser(t, &expectedResults)
}

func TestSecureBootZeroFile(t *testing.T) {

	// test a file that has zeros (secure boot is not enabled)
	secureBootFile = "test_data/misc/SecureBootZeroFile"

	expectedResults := model.HostInfo{}
	expectedResults.HardwareFeatures.UEFI.Meta.SecureBootEnabled = false

	testSecureBootParser(t, &expectedResults)

}
