/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package ta

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/hvs/constants"
	cos "github.com/intel-secl/intel-secl/v4/pkg/lib/common/os"
	"net/url"
	"strings"
	"time"

	taModel "github.com/intel-secl/intel-secl/v4/pkg/model/ta"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

var (
	defaultTimeout = 10 * time.Second
)

func NewNatsTAClient(natsServers []string, natsHostID string) (TAClient, error) {

	if len(natsServers) == 0 {
		return nil, errors.New("client/nats_client:NewNatsTAClient() At least one nats-server must be provided.")
	}

	if natsHostID == "" {
		return nil, errors.New("client/nats_client:NewNatsTAClient() The nats-host-id was not provided")
	}

	client := natsTAClient{
		natsServers: natsServers,
		natsHostID:  natsHostID,
	}

	return &client, nil
}

func (client *natsTAClient) newNatsConnection() (*nats.EncodedConn, error) {

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := cos.GetDirFileContents(constants.TrustedCaCertsDir, "*.pem")
	if err != nil {
		log.Errorf("client/nats_client:newNatsConnection() Failed to append %q to RootCAs: %v", "/etc/hvs/certs/trustedca/nats-ca.pem", err)
	}

	for _, rootCACert := range certs {
		if ok := rootCAs.AppendCertsFromPEM(rootCACert); !ok {
			log.Info("client/nats_client:newNatsConnection() No certs appended, using system certs only")
		}
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}

	conn, err := nats.Connect(strings.Join(client.natsServers, ","),
		nats.Secure(&tlsConfig),
		nats.UserCredentials(constants.NatsCredentials),
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			if s != nil {
				log.Infof("client/nats_client:newNatsConnection() NATS: Could not process subscription for subject %q: %v", s.Subject, err)
			} else {
				log.Infof("client/nats_client:newNatsConnection() NATS: Unknown error: %v", err)
			}
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Infof("client/nats_client:newNatsConnection() NATS: Client disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Infof("client/nats_client:newNatsConnection() NATS: Client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Infof("client/nats_client:newNatsConnection() NATS: Client closed")
		}))

	if err != nil {
		return nil, fmt.Errorf("Failed to create nats connection: %+v", err)
	}

	encodedConn, err := nats.NewEncodedConn(conn, "json")
	if err != nil {
		return nil, fmt.Errorf("client/nats_client:newNatsConnection() Failed to create encoded connection: %+v", err)
	}

	return encodedConn, nil
}

type natsTAClient struct {
	natsServers    []string
	natsConnection *nats.EncodedConn
	natsHostID     string
}

func (client *natsTAClient) GetHostInfo() (taModel.HostInfo, error) {
	hostInfo := taModel.HostInfo{}
	conn, err := client.newNatsConnection()
	if err != nil {
		return hostInfo, errors.Wrap(err, "client/nats_client:GetHostInfo() Error establishing connection to nats server")
	}
	defer conn.Close()

	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsHostInfoRequest), nil, &hostInfo, defaultTimeout)
	if err != nil {
		return hostInfo, errors.Wrap(err, "client/nats_client:GetHostInfo() Error getting Host Info")
	}
	return hostInfo, nil
}

func (client *natsTAClient) GetTPMQuote(nonce string, pcrList []int, pcrBankList []string) (taModel.TpmQuoteResponse, error) {
	quoteResponse := taModel.TpmQuoteResponse{}
	nonceBytes, err := base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		return quoteResponse, errors.Wrap(err, "client/nats_client:GetTPMQuote() Error decoding nonce from base64 to bytes")
	}
	quoteRequest := taModel.TpmQuoteRequest{
		Nonce:    nonceBytes,
		Pcrs:     pcrList,
		PcrBanks: pcrBankList,
	}

	conn, err := client.newNatsConnection()
	if err != nil {
		return quoteResponse, errors.Wrap(err, "client/nats_client:GetTPMQuote() Error establishing connection to nats server")
	}
	defer conn.Close()

	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsQuoteRequest), &quoteRequest, &quoteResponse, defaultTimeout)
	if err != nil {
		return quoteResponse, errors.Wrap(err, "client/nats_client:GetTPMQuote() Error getting quote")
	}
	return quoteResponse, nil
}

func (client *natsTAClient) GetAIK() ([]byte, error) {
	conn, err := client.newNatsConnection()
	if err != nil {
		return nil, errors.Wrap(err, "client/nats_client:GetAIK() Error establishing connection to nats server")
	}
	defer conn.Close()

	var aik []byte
	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsAikRequest), nil, &aik, defaultTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "client/nats_client:GetAIK() Error getting AIK")
	}
	return aik, nil
}

func (client *natsTAClient) GetBindingKeyCertificate() ([]byte, error) {
	conn, err := client.newNatsConnection()
	if err != nil {
		return nil, errors.Wrap(err, "client/nats_client:GetBindingKeyCertificate() Error establishing connection to nats server")
	}
	defer conn.Close()

	var bk []byte
	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsBkRequest), nil, &bk, defaultTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "client/nats_client:GetBindingKeyCertificate() Error getting binding key")
	}
	return bk, nil
}

func (client *natsTAClient) DeployAssetTag(hardwareUUID, tag string) error {
	var err error
	var tagWriteRequest taModel.TagWriteRequest
	tagWriteRequest.Tag, err = base64.StdEncoding.DecodeString(tag)
	if err != nil {
		return errors.Wrap(err, "client/nats_client:DeployAssetTag() Error decoding tag from base64 to bytes")
	}
	tagWriteRequest.HardwareUUID = hardwareUUID

	conn, err := client.newNatsConnection()
	if err != nil {
		return errors.Wrap(err, "client/nats_client:DeployAssetTag() Error establishing connection to nats server")
	}
	defer conn.Close()

	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsDeployAssetTagRequest), &tagWriteRequest, &nats.Msg{}, defaultTimeout)
	if err != nil {
		return errors.Wrap(err, "client/nats_client:DeployAssetTag() Error deploying asset tag")
	}
	return nil
}

func (client *natsTAClient) DeploySoftwareManifest(manifest taModel.Manifest) error {
	conn, err := client.newNatsConnection()
	if err != nil {
		return errors.Wrap(err, "client/nats_client:DeploySoftwareManifest() Error establishing connection to nats server")
	}
	defer conn.Close()

	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsDeployManifestRequest), &manifest, &nats.Msg{}, defaultTimeout)
	if err != nil {
		return errors.Wrap(err, "client/nats_client:DeploySoftwareManifest() Error deploying software flavor")
	}
	return nil
}

func (client *natsTAClient) GetMeasurementFromManifest(manifest taModel.Manifest) (taModel.Measurement, error) {
	measurement := taModel.Measurement{}
	conn, err := client.newNatsConnection()
	if err != nil {
		return measurement, errors.Wrap(err, "client/nats_client:GetMeasurementFromManifest() Error establishing connection to nats server")
	}
	defer conn.Close()

	err = conn.Request(taModel.CreateSubject(client.natsHostID, taModel.NatsApplicationMeasurementRequest), &manifest, &measurement, defaultTimeout)
	if err != nil {
		return measurement, errors.Wrap(err, "client/nats_client:GetMeasurementFromManifest() Error getting measurement from TA")
	}
	return measurement, nil
}

func (client *natsTAClient) GetBaseURL() *url.URL {
	return nil
}
