---
title: "Sectigo"
date: 2020-12-24T07:58:37-05:00
lastmod: 2021-10-08T16:25:00-05:00
description: "Interaktionen des Verzeichnisdienstes mit der Sectigo CA API"
weight: 50
---

[![Go Reference](https://pkg.go.dev/badge/github.com/trisacrypto/directory/pkg/sectigo.svg)](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)

Der TRISA-Verzeichnisdienst stellt Zertifikate mit der Sectigo-Zertifizierungsstelle über sein IoT-Portal aus. Da der Verzeichnisdienst öffentliches Schlüsselmaterial sammeln muss, um einen ersten vertrauenswürdigen Handshake für mTLS zu ermöglichen, verwendet er die Sectigo IoT Manager API als Teil des VASP-Registrierungs- und Verifizierungsprozesses. Das Paket `github.com/trisacrypto/directory/pkg/sectigo` ist eine Go-Bibliothek für die Interaktion mit der API und implementiert die vom Verzeichnisdienst benötigten Endpunkte und Methoden. Das TestNet bietet auch ein Kommandozeilenprogramm für die Interaktion mit der API zu Verwaltungs- und Debuggingzwecken. Diese Dokumentation beschreibt das Befehlszeilendienstprogramm, das auch einen Überblick darüber gibt, wie man die API direkt zum Ausstellen und Widerrufen von Zertifikaten verwendet.

Referenzmaterial:

- [Paketdokumentation](https://pkg.go.dev/github.com/trisacrypto/directory/pkg/sectigo)
- [IoT Manager API-Dokumentation](https://support.sectigo.com/Com_KnowledgeDetailPage?Id=kA01N000000bvCJ)
- [IoT-Manager-Portal](https://iot.sectigo.com)

## Erste Schritte

Um das `sectigo` CLI-Dienstprogramm zu installieren, laden Sie entweder eine vorkompilierte Binärdatei von [Releases auf GitHub](https://github.com/trisacrypto/directory/releases) herunter oder installieren Sie es lokal mit:

```
$ go get github.com/trisacrypto/directory/cmd/sectigo
```

Dies fügt den `sectigo` Befehl zu Ihrem `$PATH` hinzu.

## Authentifizierung

Der erste Schritt ist die Authentifizierung. Sie sollten Ihren Benutzernamen und Ihr Passwort in den Umgebungsvariablen `$SECTIGO_USERNAME` und `$SECTIGO_PASSWORD` setzen (alternativ können Sie sie als Parameter in der Kommandozeile übergeben). Um Ihren Authentifizierungsstatus zu überprüfen, verwenden Sie:

```
$ sectigo auth
```

Die API authentifiziert sich mit Benutzername und Passwort und gibt dann Zugriffs- und Aktualisierungs-Tokens zurück, die in einer lokalen Cache-Datei gespeichert werden. Um zu sehen, wo Ihr Cache gespeichert ist, verwenden Sie:

```
$ sectigo auth --cache
```

Wenn Sie den Status Ihrer Anmeldeinformationen überprüfen möchten, z. B. ob die Zugriffstoken gültig, auffrischbar oder abgelaufen sind, verwenden Sie:

```
$ sectigo auth --debug
```

## Autoritäten und Profile

Um mit den Zertifikaten zu arbeiten, müssen Sie die Autoritäten und Profile auflisten, auf die Ihr Benutzerkonto Zugriff hat.

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

Die Autorität zeigt die Methoden und Profile an, unter denen Zertifikate erstellt werden. Hier ist das Feld `profileId` sehr wichtig für die Verwendung in nachfolgenden Aufrufen. Sie können auch sehen, wie viele Lizenzen über alle Autoritäten hinweg wie folgt bestellt/ausgestellt wurden:

```
$ sectigo licenses
{
  "ordered": 2,
  "issued": 2
}
```

Um Detailinformationen für ein Profil zu erhalten, verwenden Sie profileId mit dem folgenden Befehl:

```
$ sectigo profiles -i 42
```

Dies gibt die rohe Profilkonfiguration zurück. Bevor Sie Zertifikate mit der Autorität erstellen, müssen Sie die erforderlichen Profilparameter kennen:

```
$ sectigo profile -i 42 --params
```

## Erstellen von Zertifikaten

Sie können die Erstellung eines Zertifikats mit den Parametern `commonName` und `pkcs12Password` wie folgt anfordern (beachten Sie, dass Sie für Profile, die andere Parameter erfordern, direkt die Codebasis verwenden und Ihre eigene Methode implementieren müssen):

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

Das Flag `-a` gibt die Autorität an, sollte aber eine Profil-ID sein. Die Domäne muss eine gültige Domäne sein. Wenn Sie kein Passwort angeben, wird eines für Sie generiert und vor dem Beenden auf der CLI ausgegeben. Das Flag `-b` gibt einen von Menschen lesbaren Namen für die Batch-Erstellung an. Die Rückgabedaten zeigen Details über den erstellten Stapelzertifikatsauftrag; Um den Status zu überprüfen, können Sie die Daten wie folgt abrufen:

```
$ sectigo batches -i 24
```

Sie können auch Informationen über die Verarbeitung des Stapels erhalten:

```
$ sectigo batches -i 24 --status
```

Sobald der Stapel erstellt ist, ist es an der Zeit, die Zertifikate in einer ZIP-Datei herunterzuladen:

```
$ sectigo download -i 24 -o certs/
```

Dadurch wird die Batch-Datei (normalerweise batchId.zip, in diesem Fall 24.zip) in das Verzeichnis `certs/` heruntergeladen. Entpacken Sie die certs und entschlüsseln Sie die Datei .pem wie folgt:

```
$ unzip certs/24.zip
$ openssl pkcs12 -in certs/example.com.p12 -out certs/example.com.pem -nodes
```

Weitere Informationen zur Arbeit mit der PKCS12-Datei finden Sie unter [Exportieren von Zertifikaten und privaten Schlüsseln aus einer PKCS#12-Datei mit OpenSSL](https://www.ssl.com/how-to/export-certificates-private-key-from-pkcs12-file-with-openssl/).

## Hochladen einer CSR

Eine Alternative zur Zertifikatserstellung ist das Hochladen eines Certificate Signing Request (CSR). Dieser Mechanismus ist oft vorzuziehen, da er bedeutet, dass kein privates Schlüsselmaterial über das Netzwerk übertragen werden muss und der private Schlüssel auf sicherer Hardware verbleiben kann.

Um eine CSR mit `openssl` auf der Kommandozeile zu erzeugen, erstellen Sie zunächst eine Konfigurationsdatei mit dem Namen `trisa.conf` in Ihrem aktuellen Arbeitsverzeichnis, wobei Sie `example.com` durch die Domain ersetzen, auf der Sie Ihren TRISA-Endpunkt hosten wollen:

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
ST = [State or Province (vollständig buchstabiert, keine Abkürzungen)]
C = [2 digit country code]
[v3ext_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = clientAuth, serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = example.com
```

Bitte füllen Sie die Konfiguration für Ihr Zertifikat sorgfältig aus. Diese Informationen müssen korrekt sein und können nicht geändert werden, ohne dass das Zertifikat neu ausgestellt wird. Achten Sie auch darauf, dass nach den Einträgen in der Konfiguration keine Leerzeichen stehen!

Führen Sie anschließend den folgenden Befehl aus und ersetzen Sie dabei `example.com` durch den Namen der Domäne, die Sie als TRISA-Endpunkt verwenden werden:

```
$ openssl req -new -newkey rsa:4096 -nodes -sha384 -config trisa.conf \
  -keyout example.com.key -out example.com.csr
```

Anschließend können Sie die CSR mithilfe des CLI-Programms wie folgt hochladen:

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

Das Flag `-p` gibt das Profil an, mit dem die CSR-Batch-Anforderung verwendet werden soll, und muss eine gültige profileId sein. Bei den hochgeladenen CSRs kann es sich um eine einzelne Textdatei mit mehreren CSRs im PEM-Format handeln, wobei Standard-BEGIN/END-Trennzeichen verwendet werden.

## Verwaltung von Zertifikaten

Sie können nach einem Zertifikat über den Namen oder die Seriennummer suchen, aber meistens sollten Sie nach der Domain oder dem Common Name suchen, um die Seriennummer zu erhalten:

```
$ sectigo find -n example.com
```

Sobald Sie die Seriennummer erhalten haben, können Sie das Zertifikat wie folgt widerrufen:

```
$ sectigo revoke -p 42 -r "cessation of operation" -s 12345
```

Dieser Befehl erwartet die Profil-ID, die das Zertifikat ausgestellt hat, mit dem Flag `-p`, einen [RFC 5280 Reason Code](https://tools.ietf.org/html/rfc5280#section-5.3.1), der über das Flag `-r` übergeben wird (standardmäßig nicht spezifiziert), und die Seriennummer des Zertifikats mit dem Flag `-s`. Wenn dieser Befehl keinen Fehler auslöst, wurde das Zertifikat erfolgreich widerrufen.

Die RFC 5280 Reason Codes sind:

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

Beachten Sie, dass bei der Angabe des Reason Codes zwischen Leerzeichen und Groß-/Kleinschreibung nicht unterschieden wird.
