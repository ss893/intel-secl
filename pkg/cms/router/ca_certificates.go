/*
 *  Copyright (C) 2020 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package router

import (
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v4/pkg/cms/controllers"
	log "github.com/sirupsen/logrus"
)

// SetCACertificatesRoutes is used to set the endpoints for CA certificate handling APIs
func SetCACertificatesRoutes(router *mux.Router) *mux.Router {
	log.Trace("router/ca_certificates:SetCACertificatesRoutes() Entering")
	defer log.Trace("router/ca_certificates:SetCACertificatesRoutes() Leaving")

	caCertController := controllers.CACertificatesController{}
	router.HandleFunc("/ca-certificates", caCertController.GetCACertificates).Methods("GET")
	return router
}
