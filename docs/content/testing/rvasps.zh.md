---
title: 机器人 VASPs
date: 2021-06-14T12:50:23-05:00
lastmod: 2021-10-13T15:27:32-05:00
description: "通过rVASP进行TestNet集成"
weight: 10
---

TestNet提供了三种方便的“robot VASPs”(rVASPs)服务，用以促进与TRISATestNet的集成和测试。这些服务如下:

- Alice (`api.alice.vaspbot.net:443`): 用于触发和接收TRISA消息的主要集成rVASP。
- Bob (`api.bob.vaspbot.net:443`): 一个用于查看与Alice进行交换的演示rVASP。
- Evil (`api.evil.vaspbot.net:443`): 一个非TRISA TestNet会员的“evil”rVASP，用于测试未经验证的交互。

注意: rVASP目前主要是为演示而配置的，为了进行集成，我们已经开始让它们更健壮；如果有任何更改，请经常查看本文档。如果您发现rVASP代码或行为中的任何错误，[请上报问题](https://github.com/trisacrypto/testnet/issues)。

## 开始使用rVASP

有两种方法可以使用rVASP来开发TRISA服务:

1. 您可以触发rVASP向您的服务发送TRISA交换消息。
2. 您可以使用有效(或无效)的rVASP客户向rVASP发送TRISA消息。

rVASP有一个内置的数据库，里面存储着带有假钱包地址的假客户。他们对TRISA消息或对触发转账的响应要求发起者/受益客户对于rVASP有效。例如，如果客户钱包地址是有效的Alice地址，并且Alice是受益人，那么rVASP将用客户的假KYC数据进行响应，但如果是无效的，它将返回TRISA错误码。

下面Alice、Bob和Evil的“客户”列表，可以作为与每个rVASP交互的参考:

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

注意，所有rVASP数据都是使用Faker工具生成的，用于产生真实/一致的测试数据，它们完全是虚构的。例如，Alice VASP(一家假的美国公司)的记录主要是在北美。

如果您是一位Traveler客户，上面粗体的地址有一些属性数据与之关联，并且它们是基于Traveler的rVASP交互的一个很好的候选对象。

### 条件

本文档假设您有一个正在运行最新`TRISANetwork`的服务，并且它已经在TRISA TestNet中注册，并且正确地安装了TestNet证书。参见[TRISA 集成概况]({{% ref "getting-started/_index.md" %}})了解更多信息。**警告**: rVASP不参与TRISA工作网络，他们将只响应已验证的TRISA TestNet mTLS连接。

要与rVASP API交互，您可以:

1. 使用`rvasp` CLI工具
2. 使用rVASP Protocol Buffers并直接与API交互

如需安装CLI工具，您可以前往[TestNet发布页面](https://github.com/trisacrypto/testnet/releases)下载`rvasp`，克隆[TestNet存储库](https://github.com/trisacrypto/testnet/)，并构建`cmd/rvasp`二进制，或者使用如下`go get`命令安装:

```
$ go get github.com/trisacrypto/testnet/...
```

如需使用[rVASP protocol buffers](https://github.com/trisacrypto/testnet/tree/main/proto/rvasp/v1)，请克隆或从TestNet存储库下载，然后使用`protoc`将它们编译成您需要的语言。

### 触发rVASP来发送一条消息

rVASP管理端点用于与rVASP直接交互，以实现开发和集成目的。注意，这个端点与上面描述的TRISA端点不同。

- Alice: `admin.alice.vaspbot.net:443`
- Bob: `admin.bob.vaspbot.net:443`
- Evil: `admin.evil.vaspbot.net:443`

如需使用命令行工具触发一条消息，请运行以下命令:

```
$ rvasp transfer -e admin.alice.vaspbot.net:443 \
        -a mary@alicevasp.us \
        -d 0.3 \
        -B trisa.example.com \
        -b cryptowalletaddress \
        -E
```

此消息使用`-e`或`--endpoint`标志向Alice rVASP发送一条消息，并使用`-a`或`--account`标志指定发起账户应为"mary@alicevasp.us"。发起账户用于确定向受益人发送什么IVMS101数据。`-d`或`--amount`标志指定要发送的“AliceCoin”的数量。

接下来的两个部分至关重要。`-E`或`--external-demo`标志告诉rVASP触发一个服务请求，而不是进行与另一个rVASP的演示交换。必须使用这个标志！最后，`-B`或`--beneficiary-vasp`标志指定rVASP将发送请求的位置。此字段应该能够在TRISA TestNet目录服务中查找；例如，如果可搜索到的话，它应该是您的通用名或VASP的名称。

注意，您可以设置`$RVASP_ADDR`和`$RVASP_CLIENT_ACCOUNT`环境变量来分别指定`-e`和`-a`标志。

如需直接使用Protocol Buffers，使用具有以下 “TransferRequest”的`TRISAIntegration`服务`Transfer` RPC:

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

这些值与命令行程序中的值具有完全相同的规范。

### 向rVASP发送一条TRISA消息

rVASP需要一个`trisa.data.generic.v1beta1.Transaction`作为事务有效负载，以及一个`ivms101.IdentityPayload`作为身份有效负载。身份有效负载受益人信息不需要填充，rVASP将响应填充受益人，但是身份有效负载不应该为空。建议您指定一些假身份数据，以利用rVASP解析和验证命令。

确保在事务有效负载中，指定与上表中rVASP受益人相匹配的受益人钱包；如使用:

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

您可以指定任何`txid`或`originator`字符串，并忽略`network`和`timestamp`字段。

使用目录服务或直接密钥交换来获取rVASP RSA公钥，并使用`AES256-GCM`和`HMAC-SHA256`作为信封加密，创建一个密封的信封。然后使用`TRISANetwork`服务`Transfer` RPC将密封的信封发送到rVASP。

待办: 很快`trisa`命令行程序将面世。我们将说明如何使用CLI程序在消息被释放后立即发送到这里。
