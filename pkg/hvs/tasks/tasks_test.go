/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDekGenerate(t *testing.T) {
	dek := "hello"
	task := CreateDek{
		DekStore: &dek,
	}
	if err := task.Validate(); err == nil {
		t.Error("first validation should not pass")
	}
	if err := task.Run(); err != nil {
		t.Error("run failed:", err.Error())
	}
	t.Log("Generated key:")
	t.Log(dek)
	if err := task.Validate(); err != nil {
		t.Error("second validation should pass:", err.Error())
	}
}

func TestDefaultFlavorGroupDes(t *testing.T) {
	// check if default flavor strings are correct
	for _, fg := range defaultFlavorGroups() {
		t.Log(fg)
	}
}

func TestDefaultFlavorTemplateDes(t *testing.T) {
	// check if default flavor templates are retrieved
	ft := CreateDefaultFlavorTemplate{
		Directory: "../../../build/linux/hvs/templates/",
	}
	flavorTemplates, err := ft.getTemplates()
	assert.NoError(t, err)
	for template, _ := range flavorTemplates {
		t.Log(template)
	}
}

func TestDefaultFlavorTemplateFaultDes(t *testing.T) {
	// check if default flavor templates are retrieved
	ft := CreateDefaultFlavorTemplate{
		Directory: "",
	}
	_, err := ft.readDefaultTemplates()
	assert.Error(t, err)
}
