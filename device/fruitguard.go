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
	return nil
}

func (m Messung) Print() {
	fmt.Printf("Flags: %d\n", m.Flags)
	fmt.Printf("NummerF: %d\n", m.NummerF)
	fmt.Printf("NummerY: %d\n", m.NummerY)
	fmt.Printf("F1: %d\n", m.F1)
	fmt.Printf("F2: %d\n", m.F2)
	fmt.Printf("Y1: %d\n", m.Y1)
	fmt.Printf("Y2: %d\n", m.Y2)
	fmt.Printf("Fo1: %d\n", m.Fo1)
	fmt.Printf("Fo2: %d\n", m.Fo2)
	fmt.Printf("Fm1: %d\n", m.Fm1)
	fmt.Printf("Fm2: %d\n", m.Fm2)
	fmt.Printf("Minute: %d\n", m.Minute)
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
		// Write Autogain trigger
		err = trigger(client, 100)
	case F_Messung:
		// Write F_Messung trigger
		err = trigger(client, 101)
	case Y_Messung:
		// Write Y_Messung trigger
		err = trigger(client, 102)
	case Reset:
		// Write Reset trigger
		err = trigger(client, 103)
	default:
		return fmt.Errorf("unknown trigger type: %v", tt)
	}
	return err
}

func trigger(client modbus.Client, address uint16) (err error) {
	res, err := client.WriteSingleRegister(address, 1)
	if err != nil {
		return fmt.Errorf("failed to write single register at address %d: %w, response: %v", address, err, res)
	}
	// Timeout for 2 seconds
	time.Sleep(2 * time.Second)

	res, err = client.WriteSingleRegister(address, 0)
	if err != nil {
		return fmt.Errorf("failed to write single register at address %d: %w, response: %v", address, err, res)
	}

	return nil
}
