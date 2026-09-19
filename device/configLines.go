package device

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// Steht für eine Linie, also, RS485 Verbindung mit mehren Geräten
type Line struct {
	ConType string
	Url     string
	Ort     string
	Id      []int
	IdCount int // Anzahl der Units in der Linie
}

type Config struct {
	Lines []Line // Linie ist
	Count int
	Base  string
}

func LoadConfig() (cfg Config, err error) {
	data, err := ini.Load("konfig.ini")
	if err != nil {
		return Config{}, fmt.Errorf("Failed to read file: %v", err)
	}

	cfg.Count, err = data.Section("").Key("LINES").Int()
	if err != nil {
		return Config{}, fmt.Errorf("Failed to read LINES: %v", err)
	}

	cfg.Base = data.Section("").Key("BASE").String()
	if cfg.Base == "" {
		fmt.Println("BASE wurde nicht gesetzt. Siehe konfig.ini ")
		cfg.Base = "DEFAULT"
	}

	cfg.Lines = make([]Line, cfg.Count)
	// loop von 1 nach 4
	for i := 1; i <= cfg.Count; i++ {
		s := fmt.Sprintf("Line%d", i)
		line := Line{
			ConType: data.Section(s).Key("CONTYPE").String(),
			Url:     data.Section(s).Key("URL").String(),
			Ort:     data.Section(s).Key("ORT").String(),
			Id:      data.Section(s).Key("ID").Ints(","),
			IdCount: len(data.Section(s).Key("ID").Ints(",")),
		}
		cfg.Lines[i-1] = line
	}

	return cfg, nil
}

func PrintConfig(cfg Config) {
	fmt.Println("Konfiguration Anlage: ", cfg.Base)
	for i, line := range cfg.Lines {
		fmt.Printf("Line %d: Ort %s - %s %s %v\n", i+1, line.Ort, line.ConType, line.Url, line.Id)
	}
}
