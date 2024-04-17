---
title: 集成概述
date: 2021-04-23T01:35:35-04:00
lastmod: 2021-10-13T15:23:16-04:00
description: "描述如何在TestNet中集成TRISA协议"
weight: 5
---

## TRISA集成概述

1.注册TRISA目录服务
2.实现TRISA网络协议
3.实现TRISA健康协议

## VASP目录服务注册

### 注册概述

在将TRISA协议集成到VASP软件之前，必须注册TRISA目录服务(DS)。TRISA DS为注册的VASP提供公钥和TRISA远程对等连接信息。

注册TRISA DS后，您将收到KYV证书。KYV证书中的公钥将通过TRISA DS提供给其他VASP。

在注册DS时，您将需要提供您的VASP实施TRISA网络服务的`address:port`端点。当您的VASP被确定为受益VASP时，该地址将被用于注册DS并由其他VASP使用。

出于集成目的，可以使用托管的TestNet TRISA DS实例进行测试。在TestNet中，注册过程非常简化，可便于快速集成。建议在与TestNet集成的同时开始DS注册。


### 目录服务注册

如需注册TRISA DS，请访问[https://vaspdirectory.net/](https://vaspdirectory.net/)。

您可以选择“注册/Register”标签开始注册。请注意，您可以使用本网站逐项输入您的注册详细信息，也可以上传包含注册详细信息的JSON文档。

此注册方式将发送电子邮件给JSON文件中指定的所有技术联系人。后续将通过电子邮件指导完成注册。一旦您完成注册，TRISA TestNet管理员将收到您的注册并进行审核。

TestNet管理员审核并批准注册后，您将通过电子邮件收到KYV证书，您的VASP将在TestNet DS中公开可见。


## 实现TRISA P2P协议


### 先决条件

开始设置前，您需要以下信息:



*   KYV证书(来自TRISA DS注册)
*   用于获取证书的CSR的公钥
*   关联的私钥
*   TRISA目录服务的主机名
*   绑定到TRISA目录服务中与VASP相关联的address:port的能力。


### 集成概述

集成TRISA协议涉及到客户端组件和服务器组件。

客户端组件将与TRISA目录服务(DS)实例进行接口，以查找集成TRISA消息传递协议的其他VASP。客户端组件用于从您的VASP发出事务，以验证接收VASP是否符合TRISA。

服务器组件接收来自集成TRISA协议的其他VASP的请求，并响应它们的请求。服务器组件提供了必须执行的回调，以便您的VASP能够返回满足TRISA网络协议的信息。

目前，在Go中可以找到TRISA网络协议的参考实现。

[https://github.com/trisacrypto/testnet/blob/main/pkg/rVASP/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rVASP/trisa.go)

集成VASP必须运行它们自己的协议实现。如果需要Go之外的语言，则可以从定义TRISA 网络协议的Protocol Buffers生成客户端库。

集成商应该集成传入的传送请求和密钥交换，也可以选择性地集成发出的传送请求和密钥交换。

### 集成注意事项

TRISA网络协议定义了数据如何在参与的VASP之间传输。用于识别信息而传输的数据的推荐格式，建议采用IVMS101。VASP有责任确保发送/接收的识别数据满足金融行动特别工作组数据转移规则(FATF Travel Rule)。

成功的TRISA交易将产生满足金融行动特别工作组数据转移规则的密钥和加密数据。TRISA并没有规定一旦获得这些数据应该如何存储。VASP有责任确保事务产生的数据的安全存储。

