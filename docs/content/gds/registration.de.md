---
title: "Registrieren"
date: 2021-07-22T09:21:59-04:00
lastmod: 2021-10-08T16:23:06-04:00
description: "Registrierung eines VASP beim Verzeichnisdienst"
weight: 20
---

Um den TRISA- oder TestNet-Netzwerken beizutreten, müssen Sie sich beim TRISA Global Directory Service (GDS) oder bei einem der gerichtsspezifischen Verzeichnisdienste registrieren. Die Registrierung beim Verzeichnisdienst umfasst zwei Arbeitsabläufe:

1. Ein KYV-Überprüfungsprozess, um sicherzustellen, dass das Netzwerk eine vertrauenswürdige Mitgliedschaft aufrechterhält
2. Ausstellung von Zertifikaten für die mTLS-Authentifizierung im Netzwerk

Demnächst: weitere Details zum Registrierungsformular, zur E-Mail-Verifizierung und zum Überprüfungsprozess.

## Ausstellung von Zertifikaten

Derzeit gibt es zwei Mechanismen, um mTLS-Zertifikate vom GDS zu erhalten, wenn Ihre Registrierung überprüft und genehmigt wurde.

1. Verschlüsselte PKCS12-Zertifikate per E-Mail
2. Antrag auf Unterzeichnung eines Zertifikats (CSR)

Sie müssen eine dieser Optionen auswählen, _wenn Sie Ihre Anmeldung einreichen_; nachdem Ihre Anmeldung eingereicht wurde, können Sie nicht mehr zwischen den Optionen wechseln.

### PKCS12-verschlüsselter E-Mail-Anhang

Der erste Mechanismus ist der einfachste &mdash: Wählen Sie bei der Registrierung einfach die E-Mail-Option und lassen Sie die CSR-Felder aus. Wenn das Registrierungsformular gültig ist, sendet das GDS ein PKCS12-Passwort zurück. **Verlieren Sie dieses Passwort nicht, es ist das einzige Mal, dass es während der Ausstellung des Zertifikats zur Verfügung gestellt wird**.

Nach der Genehmigung der Überprüfung erstellt der GDS CA ein vollständiges Zertifikat einschließlich der privaten Schlüssel und verschlüsselt es mit dem PKCS12-Passwort. Nach der Registrierung der öffentlichen Schlüssel im Verzeichnisdienst sendet der GDS das verschlüsselte Zertifikat als ZIP-Datei per E-Mail an den technischen Kontakt oder den ersten verfügbaren Kontakt im Registrierungsformular.

Nach dem Entpacken des E-Mail-Anhangs sollten Sie eine Datei mit dem Namen `<common_name>.p12` finden; Sie können diese Datei wie folgt entschlüsseln, um die Zertifikate zu extrahieren:

```
$ openssl pkcs12 -in <common_name>.p12 -out <common_name>.pem -nodes
```

Sie können die .zip-Datei auch direkt verwenden, ohne sie zu entschlüsseln oder zu extrahieren, und zwar mit dem Modul [`github.com/trisacrypto/trisa/pkg/trust`](https://pkg.go.dev/github.com/trisacrypto/trisa/pkg/trust#NewSerializer).

### Zertifikatssignierungsanfragen

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

Ihr privater Schlüssel ist nun in `example.com.key` &mdash gespeichert; **Bewahren Sie diesen privaten Schlüssel sicher auf** &mdash; er wird für mTLS-Verbindungen in Ihrem mTLS-Dienst benötigt und stellt das Vertrauen in das TRISA-Netzwerk her.

Die Datei `example.com.csr` enthält Ihre Zertifikatsignierungsanforderung. Kopieren Sie den Inhalt dieser Datei einschließlich `-----BEGIN CERTIFICATE REQUEST-----` und `-----END CERTIFICATE REQUEST-----` und fügen Sie ihn in Ihre Registrierungsanfrage ein.
