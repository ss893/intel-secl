/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package router

import (
	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v3/pkg/kbs/config"
	"github.com/intel-secl/intel-secl/v3/pkg/kbs/controllers"
)

//setSessionRoutes registers routes to perform session management operations
func setSessionRoutes(router *mux.Router, aasAPIUrl string, kbsConfig config.KBSConfig) *mux.Router {
	defaultLog.Trace("router/keys:setSessionRoutes() Entering")
	defer defaultLog.Trace("router/keys:setSessionRoutes() Leaving")

	sessionController := controllers.NewSessionController()

	router.Handle("/session",
		ErrorHandler(permissionsHandlerUsingTLSMAuth(JsonResponseHandler(sessionController.Create),
			aasAPIUrl, kbsConfig))).Methods("POST")

	return router
}
