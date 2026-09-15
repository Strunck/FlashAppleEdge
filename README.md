

## Messung

|Short|type|Beschreibung|
|-|-|-|
|F1|uint16|Fluorezenzwert rechts|
|F2|uint16|Fluorezenzwert links|
|Y1|uint16|Yieldwert rechts|
|Y2|uint16|Yieldwert links|

Es werden stündliche Messungen gemacht und dabei pro Unit in eine CSV Datei geschrieben
```CSV
file: fruitguard_ORT02_UNIT01.csv
"UTC", "Datum", "Zeit", "F1", "F2", "Y1", "Y2"
1789453566, 2026-09-15, 123,133, 890, 950
```

### Features
- Die CSV Datei kann heruntergeladen werden
- Die CSV Datei kann initialisiert werden
- Es wird immer angehängt


### Prometheus

Metric Name: ufg_messwerte
Metric Help: Fruitguard Messwert
Labels: serial, messung, ort, base
- base: Kundenkennung
- serial: Unit ID des Gerätes
- ort: Raumnummer
- messung: F1, F2, Y1, Y2

Metric Name: ufg_einstellung
Metric Help: Fruitguard Einstellungen
Labels: serial, wert, ort, base
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
    