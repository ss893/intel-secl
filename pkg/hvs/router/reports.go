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
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/domain"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/postgres"
	consts "github.com/intel-secl/intel-secl/v4/pkg/lib/common/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/validation"
)

// SetReportRoutes registers routes for reports
func SetReportRoutes(router *mux.Router, store *postgres.DataStore, hostTrustManager domain.HostTrustManager) *mux.Router {
	defaultLog.Trace("router/reports:SetReportRoutes() Entering")
	defer defaultLog.Trace("router/reports:SetReportRoutes() Leaving")

	reportStore := postgres.NewReportStore(store)
	hostStore := postgres.NewHostStore(store)
	hostStatusStore := postgres.NewHostStatusStore(store)
	reportController := controllers.NewReportController(reportStore, hostStore, hostStatusStore, hostTrustManager)

	reportIdExpr := fmt.Sprintf("%s%s", "/reports/", validation.IdReg)

	router.Handle("/reports",
		ErrorHandler(permissionsHandler(ResponseHandler(reportController.CreateSaml),
			[]string{constants.ReportCreate}))).Methods("POST").Headers("Accept", consts.HTTPMediaTypeSaml)

	router.Handle("/reports",
		ErrorHandler(permissionsHandler(JsonResponseHandler(reportController.Create),
			[]string{constants.ReportCreate}))).Methods("POST")

	router.Handle("/reports",
		ErrorHandler(permissionsHandler(ResponseHandler(reportController.SearchSaml),
			[]string{constants.ReportSearch}))).Methods("GET").Headers("Accept", consts.HTTPMediaTypeSaml)

	router.Handle(reportIdExpr,
		ErrorHandler(permissionsHandler(JsonResponseHandler(reportController.Retrieve),
			[]string{constants.ReportRetrieve}))).Methods("GET")

	router.Handle("/reports",
		ErrorHandler(permissionsHandler(JsonResponseHandler(reportController.Search),
			[]string{constants.ReportSearch}))).Methods("GET")

	return router
}
