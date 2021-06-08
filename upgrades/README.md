## Upgrades of Intel<sup>®</sup> Security Libraries for Data Center (Intel<sup>®</sup> SecL-DC)

### Intel<sup>®</sup> SecL-DC started supporting upgrades from release v3.5

Following is the matrix of upgrade support for different components in Intel<sup>®</sup> SecL-DC

Latest release: v4.0.0

| Component | Abbreviation | Supports upgrade from  |
|-----------|--------------|-----------------------|
| Certificate Management Service           | CMS         |  v3.5.0, v3.6.0 |
| Authentication and Authorization Service | AAS         |  v3.5.0, v3.6.0 |
| Workload Policy Management               | WPM         |  v3.6.0         |
| Key Broker Service                       | KBS         |  v3.5.0, v3.6.0 |
| Trust Agent                              | TA          |  v3.5.0, v3.6.0 |
| Application Agent                        | AA          |  v3.5.0, v3.6.0 |
| Workload Agent                           | WLA         |  v3.5.0, v3.6.0 |
| Host Verification Service                | HVS         |  v3.5.0, v3.6.0 |
| Integration Hub                          | iHUB        |  v3.5.0, v3.6.0 |
| Workload Service                         | WLS         |  v3.5.0, v3.6.0 |
| SGX Caching Service                      | SCS         |  v3.5.0, v3.6.0 |
| SGX Quote Verification Service           | SQVS        |  v3.5.0, v3.6.0 |
| SGX Host Verification Service            | SHVS        |  v3.5.0, v3.6.0 |
| SGX Agent                                | AGENT       |  v3.5.0, v3.6.0 |
| SKC Client/Library                       | SKC Library |  v3.5.0, v3.6.0 |


NOTE:
iHUB and WPM does not support direct upgrade from v3.5.0 to v4.0.0. As,
1. iHUB service has been changed from single instance to multi-instance support in v3.6.0.
2. WPM has changed its directory structure

For these two components, user need to upgrade to v3.6.0 first then to the latest version v4.0.0
