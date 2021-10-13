---
title: Les Robot VASPs
date: 2021-06-14T12:50:23-05:00
lastmod: 2021-10-13T14:52:08-05:00
description: "Travailler avec les rVASP pour l'intégration TestNet"
weight: 10
---

Le TestNet héberge trois services pratiques de "robot VASP" (rVASP) pour faciliter l'intégration et les tests avec le TestNet de TRISA. Ces services sont les suivants :

- Alice (`api.alice.vaspbot.net:443`) : le rVASP d'intégration primaire utilisé pour déclencher et recevoir les messages TRISA.
- Bob (`api.bob.vaspbot.net:443`) : une démo rVASP pour visualiser les échanges avec Alice.
- Evil (`api.evil.vaspbot.net:443`) : un rVASP "maléfique" qui n'est pas membre du TestNet de TRISA, utilisé pour tester les interactions non authentifiées.

Note : les rVASP sont actuellement configurés principalement pour des démonstrations et des travaux visant à les rendre plus robustes à des fins d'intégration ont été entamés ; veuillez consulter régulièrement cette documentation pour d'éventuels changements. Si vous remarquez des bogues dans le code ou le comportement des rVASP, [veuillez créer un ticket] (https://github.com/trisacrypto/testnet/issues).

## Démarrer avec les rVASP

Il y a deux façons d'utiliser les rVASP pour développer votre service TRISA :

1. Vous pouvez déclencher le rVASP pour envoyer un message d'échange TRISA à votre service.
2. Vous pouvez envoyer un message TRISA au rVASP avec un client rVASP valide (ou invalide).

Les rVASP ont une base de données intégrée de faux clients avec un faux portefeuille d'adresses. Leur réponse aux messages TRISA ou à un transfert déclenché exige que le client donneur d'ordre/bénéficiaire soit valide pour le rVASP. Par exemple, si l'adresse du portefeuille client est une adresse d'Alice valide et qu'Alice est le bénéficiaire, le rVASP répondra avec les fausses données KYC du client, mais dans le cas contraire, il renverra un code d'erreur TRISA.

Le tableau suivant des "clients" pour Alice, Bob et le Maléfique peut être utilisé comme référence pour interagir avec chaque rVASP :

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

Notez que toutes les données rVASP ont été générées à l'aide d'un outil appelé Faker pour produire des données de test et des montages réalistes/consistants et sont totalement fictives. Par exemple, les enregistrements pour Alice VASP (une société fictive américaine) se trouvent principalement en Amérique du Nord, etc.

Si vous êtes un client de Traveler, les adresses en gras ci-dessus sont associées à des données d'attribution et constituent un bon candidat pour les interactions rVASP basées sur Traveler.

### Préliminaires

Cette documentation suppose que vous avez un service qui exécute la dernière version du service `TRISANetwork` et qu'il a été enregistré dans le TestNet TRISA et que les certificats TestNet sont correctement installés. Voir [Vue d'ensemble de l'intégration de TRISA]({{< ref "integration/_index.md" >}}) pour plus d’informations. **AVERTISSEMENT**: les rVASP ne participent pas au réseau de production TRISA, ils ne répondent qu'aux connexions mTLS vérifiées de TRISA TestNet.

Pour interagir avec l'API rVASP, vous pouvez soit :

1. Utiliser l’outil d’interface de commandes `rvasp`
2. Utiliser les tampons du protocole rVASP et interagir directement avec l'API

Pour installer l'outil CLI, vous pouvez soit télécharger l'exécutable `rvasp` pour l'architecture appropriée à la [Page des versions TestNet](https://github.com/trisacrypto/testnet/releases), cloner [le répertoire TestNet](https://github.com/trisacrypto/testnet/) et construire le binaire `cmd/rvasp` ou installer avec `go get` comme suit :

```
$ go get github.com/trisacrypto/testnet/...
```

Pour utiliser le [tampon du protocole rVASP](https://github.com/trisacrypto/testnet/tree/main/proto/rvasp/v1), clonez ou téléchargez-les depuis le répertoire TestNet, puis compilez-les dans votre langue préférée en utilisant `protoc`.

### Déclencher l'envoi d'un message par un rVASP

Les terminaisons administrateurs rVASP sont utilisées pour interagir directement avec le rVASP à des fins de développement et d'intégration. Notez que ce point de terminaison est différent du point de terminaison TRISA, qui a été décrit ci-dessus.

- Alice: `admin.alice.vaspbot.net:443`
- Bob: `admin.bob.vaspbot.net:443`
- Evil: `admin.evil.vaspbot.net:443`

Pour utiliser l'outil de ligne de commande afin de déclencher un message, exécutez la commande suivante :

```
$ rvasp transfer -e admin.alice.vaspbot.net:443 \
        -a mary@alicevasp.us \
        -d 0.3 \
        -B trisa.example.com \
        -E
```

Ce message envoie un message au rVASP Alice par le champ `-e` ou `--endpoint`, et spécifie que le compte émetteur doit être "mary@alicevasp.us" par le champ `-a` ou `--account`. Le compte émetteur est utilisé pour déterminer les données IVMS101 à envoyer au bénéficiaire. Le champ `-d` ou `--amount` spécifie le montant d'"AliceCoin" à envoyer.

Les deux parties suivantes sont critiques. Le champ `-E` ou `--external-demo` indique au rVASP de déclencher une requête vers votre service plutôt que d'effectuer un échange de démo avec un autre rVASP. Ce champ est obligatoire ! Enfin, le champ `-B` ou `--beneficiary-vasp` indique où le rVASP enverra la requête. Ce champ doit pouvoir être consulté dans le service d'annuaire TestNet de TRISA ; par exemple, il doit s'agir de votre nom commun ou du nom de votre VASP s'il est consultable.

Notez que vous pouvez définir les variables d'environnement `$RVASP_ADDR` et `$RVASP_CLIENT_ACCOUNT` pour spécifier les champs `-e` et `-a` respectivement.

Pour utiliser directement les tampons de protocole, utilisez le `TRISAIntegration` service `Transfer` RPC avec la suivante `TransferRequest`:

```json
{
    "account": "mary@alicevasp.us",
    "amount": 0.3,
    "beneficiary_vasp": "trisa.example.com",
    "check_beneficiary": false,
    "external_demo": true
}
```

Ces valeurs ont la même spécification exacte que celles du CLI.

### Envoi d'un message TRISA à un rVASP

Le rVASP attend une `trisa.data.generic.v1beta1.Transaction` comme charge utile de la transaction et une `ivms101.IdentityPayload` comme charge utile de l'identité. Il n'est pas nécessaire de renseigner les informations de l'identité du bénéficiaire, le rVASP répondra en renseignant le bénéficiaire, mais l'identité du bénéficiaire ne doit pas être nulle. Il est recommandé de spécifier des données d'identité fictives pour tirer parti des commandes d'analyse et de validation du rVASP.

Assurez-vous que dans les données utiles de votre transaction, vous spécifiez un portefeuille bénéficiaire qui correspond aux bénéficiaires rVASP du tableau ci-dessus ; par exemple, utilisez :

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

Vous pouvez préciser n'importe quelle chaîne `txid` ou `originator` et les champs `network` et `timestamp` sont ignorés.

Créez une enveloppe scellée en utilisant le service d'annuaire ou l'échange direct de clés pour récupérer les clés publiques RSA de rVASP et en utilisant `AES256-GCM` et `HMAC-SHA256` comme cryptographie de l'enveloppe. Ensuite, utilisez le `TRISANetwork` service `Transfer` RPC pour envoyer l'enveloppe scellée au rVASP.

A FAIRE : Bientôt le programme de ligne de commande `trisa` sera disponible. Nous spécifierons ici comment utiliser le programme CLI pour envoyer un message dès qu'il sera disponible.

