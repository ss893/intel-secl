/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package mock

import (
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/types"
	ct "github.com/intel-secl/intel-secl/v3/pkg/model/aas"
)

type MockUserStore struct {
	CreateFunc      func(types.User) (*types.User, error)
	RetrieveFunc    func(types.User) (*types.User, error)
	RetrieveAllFunc func(types.User) (types.Users, error)
	UpdateFunc      func(types.User) error
	DeleteFunc      func(types.User) error
}

func (m *MockUserStore) Create(user types.User) (*types.User, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil, nil
}

func (m *MockUserStore) Retrieve(user types.User) (*types.User, error) {
	if m.RetrieveFunc != nil {
		return m.RetrieveFunc(user)
	}
	return nil, nil
}

func (m *MockUserStore) RetrieveAll(u types.User) (types.Users, error) {
	if m.RetrieveAllFunc != nil {
		return m.RetrieveAllFunc(u)
	}
	return nil, nil
}

func (m *MockUserStore) Update(user types.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func (m *MockUserStore) Delete(user types.User) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(user)
	}
	return nil
}

func (m *MockUserStore) GetUserRoleByID(u types.User, roleID string) (types.Role, error) {
	return types.Role{}, nil
}

func (m *MockUserStore) GetRoles(user types.User, rs *types.RoleSearch, includeID bool) ([]types.Role, error) {
	return nil, nil
}

func (m *MockUserStore) GetPermissions(user types.User, rs *types.RoleSearch) ([]ct.PermissionInfo, error) {
	return nil, nil
}

func (m *MockUserStore) AddRoles(u types.User, roleList types.Roles, mustAddAllRoles bool) error {
	return nil
}

func (m *MockUserStore) DeleteRole(u types.User, roleID string, svcFltr []string) error {
	return nil
}
