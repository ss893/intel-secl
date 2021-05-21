# HVS Installer

## Environment variables

Following env variable(s) are for configuring the behavior of hvS installer. Such variables have to be set in the shell
running installer binary.

Key         | Required | Type     | Description
----------- | -------- | -------- | --------------------------------------------------------
HVS_NOSETUP | -        | `string` | If set to `true`, installer will not perform setup tasks

## Usage

```shell
./hvs-v3.0.0.bin
```

## env file

`env file` is a file defining key-value pairs that can be automatically loaded during the setup process of hvs.
Reference [answer-file.md](../shared/setup/answer-file.md) for more information.

### Fields for HVS setup only

These fields are used for running hvs setup tasks only. These are not persisted in hvs configuration
file: `/etc/hvs/config.yml`

Field              | Required   | Type     | Description                                                 | Alternative
------------------ | ---------- | -------- | ----------------------------------------------------------- | ----------------------
BEARER_TOKEN       | `Required` | `string` | The bearer token for accessing `CMS`                        |
DB_SSL_CERT_SOURCE | -          | `string` | The source file from which to copy database SSL certificate | HVS_DB_SSL_CERT_SOURCE

### Fields for HVS configuration

Each of the following field corresponds a field in hvs configuration file.

Category  | Field                         | Required   | Type       | Default             | Alternative
--------- | ----------------------------- | ---------- | ---------- | ------------------- | ------------------------------
General   | AAS_BASE_URL                  | Required   | `string`   |                     |
General   | CMS_BASE_URL                  | Required   | `string`   |                     |
General   | CMS_CERT_DIGEST_SHA384        | `Required` | `string`   |                     |
HVS       | SERVICE_USERNAME              | `Required` | `string`   |                     |
HVS       | SERVICE_PASSWORD              | `Required` | `string`   |                     |
HVS       | DATA_ENCRYPTION_KEY           | -          | `string`   |                     | TLS                            | TLS_CERT_FILE | - | `string` |
TLS       | TLS_KEY_FILE                  | -          | `string`   |                     |
TLS       | TLS_COMMON_NAME               | -          | `string`   |                     |
TLS       | TLS_SAN_LIST                  | -          | `string`   |                     | SAN_LIST SAML                  | SAML_CERT_FILE | - | `string` |  |
SAML      | SAML_KEY_FILE                 | -          | `string`   |                     |
SAML      | SAML_COMMON_NAME              | -          | `string`   |                     |
SAML      | SAML_ISSUER_NAME              | -          | `string`   |                     |
SAML      | SAML_VALIDITY_SECONDS         | -          | `int`      | 86400               | Flavor Signing                 | FLAVOR_SIGNING_CERT_FILE | - | `string` |  |
Signing   | FLAVOR_SIGNING_KEY_FILE       | -          | `string`   |                     |
Signing   | FLAVOR_SIGNING_COMMON_NAME    | -          | `string`   |                     | Privacy CA                     | PRIVACY_CA_CERT_FILE | - | `string` |  |
CA        | PRIVACY_CA_KEY_FILE           | -          | `string`   |                     |
CA        | PRIVACY_CA_COMMON_NAME        | -          | `string`   |                     |
CA        | PRIVACY_CA_ISSUER             | -          | `string`   |                     |
CA        | PRIVACY_CA_VALIDITY_YEARS     | -          | `int`      |                     | Tag CA                         | TAG_CA_CERT_FILE | - | `string` |  |
CA        | TAG_CA_KEY_FILE               | -          | `string`   |                     |
CA        | TAG_CA_COMMON_NAME            | -          | `string`   |                     |
CA        | TAG_CA_ISSUER                 | -          | `string`   |                     |
CA        | TAG_CA_VALIDITY_YEARS         | -          | `int`      |                     |
CA        | ENDORSEMENT_CA_CERT_FILE      | -          | `string`   |                     |
CA        | ENDORSEMENT_CA_KEY_FILE       | -          | `string`   |                     |
CA        | ENDORSEMENT_CA_COMMON_NAME    | -          | `string`   |                     |
CA        | ENDORSEMENT_CA_ISSUER         | -          | `string`   |                     |
CA        | ENDORSEMENT_CA_VALIDITY_YEARS | -          | `int`      |                     |
Log       | LOG_MAX_LENGTH                | -          | `int`      |                     |
Log       | LOG_ENABLE_STDOUT             | -          | `bool`     |                     |
Log       | LOG_LEVEL                     | -          | `string`   |                     | Endorsement
Server    | SERVER_PORT                   | -          | `int`      |                     |
Server    | SERVER_READ_TIMEOUT           | -          | `Duration` |                     | HVS_SERVER_READ_TIMEOUT
Server    | SERVER_READ_HEADER_TIMEOUT    | -          | `Duration` |                     | HVS_SERVER_READ_HEADER_TIMEOUT
Server    | SERVER_WRITE_TIMEOUT          | -          | `Duration` |                     | HVS_SERVER_WRITE_TIMEOUT
Server    | SERVER_IDLE_TIMEOUT           | -          | `Duration` |                     | HVS_SERVER_IDLE_TIMEOUT
Server    | SERVER_MAX_HEADER_BYTES       | -          | `int`      |                     | HVS_SERVER_MAX_HEADER_BYTES
Database  | DB_VENDOR                     |            | `string`   |                     | HVS_DB_VENDOR
Database  | DB_HOST                       | -          | `string`   | localhost           | HVS_DB_HOSTNAME
Database  | DB_PORT                       | -          | `int`      | 5432                | HVS_DB_PORT
Database  | DB_NAME                       | -          | `string`   | localhost           | HVS_DB_NAME
Database  | DB_USERNAME                   | `Required` | `string`   |                     | HVS_DB_USERNAME
Database  | DB_PASSWORD                   | `Required` | `string`   |                     | HVS_DB_PASSWORD
Database  | DB_SSL_MODE                   | -          | `string`   | verify-full         | HVS_DB_SSL_MODE
Database  | DB_SSL_CERT                   | -          | `string`   | /etc/hvs/config.yml | HVS_DB_SSLCERT
Database  | DB_CONN_RETRY_ATTEMPTS        | -          | `int`      | 4                   |
Database  | DB_CONN_RETRY_TIME            | -          | `int`      | 1                   | HRRS                           | HRRS_REFRESH_PERIOD | - | `Duration` | 5 minutes ("5m") | VCSS | VCSS_REFRESH_PERIOD | - | `Duration` | 5 minutes ("5m") | Flavor Verification Service | FVS_NUMBER_OF_VERIFIERS | - | `int` | 20 |  | FVS_NUMBER_OF_DATA_FETCHERS | - | `int` | 20 |  | FVS_SKIP_FLAVOR_SIGNATURE_VERIFICATION | - | `bool` | false | Host Trust Manager | HOST_TRUST_CACHE_THRESHOLD | - | `int` | 100000 |
Audit Log | AUDIT_LOG_MAX_ROW_COUNT       | -          | `int`      | 10000               |
Audit Log | AUDIT_LOG_NUMBER_ROTATED      | -          | `int`      | 10                  |
Audit Log | AUDIT_LOG_BUFFER_SIZE         | -          | `int`      | 5000                |
