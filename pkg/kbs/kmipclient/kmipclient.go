/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package kmipclient

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/gemalto/kmip-go"
	"github.com/gemalto/kmip-go/kmip14"
	"github.com/gemalto/kmip-go/kmip20"
	"github.com/gemalto/kmip-go/ttlv"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/domain/models"
	commLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/pkg/errors"
)

var defaultLog = commLog.GetDefaultLogger()

type kmipClient struct {
	KMIPVersion   string
	Config        tls.Config
	requestHeader kmip.RequestHeader
	ServerIP      string
	ServerPort    string
}

func NewKmipClient() KmipClient {
	return &kmipClient{}
}

var cipherSuites = []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_RSA_WITH_AES_128_CBC_SHA256}

// InitializeClient initializes all the values required for establishing connection to kmip server
func (kc *kmipClient) InitializeClient(version, serverIP, serverPort, hostname, username, password, clientKeyFilePath, clientCertificateFilePath, rootCertificateFilePath string) error {
	defaultLog.Trace("kmipclient/kmipclient:InitializeClient() Entering")
	defer defaultLog.Trace("kmipclient/kmipclient:InitializeClient() Leaving")

	if (version != constants.KMIP_1_4) && (version != constants.KMIP_2_0) {
		return errors.Errorf("kmipclient/kmipclient:InitializeClient()Invalid Kmip version %s provided", version)
	}
	kc.KMIPVersion = version

	if serverIP == "" {
		return errors.New("kmipclient/kmipclient:InitializeClient() KMIP server address is not provided")
	}
	kc.ServerIP = serverIP

	if serverPort == "" {
		return errors.New("kmipclient/kmipclient:InitializeClient() KMIP server port is not provided")
	}
	kc.ServerPort = serverPort

	if clientCertificateFilePath == "" {
		return errors.New("kmipclient/kmipclient:InitializeClient() KMIP client certificate is not provided")
	}

	if clientKeyFilePath == "" {
		return errors.New("kmipclient/kmipclient:InitializeClient() KMIP client key is not provided")
	}

	if rootCertificateFilePath == "" {
		return errors.New("kmipclient/kmipclient:InitializeClient() KMIP root certificate is not provided")
	}

	protocolVersion := kmip.ProtocolVersion{}
	if kc.KMIPVersion == constants.KMIP_2_0 {
		protocolVersion.ProtocolVersionMajor = 2
		protocolVersion.ProtocolVersionMinor = 0
	} else {
		protocolVersion.ProtocolVersionMajor = 1
		protocolVersion.ProtocolVersionMinor = 4
	}

	kc.requestHeader.ProtocolVersion = protocolVersion
	kc.requestHeader.BatchCount = 1

	if username != "" && password != "" {
		credential := kmip.Credential{}

		credential.CredentialType = kmip14.CredentialTypeUsernameAndPassword
		credential.CredentialValue = kmip.UsernameAndPasswordCredentialValue{
			Username: username,
			Password: password,
		}
		kc.requestHeader.Authentication = &kmip.Authentication{
			Credential: []kmip.Credential{
				credential,
			},
		}
		defaultLog.Info("kmipclient/kmipclient:InitializeClient() KMIP authentication with credential type UsernameAndPassword is added")
	}

	caCertificate, err := ioutil.ReadFile(rootCertificateFilePath)
	if err != nil {
		return errors.Wrap(err, "kmipclient/kmipclient:InitializeClient() Unable to read root certificate")
	}
	defaultLog.Debugf("kmipclient/kmipclient:InitializeClient() Loaded root certificate from %s", rootCertificateFilePath)

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(caCertificate)
	certificate, err := tls.LoadX509KeyPair(clientCertificateFilePath, clientKeyFilePath)
	if err != nil {
		return errors.Wrap(err, "kmipclient/kmipclient:InitializeClient() Failed to load client key and certificate")
	}
	defaultLog.Debugf("kmipclient/kmipclient:InitializeClient() Loaded client certificate from %s", clientCertificateFilePath)
	defaultLog.Debugf("kmipclient/kmipclient:InitializeClient() Loaded client key from %s", clientKeyFilePath)

	if hostname == "" {
		hostname = kc.ServerIP
	}

	kc.Config = tls.Config{
		ServerName:               hostname,
		CipherSuites:             cipherSuites,
		PreferServerCipherSuites: true,
		RootCAs:                  rootCAs,
		Certificates:             []tls.Certificate{certificate},
		MinVersion:               tls.VersionTLS12,
	}

	defaultLog.Info("kmipclient/kmipclient:InitializeClient() Kmip client initialized")
	return nil
}

// CreateSymmetricKey creates a symmetric key on kmip server
func (kc *kmipClient) CreateSymmetricKey(length int) (string, error) {
	defaultLog.Trace("kmipclient/kmipclient:CreateSymmetricKey() Entering")
	defer defaultLog.Trace("kmipclient/kmipclient:CreateSymmetricKey() Leaving")

	var createRequestPayLoad interface{}
	if kc.KMIPVersion == constants.KMIP_2_0 {
		createRequestPayLoad = models.CreateRequestPayload{
			ObjectType: kmip20.ObjectTypeSymmetricKey,
			Attributes: models.Attributes{
				CryptographicAlgorithm: kmip14.CryptographicAlgorithmAES,
				CryptographicLength:    int32(length),
				CryptographicUsageMask: kmip14.CryptographicUsageMaskEncrypt | kmip14.CryptographicUsageMaskDecrypt,
			},
		}
	} else {
		createRequestPayLoad = kmip.CreateRequestPayload{
			ObjectType: kmip14.ObjectTypeSymmetricKey,
			TemplateAttribute: kmip.TemplateAttribute{
				Attribute: []kmip.Attribute{
					{
						AttributeName:  "Cryptographic Algorithm",
						AttributeValue: kmip14.CryptographicAlgorithmAES,
					},
					{
						AttributeName:  "Cryptographic Length",
						AttributeValue: int32(length),
					},
					{
						AttributeName:  "Cryptographic Usage Mask",
						AttributeValue: kmip14.CryptographicUsageMaskEncrypt | kmip14.CryptographicUsageMaskDecrypt,
					},
				},
			},
		}
	}

	batchItem, decoder, err := kc.SendRequest(createRequestPayLoad, kmip14.OperationCreate)
	if err != nil {
		return "", errors.Wrap(err, "failed to perform create symmetric key operation")
	}

	var respPayload models.CreateResponsePayload
	err = decoder.DecodeValue(&respPayload, batchItem.ResponsePayload.(ttlv.TTLV))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode create symmetric key response payload")
	}

	return respPayload.UniqueIdentifier, nil
}

// CreateAsymmetricKeyPair creates a asymmetric key on kmip server
func (kc *kmipClient) CreateAsymmetricKeyPair(algorithm, curveType string, length int) (string, error) {
	defaultLog.Trace("kmipclient/kmipclient:CreateAsymmetricKeyPair() Entering")
	defer defaultLog.Trace("kmipclient/kmipclient:CreateAsymmetricKeyPair() Leaving")

	var createKeyPairRequestPayLoad interface{}
	if kc.KMIPVersion == constants.KMIP_2_0 {
		createKeyPairRequestPayLoad = models.CreateKeyPairRequestPayload{
			CommonAttributes: models.CommonAttributes{
				CryptographicAlgorithm: kmip14.CryptographicAlgorithmRSA,
				CryptographicLength:    int32(length),
			},
			PrivateKeyAttributes: models.PrivateKeyAttributes{
				CryptographicUsageMask: kmip14.CryptographicUsageMaskDecrypt,
			},
			PublicKeyAttributes: models.PublicKeyAttributes{
				CryptographicUsageMask: kmip14.CryptographicUsageMaskEncrypt,
			},
		}
	} else {
		createKeyPairRequestPayLoad = kmip.CreateKeyPairRequestPayload{
			CommonTemplateAttribute: &kmip.TemplateAttribute{
				Attribute: []kmip.Attribute{
					{
						AttributeName:  "Cryptographic Algorithm",
						AttributeValue: kmip14.CryptographicAlgorithmRSA,
					},
					{
						AttributeName:  "Cryptographic Length",
						AttributeValue: int32(length),
					},
				},
			},
			PrivateKeyTemplateAttribute: &kmip.TemplateAttribute{
				Attribute: []kmip.Attribute{
					{
						AttributeName:  "Cryptographic Usage Mask",
						AttributeValue: kmip14.CryptographicUsageMaskDecrypt,
					},
				},
			},
			PublicKeyTemplateAttribute: &kmip.TemplateAttribute{
				Attribute: []kmip.Attribute{
					{
						AttributeName:  "Cryptographic Usage Mask",
						AttributeValue: kmip14.CryptographicUsageMaskEncrypt,
					},
				},
			},
		}
	}

	batchItem, decoder, err := kc.SendRequest(createKeyPairRequestPayLoad, kmip14.OperationCreateKeyPair)
	if err != nil {
		return "", errors.Wrap(err, "failed to perform create keypair operation")
	}

	var respPayload models.CreateKeyPairResponsePayload
	err = decoder.DecodeValue(&respPayload, batchItem.ResponsePayload.(ttlv.TTLV))
	if err != nil {
		return "", errors.Wrap(err, "failed to decode create keypair response payload")
	}

	return respPayload.PrivateKeyUniqueIdentifier, nil
}

// GetKey retrieves a key from kmip server
func (kc *kmipClient) GetKey(keyID, algorithm string) ([]byte, error) {
	defaultLog.Trace("kmipclient/kmipclient:GetKey() Entering")
	defer defaultLog.Trace("kmipclient/kmipclient:GetKey() Leaving")

	getRequestPayLoad := models.GetRequestPayload{
		UniqueIdentifier: kmip20.UniqueIdentifierValue{
			Text: keyID,
		},
	}

	batchItem, decoder, err := kc.SendRequest(getRequestPayLoad, kmip14.OperationGet)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform get key operation")
	}

	var respPayload models.GetResponsePayload
	err = decoder.DecodeValue(&respPayload, batchItem.ResponsePayload.(ttlv.TTLV))
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode get key response payload")
	}

	var keyValue models.KeyValue

	switch algorithm {
	case constants.CRYPTOALG_AES:
		err = decoder.DecodeValue(&keyValue, respPayload.SymmetricKey.KeyBlock.KeyValue.(ttlv.TTLV))
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode symmetric keyblock")
		}
	case constants.CRYPTOALG_RSA:
		if respPayload.ObjectType == kmip14.ObjectTypePrivateKey {
			err = decoder.DecodeValue(&keyValue, respPayload.PrivateKey.KeyBlock.KeyValue.(ttlv.TTLV))
			if err != nil {
				return nil, errors.Wrap(err, "failed to decode private keyblock")
			}
		} else {
			return nil, errors.Errorf("unsupported object type %s", respPayload.ObjectType)
		}
	default:
		return nil, errors.Errorf("unsupported %s algorithm provided", algorithm)
	}

	return keyValue.KeyMaterial, nil
}

// DeleteKey deletes a key from kmip server
func (kc *kmipClient) DeleteKey(keyID string) error {
	defaultLog.Trace("kmipclient/kmipclient:DeleteKey() Entering")
	defer defaultLog.Trace("kmipclient/kmipclient:DeleteKey() Leaving")

	deleteRequestPayLoad := models.DeleteRequest{
		UniqueIdentifier: kmip20.UniqueIdentifierValue{
			Text: keyID,
		},
	}

	_, _, err := kc.SendRequest(deleteRequestPayLoad, kmip14.OperationDestroy)
	if err != nil {
		return errors.Wrap(err, "failed to perform delete key operation")
	}

	return nil
}

// SendRequest perform send request message to kmip server and receive response messages
func (kc *kmipClient) SendRequest(requestPayload interface{}, Operation kmip14.Operation) (*kmip.ResponseBatchItem, *ttlv.Decoder, error) {
	defaultLog.Trace("kmipclient/kmipclient:SendRequest() Entering")
	defer defaultLog.Trace("kmipclient/kmipclient:SendRequest() Leaving")

	conn, err := tls.Dial("tcp", kc.ServerIP+":"+kc.ServerPort, &kc.Config)
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	_, err = conn.ConnectionState().PeerCertificates[0].Verify(x509.VerifyOptions{Roots: kc.Config.RootCAs})
	if err != nil {
		return nil, nil, err
	}

	message := kmip.RequestMessage{
		RequestHeader: kc.requestHeader,
		BatchItem: []kmip.RequestBatchItem{
			{
				Operation:      Operation,
				RequestPayload: requestPayload,
			},
		},
	}

	requestMessage, err := ttlv.Marshal(message)
	if err != nil {
		return nil, nil, err
	}

	defaultLog.Debugf("kmipclient/kmipclient:SendRequest() Request Message for operation %s \n%s", Operation.String(), requestMessage)

	_, err = conn.Write(requestMessage)
	if err != nil {
		return nil, nil, err
	}

	decoder := ttlv.NewDecoder(bufio.NewReader(conn))
	response, err := decoder.NextTTLV()
	if err != nil {
		return nil, nil, err
	}

	var responseMessage kmip.ResponseMessage
	err = decoder.DecodeValue(&responseMessage, response)
	if err != nil {
		return nil, nil, err
	}

	responseTTLV, err := ttlv.Marshal(responseMessage)
	if err != nil {
		return nil, nil, err
	}

	defaultLog.Debugf("kmipclient/kmipclient:SendRequest() Response Message for operation %s \n%s", Operation.String(), responseTTLV)

	if responseMessage.BatchItem[0].ResultStatus != kmip14.ResultStatusSuccess {
		return nil, nil, errors.Errorf("request message is failed with reason %s", responseMessage.BatchItem[0].ResultMessage)
	}
	defaultLog.Infof("kmipclient/kmipclient:SendRequest() The KMIP operation %s was executed with no errors", Operation.String())

	return &responseMessage.BatchItem[0], decoder, nil
}
