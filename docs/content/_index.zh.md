---
title: 主页
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-13T15:10:04-05:00
description: "TRISA开发人员文档"
weight: 0
---

# TRISA

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/trisa/pkg.svg)](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg)
[![Go Report Card](https://goreportcard.com/badge/github.com/trisacrypto/trisa)](https://goreportcard.com/report/github.com/trisacrypto/trisa)

Travel Rule Information Sharing Architecture (TRISA)的目标是在不修改核心区块链协议、不增加交易成本或修改虚拟货币点对点交易流的情况下，使加密货币交易身份信息符合FATF和FinCEN数据转移规则。TRISA协议和规范由[TRISA工作组](https://trisa.io)定义；欲了解更多有关规范的信息，[请阅读TRISA白皮书的当前版本](https://trisa.io/trisa-whitepaper/)。

本网站包含TRISA协议的开发人员文档和参考实现，可以在[github.com/trisacrypto/trisa](https://github.com/trisacrypto/trisa)查看。TRISA协议被定义为一个[gRPC API](https://grpc.io/)，用于促进那些必须遵守数据转移规则解决方案的虚拟资产服务提供商(VASP)之间的语言无关、高效、点对点的服务。API和消息交换格式都是通过[Protocol Buffers](https://developers.google.com/protocol-buffers)定义的，可以在存储库的[`protos`目录](https://github.com/trisacrypto/trisa/tree/main/proto)中找到。此外，在存储库的[`pkg`目录](https://github.com/trisacrypto/trisa/tree/main/proto)中提供了[Go编程语言](https://golang.org/)中的参考实现。在未来，其他实现将作为特定语言的库代码提供，可查阅存储库的[`lib`目录](https://github.com/trisacrypto/trisa/tree/main/lib)。

TRISA的v1版本正在积极开发中，更多的文档即将发布!
