---
title: "Sectigo"
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-13T14:41:50-05:00
description: "Interactions des services d'annuaire avec l'API du CA de Sectigo"
weight: 50
---

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/directory/pkg/sectigo.svg)](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)

Le service d'annuaire TRISA émet des certificats en utilisant l'autorité de certification Sectigo via son portail IdO. Comme le service d'annuaire doit collecter des clés publiques afin de faciliter la création d’une liaison de confiance initiale pour mTLS, il utilise le gestionnaire IdO de l'API de Sectigo dans le cadre du processus d'enregistrement et de vérification VASP. Le paquet `github.com/trisacrypto/directory/pkg/sectigo` est une bibliothèque Go pour interagir avec l'API, implémentant les terminaisons ainsi que les méthodes requises par le service d'annuaire. Le TestNet fournit également un utilitaire en ligne de commande pour interagir avec l'API à des fins d'administration et de débogage. Cette documentation décrit l'utilitaire de ligne de commande, qui donne également un aperçu de la façon d'utiliser l'API directement pour émettre et révoquer des certificats.

Notes de référence :

- [Documentation des paquets](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)
- [Documentation du gestionnaire IdO de l'API](https://support.sectigo.com/Com_KnowledgeDetailPage?Id=kA01N000000bvCJ)
- [Portail du Gestionnaire IdO](https://iot.sectigo.com)

## Démarrage

Pour installer l’utilitaire CLI de `sectigo`, téléchargez un binaire précompilé à partir du site [disponible sur GitHub](https://github.com/trisacrypto/directory/releases), ou bien installez le localement à l'aide du programme:

```
$ go get github.com/trisacrypto/directory/cmd/sectigo
```

Ceci ajoutera la commande `sectigo` à votre `$PATH`.

## Authentification

La première étape est l'authentification, vous devez définir votre nom d'utilisateur et votre mot de passe dans les champs variables `$SECTIGO_USERNAME` et `$SECTIGO_PASSWORD` (alternativement les passer en paramètre de ligne de commande). Pour vérifier votre statut d'authentification, vous pouvez utiliser :

```
$ sectigo auth
```

L'API authentifie par nom d'utilisateur et mot de passe, puis renvoie des jetons d'accès et d'actualisation qui sont stockés dans un fichier cache local. Pour savoir où est stocké votre cache :

```
$ sectigo auth --cache
```

Si vous souhaitez vérifier l'état de vos informations d'identification, par exemple si les jetons d'accès sont valides, actualisables ou expirés, utilisez :

```
$ sectigo auth --debug
```

## Autorités et profils

Pour commencer à interagir avec les certificats, vous devez répertorier les autorités et les profils auxquels votre compte utilisateur a accès.

```
$ sectigo authorities
[
  {
    "id": 1,
    "ecosystemId": 100,
    "signerCertificateId": 0,
    "ecosystemName": "TRISA",
    "balance": 10,
    "enabled": true,
    "profileId": 42,
    "profileName": "TRISA Profile"
  }
]
```

L'autorité affiche les méthodes et les profils par lesquels les certificats sont créés. Ici, le champ `profileId` est très important pour être utilisé dans les requêtes ultérieures. Vous pouvez également voir combien de licences ont été commandées/émises pour toutes les autorités comme suit :

```
$ sectigo licenses
{
  "ordered": 2,
  "issued": 2
}
```

Pour obtenir des informations détaillées sur un profil, utilisez le profileId avec la commande suivante :

```
$ sectigo profiles -i 42
```

Cela renverra la configuration brute du profil. Avant de créer des certificats avec l'autorité, vous devez connaître les paramètres de profil requis :

```
$ sectigo profile -i 42 --params
```

## Création des Certificats

Vous pouvez demander la création d'un certificat avec les paramètres `commonName` et `pkcs12Password` comme suit (note : pour les profils nécessitant d'autres paramètres, vous devrez utiliser directement la base de code et implémenter votre propre méthode) :

```
$ sectigo create -a 42 -d example.com -p secrtpasswrd -b "example.com certs"
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

Le champ `-a` spécifie l'autorité, mais doit être un id de profil. Le domaine doit être un domaine valide. Si vous ne spécifiez pas de mot de passe, un mot de passe est généré pour vous et imprimé sur le CLI avant la sortie. Le champ `-b` donne un nom lisible par l'homme pour la création des lots de certificats. Les données retournées donnent des détails sur la tâche de création des lots ; vous pouvez récupérer les données pour continuer à vérifier l'état comme suit :

```
$ sectigo batches -i 24
```

Vous pouvez également obtenir des informations sur le traitement du lot :

```
$ sectigo batches -i 24 --status
```

Une fois le lot créé, il ne reste plus qu'à télécharger les certificats dans un fichier ZIP :

```
$ sectigo download -i 24 -o certs/
```

Cela téléchargera le lot de fichiers (généralement batchId.zip, 24.zip dans ce cas) vers le répertoire `certs/`. Décompressez les certificats puis décryptez le fichier .pem comme suit :

```
$ unzip certs/24.zip
$ openssl pkcs12 -in certs/example.com.p12 -out certs/example.com.pem -nodes
```

Pour plus d'informations relatives au fichier PKCS12, voir [Exportation des certificats et clés privées à partir d'un fichier PKCS#12 avec OpenSSL](https://www.ssl.com/how-to/export-certificates-private-key-from-pkcs12-file-with-openssl/).

## Chargement d'un CSR

Une alternative à la création d'un certificat consiste à uploader une demande de signature de certificat (CSR). Ce mécanisme est souvent préférable car il signifie qu'aucun élément de clé privée ne doit être transmis sur le réseau et que la clé privée peut rester sur un support matériel sécurisé.

Pour générer un CSR en utilisant `openssl` en CLI, créez tout d’abord un fichier de configuration nommé `trisa.conf` dans votre répertoire courant en remplaçant `example.com` par le nom de domaine dans lequel vous souhaitez héberger votre terminal TRISA :

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
ST = [State or Province (en toutes lettres, sans abréviations)]
C = [2 digit country code]
[v3ext_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = example.com
```

Veuillez remplir soigneusement la configuration de votre certificat, ces informations doivent être correctes et ne peuvent être modifiées sans réémettre le certificat. Assurez-vous également qu'il n'y a pas d'espace après les données de la configuration !

Puis exécutez la commande suivante, en remplaçant `example.com` par le nom de domaine que vous utiliserez comme terminaison TRISA :

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

Vous pouvez ensuite charger le CSR via le CLI comme suit :

```
$ sectigo upload -p 42 <common_name>.csr
{
  "batchId": 24,
  "orderNumber": 1024,
  "creationDate": "2020-12-10T16:35:32.805+0000",
  "profile": "TRISA Profile",
  "size": 1,
  "status": "CREATED",
  "active": false,
  "batchName": "example.com certs",
  "rejectReason": "",
  "generatorParametersValues": null,
  "userId": 10,
  "downloadable": true,
  "rejectable": true
}
```

Le champ `-p` spécifie le profil à utiliser pour la requête de lot de CSR et doit être un ID de profil valide. Les CSR chargés peuvent être un fichier texte unique contenant plusieurs CSR au format PEM avec des séparateurs standards BEGIN/END.

## Gestion des Certificats

Vous pouvez rechercher un certificat par nom ou par numéro de série, mais le plus souvent on recherche par domaine ou par nom commun pour obtenir le numéro de série :

```
$ sectigo find -n example.com
```

Une fois le numéro de série obtenu, vous pouvez révoquer le certificat comme suit :

```
$ sectigo revoke -p 42 -r "cessation of operation" -s 12345
```

Cette commande vérifie l’ID du profil qui a émis le certificat avec le champ `-p` le [RFC 5280 code du motif](https://tools.ietf.org/html/rfc5280#section-5.3.1) transmis via le champ `-r` (non spécifié par défaut) et le numéro de série du certificat en utilisant le champ `-s`. Si cette commande ne présente pas d'erreur, alors le certificat a été révoqué avec succès.

Les motifs du RFC 5280 :

- "unspecified"
- "keycompromise"
- "ca compromise"
- "affiliation changed"
- "superseded"
- "cessation of operation"
- "certificate hold"
- "remove from crl"
- "privilege withdrawn"
- "aa compromise"

Notez que le motif est insensible aux espaces et à la casse.
