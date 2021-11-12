---
title: Robot VASPs
date: 2021-06-14T12:50:23-05:00
lastmod: 2021-06-14T12:50:23-05:00
description: "Working with rVASPs for TestNet integration"
weight: 10
---

The TestNet hosts three convenience "robot VASPs" (rVASPs) services to facilitate integration and testing with the TRISA TestNet. These services are as follows:

- Alice (`api.alice.vaspbot.net:443`): the primary integration rVASP used to trigger and receive TRISA messages.
- Bob (`api.bob.vaspbot.net:443`): a demo rVASP to view exchanges with Alice.
- Evil (`api.evil.vaspbot.net:443`): an "evil" rVASP that is not a TRISA TestNet member, used to test non-authenticated interactions.

Note: the rVASPs are currently primarily configured for demos and work has begun on making them more robust for integration purposes; please check in with this documentation regularly for any changes. If you notice any bugs in the rVASP code or behavior, [please open an issue](https://github.com/trisacrypto/testnet/issues).

## Getting Started with rVASPs

There are two ways that you can use the rVASPs to develop your TRISA service:

1. You can trigger the rVASP to send a TRISA exchange message to your service.
2. You can send a TRISA message to the rVASP with a valid (or invalid) rVASP customer.

The rVASPs have a built-in database of fake customers with fake wallet addresses. Their response to TRISA messages or to a triggered transfer requires the originator/beneficiary customer to be valid for the rVASP. E.g. if the customer wallet address is a valid Alice address and Alice is the beneficiary, the rVASP will respond with the customer's fake KYC data but if not, it will return an TRISA error code.

The following table of "customers" for Alice, Bob, and Evil can be used as a reference for interacting with each rVASP:

| VASP                  | "Crypto Wallet"                    | Email                 |
|-----------------------|------------------------------------|-----------------------|
| api.bob.vaspbot.net   | 18nxAxBktHZDrMoJ3N2fk9imLX8xNnYbNh | robert@bobvasp.co.uk  |
| api.bob.vaspbot.net   | 1LgtLYkpaXhHDu1Ngh7x9fcBs5KuThbSzw | george@bobvasp.co.uk  |
| api.bob.vaspbot.net   | 14WU745djqecaJ1gmtWQGeMCFim1W5MNp3 | larry@bobvasp.co.uk   |
| api.bob.vaspbot.net   | **1Hzej6a2VG7C8iCAD5DAdN72cZH5THSMt9** | fred@bobvasp.co.uk    |
| api.alice.vaspbot.net | **19nFejdNSUhzkAAdwAvP3wc53o8dL326QQ** | sarah@alicevasp.us    |
| api.alice.vaspbot.net | 1ASkqdo1hvydosVRvRv2j6eNnWpWLHucMX | mary@alicevasp.us     |
| api.alice.vaspbot.net | 1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha | alice@alicevasp.us    |
| api.alice.vaspbot.net | 14HmBSwec8XrcWge9Zi1ZngNia64u3Wd2v | jane@alicevasp.us     |
| api.evil.vaspbot.net  | **1AsF1fMSaXPzz3dkBPyq81wrPQUKtT2tiz** | gambler@evilvasp.gg   |
| api.evil.vaspbot.net  | 1PFTsUQrRqvmFkJunfuQbSC2k9p4RfxYLF | voldemort@evilvasp.gg |
| api.evil.vaspbot.net  | 172n89jLjXKmFJni1vwV5EzxKRXuAAoxUz | launderer@evilvasp.gg |
| api.evil.vaspbot.net  | 182kF4mb5SW4KGEvBSbyXTpDWy8rK1Dpu  | badnews@evilvasp.gg   |

Note that all rVASP data was generated using a Faker tool to produce realistic/consistent test data and fixtures and is completely fictional. For example, the records for Alice VASP (a fake US company) are primarily in North America, etc.

If you're a Traveler customer, the bold addresses above have some attribution data associated with them, and they're a good candidate for Traveler-based rVASP interactions.

### Preliminaries

This documentation assumes that you have a service that is running the latest `TRISANetwork` service and that it has been registered in the TRISA TestNet and correctly has TestNet certificates installed. See [ TRISA Integration Overview]({{< ref "integration/_index.md" >}}) for more information. **WARNING**: the rVASPs do not participate in the TRISA production network, they will only respond to verified TRISA TestNet mTLS connections.

To interact with the rVASP API, you may either:

1. Use the `rvasp` CLI tool
2. Use the rVASP protocol buffers and interact with the API directly

To install the CLI tool, either download the `rvasp` executable for the appropriate architecture at the [TestNet releases page](https://github.com/trisacrypto/testnet/releases), clone [the TestNet repository](https://github.com/trisacrypto/testnet/) and build the `cmd/rvasp` binary or install with `go get` as follows:

```
$ go get github.com/trisacrypto/testnet/...
```

To use the [rVASP protocol buffers](https://github.com/trisacrypto/testnet/tree/main/proto/rvasp/v1), clone or download them from the TestNet repository then compile them to your preferred language using `protoc`.

### Triggering an rVASP to send a message

The rVASP admin endpoints are used to interact with the rVASP directly for development and integration purposes. Note that this endpoint is different than the TRISA endpoint, which was described above.

- Alice: `admin.alice.vaspbot.net:443`
- Bob: `admin.bob.vaspbot.net:443`
- Evil: `admin.evil.vaspbot.net:443`

To use the command line tool to trigger a message, run the following command:

```
$ rvasp transfer -e admin.alice.vaspbot.net:443 \
        -a mary@alicevasp.us \
        -d 0.3 \
        -B trisa.example.com \
        -b cryptowalletaddress \
        -E
```

This message sends the Alice rVASP a message using the `-e` or `--endpoint` flag, and specifies that the originating account should be "mary@alicevasp.us" using the `-a` or `--account` flag. The originating account is used to determine what IVMS101 data to send to the beneficiary. The `-d` or `--amount` flag specifies the amount of "AliceCoin" to send. Finally, the `-b` or `--beneficiary` flag specifies the crypto wallet address of the beneficiary you intend as the recipient.

The next two parts are critical. The `-E` or `--external-demo` flag tells the rVASP to trigger a request to your service rather than to perform a demo exchange with another rVASP. This flag is required! Finally, the `-B` or `--beneficiary-vasp` flag specifies where the rVASP will send the request. This field should be able to be looked up in the TRISA TestNet directory service; e.g. it should be your common name or the name of your VASP if it is searchable.

Note that you can set the `$RVASP_ADDR` and `$RVASP_CLIENT_ACCOUNT` environment variables to specify the `-e` and `-a` flags respectively.

To use the protocol buffers directly, use the `TRISAIntegration` service `Transfer` RPC with the following `TransferRequest`:

```json
{
    "account": "mary@alicevasp.us",
    "amount": 0.3,
    "beneficiary": "cryptowalletaddress",
    "beneficiary_vasp": "trisa.example.com",
    "check_beneficiary": false,
    "external_demo": true
}
```

These values have the same exact specification as the ones in the command line program.

### Sending a TRISA message to an rVASP

The rVASP expects a `trisa.data.generic.v1beta1.Transaction` as the transaction payload and an `ivms101.IdentityPayload` as the identity payload. The identity payload beneficiary information does not need to be populated, the rVASP will respond populating the beneficiary, however the identity payload should not be null. It is recommended that you specify some fake identity data to take advantage of the rVASP parsing and validation commands.

Make sure that in your transaction payload you specify a beneficiary wallet that matches the rVASP beneficiaries from the table above; e.g. use:

```json
{
    "txid": "1234",
    "originator": "anydatahereisfine",
    "beneficiary": "1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha",
    "amount": 0.3,
    "network": "TestNet",
    "timestamp": "2021-06-14T16:41:52-05:00"
}
```

You may specify any `txid` or `originator` string and the `network` and `timestamp` fields are ignored.

Create a sealed envelope either using the directory service or direct key exchange to fetch the rVASP RSA public keys and using `AES256-GCM` and `HMAC-SHA256` as the envelope cryptography. Then use the `TRISANetwork` service `Transfer` RPC to send the sealed envelope to the rVASP.

TODO: Soon the `trisa` command line program will be available. We will specify how to use the CLI program to send a message here as soon as it is released.

