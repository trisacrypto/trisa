---
title: TRISA Errors and Rejections
date: 2022-07-01T16:45:58-04:00
lastmod: 2024-05-27T08:37:39-05:00
description: "Describing TRISA Errors and Rejections"
weight: 90
---

Error codes are standardized in the TRISA network to prevent confusion and allow easy identification of rejections or other problems to expedite the repair of the connection or information exchange. Then, the human-readable message stating the reason for the error should be loggable and actionable. Both standardized and unique/detail messages are acceptable. The message that caused the error should be retried with a fix, otherwise the rejection is permanent, and the request should not be retried. The payload of the additional data or reasons for the rejection, e.g., a parent error, a diff, a location for redirect, etc., should be described by the error code.

See the [error protocol buffers](https://github.com/trisacrypto/trisa/blob/main/proto/trisa/api/v1beta1/errors.proto) for more details.

## TRISA Errors Codes, Messages, and Descriptions

| ERROR CODE | ERROR MESSAGE | ERROR DESCRIPTION |
|---|---|---|
| 0 | UNHANDLED | Default error - something very bad happened. |
| 1 | UNAVAILABLE | Generic error - could not handle request, retry later. |
| 1 | SERVICE_DOWN_TIME | Generic error - could not handle request, retry later. |
| 1 | BVRC002 | Generic error - could not handle request, retry later. (Alias: Sygna BVRC Rejected Type) |
| 2 | MAINTENANCE | The service is currently in maintenance mode and cannot respond. |
| 3 | UNIMPLEMENTED | The RPC is not currently implemented by the TRISA node. |
| 49 | INTERNAL_ERROR | Request could not be processed by recipient. |
| 49 | BVRC999 | Request could not be processed by recipient. (Alias: Sygna BVRC Rejected Code) |
| 50 | REJECTED | Default rejection - no specified reason for rejecting the transaction. |
| 51 | UNKNOWN_WALLET_ADDRESS | VASP does not control the specified wallet address. |
| 52 | UNKNOWN_IDENTITY | VASP does not have KYC information for the specified wallet address. |
| 52 | UNKOWN_IDENTITY | VASP does not have KYC information for the specified wallet address. (Typo left for backwards compatibility.) |
| 53 | UNKNOWN_ORIGINATOR | Specifically, the Originator Account cannot be identified. |
| 54 | UNKOWN_BENEFICIARY | Specifically, the Beneficiary account cannot be identified. |
| 54 | BENEFICIARY_NAME_UNMATCHED | Specifically, the Beneficiary account cannot be identified. (Alias: Sygna BVRC Rejected Type) |
| 54 | BVRC007 | Specifically, the Beneficiary account cannot be identified. (Alias: Sygna BVRC Rejected Code) |
| 60 | UNSUPPORTED_CURRENCY | VASP cannot support the fiat currency or coin described in the transaction. |
| 60 | BVRC001 | VASP cannot support the fiat currency or coin described in the transaction.(Alias: Sygna BVRC Rejected Code) |
| 60 | UNSUPPORTED_NETWORK | VASP cannot support the specified network (Alias: helpful if the chain is not a currency) |
| 61 | EXCEEDED_TRADING_VOLUME | No longer able to receive more transaction inflows. |
| 61 | BVRC003 | No longer able to receive more transaction inflows. (Alias: Sygna BVRC Rejected Code) |
| 75 | UNSUPPORTED_ADDRESS_CONFIRMATION | VASP will not support the requested address confirmation mechanism. |
| 76 | CANNOT_CONFIRM_CONTROL_OF_ADDRESS | VASP cannot confirm they control the requested address, but doesn't necessarily indicate that they do not control the address. |
| 90 | COMPLIANCE_CHECK_FAIL | An internal compliance check has failed or black listing has occurred. |
| 90 | BVRC004 | An internal compliance check has failed or black listing has occurred. (Alias: Sygna BVRC Rejected Code) |
| 91 | NO_COMPLIANCE | VASP not able to implement travel rule compliance. |
| 92 | HIGH_RISK | VASP unwilling to conduct the transaction because of a risk assessment. |
| 99 | OUT_OF_NETWORK | Wallet address or transaction is not available on this network. |
| 100 | FORBIDDEN | Default access denied - no specified reason for forbidding the transaction. |
| 101 | NO_SIGNING_KEY | Could not sign transaction because no signing key is available. |
| 102 | CERTIFICATE_REVOKED | Could not sign transaction because keys have been revoked. |
| 103 | UNVERIFIED | Could not verify certificates with any certificate authority. |
| 104 | UNTRUSTED | A trusted connection could not be established. |
| 105 | INVALID_SIGNATURE | An HMAC signature could not be verified. |
| 106 | INVALID_KEY | The transaction bundle cannot be decrypted with the specified key. |
| 107 | ENVELOPE_DECODE_FAIL | Could not decode or decrypt private transaction data. (Alias: Sygna BVRC Rejected Type) |
| 107 | PRIVATE_INFO_DECODE_FAIL | Could not decode or decrypt private transaction data. (Alias: Sygna BVRC Rejected Code) |
| 107 | BVRC005 | Could not decode or decrypt private transaction data. |
| 108 | UNHANDLED_ALGORITHM | The algorithm specified by the encryption or signature is not implemented. |
| 150 | BAD_REQUEST | Generic bad request - usually implies retry when request is fixed. |
| 151 | UNPARSEABLE_IDENTITY | Could not parse the identity record; specify the type in extra. |
| 151 | PRIVATE_INFO_WRONG_FORMAT | Could not parse the identity record; specify the type in extra. (Alias: Sygna BVRC Rejected Type) |
| 151 | BVRC006 | Could not parse the identity record; specify the type in extra. (Alias: Sygna BVRC Rejected Code) |
| 152 | UNPARSEABLE_TRANSACTION | Could not parse the transaction data record; specify the type in extra. |
| 153 | MISSING_FIELDS | There are missing required fields in the transaction data, a list of these fields is specified in extra. |
| 154 | INCOMPLETE_IDENTITY | The identity record is not complete enough for compliance purposes of the receiving VASPs. Required fields or format specified in extra. |
| 155 | VALIDATION_ERROR | There was an error validating a field in the transaction data (specified in extra). |
| 198 | COMPLIANCE_PERIOD_EXCEEDED | Send by originating party if the review period has exceeded the required compliance timeline. |
| 199 | CANCELED | Cancel the ongoing TRISA exchange and do not send funds. |
| 199 | CANCEL_TRANSACTION | Cancel the ongoing TRISA exchange and do not send funds. (Alias: CANCELED) |