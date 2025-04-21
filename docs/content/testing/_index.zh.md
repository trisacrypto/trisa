---
title: TestNet
date: 2021-06-14T11:20:10-05:00
lastmod: 2021-10-13T15:26:40-05:00
description: "TRISA TestNet"
weight: 30
---

创建TRISA TestNet是为了提供TRISA点对点协议的演示，托管“robot VASP”服务，以促进TRISA集成，并且它是促进公钥交换和对等体发现的TRISA主目录服务的所在之处。

{{% figure src="/img/testnet_architecture.png" %}}

TRISA TestNet由以下服务组成:

- [TRISA Directory Service](https://trisa.directory) - 一个探索TRISA全球目录服务并注册成为TRISA会员的用户界面
- [TestNet Demo](https://vaspbot.com) - 一个演示网站，展示运行在TestNet中的“机器人”VASP之间的TRISA交互

TestNet还有三个robot VASP（简称“rVASP”），为TRISA会员集成他们的TRISA服务提供了便利。主要的rVASP是Alice，用于演示的次要rVASP是Bob，为了测试与未经验证的TRISA会员之间的交互，还有一个“Evil” rVASP。
