/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package common

import (
	"errors"
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/config"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/defender"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/domain"
	"github.com/intel-secl/intel-secl/v4/pkg/authservice/types"
	commErr "github.com/intel-secl/intel-secl/v4/pkg/lib/common/err"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
	"net/http"
	"time"
)

var defaultLog = log.GetDefaultLogger()

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

//Generates JWT token from key pair
func CreateJWTToken(keyPair nkeys.KeyPair, issuerKeyPair nkeys.KeyPair, creatorType, clientType, entityName string) (string, error) {
	defaultLog.Trace("common/common:CreateJWTToken() Entering")
	defer defaultLog.Trace("common/common:CreateJWTToken() Leaving")

	pk, err := keyPair.PublicKey()
	if err != nil {
		return "", &commErr.ResourceError{Message: "Error getting public part of account key pair"}
	}

	var claims interface{}
	var token string

	if creatorType == constants.Operator {
		// create a new operator claim
		claims = jwt.NewOperatorClaims(pk)
		claims.(*jwt.OperatorClaims).Name = entityName
		token, err = claims.(*jwt.OperatorClaims).Encode(issuerKeyPair)
	} else if creatorType == constants.Account {
		// create a new account claim
		claims = jwt.NewAccountClaims(pk)
		claims.(*jwt.AccountClaims).Name = entityName
		token, err = claims.(*jwt.AccountClaims).Encode(issuerKeyPair)
	} else if creatorType == constants.User {
		// create a new user claim
		claims = jwt.NewUserClaims(pk)
		claims.(*jwt.UserClaims).Name = entityName
		if clientType == constants.ComponentTypeHvs {
			claims.(*jwt.UserClaims).Pub.Allow = []string{"trust-agent.>"}
			claims.(*jwt.UserClaims).Sub.Allow = []string{"_INBOX.>"}
		} else if clientType == constants.ComponentTypeTa {
			claims.(*jwt.UserClaims).Pub.Deny = []string{">"}
			claims.(*jwt.UserClaims).Sub.Allow = []string{"trust-agent." + entityName + ".>"}
			claims.(*jwt.UserClaims).Resp = &jwt.ResponsePermission{
				MaxMsgs: 1,
				Expires: 0,
			}
		} else {
			return "", errors.New("invalid type provided")
		}
		claims.(*jwt.UserClaims).Expires = time.Now().Add(time.Hour * 48).Unix()
		token, err = claims.(*jwt.UserClaims).Encode(issuerKeyPair)
	}

	if err != nil {
		return "", err
	}
	return token, nil
}
