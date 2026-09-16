// Diese Datei entählt Code, zum Abspeichern in eine CSV Datei und zur Bereitstellung über prometheus.

package device

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// iso nach RFC3339
	csvHeader = "iso;timestamp;F1;F2;Y1;Y2"
)

// Dateinamen in der Form fruitguard_<ort>_ID<id>.csv , z.B fruitguard_R01_ID02.csv
func (s *State) intiFiles(folderName string) (err error) {
	// Bei Start prüfen ob der Ordner vorhanden und beschreibbar ist.
	if _, err = os.Stat(folderName); os.IsNotExist(err) {
		if err = os.Mkdir(folderName, 0755); err != nil {
			// Ordner existiert nicht, daher Fehler zurückgeben
			return fmt.Errorf("intiFiles: Ordner zum ablegen der CSV Dateien konnte nicht erstellt werden. Bitte erstellen Sie den Ordner %s: %v", folderName, err)
		}
	}

	// Dateinamen aus Config generieren und im FileManager speichern
	for l, line := range s.cfg.Lines {
		for i, id := range line.Id {

			fileName := fmt.Sprintf("fruitguard_%s_ID%02d.csv", line.Ort, id)
			filePath := filepath.Join(folderName, fileName)

			s.units[l][i].f, err = os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("intiFiles: failed to open CSV file %s: %v", filePath, err)
			}
			// Schreibe den CSV-Header, falls weniger als eine Zeile hat
			fileInfo, _ := s.units[l][i].f.Stat()
			if fileInfo.Size() == 0 {
				if _, err := s.units[l][i].f.WriteString(csvHeader + "\n"); err != nil {
					return fmt.Errorf("intiFiles: failed to write CSV header to file %s: %v", filePath, err)
				}
			}
		}
	}

	// Prüfen ob es bereits CSV Dateien gibt, und diese im Modul "Append" öffnen
	// Die FileDescriptoren sollen im struct State als Array gespeichert werden.
	// Öffne alle Dateien, die mit fruitguard_ beginnen und auf .csv enden.
	entries, err := os.ReadDir(folderName)
	if err != nil {
		return fmt.Errorf("intiFiles: Ordner konnte nicht gelesen werden: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "fruitguard") || !strings.HasSuffix(name, ".csv") {
			continue
		}

		filePath := filepath.Join(folderName, name)
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("intiFiles: failed to open CSV file %s: %v", filePath, err)
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("intiFiles: failed to close CSV file %s: %v", filePath, err)
		}
	}

	return nil
}

func (u Unit) writeLineInCSV() (err error) {
	_, err = u.f.WriteString(u.r.Mess.printCSV())
	if err != nil {
		return fmt.Errorf("writeCSV: failed to write to CSV file: %v", err)
	}
	return nil
}

func (s *State) CloseAllFiles() (err error) {
	for l, line := range s.units {
		for i, unit := range line {
			if unit.f != nil {
				if err := unit.f.Close(); err != nil {
					return fmt.Errorf("CloseAllFiles: failed to close CSV file for line %d unit %d: %v", l+1, i+1, err)
				}
			}
		}
	}
	return nil
}
