package cet

import (
	"net/http"
	"fmt"
	"net/url"
	"encoding/json"
	"strconv"
)

func WriteJson(w http.ResponseWriter, body interface{}) {
	jsonBody, err := json.Marshal(body); if err != nil {
		WriteError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonBody))
}

func WriteError(w http.ResponseWriter, errMsg string) {
	errData := map[string]string{
		"error": errMsg,
	}
	WriteJson(w, errData)
}

func ParseQuery(query string, required []string) (map[string]string, error) {
	m, err := url.ParseQuery(query); if err != nil {
		return nil, err
	}
	for _, rp := range required {
		_, ok := m[rp]; if !ok {
			return nil, fmt.Errorf("Param `%s` is missing.", rp)
		}
	}

	r := map[string]string{}
	for k, v := range m {
		if len(v) >= 1 {
			r[k] = v[0]
		}
	}
	return r, nil
}

func FindTicketHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{"province", "school", "name", "type"}
	queryParam, err := ParseQuery(r.URL.RawQuery, required); if err != nil {
		WriteError(w, err.Error())
		return
	}

	CETType, err := strconv.Atoi(queryParam["type"]); if err != nil || CETType > 2 || CETType < 1 {
		WriteError(w, "Unknow type.")
		return
	}

	ticket, err := FindTicketNumber(queryParam["province"], queryParam["school"], queryParam["name"], "", CETType); if err != nil {
		WriteError(w, err.Error())
		return
	}
	WriteJson(w, map[string]string{"ticket": ticket, })
}

func GetScoreHandler(w http.ResponseWriter, r *http.Request) {
	required := []string{"name", "ticket", "type"}
	queryParam, err := ParseQuery(r.URL.RawQuery, required); if err != nil {
		WriteError(w, err.Error())
		return
	}
	CETType, err := strconv.Atoi(queryParam["type"]); if err != nil || CETType > 2 || CETType < 1 {
		WriteError(w, "Unknow type.")
		return
	}

	score, err := GetScore(queryParam["ticket"], queryParam["name"]); if err != nil {
		WriteError(w, err.Error())
		return
	}
	WriteJson(w, score)
}

func NewCETServer(addr string, handler http.Handler) error {
	http.HandleFunc("/ticket", FindTicketHandler)
	http.HandleFunc("/score", GetScoreHandler)
	return http.ListenAndServe(addr, handler)
}

