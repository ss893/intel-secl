/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package openstackplugin

import (
	"bytes"
	"encoding/json"
	types "github.com/intel-secl/intel-secl/v4/pkg/ihub/model"
	"github.com/intel-secl/intel-secl/v4/pkg/ihub/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	openstackClient "github.com/intel-secl/intel-secl/v4/pkg/clients/openstack"
	vsPlugin "github.com/intel-secl/intel-secl/v4/pkg/ihub/attestationPlugin"
	"github.com/intel-secl/intel-secl/v4/pkg/ihub/config"
	"github.com/intel-secl/intel-secl/v4/pkg/ihub/constants"
	commonLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/saml"
	model "github.com/intel-secl/intel-secl/v4/pkg/model/openstack"
	"github.com/pkg/errors"
)

// openstackHostDetails for Openstack
type openstackHostDetails struct {
	types.HostDetails
	ResourceProviderGeneration int
	DefaultTraits              []string
	CustomTraits               []string
	CustomTeeTraits            []string
	// if the host is running Trusted Execution Environment (SGX, TDX...)
	Tee bool
}

// OpenstackDetails for requesting auth and getting host list and updating traits
type OpenstackDetails struct {
	Config             *config.Configuration
	HostDetails        []openstackHostDetails
	AllCustomTraits    []string
	OpenstackClient    *openstackClient.Client
	TrustedCAsStoreDir string
	SamlCertFilePath   string
}

var (
	log            = commonLog.GetDefaultLogger()
	osRegexEpcSize = regexp.MustCompile(constants.RegexEpcSize)
	rgx            = regexp.MustCompile(constants.RegexNonStandardChar)
)

// getHostsFromOpenstack Get Hosts from Openstack
func getHostsFromOpenstack(openstackDetails *OpenstackDetails) error {
	log.Trace("openstackplugin/openstack_plugin:getHostsFromOpenstack() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:getHostsFromOpenstack() Leaving")

	prefixURL := openstackDetails.Config.Endpoint.URL
	resourcePath := "resource_providers"

	parsedUrl, err := url.ParseRequestURI(prefixURL + resourcePath)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getHostsFromOpenstack()  Unable to parse the resource path url")
	}

	log.Debug("openstackplugin/openstack_plugin:getHostsFromOpenstack() Sending request to Openstack client to get hosts")
	res, err := openstackDetails.OpenstackClient.SendRequest(&openstackClient.RequestParams{
		Method: "GET",
		URL:    parsedUrl,
		Body:   nil,
	})
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getHostsFromOpenstack()  Error in getting the list of hosts from Openstack")
	}
	defer func() {
		derr := res.Body.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing response")
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getHostsFromOpenstack()  Error in reading the host details body")
	}

	var openStackResources model.OpenstackResources

	log.Debug("openstackplugin/openstack_plugin:getHostsFromOpenstack() Unmarshalling the Openstack resources")
	err = json.Unmarshal(body, &openStackResources)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getHostsFromOpenstack()  Error in unmarshalling the host details from Openstack")
	}

	var hostDetailsList []openstackHostDetails

	log.Debug("openstackplugin/openstack_plugin:getHostsFromOpenstack() getting host details list from resource providers")
	for _, actualObject := range openStackResources.ResourceProviders {

		hostDetails := openstackHostDetails{}
		hostDetails.HostID = actualObject.HostID
		hostDetails.HostName = actualObject.Name

		hostDetailsList = append(hostDetailsList, hostDetails)
		log.Debug("openstackplugin/openstack_plugin:getHostsFromOpenstack() Host ID : ", actualObject.HostID)

	}
	openstackDetails.HostDetails = hostDetailsList
	log.Info("openstackplugin/openstack_plugin:getHostsFromOpenstack() Retrieved the host details from Openstack")
	return nil
}

// filterHostReportsForOpenstack Get Host Reports
func filterHostReportsForOpenstack(hostDetails *openstackHostDetails, openstackDetails *OpenstackDetails) error {

	log.Trace("openstackplugin/openstack_plugin:filterHostReportsForOpenstack() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:filterHostReportsForOpenstack() Leaving")

	log.Info("openstackplugin/openstack_plugin:filterHostReportsForOpenstack() Get the host reports for Openstack")

	// split based on whether the host uses SGX/ISECL HVS
	if openstackDetails.Config.AttestationService.HVSBaseURL != "" {
		samlReport, err := vsPlugin.GetHostReports(hostDetails.HostName, openstackDetails.Config, openstackDetails.TrustedCAsStoreDir, openstackDetails.SamlCertFilePath)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:filterHostReportsForOpenstack() : Error in getting the host report")
		}
		err = getCustomTraitsFromSAMLReport(hostDetails, samlReport)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:filterHostReportsForOpenstack() : Error in generating custom traits from trust report")
		}
	}
	if openstackDetails.Config.AttestationService.SHVSBaseURL != "" {
		platformData, err := vsPlugin.GetHostPlatformData(hostDetails.HostName, openstackDetails.Config, constants.TrustedCAsStoreDir)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:filterHostReportsForOpenstack() : Error in getting the SGX Platform Data")
		}

		var sgxData types.PlatformDataSGX

		err = json.Unmarshal(platformData, &sgxData)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:filterHostReportsForOpenstack() : unmarshal SGX platform data error")
		}
		if len(sgxData) == 1 {
			// need to validate contents of EpcSize
			if !osRegexEpcSize.MatchString(sgxData[0].EpcSize) {
				log.Errorf("openstackplugin/openstack_plugin:filterHostReportsForOpenstack() Invalid EPC Size value")
				hostDetails.EpcSize = constants.SgxTraitEpcSizeNotAvailable
			} else {
				hostDetails.EpcSize = sgxData[0].EpcSize
			}
			hostDetails.FlcEnabled = sgxData[0].FlcEnabled
			hostDetails.SgxEnabled = sgxData[0].SgxEnabled
			hostDetails.SgxSupported = sgxData[0].SgxSupported
			hostDetails.TcbUpToDate = sgxData[0].TcbUpToDate
			util.EvaluateValidTo(sgxData[0].ValidTo, openstackDetails.Config.PollIntervalMinutes)
			hostDetails.ValidTo = sgxData[0].ValidTo
		} else {
			return errors.Errorf("openstackplugin/openstack_plugin:filterHostReportsForOpenstack() : SGX Platform Data response has invalid length %d", len(sgxData))
		}

		err = getCustomTraitsFromPlatformData(hostDetails)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:filterHostReportsForOpenstack() : Error in generating custom traits from SGX platform data")
		}
	}

	log.Info("openstackplugin/openstack_plugin:filterHostReportsForOpenstack() Get the custom traits from report for Openstack")

	return nil
}

// getCustomTraitsFromSAMLReport pulls custom traits from the HVS SAML report
func getCustomTraitsFromSAMLReport(hostDetails *openstackHostDetails, samlReport *saml.Saml) error {

	log.Trace("openstackplugin/openstack_plugin:getCustomTraitsFromSAMLReport() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:getCustomTraitsFromSAMLReport() Leaving")

	var customTraits []string
	trusted := false

	log.Debug("openstackplugin/openstack_plugin:getCustomTraitsFromSAMLReport() Getting traits from the report")
	for _, as := range samlReport.Attribute {

		key := as.Name
		value := as.AttributeValue
		// Asset Tags
		if strings.HasPrefix(key, "TAG") {
			log.Debugf("openstackplugin/openstack_plugin:getCustomTraitsFromSAMLReport() Constructing custom trait for Asset tag: %s - %s", key, value)
			prefix := constants.IseclTraitPrefix + constants.TraitAssetTagPrefix
			trait := getFormattedCustomTraits(prefix, key, value)
			customTraits = append(customTraits, trait)

		} else if strings.EqualFold(as.Name, "TRUST_OVERALL") && strings.EqualFold(value, "true") { //Trust tags
			log.Debug("openstackplugin/openstack_plugin:getCustomTraitsFromSAMLReport() Constructing custom trait for trust tag")
			customTraits = append(customTraits, constants.TrustedTrait)
			trusted = true
		} else if strings.HasPrefix(as.Name, "FEATURE") && !strings.EqualFold(value, "false") { //HWFeature tags
			log.Debugf("openstackplugin/openstack_plugin:getCustomTraitsFromSAMLReport() Constructing custom trait for HWFeature tag: %s", key)
			prefix := constants.IseclTraitPrefix + constants.TraitHardwareFeaturesPrefix
			trait := getFormattedCustomTraits(prefix, key, "")
			customTraits = append(customTraits, trait)
		}
	}

	if trusted {
		hostDetails.CustomTraits = customTraits
		log.Debugf("Traits for host with name %s: %v", hostDetails.HostName, customTraits)
	} else {
		log.Warnf("Host with name %s is not trusted, removing all existing custom tags", hostDetails.HostName)
	}

	hostDetails.Trusted = trusted
	return nil
}

// getCustomTraitsFromPlatformData sets the custom traits per SGX-HVS Platform Data
func getCustomTraitsFromPlatformData(hostDetails *openstackHostDetails) error {
	log.Trace("openstackplugin/openstack_plugin:getCustomTraitsFromPlatformData() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:getCustomTraitsFromPlatformData() Leaving")

	var traitSet []string

	log.Debug("openstackplugin/openstack_plugin:getCustomTraitsFromPlatformData() Getting traits from the SGX PlatformData")
	traitSet = append(traitSet, getFormattedCustomTraits(constants.IseclTraitPrefix+constants.TraitDelimiter, constants.SgxTraitEnabled, strconv.FormatBool(hostDetails.SgxEnabled)))
	traitSet = append(traitSet, getFormattedCustomTraits(constants.IseclTraitPrefix+constants.TraitDelimiter, constants.SgxTraitSupported, strconv.FormatBool(hostDetails.SgxSupported)))
	traitSet = append(traitSet, getFormattedCustomTraits(constants.IseclTraitPrefix+constants.TraitDelimiter, constants.SgxTraitTcbUpToDate, strconv.FormatBool(hostDetails.TcbUpToDate)))
	traitSet = append(traitSet, getFormattedCustomTraits(constants.IseclTraitPrefix+constants.TraitDelimiter, constants.SgxTraitFlcEnabled, strconv.FormatBool(hostDetails.FlcEnabled)))

	if hostDetails.EpcSize != "" {
		traitSet = append(traitSet, getFormattedCustomTraits(constants.IseclTraitPrefix+constants.TraitDelimiter, constants.SgxTraitEpcSize, hostDetails.EpcSize))
	} else {
		traitSet = append(traitSet, getFormattedCustomTraits(constants.IseclTraitPrefix+constants.TraitDelimiter, constants.SgxTraitEpcSize, constants.SgxTraitEpcSizeNotAvailable))
	}

	// persist custom traits
	hostDetails.CustomTeeTraits = traitSet
	// this is a  TEE enabled system
	hostDetails.Tee = true

	return nil
}

// getFormattedCustomTraits Format the custom Traits with the "CUSTOM_ISECL" prefix
func getFormattedCustomTraits(prefix string, tagKey string, tagValue string) string {

	log.Trace("openstackplugin/openstack_plugin:getFormattedCustomTraits() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:getFormattedCustomTraits() Leaving")

	log.Debug("openstackplugin/openstack_plugin:getFormattedCustomTraits() getting the formatted custom traits")
	delimiter := constants.TraitDelimiter

	newTagKey := rgx.ReplaceAllString(tagKey, delimiter)
	newTagValue := rgx.ReplaceAllString(tagValue, delimiter)

	log.Debug("openstackplugin/openstack_plugin:getFormattedCustomTraits() Add the CUSTOM_ISECL prefix to formatted string")
	formattedString := prefix + strings.ToUpper(newTagKey)

	if newTagValue != "" {
		formattedString = formattedString + delimiter + strings.ToUpper(newTagValue)
	}

	log.Info("openstackplugin/openstack_plugin:getFormattedCustomTraits() Custom traits are formatted")
	return formattedString
}

// updateOpenstackTraits Update the traits for the resources
func updateOpenstackTraits(openstackDetails *OpenstackDetails) error {

	log.Trace("openstackplugin/openstack_plugin:updateOpenstackTraits() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:updateOpenstackTraits() Leaving")

	for index := range openstackDetails.HostDetails {

		log.Debug("openstackplugin/openstack_plugin:updateOpenstackTraits() fetching all the traits for the resource")
		err := getTraitsForResource(&openstackDetails.HostDetails[index], openstackDetails)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:updateOpenstackTraits() Error in getting Traits for the resource")
		}

		if len(openstackDetails.HostDetails[index].CustomTraits) > 0 {

			log.Debug("openstackplugin/openstack_plugin:updateOpenstackTraits() creating custom traits")
			err := createCustomTraits(openstackDetails.HostDetails[index].CustomTraits, openstackDetails)
			if err != nil {
				return errors.Wrap(err, "openstackplugin/openstack_plugin:updateOpenstackTraits() Error in creating custom traits")
			}

		}

		if len(openstackDetails.HostDetails[index].CustomTeeTraits) > 0 {

			log.Debug("openstackplugin/openstack_plugin:updateOpenstackTraits() creating custom TEE traits")
			err := createCustomTraits(openstackDetails.HostDetails[index].CustomTeeTraits, openstackDetails)
			if err != nil {
				return errors.Wrap(err, "openstackplugin/openstack_plugin:updateOpenstackTraits() Error in creating custom TEE traits")
			}

		}

		log.Debug("openstackplugin/openstack_plugin:updateOpenstackTraits() Associating traits to resource")
		err = associateTraitsForResource(&openstackDetails.HostDetails[index], openstackDetails)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:updateOpenstackTraits() Error in Associating custom traits")
		}

	}

	log.Debug("openstackplugin/openstack_plugin:updateOpenstackTraits() Fetch All the custom traits")
	err := getAllCustomTraits(openstackDetails)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:updateOpenstackTraits() Error in Fetching all the custom traits for cleanup")
	}

	log.Debug("openstackplugin/openstack_plugin:updateOpenstackTraits() Delete All the Non-Associated Traits")
	err = deleteNonAssociatedTraits(openstackDetails)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:updateOpenstackTraits() Error in Deleting all the non-associated traits for cleanup")
	}

	log.Info("openstackplugin/openstack_plugin:updateOpenstackTraits() Custom traits are updated onto Openstack")

	return nil
}

// getTraitsForResource Get traits for the Openstack Resources
func getTraitsForResource(hostDetails *openstackHostDetails, openstackDetails *OpenstackDetails) error {

	log.Trace("openstackplugin/openstack_plugin:getTraitsForResource() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:getTraitsForResource() Leaving")

	log.Debug("openstackplugin/openstack_plugin:getTraitsForResource() Getting traits for the Openstack resources")
	prefixURL := openstackDetails.Config.Endpoint.URL
	resourceTraitsPath := "resource_providers/" + hostDetails.HostID.String() + "/traits"
	urlPath := prefixURL + resourceTraitsPath
	log.Debug("openstackplugin/openstack_plugin:getTraitsForResource() URL For Resource Traits : " + urlPath)
	parsedUrl, err := url.Parse(urlPath)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getTraitsForResource() Unable to parse the resource traits path url")
	}

	log.Debug("openstackplugin/openstack_plugin:getTraitsForResource() Sending request to get traits")
	res, err := openstackDetails.OpenstackClient.SendRequest(&openstackClient.RequestParams{
		Method: "GET",
		URL:    parsedUrl,
		Body:   nil,
	})

	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getTraitsForResource() : Error in getting traits for the Resources")
	}
	defer func() {
		derr := res.Body.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing response")
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getTraitsForResource() Error in reading response while getting traits for the Resources")
	}

	var openStackTrait model.OpenStackTrait

	log.Debug("openstackplugin/openstack_plugin:getTraitsForResource() Unmarshalling the Openstack traits body")
	err = json.Unmarshal(body, &openStackTrait)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getTraitsForResource(): Error in unmarshalling the Openstack custom traits body")
	}

	for _, trait := range openStackTrait.Traits {
		if !strings.HasPrefix(trait, constants.IseclTraitPrefix) {
			hostDetails.DefaultTraits = append(hostDetails.DefaultTraits, trait)
		}

	}

	hostDetails.ResourceProviderGeneration = openStackTrait.ResourceProviderGeneration
	log.Info("openstackplugin/openstack_plugin:getTraitsForResource() Retrieved the Openstack traits for the resource")

	return nil
}

// createCustomTraits Create Custom Traits
func createCustomTraits(traits []string, openstackDetails *OpenstackDetails) error {

	log.Trace("openstackplugin/openstack_plugin:createCustomTraits() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:createCustomTraits() Leaving")

	for _, trait := range traits {

		prefixURL := openstackDetails.Config.Endpoint.URL
		createTraitsPath := "traits/" + trait
		urlPath := prefixURL + createTraitsPath
		log.Debug("openstackplugin/openstack_plugin:createCustomTraits() Trait Creation URL : " + urlPath)
		parsedUrl, err := url.Parse(urlPath)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:createCustomTraits()  Unable to parse the traits path url")
		}

		log.Debug("openstackplugin/openstack_plugin:createCustomTraits() Sending request to create traits")
		res, err := openstackDetails.OpenstackClient.SendRequest(&openstackClient.RequestParams{
			Method: "PUT",
			URL:    parsedUrl,
			Body:   nil,
		})
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:createCustomTraits() Error in creation of custom traits")
		}

		log.Debug("openstackplugin/openstack_plugin:createCustomTraits() checking the response status code for creating traits")
		if res.StatusCode == http.StatusCreated {
			log.Debug("openstackplugin/openstack_plugin:createCustomTraits() Trait created :" + trait)
		} else if res.StatusCode == http.StatusNoContent {
			log.Debug("openstackplugin/openstack_plugin:createCustomTraits() Trait does not Exist :" + trait)
		} else {
			log.Debug("openstackplugin/openstack_plugin:createCustomTraits() Unable to create Trait :" + trait)
		}
	}

	log.Info("openstackplugin/openstack_plugin:createCustomTraits() Custom traits are created for Openstack")
	return nil
}

// associateTraitsForResource Associate Traits for resource
func associateTraitsForResource(hostDetails *openstackHostDetails, openstackDetails *OpenstackDetails) error {

	log.Trace("openstackplugin/openstack_plugin:associateTraitsForResource() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:associateTraitsForResource() Leaving")

	prefixURL := openstackDetails.Config.Endpoint.URL
	resourceTraitsPath := "resource_providers/" + hostDetails.HostID.String() + "/traits"
	urlPath := prefixURL + resourceTraitsPath
	var openStackTrait model.OpenStackTrait
	openStackTrait.ResourceProviderGeneration = hostDetails.ResourceProviderGeneration

	log.Debug("openstackplugin/openstack_plugin:associateTraitsForResource() Appending the default and custom traits for the resource")

	openStackTrait.Traits = hostDetails.DefaultTraits

	// host is trusted, push custom platform traits
	if hostDetails.Trusted {
		openStackTrait.Traits = append(openStackTrait.Traits, hostDetails.CustomTraits...)
	}
	// host is a TEE, push custom TEE traits
	if hostDetails.Tee {
		openStackTrait.Traits = append(openStackTrait.Traits, hostDetails.CustomTeeTraits...)
	}

	log.Debug("openstackplugin/openstack_plugin:associateTraitsForResource() Associate Trait URL :  " + urlPath)
	log.Debug("openstackplugin/openstack_plugin:associateTraitsForResource() Resource Provider generation", openStackTrait.ResourceProviderGeneration)
	log.Debug("openstackplugin/openstack_plugin:associateTraitsForResource() OpenStack traits", openStackTrait.Traits)

	log.Debug("openstackplugin/openstack_plugin:associateTraitsForResource() Marshalling the openstack traits into json body")
	jsonBody, err := json.Marshal(openStackTrait)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:associateTraitsForResource() Error in marshalling traits for the resource")
	}
	payload := bytes.NewReader(jsonBody)

	parsedUrl, err := url.Parse(urlPath)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:associateTraitsForResource()  Unable to parse the resource traits path url")
	}

	res, err := openstackDetails.OpenstackClient.SendRequest(&openstackClient.RequestParams{
		Method:            "PUT",
		URL:               parsedUrl,
		Body:              payload,
		AdditionalHeaders: map[string]string{"Content-Type": "application/json"},
	})
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:associateTraitsForResource() Error in getting response")
	}

	log.Debug("openstackplugin/openstack_plugin:associateTraitsForResource() checking the response status code for associating traits to resource")
	if res.StatusCode != http.StatusOK {
		return errors.New("openstackplugin/openstack_plugin:associateTraitsForResource() Error in associating traits for the resource :" + string(res.StatusCode))
	}

	log.Info("openstackplugin/openstack_plugin:associateTraitsForResource() Traits are associated to the resource")
	return nil
}

// getAllCustomTraits Get all the custom traits in Openstack
func getAllCustomTraits(openstackDetails *OpenstackDetails) error {
	log.Trace("openstackplugin/openstack_plugin:getAllCustomTraits() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:getAllCustomTraits() Leaving")

	prefixURL := openstackDetails.Config.Endpoint.URL
	resourceTraitsPath := "traits"
	urlPath := prefixURL + resourceTraitsPath
	log.Debug("openstackplugin/openstack_plugin:getAllCustomTraits() The URL for Getting all the traits : " + urlPath)

	parsedUrl, err := url.Parse(urlPath)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getAllCustomTraits()  Unable to parse the resource traits path url")
	}
	log.Debug("openstackplugin/openstack_plugin:getAllCustomTraits() Sending request to get all custom traits")
	res, err := openstackDetails.OpenstackClient.SendRequest(&openstackClient.RequestParams{
		Method: "GET",
		URL:    parsedUrl,
		Body:   nil,
	})
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getAllCustomTraits() : Error in retrieving all the custom traits in Openstack")
	}

	defer func() {
		derr := res.Body.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing response")
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getAllCustomTraits() : Error in reading the Openstack custom traits")
	}

	var openStackTrait model.OpenStackTrait
	log.Debug("openstackplugin/openstack_plugin:getAllCustomTraits() unmarshalling the openstack traits")
	err = json.Unmarshal(body, &openStackTrait)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:getAllCustomTraits() : Error in unmarshalling the Openstack custom traits body")
	}
	for _, trait := range openStackTrait.Traits {

		if strings.HasPrefix(trait, constants.IseclTraitPrefix) {
			openstackDetails.AllCustomTraits = append(openstackDetails.AllCustomTraits, trait)
		}

	}
	log.Debug("openstackplugin/openstack_plugin:getAllCustomTraits() The custom traits are :", openstackDetails.AllCustomTraits)

	log.Info("openstackplugin/openstack_plugin:getAllCustomTraits()  Custom traits are received from Openstack endpoint")

	return nil
}

// deleteNonAssociatedTraits Delete all non associated CustomTraits in Openstack
func deleteNonAssociatedTraits(openstackDetails *OpenstackDetails) error {

	log.Trace("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() Leaving")

	for _, trait := range openstackDetails.AllCustomTraits {

		prefixURL := openstackDetails.Config.Endpoint.URL
		resourceTraitsPath := "traits/" + trait
		urlPath := prefixURL + resourceTraitsPath
		log.Debug("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() The Url for deleting all the traits are : " + urlPath)

		parsedUrl, err := url.Parse(urlPath)
		if err != nil {
			return errors.Wrap(err, "openstackplugin/openstack_plugin:deleteNonAssociatedTraits()  Unable to parse the resource traits path url")
		}
		log.Debug("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() Sending request to delete non associated traits")
		res, err := openstackDetails.OpenstackClient.SendRequest(&openstackClient.RequestParams{
			Method:            "DELETE",
			URL:               parsedUrl,
			Body:              nil,
			AdditionalHeaders: map[string]string{"Content-Type": "application/json"},
		})
		if err != nil {
			derr := res.Body.Close()
			if derr != nil {
				log.WithError(derr).Error("Error closing response")
			}
			return errors.Wrap(err, "openstackplugin/openstack_plugin:deleteNonAssociatedTraits() : Error in deleting the non associated traits")
		}
		log.Debug("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() checking the delete status of non associated traits")
		if res.StatusCode == http.StatusConflict {
			log.Debug("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() The trait " + trait + " is in use by a resource provider. Hence, Skipping the delete.")
		} else {
			log.Debug("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() The trait " + trait + " is Deleted Successfully")
		}
		derr := res.Body.Close()
		if derr != nil {
			log.WithError(derr).Error("Error closing response")
		}
	}
	log.Info("openstackplugin/openstack_plugin:deleteNonAssociatedTraits() Non associated traits are deleted from Openstack endpoint")
	return nil
}

// SendDataToEndPoint pushes host trust data to OpenStack
func SendDataToEndPoint(openstack OpenstackDetails) error {
	log.Trace("openstackplugin/openstack_plugin:SendDataToEndPoint() Entering")
	defer log.Trace("openstackplugin/openstack_plugin:SendDataToEndPoint() Leaving")

	log.Debug("openstackplugin/openstack_plugin:SendDataToEndPoint() Fetching Hosts from Openstack")
	err := getHostsFromOpenstack(&openstack)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:SendDataToEndPoint() Error in getting Hosts from Openstack")
	}

	log.Debug("openstackplugin/openstack_plugin:SendDataToEndPoint() Filtering Hosts from Openstack")
	for index := range openstack.HostDetails {
		err := filterHostReportsForOpenstack(&openstack.HostDetails[index], &openstack)
		if err != nil {
			log.WithError(err).Errorf("openstackplugin/openstack_plugin:SendDataToEndPoint() Error in Filtering"+
				" Host details for Openstack host %s", openstack.HostDetails[index].HostID.String())
		}
	}

	log.Info("openstackplugin/openstack_plugin:SendDataToEndPoint() Updating traits to Openstack for host : ", openstack.HostDetails)
	err = updateOpenstackTraits(&openstack)
	if err != nil {
		return errors.Wrap(err, "openstackplugin/openstack_plugin:SendDataToEndPoint() Error in Filtering Host details for Openstack")
	}

	return nil
}
