---
title: TRIXO Form
date: 2022-06-28T09:15:46-04:00
lastmod: 2022-06-28T16:53:39-04:00
description: "Shows full TRIXO forms in JSON & Spreadsheet format"
weight: 21
---

The purpose of the TRIXO questionnaire is to provide a common set of questions VASPs can ask other VASPs with whom they are planning to begin exchanging Travel Rule information. It is not intended to be a due diligence questionnaire for onboarding other VASPs as customers. 

#### Spreadsheet View:

|            | VASP TRAVEL RULE INFORMATION EXCHANGE (TRIXO) QUESTIONNAIRE                                                                                                                                                                                              |   |
|------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---|
|            | "This questionnaire was derived from the draft ""Wolfsberg-style Questionnaire for VASPs"                                                                                                                                                                |   |
|            |                                                                                                                                                                                                                                                          |   |
| No #       | Question                                                                                                                                                                                                                                                 |   |
| SECTION 1. | ENTITY DETAILS                                                                                                                                                                                                                                           |   |
|     1a.    | Full Legal Name                                                                                                                                                                                                                                          |   |
|     1b.    | Doing Business As (DBA) name                                                                                                                                                                                                                             |   |
|     1c.    | Full Legal (Registered) Address                                                                                                                                                                                                                          |   |
|     1d.    | Full Primary Business Address (if different from Registered Address above)                                                                                                                                                                               |   |
|     1e.    | Date of Entity Incorporation / Establishment                                                                                                                                                                                                             |   |
|     1f.    | Incorporation Number                                                                                                                                                                                                                                     |   |
|     1g.    | Entity Identifier (e.g. Legal Entity Identifier, Employer Indentification Number), if available                                                                                                                                                          |   |
| SECTION 2. | REGULATOR & LICENSING                                                                                                                                                                                                                                    |   |
|     2a.    | Name of the Entity's primary financial regulator / supervisory authority                                                                                                                                                                                 |   |
|     2b.    | Please provide a list of national jurisdictions, other than your primary national jurisdiction, where you have been granted licenses or other approvals or have registered as required to operate, and the name of the regulator / supervisory authority |   |
|     2c.    | Entity's License / Registration Number(s) for each jurisdiction it operates in                                                                                                                                                                           |   |
|     2d.    | Is the Entity permitted to send and/or receive transfers of virtual assets in the jurisdictions in which it operates?                                                                                                                                    |   |
| SECTION 3. | CDD & TRAVEL RULE POLICIES                                                                                                                                                                                                                               |   |
|     3a.    | Does the Entity have a programme that sets minimum AML, CFT, KYC / CDD and Sanctions standards per the requirements of the jurisdiction(s) regulatory regimes where it is licensed/approved/registered?                                                  |   |
|     3b.    | Does the Entity conduct KYC / CDD before permitting its customers to send/receive virtual asset transfers?                                                                                                                                               |   |
|     3c.    |       If Yes, at what threshold does the Entity conduct KYC before permitting the customer to send/receive virtual asset transfers?                                                                                                                      |   |
|     3d.    | Is the Entity required to comply with the application of the Travel Rule standards (FATF Recommendation 16) in the jurisdiction(s) where it is licensed / approved / registered?                                                                         |   |
|     3e.    |      If Yes, please specify the applicable regulation(s)                                                                                                                                                                                                 |   |
|     3f.    | What is the minimum threshold above which the entity is required to collect/send Travel Rule information?                                                                                                                                                |   |
| SECTION 4. | DATA PROTECTION                                                                                                                                                                                                                                          |   |
|     4a.    | Is the Entity required by law to safeguard PII?                                                                                                                                                                                                          |   |
|     4b.    | Does the Entity secure and protect PII, including PII received from other VASPs under the Travel Rule?                                                                                                                                                   |   |
| SECTION 5. | TRAVEL RULE IMPLEMENTATION                                                                                                                                                                                                                               |   |
|     5a.    | Which technical solution(s) does the Entity support for sharing Travel Rule information?                                                                                                                                                                 |   |
|     5b.    | Please provide the technical details (IDs, endpoints, URLs, etc.) required to send Travel Rule information to the Entity for each solution the Entity supports.                                                                                          |   |
|     5c.    | Name, email and phone number of travel rule contact                                                                                                                                                                                                      |   |

#### JSON View

```json


 "trixo": {
  "entity": {
    "name": {
      "name_identifiers": [
        {
          "legal_person_name": "AliceCoin, Inc.",
          "legal_person_name_identifier_type": 0
        },
        {
          "legal_person_name": "Alice VASP",
          "legal_person_name_identifier_type": 2
        },
        {
          "legal_person_name": "Alice",
          "legal_person_name_identifier_type": 1
        }
      ],
      "local_name_identifiers": [],
      "phonetic_name_identifiers": []
    },
    "geographic_addresses": [
      {
        "address_type": 1,
        "address_line": [
          "23 Roosevelt Place",
          "",
          "Boston, MA 02151"
        ],
        "country": "USA"
      }
    ],
    "customer_number": "",
    "national_identification": {
      "national_identifier": "5493004YBI24IF4TIP92",
      "national_identifier_type": 8,
      "country_of_issue": "USA",
      "registration_authority": "RA000744"
    },
    "country_of_registration": "USA"
  },
  "contacts": {
    "technical": {
      "name": "Benjamin Bengfort",
      "email": "benjamin@trisa.io",
      "phone": "+1-701-456-5555"
    },
    "legal": {
      "name": "John Jefferies",
      "email": "jj@trisa.io",
      "phone": "+1-202-789-5555"
    },
    "administrative": {
      "name": "Laura Montoya",
      "email": "laura@trisa.io",
      "phone": "+1-310-123-5555"
    },
    "billing": {
      "name": "Godfried Moneybags",
      "email": "godfreid@example.com",
      "phone": "+1-555-555-5555"
    }
  },
  "trisa_endpoint": "api.alice.vaspbot.net:443",
  "common_name": "api.alice.vaspbot.net",
  "website": "https://alice.vaspbot.net/",
  "business_category": 3,
  "vasp_categories": [
    "Exchange",
    "Individual",
    "Other"
  ],
  "established_on": "2021-01-21",
  "trixo": {
    "primary_national_jurisdiction": "USA",
    "primary_regulator": "FinCEN",
    "other_jurisdictions": [],
    "financial_transfers_permitted": "no",
    "has_required_regulatory_program": "yes",
    "conducts_customer_kyc": true,
    "kyc_threshold": "1.00",
    "kyc_threshold_currency": "USD",
    "must_comply_travel_rule": true,
    "applicable_regulations": [
      "FATF Recommendation 16"
    ],
    "compliance_threshold": "3000.00",
    "compliance_threshold_currency": "USD",
    "must_safeguard_pii": true,
    "safeguards_pii": true
  }
}
```