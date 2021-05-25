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

// SetTpmEndorsementRoutes registers routes for tpm-endorsements
func SetTpmEndorsementRoutes(router *mux.Router, store *postgres.DataStore) *mux.Router {
	defaultLog.Trace("router/flavorgroups:SetTpmEndorsementRoutes() Entering")
	defer defaultLog.Trace("router/flavorgroups:SetTpmEndorsementRoutes() Leaving")

	tpmEndorsementStore := postgres.NewTpmEndorsementStore(store)
	tpmEndorsementController := controllers.TpmEndorsementController{Store: tpmEndorsementStore}
	tpmEndorsementIdExpr := fmt.Sprintf("%s%s", "/tpm-endorsements/", validation.IdReg)

	router.Handle("/tpm-endorsements",
		ErrorHandler(permissionsHandler(JsonResponseHandler(tpmEndorsementController.Create),
			[]string{constants.TpmEndorsementCreate}))).Methods("POST")

	router.Handle(tpmEndorsementIdExpr,
		ErrorHandler(permissionsHandler(JsonResponseHandler(tpmEndorsementController.Update),
			[]string{constants.TpmEndorsementStore}))).Methods("PUT")

	router.Handle("/tpm-endorsements",
		ErrorHandler(permissionsHandler(JsonResponseHandler(tpmEndorsementController.Search),
			[]string{constants.TpmEndorsementSearch}))).Methods("GET")

	router.Handle(tpmEndorsementIdExpr,
		ErrorHandler(permissionsHandler(ResponseHandler(tpmEndorsementController.Delete),
			[]string{constants.TpmEndorsementDelete}))).Methods("DELETE")

	router.Handle("/tpm-endorsements",
		ErrorHandler(permissionsHandler(ResponseHandler(tpmEndorsementController.DeleteCollection),
			[]string{constants.TpmEndorsementSearch, constants.TpmEndorsementDelete}))).Methods("DELETE")

	router.Handle(tpmEndorsementIdExpr,
		ErrorHandler(permissionsHandler(JsonResponseHandler(tpmEndorsementController.Retrieve),
			[]string{constants.TpmEndorsementRetrieve}))).Methods("GET")

	return router
}
