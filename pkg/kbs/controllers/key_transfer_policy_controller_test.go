/*
 * Copyright (C) 2020 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/controllers"
	"github.com/intel-secl/intel-secl/v4/pkg/kbs/domain/mocks"
	kbsRoutes "github.com/intel-secl/intel-secl/v4/pkg/kbs/router"
	consts "github.com/intel-secl/intel-secl/v4/pkg/lib/common/constants"
	"github.com/intel-secl/intel-secl/v4/pkg/model/kbs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyTransferPolicyController", func() {
	var router *mux.Router
	var w *httptest.ResponseRecorder
	var keyStore *mocks.MockKeyStore
	var policyStore *mocks.MockKeyTransferPolicyStore
	var keyTransferPolicyController *controllers.KeyTransferPolicyController
	BeforeEach(func() {
		router = mux.NewRouter()
		keyStore = mocks.NewFakeKeyStore()
		policyStore = mocks.NewFakeKeyTransferPolicyStore()

		keyTransferPolicyController = controllers.NewKeyTransferPolicyController(policyStore, keyStore)
	})

	// Specs for HTTP Post to "/key-transfer-policies"
	Describe("Create a new Key Transfer Policy", func() {
		Context("Provide a valid Create request", func() {
			It("Should create a new Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies", kbsRoutes.ErrorHandler(kbsRoutes.JsonResponseHandler(keyTransferPolicyController.Create))).Methods("POST")
				policyJson := `{
									"sgx_enclave_issuer_anyof": ["cd171c56941c6ce49690b455f691d9c8a04c2e43e0a4d30f752fa5285c7ee57f"],
									"sgx_enclave_issuer_product_id_anyof": [0]
							}`

				req, err := http.NewRequest(
					"POST",
					"/key-transfer-policies",
					strings.NewReader(policyJson),
				)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusCreated))
			})
		})
		Context("Provide a Create request without sgx_enclave_issuer_anyof", func() {
			It("Should fail to create new Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies", kbsRoutes.ErrorHandler(kbsRoutes.JsonResponseHandler(keyTransferPolicyController.Create))).Methods("POST")
				policyJson := `{
									"sgx_enclave_issuer_product_id_anyof": [0]
							}`

				req, err := http.NewRequest(
					"POST",
					"/key-transfer-policies",
					strings.NewReader(policyJson),
				)
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
		Context("Provide a Create request without sgx_enclave_issuer_product_id_anyof", func() {
			It("Should fail to create new Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies", kbsRoutes.ErrorHandler(kbsRoutes.JsonResponseHandler(keyTransferPolicyController.Create))).Methods("POST")
				policyJson := `{
									"sgx_enclave_issuer_anyof": ["cd171c56941c6ce49690b455f691d9c8a04c2e43e0a4d30f752fa5285c7ee57f"],
							}`

				req, err := http.NewRequest(
					"POST",
					"/key-transfer-policies",
					strings.NewReader(policyJson),
				)
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				req.Header.Set("Content-Type", consts.HTTPMediaTypeJson)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	// Specs for HTTP Get to "/key-transfer-policies/{id}"
	Describe("Retrieve an existing Key Transfer Policy", func() {
		Context("Retrieve Key Transfer Policy by ID", func() {
			It("Should retrieve a Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies/{id}", kbsRoutes.ErrorHandler(kbsRoutes.JsonResponseHandler(keyTransferPolicyController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/key-transfer-policies/ee37c360-7eae-4250-a677-6ee12adce8e2", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))
			})
		})
		Context("Retrieve Key Transfer Policy by non-existent ID", func() {
			It("Should fail to retrieve Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies/{id}", kbsRoutes.ErrorHandler(kbsRoutes.JsonResponseHandler(keyTransferPolicyController.Retrieve))).Methods("GET")
				req, err := http.NewRequest("GET", "/key-transfer-policies/e57e5ea0-d465-461e-882d-1600090caa0d", nil)
				Expect(err).NotTo(HaveOccurred())
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
	})

	// Specs for HTTP Delete to "/key-transfer-policies/{id}"
	Describe("Delete an existing Key Transfer Policy", func() {
		Context("Delete Key Transfer Policy by ID", func() {
			It("Should delete a Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies/{id}", kbsRoutes.ErrorHandler(kbsRoutes.ResponseHandler(keyTransferPolicyController.Delete))).Methods("DELETE")
				req, err := http.NewRequest("DELETE", "/key-transfer-policies/73755fda-c910-46be-821f-e8ddeab189e9", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
		})
		Context("Delete Key Transfer Policy by non-existent ID", func() {
			It("Should fail to delete Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies/{id}", kbsRoutes.ErrorHandler(kbsRoutes.ResponseHandler(keyTransferPolicyController.Delete))).Methods("DELETE")
				req, err := http.NewRequest("DELETE", "/key-transfer-policies/e57e5ea0-d465-461e-882d-1600090caa0d", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})
		Context("Delete Key Transfer Policy associated with Key", func() {
			It("Should fail to delete Key Transfer Policy", func() {
				router.Handle("/key-transfer-policies/{id}", kbsRoutes.ErrorHandler(kbsRoutes.ResponseHandler(keyTransferPolicyController.Delete))).Methods("DELETE")
				req, err := http.NewRequest("DELETE", "/key-transfer-policies/ee37c360-7eae-4250-a677-6ee12adce8e2", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

	// Specs for HTTP Get to "/key-transfer-policies"
	Describe("Search for all the Key Transfer Policies", func() {
		Context("Get all the Key Transfer Policies", func() {
			It("Should get list of all the Key Transfer Policies", func() {
				router.Handle("/key-transfer-policies", kbsRoutes.ErrorHandler(kbsRoutes.JsonResponseHandler(keyTransferPolicyController.Search))).Methods("GET")
				req, err := http.NewRequest("GET", "/key-transfer-policies", nil)
				Expect(err).NotTo(HaveOccurred())
				w = httptest.NewRecorder()
				req.Header.Set("Accept", consts.HTTPMediaTypeJson)
				router.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(http.StatusOK))

				var policies []kbs.KeyTransferPolicyAttributes
				json.Unmarshal(w.Body.Bytes(), &policies)
				// Verifying mocked data of 2 key transfer policies
				Expect(len(policies)).To(Equal(2))
			})
		})
	})
})
