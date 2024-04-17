---
title: Integrationsübersicht
date: 2021-04-23T01:35:35-04:00
lastmod: 2021-10-08T16:27:33-04:00
description: "Beschreibt, wie das TRISA-Protokoll in das TestNet integriert wird"
weight: 5
---

## Übersicht über die TRISA-Integration

1. Registrierung bei einem TRISA-Verzeichnisdienst
2. Implementierung des TRISA-Netzwerkprotokolls
3. Implementierung des TRISA-Gesundheitsprotokoll

## VASP-Verzeichnisdienst-Registrierung

### Überblick über die Registrierung

Bevor Sie das TRISA-Protokoll in Ihre VASP-Software integrieren können, müssen Sie sich bei einem TRISA Directory Service (DS) registrieren.  Der TRISA-Verzeichnisdienst stellt den öffentlichen Schlüssel und die TRISA-Verbindungsinformationen für registrierte VASPs zur Verfügung.

Sobald Sie sich beim TRISA DS registriert haben, erhalten Sie ein KYV-Zertifikat.  Der öffentliche Schlüssel im KYV-Zertifikat wird anderen VASPs über den TRISA DS zur Verfügung gestellt.

Bei der Registrierung beim DS müssen Sie den Endpunkt `address:port` angeben, an dem Ihr VASP den TRISA-Netzdienst implementiert. Diese Adresse wird im DS registriert und von anderen VASPs verwendet, wenn Ihr VASP als begünstigter VASP identifiziert wird.

Für Integrationszwecke steht eine gehostete TestNet TRISA DS-Instanz zum Testen zur Verfügung.  Der Registrierungsprozess wurde im TestNet vereinfacht, um eine schnelle Integration zu ermöglichen.  Es wird empfohlen, die Registrierung des Produktions-DS während der Integration in das TestNet zu starten.


### Verzeichnisdienst-Registrierung

Um mit der Registrierung des TRISA DS zu beginnen, besuchen Sie die Website unter [https://vaspdirectory.net/](https://vaspdirectory.net/)

Wählen Sie die Registerkarte "Registrieren", um mit der Registrierung zu beginnen. Beachten Sie, dass Sie auf dieser Website Ihre Registrierungsdaten Feld für Feld eingeben oder ein JSON-Dokument mit Ihren Registrierungsdaten hochladen können.

Diese Registrierung führt dazu, dass eine E-Mail an alle in der JSON-Datei angegebenen technischen Kontakte geschickt wird.  Die E-Mails leiten Sie durch den weiteren Verlauf des Registrierungsprozesses.  Sobald Sie die Registrierungsschritte abgeschlossen haben, erhalten die TRISA TestNet-Administratoren Ihre Registrierung zur Überprüfung.

Sobald die TestNet-Administratoren die Registrierung überprüft und genehmigt haben, erhalten Sie per E-Mail ein KYV-Zertifikat, und Ihr VASP wird im TestNet DS öffentlich sichtbar sein.


## Implementierung des TRISA P2P Protokolls


### Voraussetzungen

Um mit der Einrichtung zu beginnen, benötigen Sie Folgendes:

* KYV-Zertifikat (von der TRISA DS-Registrierung)
* Der öffentliche Schlüssel, der für die CSR verwendet wurde, um Ihr Zertifikat zu erhalten
* den zugehörigen privaten Schlüssel
* den Hostnamen des TRISA-Verzeichnisdienstes
* Die Fähigkeit, sich an address:port zu binden, die mit Ihrem VASP im TRISA-Verzeichnisdienst verbunden ist.


### Überblick über die Integration

Die Integration des TRISA-Protokolls umfasst sowohl eine Client- als auch eine Serverkomponente.

Die Client-Komponente stellt eine Schnittstelle zu einer TRISA Directory Service (DS)-Instanz her, um andere VASP zu suchen, die das TRISA-Messaging-Protokoll integrieren.  Die Client-Komponente wird für ausgehende Transaktionen von Ihrem VASP verwendet, um zu überprüfen, ob der empfangende VASP TRISA-konform ist.

Die Serverkomponente empfängt Anfragen von anderen VASPs, die das TRISA-Protokoll integrieren, und liefert Antworten auf deren Anfragen.  Die Serverkomponente bietet Rückrufe, die implementiert werden müssen, damit Ihre VASP Informationen zurückgeben kann, die dem TRISA-Netzprotokoll entsprechen.

Derzeit ist eine Referenzimplementierung des TRISA-Netzwerkprotokolls in Go verfügbar.

[https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go](https://github.com/trisacrypto/testnet/blob/main/pkg/rvasp/trisa.go)

Integrierende VASPs müssen ihre eigene Implementierung des Protokolls verwenden.  Wenn eine andere Sprache als Go erforderlich ist, können Client-Bibliotheken aus den Protokollpuffern, die das TRISA-Netzprotokoll definieren, generiert werden.

Von den Integratoren wird erwartet, dass sie eingehende Übertragungsanfragen und den Austausch von Schlüsseln integrieren und können optional auch ausgehende Übertragungsanfragen und den Austausch von Schlüsseln integrieren.

### Hinweise zur Integration

Das TRISA-Netzprotokoll legt fest, wie die Daten zwischen den teilnehmenden VASPs übertragen werden.  Das empfohlene Format für Daten, die zur Identifizierung von Informationen übertragen werden, ist das IVMS101-Datenformat.  Es liegt in der Verantwortung des implementierenden VASP, sicherzustellen, dass die gesendeten/empfangenen Identifizierungsdaten die FATF-Travel Regel erfüllen.

Das Ergebnis einer erfolgreichen TRISA-Transaktion sind ein Schlüssel und verschlüsselte Daten, die den FATF-Travel Regeln entsprechen.  TRISA legt nicht fest, wie diese Daten nach Erhalt zu speichern sind.  Es liegt in der Verantwortung des implementierenden VASP, die sichere Speicherung der aus der Transaktion resultierenden Daten zu gewährleisten.
