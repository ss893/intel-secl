/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package tasks

import (
	"errors"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/config"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/domain"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/postgres/mock"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/types"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAdmin(t *testing.T) {
	m := &mock.MockDatabase{}
	c := config.Configuration{}
	var user *types.User
	var role *types.Role
	var permission *types.Permission
	m.MockUserStore.CreateFunc = func(u types.User) (*types.User, error) {
		user = &u
		return user, nil
	}
	m.MockUserStore.RetrieveFunc = func(u types.User) (*types.User, error) {
		if user == nil {
			return nil, errors.New("Record not found")
		}
		return user, nil
	}
	m.MockRoleStore.CreateFunc = func(r types.Role) (*types.Role, error) {
		role = &r
		return role, nil
	}
	m.MockRoleStore.RetrieveFunc = func(r *types.RoleSearch) (*types.Role, error) {
		if role == nil {
			return nil, errors.New("Record not found")
		}
		return role, nil
	}
	m.MockPermissionStore.CreateFunc = func(p types.Permission) (*types.Permission, error) {
		permission = &p
		return permission, nil
	}
	m.MockPermissionStore.RetrieveFunc = func(r *types.PermissionSearch) (*types.Permission, error) {
		if permission == nil {
			return nil, errors.New("Record not found")
		}
		return permission, nil
	}

	serviceConfig := config.AASConfig{
		Username: "username",
		Password: "password",
	}

	task := Admin{
		ServiceConfigPtr: &c.AAS,
		AASConfig:        serviceConfig,
		DatabaseFactory: func() (domain.AASDatabase, error) {
			return m, nil
		},
		ConsoleWriter: os.Stdout,
	}
	err := task.Run()
	assert.NoError(t, err)
}
