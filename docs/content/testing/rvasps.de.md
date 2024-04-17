---
title: Roboter-VASPs
date: 2021-06-14T12:50:23-05:00
lastmod: 2021-10-08T16:33:29-05:00
description: "Arbeiten mit rVASPs für die TestNet-Integration"
weight: 10
---

Das TestNet beherbergt drei praktische "Roboter-VASP"-Dienste (rVASPs), die die Integration und Prüfung mit dem TRISA TestNet erleichtern. Diese Dienste sind Folgende:

- Alice (`api.alice.vaspbot.net:443`): die primäre Integration, die rVASP verwendet, um TRISA-Meldungen auszulösen und zu empfangen.
- Bob (`api.bob.vaspbot.net:443`): ein Demo-RVASP, um den Austausch mit Alice zu sehen.
- Evil (`api.evil.vaspbot.net:443`): ein "böswilliger" rVASP, der kein TRISA-TestNet-Mitglied ist und zum Testen nicht authentifizierter Interaktionen verwendet wird.

Hinweis: Die rVASPs sind derzeit in erster Linie für Demos konfiguriert, und es wird daran gearbeitet, sie für Integrationszwecke robuster zu machen; bitte schauen Sie regelmäßig in dieser Dokumentation nach, ob es Änderungen gibt. Wenn Sie irgendwelche Fehler im rVASP-Code oder im Verhalten bemerken, [reichen Sie bitte ein Anliegen ein](https://github.com/trisacrypto/testnet/issues).

## Erste Schritte mit rVASPs

Es gibt zwei Möglichkeiten, wie Sie die rVASPs für die Entwicklung Ihres TRISA-Dienstes nutzen können:

1. Sie können den rVASP veranlassen, eine TRISA-Austauschmeldung an Ihren Dienst zu senden.
2. Sie können eine TRISA-Meldung an den rVASP mit einem gültigen (oder ungültigen) rVASP-Kunden senden.

Die rVASP haben eine eingebaute Datenbank mit gefälschten Kunden mit gefälschten Wallet-Adressen. Ihre Antwort auf TRISA-Meldungen oder auf eine ausgelöste Überweisung setzt voraus, dass der Auftraggeber/Begünstigte für den rVASP gültig ist. Handelt es sich bei der Wallet-Adresse des Kunden beispielsweise um eine gültige Adresse von Alice und ist Alice der begünstigte Kunde, antwortet der rVASP mit den gefälschten KYC-Daten des Kunden; ist dies nicht der Fall, gibt er einen TRISA-Fehlercode zurück.

Die folgende Tabelle der "Kunden" für Alice, Bob und Evil kann als Referenz für die Interaktion mit jedem rVASP verwendet werden:

| VASP                  | "Crypto Wallet"                    | Email                 |
|-----------------------|------------------------------------|-----------------------|
| api.bob.vaspbot.net   | 18nxAxBktHZDrMoJ3N2fk9imLX8xNnYbNh | robert@bobvasp.co.uk  |
| api.bob.vaspbot.net   | 1LgtLYkpaXhHDu1Ngh7x9fcBs5KuThbSzw | george@bobvasp.co.uk  |
| api.bob.vaspbot.net   | 14WU745djqecaJ1gmtWQGeMCFim1W5MNp3 | larry@bobvasp.co.uk   |
| api.bob.vaspbot.net   | **1Hzej6a2VG7C8iCAD5DAdN72cZH5THSMt9** | fred@bobvasp.co.uk    |
| api.alice.vaspbot.net | **19nFejdNSUhzkAAdwAvP3wc53o8dL326QQ** | sarah@alicevasp.us    |
| api.alice.vaspbot.net | 1ASkqdo1hvydosVRvRv2j6eNnWpWLHucMX | mary@alicevasp.us     |
| api.alice.vaspbot.net | 1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha | alice@alicevasp.us    |
| api.alice.vaspbot.net | 14HmBSwec8XrcWge9Zi1ZngNia64u3Wd2v | jane@alicevasp.us     |
| api.evil.vaspbot.net  | **1AsF1fMSaXPzz3dkBPyq81wrPQUKtT2tiz** | gambler@evilvasp.gg   |
| api.evil.vaspbot.net  | 1PFTsUQrRqvmFkJunfuQbSC2k9p4RfxYLF | voldemort@evilvasp.gg |
| api.evil.vaspbot.net  | 172n89jLjXKmFJni1vwV5EzxKRXuAAoxUz | launderer@evilvasp.gg |
| api.evil.vaspbot.net  | 182kF4mb5SW4KGEvBSbyXTpDWy8rK1Dpu  | badnews@evilvasp.gg   |

Beachten Sie, dass alle rVASP-Daten mit einem Faker-Tool erzeugt wurden, um realistische/konsistente Testdaten und Vorrichtungen zu erzeugen, und dass sie völlig fiktiv sind. Die Datensätze für Alice VASP (ein gefälschtes US-Unternehmen) befinden sich beispielsweise hauptsächlich in Nordamerika usw.

Wenn Sie ein Traveler-Kunde sind, haben die fettgedruckten Adressen oben einige Attributionsdaten, die mit ihnen verbunden sind, und sie sind ein guter Kandidat für Traveler-basierte rVASP-Interaktionen.

### Präliminarien

In dieser Dokumentation wird davon ausgegangen, dass Sie einen Dienst haben, der den neuesten `TRISANetwork`-Dienst ausführt, und dass er im TRISA TestNet registriert ist und TestNet-Zertifikate korrekt installiert hat. Weitere Informationen finden Sie in der [TRISA-Integrationsübersicht](). **WARNUNG**: Die rVASPs nehmen nicht am TRISA-Produktionsnetz teil, sie antworten nur auf verifizierte TRISA TestNet mTLS-Verbindungen.

Um mit der rVASP-API zu interagieren, können Sie entweder:

1. das CLI-Tool "rvasp" verwenden
2. den rVASP-Protokollpuffer und direkte Interaktion mit der API verwenden

Um das CLI-Tool zu installieren, laden Sie entweder die ausführbare Datei `rvasp` für die entsprechende Architektur von der [TestNet-Releases-Seite] (https://github.com/trisacrypto/testnet/releases) herunter, klonen Sie [das TestNet-Repository] (https://github.com/trisacrypto/testnet/) und erstellen Sie das Binärprogramm `cmd/rvasp` oder installieren Sie es mit `go get` wie folgt:

```
$ go get github.com/trisacrypto/testnet/...
```

Um die [rVASP Protocol Buffer](https://github.com/trisacrypto/testnet/tree/main/proto/rvasp/v1) zu verwenden, klonen Sie sie oder laden Sie sie aus dem TestNet-Repository herunter und kompilieren Sie sie mit `protoc` in Ihre bevorzugte Sprache.

### Auslösen eines rVASP zum Senden einer Nachricht

Die rVASP-Admin-Endpunkte werden für die direkte Interaktion mit dem rVASP zu Entwicklungs- und Integrationszwecken verwendet. Beachten Sie, dass dieser Endpunkt sich vom TRISA-Endpunkt unterscheidet, der oben beschrieben wurde.

- Alice: `admin.alice.vaspbot.net:443`
- Bob: `admin.bob.vaspbot.net:443`
- Evil: `admin.evil.vaspbot.net:443`

Um das Kommandozeilentool zum Auslösen einer Nachricht zu verwenden, führen Sie den folgenden Befehl aus:

```
$ rvasp transfer -e admin.alice.vaspbot.net:443 \
        -a mary@alicevasp.us \
        -d 0.3 \
        -B trisa.example.com \
        -b cryptowalletaddress \
        -E
```

Diese Nachricht sendet dem Alice rVASP eine Nachricht mit dem Flag `-e` oder `--endpoint` und gibt mit dem Flag `-a` oder `--account` an, dass das Ursprungskonto "mary@alicevasp.us" sein soll. Anhand des Ursprungskontos wird bestimmt, welche IVMS101-Daten an den Empfänger gesendet werden sollen. Das `-d` oder `--amount` Flag gibt den zu sendenden Betrag an "AliceCoin" an.

Die nächsten beiden Teile sind entscheidend. Das `-E` oder `--external-demo` Flag sagt dem rVASP, dass er eine Anfrage an Ihren Dienst stellen soll, anstatt einen Demo-Austausch mit einem anderen rVASP durchzuführen. Dieses Flag ist erforderlich! Schließlich gibt das `-B` oder `--beneficiary-vasp` Flag an, wohin der rVASP die Anfrage senden wird. Dieses Feld sollte im TRISA-TestNet-Verzeichnisdienst nachgeschlagen werden können; z. B. sollte es Ihr allgemeiner Name oder der Name Ihres VASP sein, wenn dieser durchsuchbar ist.

Beachten Sie, dass Sie die Umgebungsvariablen `$RVASP_ADDR` und `$RVASP_CLIENT_ACCOUNT` setzen können, um die Flags `-e` bzw. `-a` festzulegen.

Um die Protokollpuffer direkt zu nutzen, verwenden Sie den `TRISAIntegration`-Dienst `Transfer` RPC mit dem folgenden `TransferRequest`:

```json
{
    "account": "mary@alicevasp.us",
    "amount": 0.3,
    "beneficiary": "cryptowalletaddress",
    "beneficiary_vasp": "trisa.example.com",
    "check_beneficiary": false,
    "external_demo": true
}
```

Diese Werte haben genau die gleiche Spezifikation wie die im Befehlszeilenprogramm.

### Senden einer TRISA-Nachricht an einen rVASP

Der rVASP erwartet eine `trisa.data.generic.v1beta1.Transaction` als Transaktions-Payload und einen `ivms101.IdentityPayload` als Identitäts-Payload. Die Begünstigten-Informationen der Identitäts-Payload müssen nicht ausgefüllt werden, rVASP antwortet mit der Ausfüllung des Begünstigten, die Identitäts-Payload sollte jedoch nicht null sein. Es wird empfohlen, einige gefälschte Identitätsdaten anzugeben, um die Vorteile der rVASP-Parsing- und Validierungsbefehle zu nutzen.

Vergewissern Sie sich, dass Sie in Ihrer Transaktionsnutzlast eine begünstigte Wallet angeben, die mit den rVASP-Begünstigten aus der obigen Tabelle übereinstimmt, z. B. verwenden Sie:

```json
{
    "txid": "1234",
    "originator": "anydatahereisfine",
    "beneficiary": "1MRCxvEpBoY8qajrmNTSrcfXSZ2wsrGeha",
    "amount": 0.3,
    "network": "TestNet",
    "timestamp": "2021-06-14T16:41:52-05:00"
}
```

Sie können einen beliebigen `txid`- oder `originator`-String angeben, und die Felder `network` und `timestamp` werden ignoriert.

Erstellen Sie einen versiegelten Umschlag, indem Sie entweder den Verzeichnisdienst oder den direkten Schlüsselaustausch verwenden, um die öffentlichen RSA-Schlüssel des rVASP abzurufen, und verwenden Sie `AES256-GCM` und `HMAC-SHA256` als Umschlagkryptographie. Anschließend wird der versiegelte Umschlag mit dem RPC-Dienst `Transfer` des `TRISANetwork`-Dienstes an den rVASP gesendet.

TO-DO: Bald wird das `trisa` Kommandozeilenprogramm verfügbar sein. Wir werden hier angeben, wie man das CLI-Programm verwendet, um eine Nachricht zu senden, sobald es veröffentlicht ist.

