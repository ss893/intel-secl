# Setup Library - Answer file

## Functions

### `ReadAnswerFileToEnv`

```go
func ReadAnswerFileToEnv(filename string) error
```

`ReadAnswerFileToEnv` reads in the answer file and dump its contents into the
environment variables. 

#### Answer file

The format of answer files are defined as followed
- Lines started with pound sign (#) are comments thus ignored
  - **No** inline comments support
- All empty lines are ignored
- Key-value pair for environment variables should match the format `<key>=<value>`
  - The line is split at the first equal sign (`=`)
  - Everything before the equal sign is set as key, while everything after is value
    - Mal-formatted includes
      1. Either key or value is empty
      2. No equal sign found in line
    - Any `#` in a valid line is considered a part of key or value
  - All mal-formatted lines are ignored and will not return error

##### Example

- Unit test file
```text
# This is the answer file for testing

# It should export all the following environment variables
# with given keys and values
TEST_ENV_ONE=1
TEST_ENV_TWO=12
TEST_ENV_THREE=123
TEST_ENV_FOUR=1234
TEST_ENV_FIVE=12345

# THIS_LINE_SHOULD_BE_IGNORED=qwert

# following line should be ignored and not cause error
INVALID_LINE=
=invalid-line
```

- An example of an answer file for setting up HVS:

```text
CMS_BASE_URL=https://<cms-url>:<cms-port>/cms/v1/
CMS_TLS_CERT_SHA384=<cert-digest>
AAS_API_URL=https://<aas-url>:<aas-port>/aas/v1

BEARER_TOKEN=<token>

TLS_COMMON_NAME=Mt Wilson TLS Certificate
TLS_SAN_LIST=<hvs-ip>,127.0.0.1,localhost

HVS_DB_HOSTNAME=<host-name-of-db>
HVS_DB_NAME=hvsdb
HVS_DB_SSLCERTSRC=/etc/postgresql/server.crt
HVS_DB_PORT="5432"
HVS_DB_SSLMODE=verify-full
HVS_DB_USERNAME=vsdbuser
HVS_DB_PASSWORD=dbpassword
HVS_DB_NAME=vsdb
HVS_DB_CONN_RETRY_ATTEMPTS=5
HVS_DB_CONN_RETRY_TIME=100

SERVER_PORT=4567

HVS_SERVICE_USERNAME=<HVS username in AAS>
HVS_SERVICE_PASSWORD=<HVS password in AAS>
```