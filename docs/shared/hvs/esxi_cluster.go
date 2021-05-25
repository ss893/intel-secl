/*
 *  Copyright (C) 2020 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package hvs

import "github.com/intel-secl/intel-secl/v4/pkg/model/hvs"

// ESXiCluster API request payload
// swagger:parameters ESXiClusterCreateRequest
type ESXiClusterCreateRequest struct {
	// in:body
	Body hvs.ESXiClusterCreateRequest
}

// ESXiCluster API response payload
// swagger:parameters ESXiCluster
type ESXiCluster struct {
	// in:body
	Body hvs.ESXiCluster
}

// ESXiClusterCollection response payload
// swagger:parameters ESXiClusterCollection
type ESXiClusterCollection struct {
	// in:body
	Body hvs.ESXiClusterCollection
}

// ---
//
// swagger:operation GET /esxi-cluster ESXi-Cluster Search-ESXi-cluster-records
// ---
//
// description: |
//
//   This API resource allows entire Clusters of VMWare  ESXi hosts to be managed as a group, using the vCenter Cluster object. When a Cluster is registered to the Host Verification Service, the Host Verification Service will automatically mirror the Cluster object in vCenter, automatically registering any ESXi hosts currently in the Cluster in vCenter.  As additional ESXi hosts are added or removed from the Cluster object in vCenter, the Host Verification Service will also register or remove the ESXi hosts from its own database.
//
//   Searches for ESXi clusters. Only one identifying parameter can be specified to search ESXi clusters which will return ESXi cluster collection as a result.
//
// x-permissions: esxi_clusters:search
// security:
//  - bearerAuth: []
// produces:
//  - application/json
// parameters:
// - name: id
//   description: ESXi cluster ID
//   in: query
//   type: string
//   format: uuid
//   required: false
// - name: clusterName
//   description: ESXi cluster name.
//   in: query
//   type: string
//   required: false
// - name: Accept
//   description: Accept header
//   in: header
//   type: string
//   required: true
//   enum:
//     - application/json
// responses:
//   '200':
//     description: Successfully retrieved the ESXi clusters.
//     content:
//       application/json
//     schema:
//       $ref: "#/definitions/ESXiClusterCollection"
//   '400':
//     description: Invalid search criteria provided
//   '415':
//     description: Invalid Accept Header in Request
//   '500':
//     description: Internal server error
//
// x-sample-call-endpoint: https://hvs.com:8443/hvs/v2/esxi-cluster?clusterName=Cluster-name
// x-sample-call-output: |
//      {
//          "esxi_clusters": [
//          {
//              "id": "9519febc-2c8d-4bb0-afec-b7a23db5735a",
//              "connection_string": "https://vCenter-url:443/sdk",
//              "cluster_name": "Cluster name",
//              "hosts": [
//                  "host.ip1",
//                  "host.ip2"
//              ]
//          } ]
//      }

// ---

// swagger:operation POST /esxi-cluster ESXi-Cluster Create-ESXi-cluster-record
// ---
//
// description: |
//   Creates a ESXi cluster record in database.
//
//   The serialized ESXi cluster Go struct object represents the content of the request body.
//
//    | Attribute                      | Description                                     |
//    |--------------------------------|-------------------------------------------------|
//    | cluster_name                   | Name of the vCenter cluster. The name needs to be exactly as it appears in vCenter. |
//    | connection_string 			   | The connection string is of the form <b>https://vCenter-url:443/sdk;u=vCenter-username;p=password"</b>. This is used to connect to vCenter and get the cluster information. |
//
// x-permissions: esxi_clusters:create
// security:
//  - bearerAuth: []
// produces:
// - application/json
// consumes:
// - application/json
// parameters:
// - name: request body
//   required: true
//   in: body
//   schema:
//    "$ref": "#/definitions/ESXiClusterCreateRequest"
// - name: Content-Type
//   description: Content-Type header
//   in: header
//   type: string
//   required: true
//   enum:
//     - application/json
// - name: Accept
//   description: Accept header
//   in: header
//   type: string
//   required: true
//   enum:
//     - application/json
// responses:
//   '201':
//     description: Successfully created the ESXi cluster record.
//     content:
//       application/json
//     schema:
//       $ref: "#/definitions/ESXiCluster"
//   '400':
//     description: Invalid request body provided
//   '415':
//     description: Invalid Content-Type/Accept Header in Request
//   '500':
//     description: Internal server error
//
// x-sample-call-endpoint: https://hvs.com:8443/hvs/v2/esxi-cluster
// x-sample-call-input: |
//      {
//          "connection_string" : "https://vCenter-url:443/sdk;u=vCenter-username;p=password",
//          "cluster_name" : "Cluster name"
//      }
// x-sample-call-output: |
//      {
//          "id": "9519febc-2c8d-4bb0-afec-b7a23db5735a",
//          "connection_string": "https://vCenter-url:443/sdk",
//          "cluster_name": "CSS-Attestation"
//      }

// ---

// swagger:operation GET /esxi-cluster/{esxi-cluster_id} ESXi-Cluster Retrieve-ESXi-cluster-record
// ---
//
// description: |
//   Retrieves an ESXi cluster.
//   Returns - The serialized ESXi cluster Go struct object that was retrieved and a list of associated host names
// x-permissions: esxi_clusters:retrieve
// security:
//  - bearerAuth: []
// produces:
// - application/json
// parameters:
// - name: esxi-cluster_id
//   description: Unique ID of the ESXi cluster.
//   in: path
//   required: true
//   type: string
//   format: uuid
// - name: Accept
//   description: Accept header
//   in: header
//   type: string
//   required: true
//   enum:
//     - application/json
// responses:
//   '200':
//     description: Successfully retrieved the ESXi cluster.
//     content:
//       application/json
//     schema:
//       $ref: "#/definitions/ESXiCluster"
//   '404':
//     description: No relevant cluster records found.
//   '415':
//     description: Invalid Accept Header in Request
//   '500':
//     description: Internal server error.
//
// x-sample-call-endpoint: https://hvs.com:8443/hvs/v2/esxi-cluster/9519febc-2c8d-4bb0-afec-b7a23db5735a
// x-sample-call-output: |
//          {
//              "id": "9519febc-2c8d-4bb0-afec-b7a23db5735a",
//              "connection_string": "https://vCenter-url:443/sdk",
//              "cluster_name": "Cluster name",
//              "hosts": [
//                  "host.ip1",
//                  "host.ip2"
//              ]
//          }

// ---

// swagger:operation DELETE /esxi-cluster/{esxi-cluster_id} ESXi-Cluster Delete-ESXi-cluster-record
// ---
//
// description: |
//   Deletes an ESXi cluster. If the cluster is still associated with any hosts, it will delete the host associations and the host records as well from database
// x-permissions: esxi_clusters:delete
// security:
//  - bearerAuth: []
// parameters:
// - name: esxi-cluster_id
//   description: Unique ID of the ESXi cluster.
//   in: path
//   required: true
//   type: string
//   format: uuid
// responses:
//   '204':
//     description: Successfully deleted the ESXi cluster.
//   '404':
//     description: The cluster to be deleted was not found.
//   '500':
//     description: Internal server error
// x-sample-call-endpoint: https://hvs.com:8443/hvs/v2/esxi-cluster/9519febc-2c8d-4bb0-afec-b7a23db5735a

// ---
