package device

import (
	"encoding/binary"
	"fmt"
	"time"

	modbus "github.com/goburrow/modbus"
)

// Messung 12 Register
type Messung struct {
	Flags   uint16 `register:"10"`
	NummerF uint16 `register:"11"`
	NummerY uint16 `register:"12"`
	F1      uint16 `register:"13"`
	F2      uint16 `register:"14"`
	Y1      uint16 `register:"15"`
	Y2      uint16 `register:"16"`
	Fo1     uint16 `register:"17"`
	Fo2     uint16 `register:"18"`
	Fm1     uint16 `register:"19"`
	Fm2     uint16 `register:"20"`
	Minute  uint16 `register:"21"`
	gauges  MessungGauges
}

type Trigger struct {
	Autogain  uint16 `register:"100"`
	F_Messung uint16 `register:"101"`
	Y_Messung uint16 `register:"102"`
	Reset     uint16 `register:"103"`
}

type TriggerType int

const (
	Autogain TriggerType = iota
	F_Messung
	Y_Messung
	Reset
)

// Messung
type Einstellungen struct {
	ModusF          uint16 `register:"200"` //Modus F-Messung
	IntervallF      uint16 `register:"201"` //Intervall F-Messung
	ModusY          uint16 `register:"202"` //Modus Y-Messung
	IntervallY      uint16 `register:"203"` //Intervall Y-Messung
	ZeitpunktY      uint16 `register:"204"` //Zeitpunkt Y-Messung
	DAC1            uint16 `register:"205"` //DAC #1
	DAC2            uint16 `register:"206"` //DAC #2
	Messfrequenz    uint16 `register:"207"` //Messfrequenz
	IntegratZeitF   uint16 `register:"208"` //Integrationszeit F
	SAT1            uint16 `register:"209"` //SAT #1
	SAT2            uint16 `register:"210"` //SAT #2
	MinutenDesTages uint16 `register:"211"` //Minuten des Tages
}

type Register struct {
	Mess Messung
	Trig Trigger
	Eins Einstellungen
}

func (m *Messung) Read(cl modbus.Client) (err error) {
	results, err := cl.ReadHoldingRegisters(10, 12)
	if err != nil {
		return fmt.Errorf("failed to read holding registers 10-21: %w", err)
	}
	// Parse results into the Messung struct fieldsModus F-Messung
	m.Flags = binary.BigEndian.Uint16(results[0:2])
	m.NummerF = binary.BigEndian.Uint16(results[2:4])
	m.NummerY = binary.BigEndian.Uint16(results[4:6])
	m.F1 = binary.BigEndian.Uint16(results[6:8])
	m.F2 = binary.BigEndian.Uint16(results[8:10])
	m.Y1 = binary.BigEndian.Uint16(results[10:12])
	m.Y2 = binary.BigEndian.Uint16(results[12:14])
	m.Fo1 = binary.BigEndian.Uint16(results[14:16])
	m.Fo2 = binary.BigEndian.Uint16(results[16:18])
	m.Fm1 = binary.BigEndian.Uint16(results[18:20])
	m.Fm2 = binary.BigEndian.Uint16(results[20:22])
	m.Minute = binary.BigEndian.Uint16(results[22:24])
	m.setGauges()
	return nil
}

func (m Messung) Print() {
	/*
		----------------------- Messung ---------------------
		Zeit: 2026-09-01 12:00:00
		F1: 123 - F2: 456 - Y1: 789 - Y2: 012
		Flags: 4 - NummerF/Y: 1/2 - Fo: 123/123 Fm: 123/123
		Minuten des Tages: 123
		-----------------------------------------------------
	*/
	fmt.Printf("Zeit: %s\n", NettesDatum())
	fmt.Printf("F1: %d - F2: %d - Y1: %d - Y2: %d\n", m.F1, m.F2, m.Y1, m.Y2)
	fmt.Printf("Flags: %d - NummerF/Y: %d/%d - Fo: %d/%d Fm: %d/%d\n", m.Flags, m.NummerF, m.NummerY, m.Fo1, m.Fo2, m.Fm1, m.Fm2)
	fmt.Printf("Minuten des Tages: %d\n", m.Minute)
}

func (r Register) Print(i, j int) {
	fmt.Printf("-------- Line: %d, Slave: %d -----------------\n", i+1, j+1)
	r.Mess.Print()
	fmt.Printf("--------------------------------------------\n")
}

func (m Messung) printCSV() string {
	return fmt.Sprintf("%s;%d;%d;%d;%d;%d", time.Now().Format(time.RFC3339), time.Now().Unix(), m.F1, m.F2, m.Y1, m.Y2)
}

func (e *Einstellungen) Read(client modbus.Client) (err error) {
	results, err := client.ReadHoldingRegisters(200, 12)
	if err != nil {
		return fmt.Errorf("failed to read holding registers 200-211: %w", err)
	}
	e.ModusF = binary.BigEndian.Uint16(results[0:2])
	e.IntervallF = binary.BigEndian.Uint16(results[2:4])
	e.ModusY = binary.BigEndian.Uint16(results[4:6])
	e.IntervallY = binary.BigEndian.Uint16(results[6:8])
	e.ZeitpunktY = binary.BigEndian.Uint16(results[8:10])
	e.DAC1 = binary.BigEndian.Uint16(results[10:12])
	e.DAC2 = binary.BigEndian.Uint16(results[12:14])
	e.Messfrequenz = binary.BigEndian.Uint16(results[14:16])
	e.IntegratZeitF = binary.BigEndian.Uint16(results[16:18])
	e.SAT1 = binary.BigEndian.Uint16(results[18:20])
	e.SAT2 = binary.BigEndian.Uint16(results[20:22])
	e.MinutenDesTages = binary.BigEndian.Uint16(results[22:24])
	return nil
}

func (t *Trigger) Write(client modbus.Client, tt TriggerType) (err error) {

	switch tt {
	case Autogain:
		fmt.Println("Autogain trigger, bitte 2 sek warten")
		err = trigger(client, 100, 2*time.Second)
	case F_Messung:
		fmt.Println("F_Messung trigger, bitte 6 sek warten")
		err = trigger(client, 101, 6*time.Second)
	case Y_Messung:
		fmt.Println("Y_Messung trigger, bitte 8 sek warten")
		err = trigger(client, 102, 8*time.Second)
	case Reset:
		fmt.Println("Reset trigger, bitte 1 sek warten")
		err = trigger(client, 103, 1*time.Second)
	default:
		return fmt.Errorf("unknown trigger type: %v", tt)
	}

	return err
}

func trigger(client modbus.Client, address uint16, sek time.Duration) (err error) {
	res, err := client.WriteSingleRegister(address, 1)
	if err != nil {
		return fmt.Errorf("failed to write single register at address %d: %w, response: %v", address, err, res)
	}
	time.Sleep(sek)

	return nil
}

func NettesDatum() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
