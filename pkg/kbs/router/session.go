/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package router

import (
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/config"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/controllers"
)

//setSessionRoutes registers routes to perform session management operations
func setSessionRoutes(router *mux.Router, kbsConfig *config.Configuration) *mux.Router {
	defaultLog.Trace("router/keys:setSessionRoutes() Entering")
	defer defaultLog.Trace("router/keys:setSessionRoutes() Leaving")

	sessionController := controllers.NewSessionController(kbsConfig, constants.TrustedCaCertsDir)

	router.Handle("/session",
		ErrorHandler(permissionsHandlerUsingTLSMAuth(JsonResponseHandler(sessionController.Create),
			kbsConfig.AASApiUrl, kbsConfig.KBS))).Methods("POST")
	return router
}
