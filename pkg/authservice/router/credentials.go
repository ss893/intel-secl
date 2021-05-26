/*
 *  Copyright (C) 2021 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package router

import (
	"github.com/gorilla/mux"
	consts "github.com/intel-secl/intel-secl/v4/pkg/authservice/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/controllers"
)

func SetCredentialsRoutes(r *mux.Router, username string) *mux.Router {
	defaultLog.Trace("router/credentials_controller:SetCredentialsRoutes() Entering")
	defer defaultLog.Trace("router/jwt_certificate:SetCredentialsRoutes() Leaving")

	controller := controllers.CredentialsController{Username: username}
	r.Handle("/credentials", ErrorHandler(permissionsHandler(ResponseHandler(controller.CreateCredentials,
		"text/plain"), []string{consts.CredentialCreate}))).Methods("POST")

	return r
}
