/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package common

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/config"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/defender"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/domain"
	"github.com/intel-secl/intel-secl/v3/pkg/authservice/types"
	"net/http"
	"time"
)

var defend *defender.Defender

func init() {
	c, _ := config.LoadConfiguration()

	defend = defender.New(c.AuthDefender.MaxAttempts,
		time.Duration(c.AuthDefender.IntervalMins)*time.Minute,
		time.Duration(c.AuthDefender.LockoutDurationMins)*time.Minute)
	quit := make(chan struct{})

	go defend.CleanupTask(quit)

}

func HttpHandleUserAuth(u domain.UserStore, username, password string) (int, error) {
	// first let us make sure that this is not a user that is banned

	foundInDefendList := false
	// check if we have an entry for the client in the defend map.
	// There are several scenarios in this case
	if client, ok := defend.Client(username); ok {
		foundInDefendList = true
		if client.Banned() {
			// case 1. Client is banned - however, the ban expired but cleanup is not done.
			// just delete the client from the map
			if client.BanExpired() {
				defend.RemoveClient(client.Key())
			} else {
				return http.StatusTooManyRequests, fmt.Errorf("Maximum login attempts exceeded for user : %s. Banned !", username)
			}
		}
	}

	// fetch by user
	user, err := u.Retrieve(types.User{Name: username})
	if err != nil {
		return http.StatusUnauthorized, fmt.Errorf("BasicAuth failure: could not retrieve user: %s error: %s", username, err)
	}
	if err := user.CheckPassword([]byte(password)); err != nil {
		if defend.Inc(username) {
			return http.StatusTooManyRequests, fmt.Errorf("Authentication failure - maximum login attempts exceeded for user : %s. Banned !", username)
		}
		return http.StatusUnauthorized, fmt.Errorf("BasicAuth failure: password mismatch, user: %s, error : %s", username, err)
	}
	// If we found the user earlier in the defend list, we should now remove as user is authorized
	if foundInDefendList {
		if client, ok := defend.Client(username); ok {
			defend.RemoveClient(client.Key())
		}
	}
	return 0, nil
}
