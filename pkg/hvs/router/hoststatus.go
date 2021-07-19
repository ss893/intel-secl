/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package router

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/controllers"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/postgres"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/validation"
)

// SetHostStatusRoutes registers routes for host-status APIs
func SetHostStatusRoutes(router *mux.Router, store *postgres.DataStore) *mux.Router {
	defaultLog.Trace("router/hoststatus:SetHostStatusRoutes() Entering")
	defer defaultLog.Trace("router/hoststatus:SetHostStatusRoutes() Leaving")

	hoststatusStore := postgres.NewHostStatusStore(store)
	hoststatusController := controllers.HostStatusController{Store: hoststatusStore}

	router.Handle("/host-status", ErrorHandler(permissionsHandler(JsonResponseHandler(hoststatusController.Search),
		[]string{constants.HostStatusSearch}))).Methods("GET")

	hostStatusIdExpr := fmt.Sprintf("%s%s", "/host-status/", validation.IdReg)
	router.Handle(hostStatusIdExpr, ErrorHandler(permissionsHandler(JsonResponseHandler(hoststatusController.Retrieve),
		[]string{constants.HostStatusRetrieve}))).Methods("GET")

	return router
}
