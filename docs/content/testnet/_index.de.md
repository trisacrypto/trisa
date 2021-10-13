---
title: Das TRISA TestNet
date: 2021-06-14T11:20:10-05:00
lastmod: 2021-10-08T16:34:11-05:00
description: "Das TRISA TestNet"
weight: 30
---

Das TRISA TestNet wurde eingerichtet, um eine Demonstration des TRISA Peer-to-Peer-Protokolls zu ermöglichen, "Roboter-VASP"-Dienste zur Erleichterung der TRISA-Integration zu beherbergen und den primären TRISA-Verzeichnisdienst zu betreiben, der den Austausch öffentlicher Schlüssel und die Ermittlung von Peers erleichtert.

{{% figure src="/img/testnet_architecture.png" %}}

Das TRISA TestNet besteht aus den folgenden Diensten:

- [TRISA Directory Service](https://vaspdirectory.net) - eine Benutzeroberfläche, um den TRISA Global Directory Service zu erkunden und sich als TRISA-Mitglied zu registrieren
- [TestNet Demo](https://vaspbot.net) - eine Demoseite, die TRISA-Interaktionen zwischen "Roboter"-VASPs zeigt, die im TestNet laufen

Das TestNet beherbergt auch drei Roboter-VASPs oder rVASPs, die als Annehmlichkeit für TRISA-Mitglieder implementiert wurden, um ihre TRISA-Dienste zu integrieren. Der primäre rVASP ist Alice, ein sekundärer für Demozwecke ist Bob, und um Interaktionen mit nicht verifizierten TRISA-Mitgliedern zu testen, gibt es auch einen "bösartigen" rVASP.
