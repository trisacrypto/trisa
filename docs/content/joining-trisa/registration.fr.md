---
title: "Enregistrement"
date: 2021-07-22T09:21:59-04:00
lastmod: 2021-10-13T14:43:51-04:00
description: "Enregistrement d'un fournisseur de services d'actifs virtuels (VASP) grâce au service d’annuaire"
weight: 20
---

Pour rejoindre les réseaux TRISA ou TestNet, vous devez vous inscrire auprès du service d'annuaire mondial (GDS) de TRISA ou de l'un des services d'annuaire propres à chaque juridiction. L'enregistrement auprès du service d'annuaire implique deux flux de travail :

1. Un processus de révision de la Connaissance Client (KYV) pour s'assurer que le réseau conserve une adhésion de confiance.
2. La délivrance d'un certificat pour l'authentification mTLS dans le réseau.

À venir : plus de détails sur le formulaire d'inscription, la vérification par e-mail et le processus de révision.

## Délivrance des certificats

Il existe actuellement deux mécanismes pour recevoir des certificats mTLS du GDS lorsque votre inscription a été examinée et approuvée.

1. Certificats cryptés PKCS12 envoyés par e-mail.
2. Demande de signature de certificat (Certificate Signing Request - CSR)

Vous devez sélectionner l'une de ces options _lorsque vous soumettez_ votre inscription ; après la soumission de votre inscription, vous ne pourrez pas passer d'une option à l'autre.

### Pièce jointe du courriel crypté PKCS12

Le premier mécanisme est le plus simple &mdash; il suffit de sélectionner l'option e-mail lors de l'inscription et d'omettre les champs CSR. Si le formulaire d'enregistrement est valide, le GDS renverra un mot de passe PKCS12. **Ne perdez pas ce mot de passe, c'est la seule fois où il est mis à disposition pendant le processus d'émission du certificat**.

Après vérification et approbation, l’autorité de Certification de GDS générera un certificat complet incluant les clés privées et le crypter en utilisant le mot de passe PKCS12. Après avoir enregistré les clés publiques dans le service d'annuaire, le GDS enverra le certificat crypté sous forme de fichier ZIP au contact technique, ou au premier contact disponible sur le formulaire d'enregistrement.

Après avoir décompressé la pièce jointe de l'e-mail, vous devriez trouver un fichier nommé `<common_name>.p12`; vous pouvez décrypter ce fichier pour extraire les certificats comme suit :

```
$ openssl pkcs12 -in <common_name>.p12 -out <common_name>.pem -nodes
```

Vous pouvez également utiliser directement le fichier .zip sans le décrypter ou l'extraire grâce au module [`github.com/trisacrypto/trisa/pkg/trust`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#NewSerializer).

### Demande de signature de certificat

Une alternative à la création d'un certificat consiste à télécharger une demande de signature de certificat (CSR). Ce mécanisme est souvent préférable car il signifie qu'aucun élément de la clé privée ne doit être transmis sur le réseau et que la clé privée peut rester sur un matériel sécurisé.

Pour générer un CSR en utilisant `openssl` en ligne de commande, créez d'abord un fichier de configuration nommé `trisa.conf` dans votre répertoire de travail courant, en remplaçant `example.com` par le domaine sur lequel vous prévoyez d'héberger votre terminaison TRISA:

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

Puis, exécutez la commande suivante, en remplaçant `example.com` par le nom de domaine que vous utiliserez comme terminaison TRISA :

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

Votre clé privée est maintenant stockée dans `example.com.key` &mdash; ** gardez soigneusement cette clé privée ** &mdash; elle est requise pour les connexions mTLS dans votre service mTLS et établit la confiance sur le réseau TRISA.

Le fichier `example.com.csr` contient votre demande de signature de certificat. Copiez et collez le contenu de ce fichier, y compris `-----BEGIN CERTIFICATE REQUEST-----` et `-----END CERTIFICATE REQUEST-----` dans votre demande d'inscription.
