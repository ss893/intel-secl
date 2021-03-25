/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package imageflavor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateImageFlavor(t *testing.T) {
	imageFlavor, err := CreateImageFlavor("label", "", "cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", false)
	assert.NotNil(t, err)
	assert.Equal(t, imageFlavor, "")
}

func TestCreateImageFlavorToFile(t *testing.T) {
	imageFlavor, err := CreateImageFlavor("label", "image_flavor.txt", "cirros-x86.qcow2", "cirros-x86.qcow2_enc", "", false)
	assert.NotNil(t, err)
	assert.Equal(t, imageFlavor, "")
}
