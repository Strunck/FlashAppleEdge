

## Messung
F1 uint16   //Fluorezenzwert rechts
F2 uint16   //Fluorezenzwert links
Y1 uint16   //Yieldwert rechts
Y2 uint16   //Yieldwert links

### Prometheus

Metric Name: ufg_messwerte
Metric Help: Fruitguard Messwert
Labels: serial, messung, ort
- serial: Unit ID des Gerätes
- ort: Raumnummer
- messung: F1, F2, Y1, Y2
