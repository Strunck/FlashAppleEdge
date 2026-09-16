

## Messung

|Short|type|Beschreibung|Dauer|
|-|-|-|-|
|F1|uint16|Fluorezenzwert rechts|
|F2|uint16|Fluorezenzwert links|
|Y1|uint16|Yieldwert rechts|
|Y2|uint16|Yieldwert links|


### Trigger
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
