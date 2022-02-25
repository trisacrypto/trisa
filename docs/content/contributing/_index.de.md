---
title: Beitragender
date: 2021-06-14T11:34:06-05:00
lastmod: 2022-02-23T10:30:47-05:00
description: "Zum Open-Source-Projekt beitragen"
weight: 110
---

TRISA ist ein Open-Source-Projekt und freut sich über Beiträge!

Wenn Sie ein Entwickler sind, dessen Organisation das TRISA-Protokoll verwendet (oder plant, es zu übernehmen), ist dieser Abschnitt für Sie!

## Navigieren im Repository

Dieses Repository enthält eine gRPC-Implementierung des TRISA-Protokolls, wie im [Whitepaper](https://trisa.io/trisa-whitepaper/) beschrieben, das [protocol buffers](https://grpc.io/) und Golang nutzt.

Der Ordner `proto` enthält die wichtigsten RPC-Definitionen, einschließlich:
- die Nachrichtendefinitionen des interVASP Messaging Standard (IVMS), die als Grundlage dafür dienen, wie zwei VASP-Peers gemeinsam Entitäten beschreiben sollten, die an kryptografischen Übertragungen beteiligt sind, einschließlich Namen, Standorte und Regierungskennungen. Dies ist die Spezifikation, die es den Originatoren ermöglicht, sich gegenüber den Begünstigten zu identifizieren und von diesen Begünstigten Informationen anzufordern, um die gesetzlichen Anforderungen ihrer Regulierungsbehörden zu erfüllen.
- die Service-Definitionen des TRISA-Netzwerks, im Wesentlichen wie die verschiedenen Teile der API funktionieren &mdash;  vom Austausch von Schlüsseln (um sicherzustellen, dass beide Peers über die erforderlichen Informationen verfügen) bis hin zur Übertragung von "sicheren Umschlägen" (kryptografisch versiegelte protocol buffer, die nur vom beabsichtigten Empfänger entschlüsselt werden können). Der  Unterordner `trisa` enthält auch generische Nachrichtentypen für Transaktionen, die maximale Flexibilität für eine Vielzahl von TRISA-Anwendungsfällen bieten sollen.

Der Ordner `pkg` enthält den Referenzimplementierungscode, einschließlich kompiliertem Code, der aus den protocol buffer im  Ordner 'proto' generiert wurde[^1].
 - Der Folder `iso3166` enthält Sprachcodes.
 - Der Ordner `ivms101` erweitert den generierten Protobuf-Code um JSON-Ladewerkzeuge, Validierungshelfer, kurze Konstanten usw.
 - Der Ordner `trisa` enthält Strukturen und Methoden für eine Reihe von TRISA-bezogenen Aufgaben, z. B. das Ausführen von Kryptographie und das Herstellen von mTLS-Verbindungen .

Der Ordner `lib` soll Hilfsprogrammcode ähnlich dem im Ordner `pkg` anzeigen  , jedoch für andere Sprachen als Golang. Wenn Sie in einer anderen Sprache als Golang arbeiten, wäre dies ein großartiger Ort, um Ihren Beitrag zu beginnen!

[^1]: Beachten Sie, dass diese kompilierten Dateien für Golang kompiliert werden; aber dies ist sicherlich nicht die einzige Option. Diejenigen, die Implementierungscode in einer anderen Sprache erstellen möchten, sollten sich den  Ordner `lib` ansehen  , der derzeit Platzhalterordner enthält, aber solche anderen Implementierungen (einschließlich kompiliertem protocol buffer für diese anderen Sprachen) präsentieren soll.

## Der globale Verzeichnisdienst

Ein weiterer integraler Bestandteil des TRISA-Protokolls ist der Global Directory Service, der als Suchwerkzeug für TRISA-Mitglieder dient, um Peers zu identifizieren, mit denen sie Informationen austauschen möchten. RPC-Definitionen und Implementierungscode im Zusammenhang mit dem globalen Verzeichnisdienst finden Sie im zugehörigen [`directory` repository](https://github.com/trisacrypto/directory).  Um mehr darüber zu erfahren, wie Sie Mitglied des Verzeichnisses werden können, besuchen Sie [vaspdirectory.net](https://vaspdirectory.net/).

## Übersetzungen

Übersetzungen der Dokumentation über trisa.dev werden regelmäßig von menschlichen Übersetzern durchgeführt und können nicht mehr synchron sein oder Fehler widerspiegeln. Wenn Sie einen Fehler bemerken, öffnen Sie bitte einen [Fehlerbericht](https://github.com/trisacrypto/trisa/issues/new), um uns zu benachrichtigen.


