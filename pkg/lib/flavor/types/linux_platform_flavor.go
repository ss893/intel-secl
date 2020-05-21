/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package types

import (
	"crypto/rsa"
	"encoding/xml"
	"fmt"
	cf "github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/common"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/constants"
	cm "github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/model"
	"github.com/intel-secl/intel-secl/v3/pkg/lib/flavor/util"
	hcTypes "github.com/intel-secl/intel-secl/v3/pkg/lib/host-connector/types"
	"github.com/intel-secl/intel-secl/v3/pkg/model/hvs"
	taModel "github.com/intel-secl/intel-secl/v3/pkg/model/ta"
	"github.com/pkg/errors"
	"strings"
)

/**
 *
 * @author mullas
 */

// LinuxPlatformFlavor is used to generate various Flavors for a Intel-based Linux host
type LinuxPlatformFlavor struct {
	HostManifest   *hcTypes.HostManifest        `json:"host_manifest"`
	HostInfo       *taModel.HostInfo            `json:"host_info"`
	TagCertificate *cm.X509AttributeCertificate `json:"tag_certificate"`
}

var (
	platformModules = []string{"LCP_DETAILS_HASH", "BIOSAC_REG_DATA", "OSSINITDATA_CAP_HASH", "STM_HASH",
		"MLE_HASH", "NV_INFO_HASH", "tb_policy", "CPU_SCRTM_STAT", "HASH_START", "SINIT_PUBKEY_HASH",
		"LCP_AUTHORITIES_HASH", "EVTYPE_KM_HASH", "EVTYPE_BPM_HASH", "EVTYPE_KM_INFO_HASH", "EVTYPE_BPM_INFO_HASH",
		"EVTYPE_BOOT_POL_HASH"}
	osModules         = []string{"vmlinuz"}
	hostUniqueModules = []string{"initrd", "LCP_CONTROL_HASH"}
)

var pfutil util.PlatformFlavorUtil
var sfutil util.SoftwareFlavorUtil

// NewRHELPlatformFlavor returns an instance of LinuxPlatformFlavor
func NewRHELPlatformFlavor(hostReport *hcTypes.HostManifest, tagCertificate *cm.X509AttributeCertificate) PlatformFlavor {
	return LinuxPlatformFlavor{
		HostManifest:   hostReport,
		HostInfo:       &hostReport.HostInfo,
		TagCertificate: tagCertificate,
	}
}

// GetFlavorPartRaw extracts the details of the flavor part requested by the
// caller from the host report used during the creation of the PlatformFlavor instance
func (rhelpf LinuxPlatformFlavor) GetFlavorPartRaw(name cf.FlavorPart) ([]string, error) {
	switch name {
	case cf.Platform:
		return rhelpf.getPlatformFlavor()
	case cf.Os:
		return rhelpf.getOsFlavor()
	case cf.AssetTag:
		return rhelpf.getAssetTagFlavor()
	case cf.HostUnique:
		return rhelpf.getHostUniqueFlavor()
	case cf.Software:
		return rhelpf.getDefaultSoftwareFlavor()
	}
	return nil, cf.UNKNOWN_FLAVOR_PART()
}

// GetFlavorPartNames retrieves the list of flavor parts that can be obtained using the GetFlavorPartRaw function
func (rhelpf LinuxPlatformFlavor) GetFlavorPartNames() ([]cf.FlavorPart, error) {
	flavorPartList := []cf.FlavorPart{cf.Platform, cf.Os, cf.HostUnique, cf.Software}

	// For each of the flavor parts, check what PCRs are required and if those required PCRs are present in the host report.
	for i := 0; i < len(flavorPartList); i++ {
		flavorPart := flavorPartList[i]
		pcrList := rhelpf.getPcrList(flavorPart)
		pcrExists := pfutil.PcrExists(rhelpf.HostManifest.PcrManifest, pcrList)
		if !pcrExists {
			// remove the non-existent FlavorPart from list
			flavorPartList = append(flavorPartList[:i], flavorPartList[i+1:]...)
		}
	}

	// Check if the AssetTag flavor part is present by checking if tagCertificate is present
	if rhelpf.TagCertificate != nil {
		flavorPartList = append(flavorPartList, cf.AssetTag)
	}
	return flavorPartList, nil
}

// GetPcrList Helper function to calculate the list of PCRs for the flavor part specified based
// on the version of the TPM hardware.
func (rhelpf LinuxPlatformFlavor) getPcrList(flavorPart cf.FlavorPart) []int {
	var pcrSet = make(map[int]bool)
	var pcrs []int
	var isTboot bool

	hostInfo := *rhelpf.HostInfo

	isTboot = hostInfo.TbootInstalled

	switch flavorPart {
	case cf.Platform:
		pcrSet[0] = true
		// check if CBNT is enabled
		if isCbntMeasureProfile(hostInfo.HardwareFeatures.CBNT) {
			pcrSet[7] = true
		}
		// check if SUEFI is enabled
		if hostInfo.HardwareFeatures.SUEFI != nil {
			if hostInfo.HardwareFeatures.SUEFI.Enabled {
				for _, pcrx := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
					pcrSet[pcrx] = true
				}
			}
		}

		// check if TBOOT is enabled
		if isTboot {
			for _, pcrx := range []int{17, 18} {
				pcrSet[pcrx] = true
			}
		}
	case cf.Os:
		// check if TBOOT is enabled
		if isTboot {
			pcrSet[17] = true
		}
	case cf.HostUnique:
		// check if TBOOT is enabled
		if isTboot {
			for _, pcrx := range []int{17, 18} {
				pcrSet[pcrx] = true
			}
		}
	case cf.Software:
		pcrSet[15] = true
	}

	// convert set back to list
	for k := range pcrSet {
		pcrs = append(pcrs, k)
	}
	return pcrs
}

func isCbntMeasureProfile(cbnt *taModel.CBNT) bool {
	if cbnt != nil {
		return cbnt.Enabled && cbnt.Meta.Profile == cf.BootGuardProfile5().Name
	}
	return false
}

// eventLogRequired Helper function to determine if the event log associated with the PCR
// should be included in the flavor for the specified flavor part
func (rhelpf LinuxPlatformFlavor) eventLogRequired(flavorPartName cf.FlavorPart) bool {
	// defaults to false
	var eventLogRequired bool

	switch flavorPartName {
	case cf.Platform:
		eventLogRequired = true
	case cf.Os:
		eventLogRequired = true
	case cf.HostUnique:
		eventLogRequired = true
	case cf.Software:
		eventLogRequired = true
	}
	return eventLogRequired
}

// getPlatformFlavor returns a json document having all the good known PCR values and
// corresponding event logs that can be used for evaluating the PLATFORM trust of a host
func (rhelpf LinuxPlatformFlavor) getPlatformFlavor() ([]string, error) {
	var errorMessage = "Error during creation of PLATFORM flavor"
	var platformFlavors []string
	var platformPcrs = rhelpf.getPcrList(cf.Platform)
	var includeEventLog = rhelpf.eventLogRequired(cf.Platform)
	var allPcrDetails = pfutil.GetPcrDetails(
		rhelpf.HostManifest.PcrManifest, platformPcrs, includeEventLog)
	var filteredPcrDetails = pfutil.IncludeModulesToEventLog(
		allPcrDetails, platformModules)

	newMeta, err := pfutil.GetMetaSectionDetails(rhelpf.HostInfo, rhelpf.TagCertificate, "", cf.Platform, "")
	if err != nil {
		err = errors.Wrap(err, errorMessage+" - failure in Meta section details")
		return nil, err
	}
	newBios := pfutil.GetBiosSectionDetails(rhelpf.HostInfo)
	if newBios == nil {
		err = fmt.Errorf(errorMessage + " - failure in Bios section details")
		return nil, err
	}
	newHW := pfutil.GetHardwareSectionDetails(rhelpf.HostInfo)
	if newHW == nil {
		err = fmt.Errorf(errorMessage + " - failure in Hardware section details")
		return nil, err
	}

	// Assemble the Platform Flavor
	fj, err := hvs.NewFlavorToJson(newMeta, newBios, newHW, filteredPcrDetails, nil, nil, errorMessage)
	if err != nil {
		return nil, errors.Wrap(err, errorMessage+" - JSON marshal failure")
	}
	// return JSON
	platformFlavors = append(platformFlavors, fj)
	return platformFlavors, nil
}

// getOsFlavor Returns a json document having all the good known PCR values and
// corresponding event logs that can be used for evaluating the OS Trust of a host
func (rhelpf LinuxPlatformFlavor) getOsFlavor() ([]string, error) {
	var errorMessage = "Error during creation of OS flavor"
	var err error
	var osFlavors []string
	var osPcrs = rhelpf.getPcrList(cf.Os)
	var includeEventLog = rhelpf.eventLogRequired(cf.Os)
	var allPcrDetails = pfutil.GetPcrDetails(
		rhelpf.HostManifest.PcrManifest, osPcrs, includeEventLog)
	var filteredPcrDetails = pfutil.IncludeModulesToEventLog(
		allPcrDetails, osModules)

	newMeta, err := pfutil.GetMetaSectionDetails(rhelpf.HostInfo, rhelpf.TagCertificate, "", cf.Os, "")
	if err != nil {
		err = errors.Wrap(err, errorMessage+" Failure in Meta section details")
		return nil, err
	}
	newBios := pfutil.GetBiosSectionDetails(rhelpf.HostInfo)
	if newBios == nil {
		err = fmt.Errorf("%s Failure in Bios section details", errorMessage)
		return nil, err
	}

	// Assemble the OS Flavor
	fj, err := hvs.NewFlavorToJson(newMeta, newBios, nil, filteredPcrDetails, nil, nil, errorMessage)
	if err != nil {
		return nil, err
	}
	// return JSON
	osFlavors = append(osFlavors, fj)
	return osFlavors, nil
}

// getHostUniqueFlavor Returns a json document having all the good known PCR values and corresponding event logs that
// can be used for evaluating the unique part of the PCR configurations of a host. These include PCRs/modules getting
// extended to PCRs that would vary from host to host.
func (rhelpf LinuxPlatformFlavor) getHostUniqueFlavor() ([]string, error) {
	var errorMessage = "Error during creation of HOST_UNIQUE flavor"
	var err error
	var hostUniqueFlavors []string
	var hostUniquePcrs = rhelpf.getPcrList(cf.HostUnique)
	var includeEventLog = rhelpf.eventLogRequired(cf.HostUnique)
	var allPcrDetails = pfutil.GetPcrDetails(
		rhelpf.HostManifest.PcrManifest, hostUniquePcrs, includeEventLog)
	var filteredPcrDetails = pfutil.IncludeModulesToEventLog(
		allPcrDetails, hostUniqueModules)

	newMeta, err := pfutil.GetMetaSectionDetails(rhelpf.HostInfo, rhelpf.TagCertificate, "", cf.HostUnique, "")
	if err != nil {
		err = errors.Wrap(err, errorMessage+" Failure in Meta section details")
		return nil, err
	}
	newBios := pfutil.GetBiosSectionDetails(rhelpf.HostInfo)
	if newBios == nil {
		err = errors.Wrap(err, errorMessage+" Failure in Bios section details")
		return nil, err
	}

	// Assemble the Host Unique Flavor
	fj, err := hvs.NewFlavorToJson(newMeta, newBios, nil, filteredPcrDetails, nil, nil, errorMessage)
	if err != nil {
		return nil, err
	}
	// return JSON
	hostUniqueFlavors = append(hostUniqueFlavors, fj)
	return hostUniqueFlavors, nil
}

// getAssetTagFlavor Retrieves the asset tag part of the flavor including the certificate and all the key-value pairs
// that are part of the certificate.
func (rhelpf LinuxPlatformFlavor) getAssetTagFlavor() ([]string, error) {
	var errorMessage = "Error during creation of ASSET_TAG flavor"
	var err error
	var assetTagFlavors []string
	if rhelpf.TagCertificate == nil {
		return nil, fmt.Errorf("%s - %s", errorMessage, cf.FLAVOR_PART_CANNOT_BE_SUPPORTED().Message)
	}

	// create meta section details
	newMeta, err := pfutil.GetMetaSectionDetails(rhelpf.HostInfo, rhelpf.TagCertificate, "", cf.AssetTag, "")
	if err != nil {
		err = errors.Wrap(err, errorMessage+" Failure in Meta section details")
		return nil, err
	}
	// create bios section details
	newBios := pfutil.GetBiosSectionDetails(rhelpf.HostInfo)
	if newBios == nil {
		err = fmt.Errorf("%s - Failure in Bios section details", errorMessage)
		return nil, err
	}
	// create external section details
	newExt, err := pfutil.GetExternalConfigurationDetails(rhelpf.TagCertificate)
	if err != nil {
		err = errors.Wrap(err, errorMessage+" Failure in External configuration section details")
		return nil, err
	}

	// Assemble the Asset Tag Flavor
	fj, err := hvs.NewFlavorToJson(newMeta, newBios, nil, nil, newExt, nil, errorMessage)
	if err != nil {
		return nil, err
	}
	// return JSON
	assetTagFlavors = append(assetTagFlavors, fj)
	return assetTagFlavors, nil
}

// getDefaultSoftwareFlavor Method to create a software flavor. This method would create a software flavor that would
// include all the measurements provided from host.
func (rhelpf LinuxPlatformFlavor) getDefaultSoftwareFlavor() ([]string, error) {
	var softwareFlavors []string
	var errorMessage = cf.SOFTWARE_FLAVOR_CANNOT_BE_CREATED().Message

	if rhelpf.HostManifest != nil && rhelpf.HostManifest.MeasurementXmls != nil {
		measurementXmls, err := rhelpf.getDefaultMeasurement()
		if err != nil {
			return nil, errors.Wrapf(err, errorMessage)
		}

		for _, measurementXml := range measurementXmls {
			var softwareFlavor = NewSoftwareFlavor(measurementXml)
			swFlavorStr, err := softwareFlavor.GetSoftwareFlavor()
			if err != nil {
				return nil, err
			}
			softwareFlavors = append(softwareFlavors, swFlavorStr)
		}
	}
	return softwareFlavors, nil
}

// getDefaultMeasurement returns a default set of measurements for the Platform Flavor
func (rhelpf LinuxPlatformFlavor) getDefaultMeasurement() ([]string, error) {
	var measurementXmlCollection []string
	var err error

	for _, measurementXML := range rhelpf.HostManifest.MeasurementXmls {
		var measurement taModel.Measurement
		err = xml.Unmarshal([]byte(measurementXML), &measurement)
		if err != nil {
			err = errors.Wrapf(err, "Error unmarshalling measurement XML: %s", err.Error())
			return nil, err
		}
		if strings.Contains(measurement.Label, constants.DefaultSoftwareFlavorPrefix) ||
			strings.Contains(measurement.Label, constants.DefaultWorkloadFlavorPrefix) {
			measurementXmlCollection = append(measurementXmlCollection, measurementXML)
		}
	}
	return measurementXmlCollection, nil
}

// GetFlavorPart extracts the details of the flavor part requested by the caller from
// the host report used during the creation of the PlatformFlavor instance and it's corresponding signature.
func (rhelpf LinuxPlatformFlavor) GetFlavorPart(part cf.FlavorPart, flavorSigningPrivateKey *rsa.PrivateKey) ([]hvs.SignedFlavor, error) {
	var flavors []string
	var err error

	// validate private key
	if flavorSigningPrivateKey != nil {
		err := flavorSigningPrivateKey.Validate()
		if err != nil {
			return nil, errors.Wrap(err, "signing key validation failed")
		}
	}

	// get flavor
	flavors, err = rhelpf.GetFlavorPartRaw(part)
	if err != nil {
		return nil, err
	}

	sfList, err := pfutil.GetSignedFlavorList(flavors, flavorSigningPrivateKey)
	if err != nil {
		return []hvs.SignedFlavor{}, err
	}
	return *sfList, nil
}