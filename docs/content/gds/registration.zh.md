---
title: "注册"
date: 2021-07-22T09:21:59-04:00
lastmod: 2021-10-13T15:18:51-04:00
description: "为VASP注册目录服务"
weight: 20
---

如需加入TRISA或TestNet网络，您必须注册TRISA全球目录服务(GDS)或司法管辖区指定的目录服务。注册目录服务包括两个工作流程:

1. KYV审查过程，以确保网络维持可信任的成员
2. 颁发用于在网络中进行mTLS认证的证书

即将上线: 关于注册表格、邮件验证和审核过程的更多信息。

## 证书颁发

当您的注册通过审核和批准后，目前有两种机制可以从GDS获得mTLS证书。

1. 电子邮件发送的PKCS12加密证书
2. 证书签名请求(CSR)

_当您提交_注册时，您必须选择其中一个选项；在您提交注册后，您将无法更改选项。

### PKCS12加密电子邮件附件

第一种机制是最简单的&mdash; 只需在注册期间选择电子邮件选项，并忽略CSR字段。如果注册表格有效，GDS将返回一个PKCS12密码。**请勿遗忘此密码，它在证书颁发过程中仅显示一次**。

审核通过后，GDS CA将生成包含私钥的完整证书，并使用PKCS12密码对其进行加密。在目录服务中注册公钥之后，GDS将把加密的证书作为ZIP文件通过电子邮件发送给技术联系人，或注册表单上第一个联系人。

解压电子邮件附件后，您会看到一个名为`<common_name>.p12`的文件; 您可以使用如下命令对该文件进行解密以提取证书:

```
$ openssl pkcs12 -in <common_name>.p12 -out <common_name>.pem -nodes
```

您还可以直接使用.zip文件，而无需通过[`github.com/trisacrypto/trisa/pkg/trust`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#NewSerializer)模块进行解密或解压。

### 证书签名请求

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

请仔细填写您的证书的配置，信息必须正确无误，所填信息无法在不重新颁发证书的情况下更改。还要确保配置中的条目后面没有空格!

然后将`example.com`替换为您将使用的TRISA端点的域名，并运行以下命令:

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

您的私钥现在存储在`example.com.key` &mdash; **确保私钥安全** &mdash;它在您mTLS服务中的mTLS连接时必须使用，并将在TRISA网络上建立信任。

`example.com.csr`文件包含您的证书签名请求。复制并粘贴此文件的内容，在您的注册请求中包含`-----BEGIN CERTIFICATE REQUEST-----` 和`-----END CERTIFICATE REQUEST-----`。
