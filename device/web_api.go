package device

import (
	"encoding/json"
	"net/http"
)

func (m *Messung) toJSON() (res []byte, err error) {
	if res, err = json.Marshal(m); err != nil {
		return nil, err
	}
	return res, nil
}

func (e *Einstellungen) toJSON() (res []byte, err error) {
	res, err = json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u Unit) toJSON() (res []byte, err error) {
	res, err = json.Marshal(u)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ulist State) allUnitsToJSON() (res []byte, err error) {
	res, err = json.Marshal(ulist.Units)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func HandleAPI(ds State) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implement your API handling logic here
		w.Header().Set("Content-Type", "application/json")
		data, err := ds.allUnitsToJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
}

func HandleGetConfig(ds State) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := ds.Cfg.toJSON()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
}
