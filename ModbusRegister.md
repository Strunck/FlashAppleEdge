
# Fruitguard Modbus Register

## Messwerte

|Register   |Funktion                     |Bereich             |
|-----------|-----------------------------|--------------------|
|10         |Flags                        |                    |
|11         |fortlaufende Nummer F-Messung|0 – 65535           |
|12         |fortlaufende Nummer Y-Messung|0 – 65535           |
|13         |F #1                         |0 – 5000            |
|14         |F #2                         |0 – 5000            |
|15         |Y #1                         |0 – 1000, entspricht|
|16         |Y #2                         |0.000 – 1.000       |
|17         |Fo #1                        |0 – 5000            |
|18         |Fo #2                        |0 – 5000            |
|19         |Fm #1                        |0 – 5000            |
|20         |Fm #2                        |0 – 5000            |
|21         |Minuten des Tages            |0 – 1439            |

Als Zeitbasis Minuten, da die Modbus Register 16 Bit breit sind und 86400 Sekunden da nicht hineinpassen.
Über die „fortlaufend Nummer“ kann man erkennen, dass ein neues Messergebnis vorhanden ist, die läuft immer im Kreis herum, auf 65535 folgt wieder 0.
Flags enthalten Fehler und Status Meldungen, z.B. Signal Low, muss ich noch festlegen.
Man würde also Register 10 bis 17 lesen.
Die Uhrzeit wird nicht sehr stabil sein, wenn man das wollte müsste man einen Uhrenchip einbauen, daher macht es vermutlich Sinn, sie einmal am Tag zu synchronisieren.

## Trigger, zum Auslösen 1 schreiben, lesen immer 0

|Register   |Funktion                |
|-----------|------------------------|
|100        | Autogain               |
|101        | F-Messung              |
|102        | Y-Messung              |
|103        | Reset Default Settings |

Der Autogain misst die korrekte Messlicht Intensität ein, das würde man bei der Inbetriebnahme mit neuen Äpfeln machen.
Es war noch zu testen, ob der Autogain in der Praxis sinnvoll ist.
Alternativ ist die Messlicht Intensität auch manuell einstellbar.
F/Y-Messung könnten manuell getriggert werden.

## Einstellungen (lesen/schreiben):

|Register|Funktion   |Werte|
|--------|-------------------|-----------------------------------------------|
|200     |Modus F-Messung    | 0 = Intervall, 1 = getriggert, 2 = fortlaufend|
|201     |Intervall F-Messung| in Minuten                                    |
|202     |Modus Y-Messung    | 0 = Intervall, 1 = getriggert, 2 = Zeitpunkt  |
|203     |Intervall Y-Messung| in Minuten                                    |
|204     |Zeitpunkt Y-Messung| Uhrzeit in Minuten des Tages (0 – 1439)       |
|205     |DAC #1             | Messlicht Intensität, 0 – 4095, vom Autogain  |
|206     |DAC #2             | eingemessen, könnte man auch manuell setzen   |
|207     |Messfrequenz       | 1 – 5, Häufigkeit der Messpulse               |
|208     |Integrationszeit F |                                               |
|209     |SAT #1             | Satpuls Intensität                            |
|210     |SAT #2             |                                               |
|211     |Minuten des Tages  | 0 - 1439                                      |

Die F-Messung würde im Intervall-Modus laufen, mit 60 Minuten Intervall
Die Y-Messung würde zu einem bestimmten Zeitpunkt, einmal pro Tag ausgeführt.
Mit Register 211 kann man die Uhrzeit synchronisieren.

## Kommandozeile

Unter Linux gibt es den Befehl mbpoll. Hiermit können Low Level Test gemacht werden, ob die Verbindung zu den Sensoren funtioniert.

`mbpoll -P none -m rtu -a 33 -t 3   -r 10 -c 8  /dev/ttyS1`

Erläuterung
- `-P none` Parität auf "none" einstellen, Wichtig, das sonst ein CRC Fehler kommt
- `-m rtu` stellt auf RS485-Modus um. 
- `-a 33` stellt die Unit ID ein,  bzw Modbus-Adresse 
- `-t 3` stellt sicher, das ein 16 bit Integer abgefragt wird
- `-r 10` starte bei Register 10
- `-c 8` frage 8 Register ab
- `/dev/ttyS1` gibt den Serial Port an 


Version 0.2 / 19.08.2026
