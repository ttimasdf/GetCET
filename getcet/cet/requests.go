package cet

import (
	"strings"
	"math/rand"
	"net/http"
	"fmt"
	"bytes"
	"io/ioutil"
	"net/url"
)

const (
	CET4 = 1
	CET6 = 2
)

type Score struct {
	Name      string    `json:"name"`
	School    string    `json:"school"`
	Listening string    `json:"listening"`
	Reading   string    `json:"reading"`
	Writing   string    `json:"writing"`
	Total     string    `json:"total"`
}

func random(min int, max int) int {
	return rand.Intn(max - min) + min
}

func randomMAC() string {
	randUnit := [6]string{}
	for i := 0; i < 6; i++ {
		randUnit[i] = fmt.Sprintf("%.2X", random(0, 16))
	}
	return strings.Join(randUnit[:], "-")
}

func FindTicketNumber(province string, school string, name string, examroom string, CETType int) (string, error) {
	provinceId, ok := PROVINCE[province]; if !ok {
		return "", fmt.Errorf("Not supported province: %s.", province)
	}
	payload := fmt.Sprintf("type=%d&provice=%d&school=%s&name=%s&examroom=%s&m=%s", CETType, provinceId, school, name, examroom, randomMAC())
	gbkPayload, err := Utf8ToGbk([]byte(payload)); if err != nil {
		return "", err
	}

	cipher := NewCETCipher(TICKET_KEY, REQUEST_KEY)
	encPayload, err := cipher.EncryptRequest([]byte(gbkPayload))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", SEARCH_URL, bytes.NewReader(encPayload)); if err != nil {
		return "", err
	}
	defer req.Body.Close()

	req.Header.Set("User-Agent", USER_AGENT)

	client := &http.Client{}
	resp, err := client.Do(req); if err != nil {
		return "", err
	}
	defer req.Body.Close()

	respContent, err := ioutil.ReadAll(resp.Body); if err != nil || len(respContent) < 2 {
		return "", fmt.Errorf("Can't find your ticket.")
	}

	respContent = respContent[2:] // Remove Magic Number
	ticket, err := cipher.DecryptTicket(respContent); if err != nil {
		return "", err
	}

	return string(ticket[:]), nil
}

func GetScore(ticket string, name string) (*Score, error) {
	gbkName, err := Utf8ToGbk([]byte(name)); if err != nil {
		return nil, err
	}
	payload := url.Values{"id": {ticket}, "name": {string(gbkName)}}

	req, err := http.NewRequest("POST", SCORE_URL, strings.NewReader(payload.Encode())); if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Referer", "http://cet.99sushe.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req); if err != nil {
		return nil, fmt.Errorf("Get Score Failed.")
	}
	defer resp.Body.Close()

	scoreData, err := ioutil.ReadAll(resp.Body); if err != nil {
		return nil, err
	}
	scoreData, err = GbkToUtf8(scoreData); if err != nil {
		return nil, err
	}
	scoreArr := strings.Split(string(scoreData), ","); if len(scoreArr) != 7 {
		return nil, fmt.Errorf("Parse Score Data Failed.")
	}

	return &Score{
		Name:scoreArr[6],
		School:scoreArr[5],
		Listening:scoreArr[1],
		Reading:scoreArr[2],
		Writing:scoreArr[3],
		Total:scoreArr[4],
	}, nil
}