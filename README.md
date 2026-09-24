# Flash Apple Edge
Hier ist ein Programm, welches per Modbus periodisch Daten eines Walz Fruitguard Sensors abfragt. Die ausgelesen Daten werden pro Sensor in eine CSV geschrieben.

Die Sensoren werden über die die Datei konfig.ini konfiguriert, die im selben Ordner, wie die ausführbare Datei liegen muss.
In der Konfigdatei gibt man an, wie welche Sensoren angeschlossen sind. In der Datei sind Erläuterungen, was man einstellen kann.

Das Programm funktioniert sowohl auf Windows 64 bit und Linux aarch64 und x86_64. 


## Inbetriebnahme

1. Autogain ausführen
2. Yield Messung ausführen
2. Der Yield Wert kann max. 1000 werden. Je höher desto besser. 
3. Den Sensor so lange plazieren, bis ein maximaler Yield Wert erreicht wird.

## Messung
Eine F_Messung findet stündlich statt. Hierbei wird der Fluorezenzwert und der Yield Wert in eine CSV Datei geschrieben.

|Short|type|Beschreibung|Dauer|
|-|-|-|-|
|F1|uint16|Fluorezenzwert rechts|
|F2|uint16|Fluorezenzwert links|
|Y1|uint16|Yieldwert rechts|
|Y2|uint16|Yieldwert links|


### Trigger

Es gibt 4 Aktionen, die man auslösen kann. Die sind zur Inbetriennahme.

|Register|Zweck|Short|Dauer|
|-|-|-|-|
|100|Inbetriebnahme, Intensität ermitteln|Autogain|2 sek|
|101|Fluorezenzmessung stündlich|F_Messung|6 sek|
|102|Yield Messung täglich |Y_Messung|8 sek|
|103|Werkseinstellungen|Reset|2 sek|

## Persitenz CSV

Es werden stündliche Messungen (F_Messung) gemacht und dabei pro Unit in eine CSV Datei geschrieben
```CSV
file: fruitguard_ORT02_UNIT01.csv
"iso", "timestamp", "Zeit", "F1", "F2", "Y1", "Y2"
2026-09-16T09:08:52+02:00; 1789453566, 2026-09-15, 123,133, 890, 950
```
iso entspricht RFC3339, oder ISO 8601 (vermutlich...)
timestamp ist der UNIX timestamp, 

### Features
- Die CSV Datei kann heruntergeladen werden
- Die CSV Datei kann initialisiert werden
- Es wird immer angehängt


## Prometheus
Ein Prometheus Endpunkt wird zum Abspeichern der Messwerte und Einstellungen zu Verfügung gestellt

#### ufg_messwerte
- Metric Help: Fruitguard Messwert
- Labels: 
    - base: Kundenkennung
    - serial: Unit ID des Gerätes
    - ort: Raumnummer
    - messung: F1, F2, Y1, Y2

#### ufg_einstellung
- Metric Help: Fruitguard Einstellungen
- Labels: serial, wert, ort, base
- base: Kundenkennung
    - serial: Unit ID des Gerätes
    - ort: Raumnummer
    - short: 
        - ModusF         
        - IntervallF     
        - ModusY         
        - IntervallY     
        - ZeitpunktY     
        - DAC1           
        - DAC2           
        - Messfrequenz   
        - IntegratZeitF  
        - SAT1           
        - SAT2           
        - MinutenDesTages

## Aufruf Hierachy

```
package main
    |
    package device (notify changes, recieves trigger)
    |
    package web (refresh on change, transmit trigger)
        |  
        call device methods

pollClients
    Vorraussetzung: initDone
    Loop units 
        1. Trigger F_Messung unit
        2. Read Messwerte (m.Read)
        3. Upload Metrics (prometheus.Gauge.Set(float)) 
            (Voraussetung webServer running -> handerPrometheus OK)
        4. Refresh unit-html, 
                falls User Online
        5. Publish Nats-Topic

```


# Programmatische Abstimmung

### Probleme:

1. Sicherstellen, das bei der Initialisierung alles geladen ist und nichts abstürzt. Dazu muss eine Reihenfolge eingehalten werden.
    - KONFIG_FILE
    - INIT_NATS
    - INIT_WEBSERVER
    - INIT_UNITS        
    - INIT_METRIC       - Fehlerfälle?
    - READY_TO_POLL     - Fehlerfälle?
    - INIT_SCHEDULER    - Fehlerfälle?
    - RUNNING           - msg to topic:Errors || Data
    - SHUTDOWN

Im Fehlerfall - Logdatei-Eintrag, UI-Hinweistext


2. Neue Daten signalisieren und damit Aktionen auslösen
    - UI aktulisieren
    - Auf Interaktion reagieren
     - Fehler behandeln

### Pipeline Channel (Init-Aktion)
Über die Pipeline kommunizieren die Programbausteine untereinander. 
Sie kommunizieren was als nächstes zu tun ist. 
D.h, welche Programmbaustein als nächstes aufgerufen wird.
Ein switch-case ruft den nächsten Baustein auf.

### Message Bus (Betriebs-Aktionen)
Der Message-Bus sendet Benutzerinteraktion, neue Daten und Fehler.
Andere Programmbausteine sollen dann Aktionen auslösen. In etwa UI aktualisieren, metrics aktualierem, csv Datei schreiben.