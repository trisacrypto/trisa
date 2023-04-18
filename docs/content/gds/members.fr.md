---
title: "Les Membres Entrants"
date: 2022-02-21T09:34:44-05:00
lastmod: 2022-02-21T09:34:44-05:00
description: "Connexion à d'autres membres de TRISA"
weight: 70
---

Le service `TRISAMembers` est un service expérimental qui fournit un accès supplémentaire et sécurisé au service d'annuaire pour les membres approuvés du réseau TRISA. Seuls les membres du réseau TRISA (c'est-à-dire les membres qui ont reçu des certificats d'identité valides) peuvent accéder aux données détaillées du service d'annuaire concernant le réseau en utilisant ce service. *Remarque : une fois validé, ce service sera intégré à la spécification officielle de TRISA.*

## Le Service `TRISAMembers`

Cette section décrit les protocol buffers pour la terminaison `TRISAMembers`, qui peut être trouvée [ici](https://github.com/trisacrypto/directory/blob/main/proto/gds/members/v1alpha1/members.proto). Ce fichier peut être compilé dans le langage de votre choix ([exemple pour le Golang](https://github.com/trisacrypto/directory/tree/main/pkg/gds/members/v1alpha1)). *Remarque : Vous devrez télécharger et installer le compilateur du protocol buffer si vous ne l'avez pas encore.*

Actuellement, le service `TRISAMembers` ne comprend qu'une seule RPC &mdash; le `List` RPC. D'autres RPCs expérimentaux et sécurisés pourront être mis à disposition à l'avenir.

```proto
service TRISAMembers {
    // List all verified VASP members in the Directory Service.
    rpc List(ListRequest) returns (ListReply) {};
}
```

## Liste des membres vérifiés

Le `List` RPC renvoie une liste paginée de tous les membres vérifiés de TRISA afin de faciliter la recherche des pairs dans le réseau TRISA. Le RPC reçoit en entrée une requête de type `ListRequest` et renvoie un `ListReply`.

### `ListRequest`

La fonction `ListRequest` peut être utilisée pour gérer la pagination de la demande de liste VASP (VASP - Fournisseur de Services d’Actifs Virtuels). S'il y a plus de résultats que la taille de page spécifiée, la fonction `ListReply` renvoie un jeton de page qui peut être utilisé pour récupérer la page suivante (tant que les paramètres de la requête originale ne sont pas modifiés, par exemple les filtres ou les paramètres de pagination).

La variable `page_size` spécifie le nombre de résultats par page et ne peut pas changer entre les requêtes de page ; sa valeur par défaut est 100. La variable `page_token` spécifie le jeton de page permettant d'extraire la page de résultats suivante.

```proto
message ListRequest {
    int32 page_size = 1;
    string page_token = 2;
}
```

### `ListReply`

Un `ListReply` renvoie une liste abrégée des détails du VASP, destinée à faciliter les échanges de clés entre pairs ou les recherches plus détaillées dans le service d'annuaire.

Les `vasps` sont une liste de VASP (voir la déclaration de `VASPMember` ci-dessous), et si le paramètre `next_page_token`, est spécifié, il indique qu'une autre page de résultats existe.

```proto
message ListReply {
    repeated VASPMember vasps = 1;
    string next_page_token = 2;
}
```

### `VASPMember`

Un `VASPMember` contient suffisamment d'informations pour faciliter les échanges d'égal à égal ou des recherches plus détaillées dans le service d'annuaire. Le `ListReply` contiendra une liste d'aucun, d'un ou de plusieurs `VASPMembers`.

```proto
message VASPMember {
    // Les éléments d'identification unique du VASP dans le service d'annuaire
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;

    // Adresse à laquelle se connecter au VASP distant pour effectuer une demande TRISA
    string endpoint = 4;

    // Détails supplémentaires utilisés pour faciliter les recherches et les correspondances
    string name = 5;
    string website = 6;
    string country = 7;
    trisa.gds.models.v1beta1.BusinessCategory business_category = 8;
    repeated string vasp_categories = 9;
    string verified_on = 10;
}
```

## Connexion via authentification mTLS

Pour utiliser le service `TRISAMembers`, vous devez vous authentifier à l’aide du protocole [mTLS](https://grpc.io/docs/guides/auth/) en utilisant les certificats d'identité TRISA qui vous ont été attribués lors de votre inscription.

La documentation gRPC sur [l’authentification](https://grpc.io/docs/guides/auth) fournit des exemples de code permettant de se connecter avec mTLS dans différents langages, notamment [Java](https://grpc.io/docs/guides/auth/#java), [C++](https://grpc.io/docs/guides/auth/#c), [Golang](https://grpc.io/docs/guides/auth/#go), [Ruby](https://grpc.io/docs/guides/auth/#ruby), et [Python](https://grpc.io/docs/guides/auth/#python).

Par exemple, si vous utilisez Golang pour vous connecter au service d'annuaire, vous utiliserez les bibliothèques [`tls`](https://pkg.go.dev/crypto/tls), [`x509`](https://pkg.go.dev/crypto/x509), et [`credentials`](https://pkg.go.dev/google.golang.org/grpc/credentials) pour charger vos certificats d'identité TRISA depuis leur emplacement sécurisé sur votre ordinateur et construire des informations d'identification TLS pour valider mutuellement la connexion avec le serveur. Enfin, vous utiliserez le code compilé des membres protocol buffer pour créer un client de membres. *Remarque : les déclarations du protocol buffer sont décrites plus en détail plus haut dans cette page.*

```golang
import (
    "crypto/tls"
    "crypto/x509"

    members "github.com/trisacrypto/directory/pkg/gds/members/v1alpha1"
    "google.golang.org/grpc/credentials"
)

func (m *MyProfile) Connect() (_ members.TRISAMembersClient, err error){
    config := &tls.Config{
		Certificates: []tls.Certificate{m.Cert}, // m.Cert is your TRISA certificate parsed into a *x509.Certificate
		MinVersion:   tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		ClientAuth: tls.RequireAndVerifyClientCert, // this will ensure an mTLS connection
		ClientCAs:  m.Pool, // m.Pool is a *x509.CertPool that must contain the RA and IA public certs from your TRISA certificate chain
	}

    var creds *credentials.TransportCredentials
    if creds, err = credentials.NewTLS(config); err != nil {
        return nil, err
    }

    var cc *grpc.ClientConn
    if cc, err = grpc.Dial(m.Endpoint, grpc.WithTransportCredentials(creds)); err != nil {
        return nil, err
    }

    return members.NewTRISAMembersClient(cc), nil
}
```

*Notez qu'il existe actuellement deux répertoires TRISA : le TRISA [TestNet](https://trisatest.net/), qui permet aux utilisateurs d'expérimenter les interactions TRISA, et le [Répertoire VASP](https://vaspdirectory.net/), qui est le réseau de production pour les transactions TRISA. Si vous vous êtes inscrit au TestNet et que vous disposez de certificats TestNet, le point de terminaison que vous passerez dans la fonction de numérotation sera `members.trisatest.net:443`. Alternativement, si vous souhaitez accéder aux membres de l'annuaire VASP et que vous êtes déjà enregistré, vous utiliserez le terminal suivant `members.vaspdirectory.net:443`.*
