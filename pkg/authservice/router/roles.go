/*
 *  Copyright (C) 2020 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package router

import (
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/controllers"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/domain"
)

func SetRolesRoutes(r *mux.Router, db domain.AASDatabase) *mux.Router {
	defaultLog.Trace("router/roles:SetRolesRoutes() Entering")
	defer defaultLog.Trace("router/roles:SetRolesRoutes() Leaving")

	controller := controllers.RolesController{Database: db}

	r.Handle("/roles", ErrorHandler(ResponseHandler(controller.CreateRole, "application/json"))).Methods("POST")
	r.Handle("/roles", ErrorHandler(ResponseHandler(controller.QueryRoles, "application/json"))).Methods("GET")
	r.Handle("/roles/{id}", ErrorHandler(ResponseHandler(controller.DeleteRole, ""))).Methods("DELETE")
	r.Handle("/roles/{id}", ErrorHandler(ResponseHandler(controller.GetRole, "application/json"))).Methods("GET")
	r.Handle("/roles/{id}", ErrorHandler(ResponseHandler(controller.UpdateRole, ""))).Methods("PATCH")
	return r
}
