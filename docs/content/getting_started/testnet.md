---
title: "TRISA TESTNET"
draft: false
weight: 40
---

## Introduction

For testing purposes to develop new clients and integrations, we have a test PKI available and some VASP nodes
running the latest version of TRISA server. This allows for full integration testing without disturbing the
production TRISA network mesh.

## Test VASP nodes

The following VASP nodes are available for testing:

* vasp1.testnet.trisa.io
* vasp2.testnet.trisa.io
* vasp3.testnet.trisa.io

Each VASP uses port _**8888**_ for the _**gRPC peer-to-peer**_ communication and port _**9999**_ for its _**admin**_ endpoints.


## Connecting your own node to the TRISA TESTNET

### Configuration

The TRISA server node accepts a `--config` command line argument with the path and filename of a configuration file which
specifies where the node can find its private key, certificate and the trust chain as well as what ports and
hostname to use.

Here is a sample configuration file: **[config.yaml](/testnet/config/config.yaml)**

```yaml
tls:
  privateKeyFile: server.key
  certificateFile: server.crt
  trustChain: trust.chain
server:
  listenAddress: :8888
  listenAddressAdmin: :9999
  hostname: myOwnVasp
```

The tls entries are all reference to filenames in which to find the artificats. We suggest creating a small directory
for this purpose for instance under the subdirectory `artifacts/testnet/myOwnVasp` to hold all these files and this config.yaml.

Though strictly not required to connect to the TESTNET, it is suggested to use the same gRPC and admin ports than those
that the TESTNET nodes use, namely 8888 for gRPC traffic and 9999 for admin HTTPS calls.

The hostname is not relevant for connecting to the TESTNET.




### Test certificate, private key and chain of trust

Note that the TESTNET PKI is currently in flux, so the root and issuing CA's may still change. To retrieve your
test certificate and private key, visit our TESTNET certificate portal page:

> [TESTNET Portal Page](http://testnet.trisa.io)

You will need to login using your Github account. You can generate as many test certificates as you want,
there are currently no rate limits implemented. Please read the rest of this section before generating your
key and certificate.

> ##### NOTE
> TRISA **PRODUCTION** SSL certificates will be issued following standard industry practices and artifacts.
> I.e. VASPs will have to submit a SSL CSR as part of their TRISA certification process.
> This TESTNET setup is merely a convenience to get developers up and running as quickly as possible.


#### server.crt

Create a certificate file, e.g. **server.crt** in the before established directory and copy everything between
```
-----BEGIN CERTIFICATE-----
```
and
```
-----END CERTIFICATE-----
```
including both tags from the newcert portal page and paste it into this certificate file.

#### server.key

Create a key file, e.g. **server.key** in the before established directory and copy everything between
```
-----BEGIN RSA PRIVATE KEY-----
```
and
```
-----END RSA PRIVATE KEY-----
```
including both tags from the newcert portal page and paste it into this key file.

#### trust.chain

Lastly, your TRISA node needs to know who to trust. For the production systems, the chain of trust will be retrieved
from embedded root CAs and external services they host. However for development purpose the chain of trust is passed
to the node in the form of a file containing the certificates of the trusted CAs.

[**Click here to download and save the trust.chain file**](/testnet/config/trust.chain) containing the list of the 3 CA certificates that are being used for the TRISA TESTNET and save it to your VASP configuration directory.


### Verifying the configuration

At this point the VASP configuration directory, e.g. `artifacts/testnet/myOwnVasp`, you created should have the following files:
```bash
config.yaml
server.crt
server.key
trust.chain
```

And each of them should have their respective, configuration, certificate, private key and chain of trust.




## Running your VASP node

Running the following command in the root of the cloned repo and assuming you created the above files in the
`artifacts/testnet/myOwnVasp` subdirectory, should start your VASP node and allow it to connect to the TRISA
TESTNET:

```bash
bazel run --run_under="cd $PWD/artifacts/testnet/myOwnVasp && " //cmd/trisa -- server --config config.yaml
```

The output should look something like this:

```
INFO: Build option --run_under has changed, discarding analysis cache.
INFO: Analyzed target //cmd/trisa:trisa (0 packages loaded, 7874 targets configured).
INFO: Found 1 target...
Target //cmd/trisa:trisa up-to-date:
  bazel-bin/cmd/trisa/darwin_amd64_stripped/trisa
INFO: Elapsed time: 0.538s, Critical Path: 0.02s
INFO: 0 processes.
INFO: Build completed successfully, 1 total action
INFO: Running command line: /bin/bash -c 'cd /Users/frank/go/src/github.com/trisacrypto/trisa/artifacts/testnet/myOwnVasp &&  bazel-bINFO: Build completed successfully, 1 total action
INFO[0000] starting TRISA admin server                   component=admin port=":9999" tls=listening
INFO[0000] starting TRISA server                         component=grpc port=":8888" tls=listening
```



## Testing your TESTNET connected VASP
While your node is running, let's see if we can connect to it and have it send a command to one of the TRISA
TESTNET vasps. Open another terminal window and execute the following command:

```bash
curl -ks "https://127.0.0.1:9999/send?target=vasp3.testnet.trisa.io:8888"
```
This command will connect to your local machine on port 9999 (the admin port) and send a test message to vasp3.testnet.trisa.io
using gRPC port 8888. The curl command output will be a `'.%'` if everything went fine.

Your node should have logged some output like this:

```
INFO[1028] sent transaction dadb719d-261f-47b1-a5b7-99326b90419a to vasp3.testnet.trisa.io:8888  identity="first_name:\"John\" last_name:\"Doe\" ssn:\"001-0434-4983\" state:\"CA\" driver_license:\"FA-387463\" " identity-type=trisa.identity.us.v1alpha1.Identity network="source:\"ae8a7287-78ef-4b09-a71f-71d38f929127\" destination:\"9baee60d-36b6-4dbc-8c9a-6216bb04fadf\" " network-type=trisa.data.bitcoin.v1alpha1.Data
INFO[1028] protocol envelope for incomingtransaction dadb719d-261f-47b1-a5b7-99326b90419a  direction=incoming enc_algo=AES256_GCM enc_blob="[206 225 79 54 251 177 103 181 203 81 108 253 28 45 93 142 203 111 242 142 104 29 168 106 81 73 14 229 47 206 247 163 227 108 79 165 3 217 190 188 185 13 175 97 6 121 181 87 200 245 138 91 79 216 193 92 232 76 15 214 211 66 157 205 61 100 79 86 69 144 130 107 186 31 155 118 240 183 192 113 136 250 41 169 131 112 192 183 30 88 107 144 125 60 8 242 204 223 109 57 240 70 249 180 142 95 92 141 227 109 18 32 26 175 9 223 31 68 22 40 100 252 237 120 6 136]" hmac="[223 69 96 218 116 1 202 106 4 50 98 71 166 254 125 210 18 17 184 212 147 204 160 67 74 115 227 192 134 99 120 90]" hmac_algo=HMAC_SHA256
INFO[1028] received transaction confirmation for dadb719d-261f-47b1-a5b7-99326b90419a  identity="first_name:\"Jane\" last_name:\"Foe\" national_number:\"109-800211-69\" city_of_birth:\"Zwevezele\" " identity-type=trisa.identity.be.v1alpha1.Identity
```

If you see some output like this, your local VASP node is able to connect to the TRISA TESTNET and you can start developing.
