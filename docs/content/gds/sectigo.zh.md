---
title: "Sectigo"
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-13T15:16:54-05:00
description: "目录服务与Sectigo CA API的交互"
weight: 50
---

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/directory/pkg/sectigo.svg)](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)

TRISA目录服务通过IoT门户使用Sectigo证书颁发机构颁发证书。由于目录服务必须收集公钥信息，以推进mTLS的初始信任握手，因此它使用Sectigo IoT Manager API作为VASP注册和验证过程的一部分。`github.com/trisacrypto/directory/pkg/sectigo`包是一个用于与API交互以及实现目录服务所需的端点和方法的Go库。TestNet还提供了一个命令行公用程序，用于与API进行交互，以实现管理和调试目的。本文档描述了命令行公用程序，还概述了如何直接使用API来颁发和撤销证书。

参考资料:

- [软件包文档](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)
- [IoT Manager API 文档](https://support.sectigo.com/Com_KnowledgeDetailPage?Id=kA01N000000bvCJ)
- [IoT Manager 门户网站](https://iot.sectigo.com)

## 开始

如需安装`sectigo` CLI实用程序，可从[GitHub发布的文件](https://github.com/trisacrypto/directory/releases)下载预编译的二进制文件，或者使用以下命令在本地安装:

```
$ go get github.com/trisacrypto/directory/cmd/sectigo
```

此操作将添加`sectigo`命令到您的`$PATH`。

## 身份验证

第一步是身份验证，您应该在`$SECTIGO_USERNAME`和`$SECTIGO_PASSWORD`环境变量中设置您的用户名和密码(或者您可以在命令行中将它们作为参数传递)。如需验证您的身份验证状态，您可以使用以下命令:

```
$ sectigo auth
```

API通过用户名和密码进行身份验证，然后返回访问权限以及存储在本地缓存文件中的刷新令牌。查看缓存的存储位置:

```
$ sectigo auth --cache
```

如果您想检查凭据状态，例如访问令牌是否有效、可刷新或已过期，请使用:

```
$ sectigo auth --debug
```

## 权限和配置文件

开始与证书交互前，您需要列出您的用户账户可以访问的权限和配置文件。

```
$ sectigo authorities
[
  {
    "id": 1,
    "ecosystemId": 100,
    "signerCertificateId": 0,
    "ecosystemName": "TRISA",
    "balance": 10,
    "enabled": true,
    "profileId": 42,
    "profileName": "TRISA Profile"
  }
]
```

权限显示创建证书所依附的方法和概要文件。这里的`profileId`字段对于在后续调用中的使用非常重要。您还可以通过如下命令查看权限中已经订购/颁发了多少许可证:

```
$ sectigo licenses
{
  "ordered": 2,
  "issued": 2
}
```

如需获取配置文件的详细信息，请通过以下命令来使用profileId:

```
$ sectigo profiles -i 42
```

此操作将返回原始配置文件的配置。在使用权限创建证书之前，您需要知道所需的配置文件参数:

```
$ sectigo profile -i 42 --params
```

## 创建证书

您可以请求用`commonName`和`pkcs12Password`参数创建一个证书，如下所示(注意，对于需要其他参数的配置文件，您必须直接使用代码库并实现自己的方法):

```
$ sectigo create -a 42 -d example.com -p secrtpasswrd -b "example.com certs"
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

`-a`标志指定权限，但它应该是配置文件id。域名必须有效。如果您没有指定密码，则会为您生成一个密码，并在退出之前显示在CLI中。`-b`标志为批处理创建提供了一个人类可读的名称。返回数据显示创建的批处理证书作业的详细信息；您可以使用以下命令获取数据，以继续查看状态:

```
$ sectigo batches -i 24
```

您还可以获取批处理信息:

```
$ sectigo batches -i 24 --status
```

一旦创建批处理，就该下载ZIP文件中的证书了:

```
$ sectigo download -i 24 -o certs/
```

此操作将下载批处理文件(通常为batchId.zip，在本例中是24.zip)到`certs/`目录。解压缩certs，然后通过以下命令对.pem文件进行解密:

```
$ unzip certs/24.zip
$ openssl pkcs12 -in certs/example.com.p12 -out certs/example.com.pem -nodes
```

有关使用PKCS12文件的更多信息，请参见[使用OpenSSL从PKCS#12文件导出证书和私钥](https://www.ssl.com/how-to/export-certificates-private-key-from-pkcs12-file-with-openssl/)。

## 上传CSR

创建证书的另一种方法是上传证书签名请求(CSR)。这种机制通常更可取，因为它意味着不需要通过网络传输私钥信息，而且私钥可以保存在安全的硬件中。

如需在命令行中使用`openssl`生成CSR，首先在您的当前工作目录中创建一个名为`trisa.conf`的配置文件，通过以下命令用您打算托管TRISA端点的域名替换`example.com`:

```conf
[req]
distinguished_name = dn_req
req_extensions = v3ext_req
prompt = no
default_bits = 4096
[dn_req]
CN = example.com
O = [Organization]
L = [City]
ST = [State or Province (完整拼出，勿使用缩写)]
C = [2 digit country code]
[v3ext_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = example.com
```

请仔细填写您的证书配置，信息必须正确无误，所填信息无法在不重新颁发证书的情况下更改。还要确保配置中的条目后面没有空格!

然后将`example.com`替换为您将使用的TRISA端点的域名，并运行以下命令:

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

然后可以通过如下命令使用CLI程序上传CSR:

```
$ sectigo upload -p 42 <common_name>.csr
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

`-p`标志指定使用CSR批处理请求的配置文件，并且必须是一个有效的profileId。上传的CSR可以是一个文本文件，包含多个使用标准BEGIN/END分隔符的PEM格式的CSR。

## 管理证书

您可以按名称或序列号搜索证书，但通常情况下，您按域名或通用名称搜索来获得序列号:

```
$ sectigo find -n example.com
```

一旦您获得序列号，您可以通过以下命令撤销证书:

```
$ sectigo revoke -p 42 -r "cessation of operation" -s 12345
```

此命令需要使用`-p`标志颁发证书的配置文件id，一个通过`-r`标志传递(默认未指定)的[RFC 5280 原因码](https://tools.ietf.org/html/rfc5280#section-5.3.1)，以及使用`-s`标志的证书序列号。如果此命令没有错误，则证书已被成功撤销。

RFC 5280原因包括:

- "unspecified"
- "keycompromise"
- "ca compromise"
- "affiliation changed"
- "superseded"
- "cessation of operation"
- "certificate hold"
- "remove from crl"
- "privilege withdrawn"
- "aa compromise"

请注意，原因不区分空格和大小写。
