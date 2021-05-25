/*
 *  Copyright (C) 2021 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package wpm

import (
	"fmt"
	"github.com/intel-secl/intel-secl/v4/pkg/wpm/version"
)

const helpStr = `
Usage:
    wpm <command> [arguments]

Available Commands:
    -h|--help                        Show this help message
    -v|--version                     Print version/build information
    create-image-flavor              Create VM image flavors and encrypt the image
    create-container-image-flavor    Create container image flavors and encrypt the container image
    get-container-image-id           Fetch the container image ID given the sha256 digest of the image
    unwrap-key                       Unwraps the image encryption key fetched from KBS
    fetch-key                        Fetches the image encryption key with associated tags from KBS
    uninstall [--purge]              Uninstall wpm. --purge option needs to be applied to remove configuration and data files
    setup                            Run workload-policy-manager setup tasks

Setup command usage:     wpm setup [task] [--force]

Available tasks for setup:
   all                                         Runs all setup tasks
                                               Required env variables:
                                                   - get required env variables from all the setup tasks
                                               Optional env variables:
                                                   - get optional env variables from all the setup tasks

   download-ca-cert                            Download CMS root CA certificate
                                               - Option [--force] overwrites any existing files, and always downloads new root CA cert
                                               Required env variables specific to setup task are:
                                                   - CMS_BASE_URL=<url>                              : for CMS API url
                                                   - CMS_TLS_CERT_SHA384=<CMS TLS cert sha384 hash>  : to ensure that WPM is talking to the right CMS instance

   download-cert-flavor-signing                Generates Key pair and CSR, gets it signed from CMS
                                               - Option [--force] overwrites any existing files, and always downloads newly signed WPM Flavor Signing cert
                                               Required env variables specific to setup task are:
                                                   - CMS_BASE_URL=<url>                       : for CMS API url
                                                   - BEARER_TOKEN=<token>                     : for authenticating with CMS
                                               Optional env variables specific to setup task are:
                                                   - FLAVOR_SIGNING_CERT_FILE    The file to which certificate is saved
                                                   - FLAVOR_SIGNING_KEY_FILE     The file to which private key is saved
                                                   - FLAVOR_SIGNING_COMMON_NAME  The common name of signed certificate

   create-envelope-key                           Creates the key pair required to securely transfer key from KBS
                                               - Option [--force] overwrites existing envelope key pairs
`

func (a *App) printUsage() {
	fmt.Fprintln(a.consoleWriter(), helpStr)
}

func (a *App) printVersion() {
	fmt.Fprintf(a.consoleWriter(), version.GetVersion())
}

func (a *App) printUsageWithError(err error) {
	fmt.Fprintln(a.errorWriter(), "Application returned with error:", err.Error())
	fmt.Fprintln(a.errorWriter(), helpStr)
}

func (a *App) printContainerFlavorUsage() {
	log.Trace("app:printContainerFlavorUsage() Entering")
	defer log.Trace("app:printContainerFlavorUsage() Leaving")

	fmt.Fprintf(a.consoleWriter(), "usage: wpm create-container-image-flavor -i img-name [-t tag] [-f dockerFile] [-d build-dir] [-k keyId]\n"+
		"                            [-e] [-s] [-n notaryServer] [-o out-file]\n"+
		"\t  -i, --img-name                  container image name\n"+
		"\t  -t, --tag                       (optional) container image tag name\n"+
		"\t  -f, --docker-file               (optional) container file path\n"+
		"\t                                  to build the container image\n"+
		"\t  -d, --build-dir                 (optional) build directory to build the\n"+
		"\t                                  container image. To be provided when container\n"+
		"\t                                  file path [-f] is provided as parameter\n"+
		"\t  -k, --key-id                    (optional) existing key ID\n"+
		"\t                                  if not specified, a new key is generated\n"+
		"\t  -e, --encryption-required       (optional) boolean parameter specifies if\n"+
		"\t                                  container image needs to be encrypted\n"+
		"\t  -s, --integrity-enforced        (optional) boolean parameter specifies if\n"+
		"\t                                  container image should be signed\n"+
		"\t  -n, --notary-server             (optional) specify notary server url\n"+
		"\t  -o, --out-file                  (optional) specify output file name")

}

// fetch-key command usage string
func (a *App) printFetchKeyUsage() {
	log.Trace("app:printFetchKeyUsage() Entering")
	defer log.Trace("app:printFetchKeyUsage() Leaving")

	fmt.Fprintf(a.consoleWriter(), "usage: wpm fetch-key [-k key]\n"+
		"\t  -k, --key       (optional) existing key ID\n"+
		"\t                  if not specified, a new key is generated\n"+
		"\t  -t, --asset-tag (optional) asset tags associated with the new key\n"+
		"\t                  tags are key:value separated by comma\n"+
		"\t  -a, --asymmetric (optional) specify to use asymmetric encryption\n"+
		"\t                  currently only supports RSA")
}

// unwrap-key command usage string
func (a *App) printUnwrapKeyUsage() {
	log.Trace("app:printUnwrapKeyUsage() Entering")
	defer log.Trace("app:printUnwrapKeyUsage() Leaving")

	fmt.Fprintf(a.consoleWriter(), "usage: unwrap-key [-i |--in] <wrapped key file path>")
}

// get-container-image-id command usage string
func (a *App) printGetContainerImageIdUsage() {
	log.Trace("app:printGetContainerImageIdUsage() Entering")
	defer log.Trace("app:printGetContainerImageIdUsage() Leaving")

	fmt.Fprintf(a.consoleWriter(), "usage: get-container-image-id [<sha256 digest of image>]")
}

// create-image-flavor command usage
func (a *App) printImageFlavorUsage() {
	log.Trace("main:imageFlavorUsage() Entering")
	defer log.Trace("main:imageFlavorUsage() Leaving")

	fmt.Fprintf(a.consoleWriter(), "usage: wpm create-image-flavor [-l label] [-i in] [-o out] [-e encout] [-k key]\n"+
		"\t  -l, --label     image flavor label\n"+
		"\t  -i, --in        input image file name\n"+
		"\t  -o, --out       (optional) output image flavor file name\n"+
		"\t                  if not specified, will print to the console\n"+
		"\t  -e, --encout    (optional) output encrypted image file name\n"+
		"\t                  if not specified, encryption is skipped\n"+
		"\t  -k, --key       (optional) existing key ID\n"+
		"\t                  if not specified, a new key is generated\n\n")
}
