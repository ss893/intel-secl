Test files in this directory are from java HVS.  The following list 
describes where they came from and any alterations.

- host_manifest.json: extracted from 'mtwilson/v2/host-status?filter=false'
  - Removed the 'wrapper' json so that the json reflects a single manifest.
  - Altered Hardware features to have bools with quotes.
  - The PCR 15 event log hashes were manually updated due an issue in 
    application-agent.

- signed_flavor.json: imported flavors from the trust-agent via 'mtwilson/v2/flavors'

- trust_report.json: extracted from 'mtwilson/v2/reports/'.  
  - This file was reduced to a map of flavopart to trust reports.

- cms-ca-cert.pem: from /opt/mtwilson/configuration

- flavor-signer.crt.pem: from /opt/mtwilson/configuration

- PrivacyCA.pem: from /opt/mtwilson/configuration

- tag-cacerts.pem: from /opt/mtwilson/configuration