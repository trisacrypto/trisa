---
title: 贡献
date: 2021-06-14T11:34:06-05:00
lastmod: 2022-02-19T13:41:33-05:00
description: "为开源项目做出贡献"
weight: 110
---

TRISA是一个开源项目，欢迎大家做出贡献。

如果您是一名开发人员，您的公司正在使用（或计划采用）TRISA协议，那么这个就是为您准备的！

## 浏览存储库

这个存储库包含了[白皮书](https://trisa.io/trisa-whitepaper/) 中描述的TRISA协议中gRPC的实现，它利用了[protocol buffers](https://grpc.io/) 和Golang。

The `proto` folder contains the core RPC definitions, including:
“proto”文件夹包含核心RPC定义，其中包括：
 - interVASP消息标准（IVMS）消息定义，作为两个VASP双方如何相互描述加密传输中涉及的实体的基础  ，其中包括名称、位置和政府识别码。该规范将允许发起者向受益人表明自己的身份，并要求受益人提供信息来满足监管机构的法律要求。
 - TRISA网络的服务定义，本质上是API的不同部分是如何工作和 &mdash; 从密钥的交换（确保双方都有交换信息所需的详细信息）到“安全信封”的传输（加密密封的protocol buffer信息，只能由预定的接收者解密）。`trisa`子文件夹还包含用于交易的通用消息类型，旨在为广泛的TRISA用例提供最大的灵活性。

`pkg`文件夹包含参考实现代码，包括从`proto`文件夹[^1] 中protocol buffer定义生成的编译代码。
 - `iso3166`文件夹包含语言代码
 - `ivms101`文件夹用JSON加载工具、验证助手、短常数等扩展了生成的protobuf代码。
 - `trisa`文件夹包含了一系列用于TRISA相关任务的结构和方法，如执行加密以及建立mTLS连接。

`lib`文件夹的目的是展示类似于“pkg”文件夹中的实用程序代码，但针对的是Golang以外的语言。如果你的工作语言不是Golang，你可以开始你的贡献啦！

[^1]: 请注意，这些编译的文件是为Golang编译的; 但这肯定不是唯一的选择。 那些对用不同语言构建实现代码感兴趣的人可以看看“lib”文件夹，该文件夹目前包含占位符文件夹，但目的是展示其他实现（包括为这些其他语言编译的protocol buffer代码）。

## 全球目录服务


TRISA协议的另一个组成部分是全球目录服务，它作为TRISA成员的查询工具，用于识别他们希望与之交换信息的双方。关于与全球目录服务有关的RPC定义和实现代码，请访问配套的[`directory` repository](https://github.com/trisacrypto/directory)。有意了解更多关于如何成为目录的成员，请访问[vaspdirectory.net](https://vaspdirectory.net/)。

## 翻译

在trisa.dev文档的翻译是由人工翻译完成的，可能会出现不同步或有个别错误的情况。如果您发现了问题，请开启[错误报告](https://github.com/trisacrypto/trisa/issues/new) 来通知我们