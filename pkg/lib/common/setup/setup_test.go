/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package setup_test

import (
	"io"
	"testing"

	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/setup"
	"github.com/pkg/errors"
)

var testConfigMap = map[string]string{
	"TEST_ARG_KEY_ONE": "test-arg-val-one",
	"TEST_ARG_KEY_TWO": "test-arg-val-two",
}

type testTaskOne struct {
	Arg string

	hasConfig bool
	hasRun    bool
}

func (t *testTaskOne) Run() error {
	t.hasRun = true
	return nil
}

func (t *testTaskOne) Validate() error {
	if t.hasRun &&
		t.Arg == "test-arg-val-one" {
		return nil
	}
	return errors.New("validation failed")
}

type testTaskTwo struct {
	Arg string

	hasConfig bool
	hasRun    bool
}

func (t *testTaskTwo) Run() error {
	t.hasRun = true
	return nil
}

func (t *testTaskTwo) Validate() error {
	if t.hasRun &&
		t.Arg == "test-arg-val-two" {
		return nil
	}
	return errors.New("validation failed")
}

func (t *testTaskOne) PrintHelp(w io.Writer) {}
func (t *testTaskTwo) PrintHelp(w io.Writer) {}

func (t *testTaskOne) SetName(string, string) {}
func (t *testTaskTwo) SetName(string, string) {}

func TestSetupRunner(t *testing.T) {
	runner := setup.NewRunner()
	runner.AddTask("task-1", "", &testTaskOne{
		Arg: testConfigMap["TEST_ARG_KEY_ONE"],
	})
	runner.AddTask("task-1", "", &testTaskTwo{
		Arg: testConfigMap["TEST_ARG_KEY_TWO"],
	})
	if err := runner.RunAll(true); err != nil {
		t.Error("Failed to run all tasks:", err.Error())
	}
}
