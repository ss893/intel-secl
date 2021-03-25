/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package containerimageflavor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateContainerImageFlavor(t *testing.T) {
	imageFlavor, err := CreateContainerImageFlavor("hello-world", "latest", "", "", "", false, false, "", "")
	assert.NotNil(t, err)
	assert.Equal(t, imageFlavor, "")
}

func TestCreateContainerImageFlavorToFile(t *testing.T) {
	imageFlavor, err := CreateContainerImageFlavor("hello-world", "latest", "", "", "", false, false, "", "container_flavor.txt")
	assert.NotNil(t, err)
	assert.Equal(t, imageFlavor, "")
}
