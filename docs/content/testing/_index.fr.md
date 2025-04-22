---
title: Le TestNet
date: 2021-06-14T11:20:10-05:00
lastmod: 2021-10-13T14:50:58-05:00
description: "Le TestNet TRISA"
weight: 30
---

Le TestNet TRISA a été mis en place pour fournir une démonstration du protocole TRISA P2P, héberger des services "robot VASP" pour faciliter l'intégration de TRISA, et est l'emplacement du principal service d'annuaire TRISA qui facilite l'échange de clés publiques et la découverte de pairs.

{{% figure src="/img/testnet_architecture.png" %}}

Le TestNet TRISA est composé des services suivants :

- [TRISA Directory Service](https://trisa.directory) - une interface utilisateur permettant d'explorer le service d'annuaire global TRISA et de se faire enregistrer pour devenir membre de TRISA
- [TestNet Demo](https://vaspbot.com) - un site de démonstration des interactions TRISA entre les VASP "robots" qui fonctionnent dans le TestNet.

Le TestNet héberge également trois robots VASP ou rVASP qui ont été mis en place pour permettre aux membres de TRISA d'intégrer leurs services TRISA. Le premier rVASP est Alice, le second est Bob et pour tester les interactions avec des membres TRISA non vérifiés, il existe également un rVASP "maléfique".
