---
title: "Zugriff auf Mitglieder"
date: 2022-01-03T14:02:37-05:00
lastmod: 2022-02-23T10:35:11-05:00
description: "Zugriff auf andere TRISA-Mitglieder"
weight: 70
---

Der Dienst `TRISAMembers` ist ein experimenteller Dienst, der verifizierten TRISA-Mitgliedern einen zusätzlichen, sicheren Zugriff auf den Verzeichnisdienst bietet. Nur Mitglieder des TRISA-Netzwerks (z. B. Mitglieder, denen gültige Identitätszertifikate ausgestellt wurden) können über diesen Dienst auf detaillierte Verzeichnisdienstdaten über das Netzwerk zugreifen. *Hinweis: Nach der Validierung wird dieser Dienst in die offizielle TRISA-Spezifikation verschoben.*

## Der `TRISAMembers` Service

Dieser Abschnitt beschreibt die protocol buffers für den `TRISAMembers`-Endpunkt, der [hier](https://github.com/trisacrypto/directory/blob/main/proto/gds/members/v1alpha1/members.proto) zu finden  ist. Diese Datei kann in die Sprache Ihrer Wahl kompiliert werden ([Beispiel für Golang](https://github.com/trisacrypto/directory/tree/main/pkg/gds/members/v1alpha1)). *Hinweis: Sie müssen den protocol buffers herunterladen und installieren  , falls Sie dies noch nicht getan haben.*

Derzeit besteht der Service `TRISAMembers` nur aus einem einzigen RPC &mdash;  der RPC `List`. Andere experimentelle, sichere RPCs könnten in Zukunft zugänglich gemacht werden.

```proto
service TRISAMembers {
    // Listen Sie alle verifizierten VASP-Mitglieder im Verzeichnisdienst auf.
    rpc List(ListRequest) returns (ListReply) {};
}
```

## Auflisten von verifizierten Mitgliedern

Der RPC `List` gibt eine paginierte Liste aller _ve rified_ TRISA-Mitglieder zurück, um die SUCHE NACH TRISA-Peers zu erleichtern. Der RPC erwartet als Eingabe eine `ListRequest` und gibt eine `ListReply` zurück.

### `ListRequest`

Ein `ListRequest` kann verwendet werden, um die Paginierung für die VASP-Listenanforderung zu verwalten. Wenn es mehr Ergebnisse als die angegebene Seitengröße gibt die `ListReply` ein Seitentoken zurück, das zum Abrufen der nächsten Seite verwendet werden kann (solange die Parameter der ursprünglichen Anforderung nicht geändert werden , z. B. Filter oder Paginierungsparameter).

Die `page_size` gibt die Anzahl der Ergebnisse pro Seite an und kann nicht zwischen Seitenanforderungen wechseln.  Der Standardwert ist 100. Die `page_token` gibt das Seitentoken an, um die nächste Seite mit Ergebnissen abzurufen.

```proto
message ListRequest {
    int32 page_size = 1;
    string page_token = 2;
}
```

### `ListReply`

Ein `ListReply` gibt eine abgekürzte Auflistung von VASP-Details zurück, die den Peer-to-Peer-Schlüsselaustausch oder detailliertere Suchen gegen den Verzeichnisdienst erleichtern sollen.

Die `vasps` sind eine Liste von VASPs (siehe Definition von `VASPMember` unten), und die `next_page_token`, falls angegeben, ist die Benachrichtigung, dass eine andere Seite mit Ergebnissen existiert.

```proto
message ListReply {
    repeated VASPMember vasps = 1;
    string next_page_token = 2;
}
```

### `VASPMember`

Ein `VASPMember` enthält genügend Informationen, um den Peer-to-Peer-Austausch oder detailliertere Suchen gegen den Verzeichnisdienst zu erleichtern. Die `ListReply` enthält eine Liste mit keinem, einem oder mehreren `VASPMembers`.

```proto
message VASPMember {
    // Die eindeutig identifizierenden Komponenten des VASP im Verzeichnisdienst
    string id = 1;
    string registered_directory = 2;
    string common_name = 3;

    // Adresse für die Verbindung mit dem Remote-VASP, um eine TRISA-Anforderung auszuführen
    string endpoint = 4;

    //Zusätzliche Details zur Erleichterung der Suche und des Abgleichs
    string name = 5;
    string website = 6;
    string country = 7;
    trisa.gds.models.v1beta1.BusinessCategory business_category = 8;
    repeated string vasp_categories = 9;
    string verified_on = 10;
}
```

## Verbinden mit mTLS

Um den Dienst `TRISAMembers` nutzen zu  können, müssen Sie sich mit [mTLS](https://grpc.io/docs/guides/auth/) mit den TRISA-Identitätszertifikaten authentifizieren, die Ihnen bei der Registrierung erteilt wurden.

Die gRPC-Dokumentation zu [authentication](https://grpc.io/docs/guides/auth) enthält Codebeispiele für die Verbindung über mTLS in einer Vielzahl von Sprachen, einschließlich [Java](https://grpc.io/docs/guides/auth/#java), [C++](https://grpc.io/docs/guides/auth/#c), [Golang](https://grpc.io/docs/guides/auth/#go), [Ruby](https://grpc.io/docs/guides/auth/#ruby) und [Python](https://grpc.io/docs/guides/auth/#python).

Wenn Sie beispielsweise Golang verwenden, um eine Verbindung mit dem Verzeichnisdienst herzustellen, würden Sie die Bibliotheken [`tls`](https://pkg.go.dev/crypto/tls), [`x509`](https://pkg.go.dev/crypto/x509) und [`credentials`](https://pkg.go.dev/google.golang.org/grpc/credentials)  verwenden, um Ihre TRISA-Identitätszertifikate von ihrem sicheren Speicherort auf Ihrem Computer zu laden und TLS-Anmeldeinformationen zu erstellen, um sie gegenseitig zu überprüfen. die Verbindung mit dem Server. Schließlich würden  Sie den kompilierten protocol buffer verwenden  , um einen Member-Client zu erstellen. *Hinweis: Die protocol buffer werden weiter oben auf dieser Seite ausführlicher beschrieben.*

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
    if cc, err = grpc.NewClient(m.Endpoint, grpc.WithTransportCredentials(creds)); err != nil {
        return nil, err
    }

    return members.NewTRISAMembersClient(cc), nil
}
```

*Beachten Sie, dass es derzeit zwei TRISA-Verzeichnisse gibt; das TRISA [TestNet](https://testnet.directory/),, mit dem Benutzer mit den TRISA-Interaktionen experimentieren können, und das [VASP-Verzeichnis](https://trisa.directory/),, das das Produktionsnetzwerk für TRISA-Transaktionen ist. Wenn Sie sich für das TestNet registriert haben  und über TestNet-Zertifikate verfügen, ist der Endpunkt, den Sie an die Wählfunktion übergeben, `members.testnet.directory:443`. Wenn Sie auf Mitglieder des VASP-Verzeichnisses zugreifen möchten und bereits ein erneutes Mitglied sind, verwenden Sie alternativ den Endpunkt `members.trisa.directory:443`.*
