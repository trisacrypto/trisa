---
title: "Hackathon"
draft: false
weight: 30
---

{{% button icon="fas fa-download" url="/conference-2019/trisa-conference-hackathon-2019.pdf" %}}Hackathon Slide Deck{{% /button %}}

## Projects & Slide Decks

* [VASP Adress Resolution](/conference-2019/trisa_vasp_address_resolution.pdf)
* [Enhancement of PII privacy for TRISA transactions](/conference-2019/trisa_pii_enhancement.pdf)
* [Red Flagging](/conference-2019/trisa_redflagging.pdf)
* [Travel Rule Extension TrisaScan](/conference-2019/trisa_travel_rule_extension_trisascan.pdf)
* [Zero Knowledge TRISA](/conference-2019/trisa_zk.pdf)

### VASP Adress Resolution

VASPs need to transfer some required data when sending value to another VASP.
One problem that the current implementation of TRISA does not address is recognizing or 
identifying an address as belonging to a VASP. This proposal complements the TRISA protocol.

### Enhancement of PII privacy for TRISA transactions

TRISA suggests a VASP to expose the PII of any account holder listed in a proposed transaction.
This runs directly contrary to one of the key benefits of digital currency which is the ability
execute low-value transactions with low friction.

### Red Flagging

Determine and return transactions (incl participants) where the use of mixers or other obfuscation techniques (e.g. transaction segmentation, 'hawallah', etc.) are detected - which we refere to as "red flagged" transactions.

### Travel Rule Extension TrisaScan

We want to make Travel Rule benefit to user, make it more simple for user to see
VASP, individual information in etherscan. This solution integrates with TRISA library.

### Zero Knowledge TRISA

Our constraints:

* We want all VASPs to use TRISA
* Travel Rule: Counterparty PII recorded before a Tx is sent
* All ICOs accepting Fiat etc. are VASPs
* … How many ICOs in 2018 were Scams...

So we’s have to bottleneck access to TRISA, but
* this limits adoption of TRISA
* This delegates trust to the TRISA gatekeepers
* who may have different incentives than users & VASPs
* Small Exchanges can be sold...

## Results


Team | Score 
------|-----
VASP Adress Resolution | 45.5
Enhancement of PII privacy for TRISA transactions | 52
Red Flagging | 48
Travel Rule Extension TrisaScan | 45
Zero Knowledge TRISA | 49.5


## Follow-up

* Monthly online meetup
* Elect board members
* Bylaws proposal
* Setup trisacrypto organization
* Workout production PKI
* Formalize design (issue tracking)
