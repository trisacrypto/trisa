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

Here is a sample configuration file: **config.yaml**

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

Below is the list of the 3 CA certificates that are being used for the TRISA TESTNET. Copy everything between the **first**
```
-----BEGIN CERTIFICATE-----
```
and **last**
```
-----END CERTIFICATE-----
```
including both tags and paste it into a chain of trust file e.g. **trust.chain**.


```sh
-----BEGIN CERTIFICATE-----
MIIFQDCCAyigAwIBAgIUE87m4cufCbdffDpup1WeHOZxIGkwDQYJKoZIhvcNAQEN
BQAwODEWMBQGA1UEChMNVFJJU0EgVGVzdE5ldDEeMBwGA1UEAxMVVFJJU0EgVGVz
dG5ldCBSb290IENBMB4XDTIwMDQwNjAyNDcwMFoXDTMwMDQwNDAyNDcwMFowODEW
MBQGA1UEChMNVFJJU0EgVGVzdE5ldDEeMBwGA1UEAxMVVFJJU0EgVGVzdG5ldCBS
b290IENBMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAtWawb42ZnPZi
yD8G8F8dakNJaHXWEYZL/Sx/yDNeWgmeDu74Q5mEB/H9R85yyf+M9mWHT13F2SGf
akHnPdGolPUqBKdrCunke/cKsDjvytMBUUqWDh/jQgzZLwdc4WLUkYLXwJsjgSjj
zVcBeXJn4xk0RYiLTN6RDbgS7BOm458U3b2DvFWjoPa94NUj3jLMDnnHmPLkgw6M
FPLGnYVPRkJbaiv95sws4Jigy0iBdqtnB3Twoxml+l4vEqiQEu08gz+ujXz3i7Zj
Q2DkGob5y52DcskGH9H4tZAkXwS6gyL5pkEpdcFuv9bbSWznxQxJNfgmsGlm3Abd
flHFaoxzrl63dxh5Ea/hNl52/dz7qZ3MBrZqGsinL0jvdruO4vdhnB/s1aotOD+h
r17iXP+tfFmmpWjFjG0o+VUpWO6OmbKF72SGvRxoFWPTgKMLU+a75NwQTfSDTwME
EsqQgKq20lRc6oOi5niZbqlQZKlUuPQhHLsJIDu3GlANd7/lE2PMUiyzeP6xq0aN
6ye9sS63+AIqArFUF2nx8Lg2pZNHXWUriwTlCObf01/RukGE5f4pz62PEmO7DlT4
58xX/d0TUxVH8+6+638B+tMBiRK+0CAeNIiculNjMTrH+IYu6sXabzay10qubR0t
ZHB6Ro1cmz4r4B98Lo7+SoVPxPMTjWUCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgEG
MA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHug4sSTGfPzmp+KPhtpjWgmISB3
MA0GCSqGSIb3DQEBDQUAA4ICAQAkhs0kURn9Jsv/0G7npna46Xawo7WbMRPsNek4
4ueeG+7wMGCjOc8W3+wGMgqI0lqCiUDRh69chYmUPZX0Husg0cQL9rZebL9Fsg2q
cGkOFnDJJcnoOQaiJpqsLbq+aL/AZun22DDdik8pZqEtCK0TYK4wsP/F3UrgJDTY
vxz6IiDE2PrWA/IevIk5gskowhg0ZVXlGcADT7HQ42d0SARdJ46c0BvXhn8V1B4K
udSv+rJAfzAFJsvX9VGB0wSuqMqd91nkGieeS7Wddw5gj0smr+XtcVz7YKOtV5Tc
yQIvva7HMYK8wWsIqURn3Z14LhxtPa2CdWeJ2fEo3YAhUxuW1F4OwsemwGy6erx8
ER80Y4HJrVY/cDpGjhDoVAFfOboeAiICdYJFBHWqiWTUEHJyIHR4QUOyLEsyMByL
Adhd+4fZGW9CbXxRD6iTVWCv9Tv5oOZjXeqwuKZ7aaBYspuNPKqY5yKip/dMpqB7
X/8v5YHBTkUAylcv46XWahfON39KYer05757hXJXIqwHklOxQPOGyHi1sr6XxYer
FMYoR9O392e2YDx0FOrFFzJFsYrBPchB8VttDJAMmIqHCxIYG396xNNfy86cK3rE
hJaNebXn1lI1+L26fkCsvAksw6v0ShGkphOC09Q/S7zO/chAv2F4DqYnjrMCwcOL
affP/w==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFZjCCA06gAwIBAgIUDymOMrDn6XjmlRj/IT6fZtkT1EQwDQYJKoZIhvcNAQEN
BQAwODEWMBQGA1UEChMNVFJJU0EgVGVzdE5ldDEeMBwGA1UEAxMVVFJJU0EgVGVz
dG5ldCBSb290IENBMB4XDTIwMDQwNjAyNDcwMFoXDTIxMDQwNjAyNDcwMFowPTEW
MBQGA1UEChMNVFJJU0EgVGVzdE5ldDEjMCEGA1UEAxMaVFJJU0EgVGVzdE5ldCBJ
c3N1aW5nIENBIDEwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC0Vf+r
xUmMKRAgr8xR6BfdDeSn/lrbIryCSPJTIBB/NjGluvmHoH29gVrvU99aTKwkbL0W
Qw04jS2eGuMUM/dJnIT9g2B0Fa5eP1rFTB/zE5YFExySEvY4h6JjZ/ab4XybZIVS
YnYEudbZ2jZQBWMM+JWdgTO6HWi3yBk8siODmkTDkS5ZBNbqnEtNPYgi1qUHfMMb
kJ/BIGGevv2gOy+hrnjMpV6AXexGifdzDCoYIA6WIlXRHEh6CqyRGnvssYswWAER
uxLEf+HTPMX2N0cRAbg4wcPrYgJmnTVF/hzn54raajgB/6fVICC5nREOeszk75d+
z/jM9smVPDhDdyFzbZED6wt9Yf5g4TL9WHZGmhakO9Dy6BXt/fAhLwiVxm3nSYM4
iFKHnJYN9juN0A3bX0JMuoRm1XvVUOVoUSULsLYDzmEv6dYkgiDQrnjNsSw6Z/1O
Or78rsFIb6wsTwlystcdHzoQYjY7gXjp85204EKIhEGNMMirfKbtFNdV+pDq10gD
JqkFTSC7yoIoFkyZ1U1tBr7fCYqdymAHGzQaqmRWmvrWGu049upYv5Es57dIJ/o4
b8mspSLE835Mjk7jy1MmmksP3+ZeGV6rJzz8r4PWrI9v/y83zNBWJtc6YuTU87om
H4RQAyyAGePBeKWWTzz2YCtAefuFi5zGkp84IwIDAQABo2MwYTAOBgNVHQ8BAf8E
BAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU/jo1K6Uqh/kM8tPwGeOn
mEEqLYQwHwYDVR0jBBgwFoAUe6DixJMZ8/Oan4o+G2mNaCYhIHcwDQYJKoZIhvcN
AQENBQADggIBACYv/kFjfgauDPWJXSb+XMPT6bfsJztC55bX2tpaLpV9u6QNH72r
t79PvOSQQ3eV72Gw1tfaiWbPRKIY9SGELjp/G7vcOS4pF8BTesJnHva7aEsHqtLL
fqw+ZkrKW+ZyvY0vefdBh5hrJIMKDgGkuKayiaIzHzkBBSdBPJEOg818cVErmNX2
dLAu3o5xl7v20sQmK9j0gEDYxj7tuhOWVYZmcog8bQWxXtCs14znNut00Gn6CSRI
CsdwWFvrbZlfiv4OqTEZjsmqx/aSf6IwjfUYjATGVaX9Uwb+YMjp+PyADMbaj81K
XKP0gBvaoBiKXOmRl3/xQmVgYZfV24UEv1tPlY5ueKw+3zVZei/OQG1gpdieK4Tm
fq+yRkR23PLXcrGoUWw519+AJHu0Uiq1ZFf7UTZIq8Y/4ZPDyIkgP2TycxcBJEEy
vwMdsbchzpGbd7OqWFEWIjAAza0zc9kRkqVuygbTNYSC5PzeeolcS2l7+BG5lyG2
7rXed0PxMbz4wHYclHuuI6EKYJRVWpm8Uku60mgzSDMdVVxy5ccYCTvJ/VuW0Xff
Cluk2LyVAvi+DoViS5QABDA1duMvYs8YmwwYo00xhSsjsTg+BVlULqsQH/AMqX21
iFzTSfl43xduyUb7zT2mxq0iVUP0K16jXgVa10CCW4dcrK7hPpXGd9A9
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFZjCCA06gAwIBAgIUO+nNG89+qVKI2E45gcKBehhWVaowDQYJKoZIhvcNAQEN
BQAwODEWMBQGA1UEChMNVFJJU0EgVGVzdE5ldDEeMBwGA1UEAxMVVFJJU0EgVGVz
dG5ldCBSb290IENBMB4XDTIwMDQwNjAyNDcwMFoXDTIxMDQwNjAyNDcwMFowPTEW
MBQGA1UEChMNVFJJU0EgVGVzdE5ldDEjMCEGA1UEAxMaVFJJU0EgVGVzdE5ldCBJ
c3N1aW5nIENBIDIwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDLzfww
hDHETRYL2/+nfcthBju09h9i2T72w+VAEFLMwDVQvXfnEGdToznx8Robt1SCRglj
Y4zN7vaqWHEtnz/76UMIinLzjz3BGaMsyiFeQ7q/ZEoQ7SsHOJzuHaUVTFydtQfQ
jWMPn/5GqfUGiftEvtwhsMnFaQ7CbeucV6KwfU7nuZBHBxZKFDxclZmfFxLyAdAp
mSo4Y9/XRnfGv3hAl8+1VhhYmZl7pSC2uNnPy2cply87J7EChUpIFhL5PplEQa7n
Hj95AuEv8Axj1iHw/W4wwfG4JnDdJPVDaXnZuEfHuAHxMd+ZgId/vRjkYQtk6wEF
3b4/HSMgwV2HOr4HnsBuSu5hw5yhTjVDpCALbgSg3DojzeZiftJtkJD2zcDu2WYU
KZxLSuu0iMw8E2xbptdx4AoOeTkJxvcBqGFqvdATFwbn1AoRmN2NDuANlqWvInMz
JEa8HMmBAhBje9WNc+hOnRVcmzIDDm6U8JFIUz+A20x6NpmYKUkIJa0KxAl19G8T
rdhIXCyW2lGYKXc37JJQbK7CtS3O7Ba/Cdrse4WocQmu/XtO5oPVI3nDmbwstxFN
yh2++ZZW4l72eqStW12PUCzF1eok23T7ppNkj64FZZwTIxygTVNeXEdQTp4axjM+
iP6X8XOCABCtJMBy7t7XMsI6jjZ+PzdM0J8bjQIDAQABo2MwYTAOBgNVHQ8BAf8E
BAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUev4c5lKpxHbn8iZ7R75+
tS4daLAwHwYDVR0jBBgwFoAUe6DixJMZ8/Oan4o+G2mNaCYhIHcwDQYJKoZIhvcN
AQENBQADggIBAFXWbKI0lyQ7Jx9iNXXAuGG8gena7HP814PksJQi7xqwITR5pLy2
gcenMSWmLQc4678zjt9HyhjDkkwwxYlCFSP4Kq2G0m+qSb8tCc4HwvY3P8o2EqgE
mcy5Z4Lbm8uPJqh2BW6fwlSXxYH/NRoi170VgEDW1DVbnjBuxvOtp66U5mFlM9MB
DTKWfb6SyIyGwN7AzqZE5F1FdNkZskJQB67kUY3faJ0XdK8yEBczL71TZM9nPVv8
hCmS3SJCCkv/6pt4q6/WdowAYFvae3rIJDr70zivbisqRLr/wiNKwFugSsYHA/6M
4uCV6DdM6tibAFm+C/ezpdlX47rnzzs65yUR0diiSzK5KH/zhEqGGVOquzo8NOUW
7rMs/oBdJlNwDKm26ypUx0wz/1CW4MPezqG8/rRUjT+My+uzoqK2RHZKx8CWZqYy
HsAoog69KjlJ3CbXZaP+4zyZvDYjsu45Zbg+7z2wrljKhthkR8iBy+yEsn/f1D+B
4ZDZgzBhBSpdIydLxNo+q9xNUm3WjkX1ZLJpweNgTMbD0Iysx5C9U2hF2B0jjwRB
t7XgwQbOEHdoyDMTFnvdoGqT6lWZMmtvBlCD7XK1OaPZk6K9ZeAUbVz1H2L1Lc7y
37LmfptkDXVthqPXuLgrFoKRPN4F/sQ/5VJ98cZmrw7mt2utidjxYltu
-----END CERTIFICATE-----
```

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
