---
title: Contribution
date: 2021-06-14T11:34:06-05:00
lastmod: 2022-02-21T09:33:45-05:00
description: "Contribution au projet Open Source"
weight: 110
---

TRISA est un projet open source et donc, les contributions sont les bienvenues !

Si vous êtes un développeur dont l'organisation utilise (ou prévoit d'adopter) le protocole TRISA, cette section est pour vous !

## Navigation dans le référentiel

Ce référentiel contient une implémentation gRPC du protocole TRISA tel que décrit dans le [livre blanc](https://trisa.io/trisa-whitepaper/), qui tire parti des [protocol buffers](https://grpc.io/) et du Golang.

Le dossier `proto` contient les définitions RPC de base, y compris :
 -  les définitions des messages de la norme de messagerie interVASP (IVMS), qui constituent la base de la manière dont deux pairs VASP doivent décrire mutuellement les entités impliquées dans les transferts cryptographiques, notamment les noms, les emplacements et les identifiants gouvernementaux. C'est la spécification qui permettra aux expéditeurs de s'identifier auprès des bénéficiaires et de demander des informations à ces derniers pour répondre aux exigences légales de leurs organes de réglementation.
 - les définitions des services du réseau TRISA, qui expliquent essentiellement le fonctionnement des différentes parties de l'API &mdash; de l'échange de clés ( afin de s'assurer que les deux homologues disposent des informations nécessaires à l'échange) au transfert des « enveloppes sécurisées » (messages de type protocol buffer cryptographiquement scellés et qui ne peuvent être déchiffrés que par le destinataire prévu). Le sous-dossier `trisa` contient également des types de messages génériques à utiliser pour les transactions visant à fournir une flexibilité maximale pour un large éventail de cas d'utilisation de TRISA.

Le dossier `pkg` contient le code d'implémentation de référence, y compris le code compilé généré à partir des définitions du protocol buffer dans le dossier `proto`[^1].
 - Le dossier `iso3166` contient les codes du langage.
 - Le dossier `ivms101` étend le code protobuf généré avec des utilitaires de chargement JSON, des aides à la validation, des constantes de longueurs courtes, etc...
 - Le dossier `trisa` contient des structures et des méthodes pour une série de tâches liées à TRISA, comme la cryptographie et les connexions mTLS.

 Le dossier `lib` présente un code utilitaire similaire à celui du dossier `pkg` mais pour des langages autres que Golang. Si vous travaillez dans un langage autre que Golang, ce serait un bon endroit pour commencer votre contribution !

[^1]: Veuillez noter que ces fichiers compilés sont réalisés pour le langage Golang ; mais ce n'est certainement pas la seule option. Ceux qui sont intéressés par la construction de code d'implémentation dans un langage autre devraient consulter le dossier `lib` qui contient actuellement des dossiers placeholder mais qui est destiné à présenter de telles autres implémentations (y compris le code compilé du protocol buffer pour ces autres langages).

## Le Service d'Annuaire Global

Une autre partie intégrante du protocole TRISA est le service d'annuaire global, qui sert d'outil de recherche aux membres de TRISA pour identifier les pairs avec lesquels ils souhaitent échanger des informations. Pour obtenir les définitions des RPC et le code d'implémentation relatifs au service d'annuaire mondial, consultez le référentiel [directory repository](https://github.com/trisacrypto/directory). Pour en savoir plus sur la façon de devenir membre de l'annuaire, consultez [vaspdirectory.net](https://vaspdirectory.net/).

## Les Traductions

Les traductions de la documentation sur trisa.dev sont effectuées périodiquement par des traducteurs humains, et peuvent être désynchronisées ou contenir des erreurs. Si vous remarquez une erreur, veuillez ouvrir une [demande de bogue](https://github.com/trisacrypto/trisa/issues/new) pour nous en informer.
