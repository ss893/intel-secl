## Upgrades of Intel<sup>®</sup> Security Libraries for Data Center (Intel<sup>®</sup> SecL-DC)

### Intel<sup>®</sup> SecL-DC started supporting upgrades from release v3.5

Following is the matrix of upgrade support for different components in Intel<sup>®</sup> SecL-DC

Latest release: v4.0.0

#### Compatibility Matrix:
| Component (v4) |  CMS                | AAS                    | WPM | KBS                    | TA     | AA | WLA | HVS                    | iHUB | WLS                    |
|-------------|------------------------|------------------------|-----|------------------------|--------|----|-----|------------------------|------|------------------------|
| CMS         | NA                     | NA                     | NA  | NA                     | NA     | NA | NA  | NA                     | NA   | NA                     |
| AAS         | v3.5.0, v3.6.0, v4.0.0 | NA                     | NA  | NA                     | NA     | NA | NA  | NA                     | NA   | NA                     |
| WPM         | v3.5.0, v3.6.0, v4.0.0 | v3.5.0, v3.6.0, v4.0.0 | NA  | v3.5.0, v3.6.0, v4.0.0 | NA     | NA | NA  | NA                     | NA   | NA                     |
| KBS         | v3.5.0, v3.6.0, v4.0.0 | v3.5.0, v3.6.0, v4.0.0 | NA  | NA                     | NA     | NA | NA  | NA                     | NA   | NA                     |
| TA          | v3.5.0, v3.6.0, v4.0.0 | v4.0.0                 | NA  | NA                     | NA     | NA | NA  | v4.0.0                 | NA   | NA                     |
| AA          | NA                     | NA                     | NA  | NA                     | NA     | NA | NA  | NA                     | NA   | NA                     |
| WLA         | v3.5.0, v3.6.0, v4.0.0 | v3.5.0, v3.6.0, v4.0.0 | NA  | NA                     | v4.0.0 | NA | NA  | NA                     | NA   | v3.5.0, v3.6.0, v4.0.0 |
| HVS         | v3.5.0, v3.6.0, v4.0.0 | v4.0.0                 | NA  | NA                     | v4.0.0 | NA | NA  | NA                     | NA   | NA                     |
| iHUB        | v3.5.0, v3.6.0, v4.0.0 | v3.5.0, v3.6.0, v4.0.0 | NA  | NA                     | NA     | NA | NA  | v3.5.0, v3.6.0, v4.0.0 | NA   | NA                     |
| WLS         | v3.5.0, v3.6.0, v4.0.0 | v3.5.0, v3.6.0, v4.0.0 | NA  | v3.5.0, v3.6.0, v4.0.0 | NA     | NA | NA  | NA                     | NA   | NA                     |

#### Supported upgrade path:

Binary deployment:

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
WPM does not support direct upgrade from v3.5.0 to v4.0.0. As we have changed directory structure of WPM in v3.6

For WPM, user need to upgrade to v3.6.0 first then to the latest version v4.0.0

Container deployment:

| Component | Abbreviation | Supports upgrade from  |
|-----------|--------------|-----------------------|
| Certificate Management Service           | CMS         |  v3.6.0 |
| Authentication and Authorization Service | AAS         |  v3.6.0 |
| Workload Policy Management               | WPM         |  v3.6.0 |
| Key Broker Service                       | KBS         |  v3.6.0 |
| Trust Agent                              | TA          |  v3.6.0 |
| Application Agent                        | AA          |  v3.6.0 |
| Workload Agent                           | WLA         |  v3.6.0 |
| Host Verification Service                | HVS         |  v3.6.0 |
| Integration Hub                          | iHUB        |  v3.6.0 |
| Workload Service                         | WLS         |  v3.6.0 |
| SGX Caching Service                      | SCS         |  v3.6.0 |
| SGX Quote Verification Service           | SQVS        |  v3.6.0 |
| SGX Host Verification Service            | SHVS        |  v3.6.0 |
| SGX Agent                                | AGENT       |  v3.6.0 |
| SKC Client/Library                       | SKC Library |  v3.6.0 |

##### Upgrade to v3.6.0:
*iHUB* :
iHUB in v3.6.0, has added multi instance installation support. Hence, it requires following ENV variables for the upgrade,

```shell
HVS_BASE_URL
SHVS_BASE_URL
```

##### Upgrade to v4.0.0:
*TA and WLA* :
In v4.0.0, TA has modified policies on TPM and NVRAM and it requires to re-provision itself with HVS. This would need following 
ENV variable for the upgrade. Also, WLA would need to recreate keys as Binding Key gets updated after re-provisioning.

```shell
BEARER_TOKEN
```

NOTE:
If in case some components needs to get upgraded from v3.5.0, directly to the latest version then it would require ENV variables 
if they are mentioned above and comes in the upgrade path.
e.g if iHUB needs upgrade from v3.5.0 to v4.0.0 then it would require following ENV variables,

```shell
HVS_BASE_URL
SHVS_BASE_URL
```