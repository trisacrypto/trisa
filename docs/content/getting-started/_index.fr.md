---
title: L'intégration
date: 2021-04-23T01:35:35-04:00
lastmod: 2021-10-13T14:47:18-04:00
description: "Décrit comment intégrer le protocole TRISA dans le TestNet."
weight: 5
---

## Aperçu de l'intégration TRISA

1. S'inscrire auprès d'un service d'annuaire TRISA
2. Mettre en œuvre le protocole Réseau TRISA
3. Mettre en œuvre le protocole de Santé TRISA

## Enregistrement du service d'annuaire VASP (Fournisseur de Services d'Actifs Virtuels)

### Présentation de l'enregistrement

Avant de pouvoir intégrer le protocole TRISA dans votre logiciel VASP, vous devez vous faire enregistrer dans un service d'annuaire TRISA (DS).  Le DS TRISA fournit la clé publique ainsi que les informations de connexion aux pairs TRISA à distants pour les VASP enregistrés.

Une fois que vous vous êtes enregistré auprès du DS TRISA, vous recevez un certificat KYV.  La clé publique du certificat KYV (Know Your sera mise à la disposition des autres VASP via le TRISA DS.

Lors de l'enregistrement auprès du DS, vous devrez fournir les éléments `address:port` du terminal où votre VASP exécute le service Réseau TRISA. Cette adresse sera enregistrée auprès du DS et utilisée par les autres VASP lorsque votre VASP sera identifié comme étant le VASP bénéficiaire.

À des fins d'intégration, une instance hébergée de TestNet TRISA DS est disponible pour des tests.  Le processus d'enregistrement est simplifié dans le TestNet pour favoriser une intégration rapide.  Il est recommandé de commencer l'enregistrement du DS de production lors de l'intégration avec le TestNet.


### Enregistrement au Service d'Annuaire

Pour commencer l'enregistrement avec le DS de TRISA, visitez le site web à l'adresse suivante [https://vaspdirectory.net/](https://vaspdirectory.net/)

Vous pouvez sélectionner l'onglet "Enregistrer" pour commencer l'enregistrement. Notez que vous pouvez utiliser ce site web pour saisir les détails de votre enregistrement, champ par champ, ou pour télécharger un document JSON contenant les détails de votre enregistrement.

L'enregistrement aboutira à l'envoi d'un courriel à tous les contacts techniques spécifiés dans le fichier JSON.  Ces courriels vous guideront dans la suite de la procédure d'enregistrement.  Une fois que vous aurez terminé les étapes d'enregistrement, les administrateurs TestNet de TRISA recevront votre enregistrement pour examen.

Une fois que les administrateurs de TestNet auront examiné et approuvé l'enregistrement, vous recevrez un certificat KYV par courrier électronique et votre VASP sera publiquement visible dans le DS TestNet.


## Implementation du Protocole P2P TRISA


### Prérequis

Pour commencer l'installation, vous aurez besoin des éléments suivants :

*   Le certificat KYV (de l’enregistrement au DS de TRISA)
*   La clé publique utilisée pour la Demande de Signature de Certificat (Certificate Signing Request - CSR) afin d'obtenir votre certificat.
*   La clé privée associée
*   Le nom d'hôte du service d'annuaire TRISA
*   Possibilité de lier l’address:port associé à votre VASP dans le service d'annuaire TRISA.


### Présentation de l'intégration

L'intégration du protocole TRISA implique à la fois un composant client et un composant serveur.

Le client s'interface avec une instance du Service d'annuaire TRISA (DS) pour rechercher d'autres VASP qui intègrent le protocole de messagerie TRISA.  Le client est utilisé pour les transactions sortantes de votre VASP afin de vérifier que le VASP destinataire est conforme à TRISA.

Quant au serveur, il reçoit les demandes d'autres VASP qui intègrent le protocole TRISA et fournit des réponses à leurs demandes.  Le serveur émet des rappels qui doivent être implémentés pour que votre VASP puisse renvoyer des informations conformes au protocole Réseau TRISA.

Actuellement, une implémentation de référence du protocole Réseau TRISA est disponible en Go.

[https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go)

Les VASP intégrants doivent exécuter leur propre mise en œuvre du protocole.  Si un langage autre que Go est requis, des bibliothèques client peuvent être générées à partir des tampons de protocole qui définissent le protocole du réseau TRISA.

Ils devraient intégrer les demandes entrantes de transfert et d'échange de clés et peuvent, en option, intégrer les demandes sortantes de transfert et d'échange de clés.

### Notes sur l'intégration

Le protocole Réseau TRISA définit la manière dont les données sont transférées entre les VASP participants.  Le format recommandé pour les données transférées pour les informations d'identification est le format de données IVMS101.  Il est de la responsabilité du VASP qui l'implémente de s'assurer que les données d'identification envoyées/reçues sont conformes à la règle de voyage du GAFI.

Le résultat d'une transaction TRISA réussie est une clé et des données cryptées qui répondent à la règle de voyage du GAFI.  TRISA ne définit pas comment ces données doivent être stockées une fois obtenues.  Il incombe au VASP chargé de la mise en œuvre de gérer le stockage sécurisé des données résultant de la transaction.

