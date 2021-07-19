/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package router

import (
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/controllers"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/domain/models"
)

func SetCreateCaCertificatesRoutes(router *mux.Router, certStore *models.CertificatesStore) *mux.Router {
	defaultLog.Trace("router/create_ca_certificates:SetCreateCaCertificatesRoutes() Entering")
	defer defaultLog.Trace("router/create_ca_certificates:SetCreateCaCertificatesRoutes() Leaving")

	caCertController := controllers.CaCertificatesController{CertStore: certStore}

	router.Handle("/ca-certificates",
		ErrorHandler(permissionsHandler(JsonResponseHandler(caCertController.Create),
			[]string{constants.CaCertificatesCreate}))).Methods("POST")
	return router
}
