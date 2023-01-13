package modules

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func WriteResponseContent(writer http.ResponseWriter, status int, content []byte) error {
	writer.WriteHeader(status)
	_, err := writer.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func WriteResponseHTML(writer http.ResponseWriter, status int) error {
	pageContent, err := GetPageByStatus(status)
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", "text/html")
	err = WriteResponseContent(writer, status, pageContent)
	if err != nil {
		return err
	}
	return nil
}

func WriteResponseJSON(writer http.ResponseWriter, status int) error {
	statusText := http.StatusText(status)
	writer.Header().Set("Content-Type", "application/json")
	body := struct {
		Code    int    `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}{
		Code:    0,
		Type:    statusText,
		Message: strconv.Itoa(status) + " " + statusText,
	}
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = WriteResponseContent(writer, status, bodyJson)
	if err != nil {
		return err
	}
	return nil
}

func WriteResponse(request *http.Request, writer http.ResponseWriter, status int) error {

	isJson := strings.Contains(request.Header.Get("Accept"), "application/json")

	var err error

	if isJson {
		err = WriteResponseJSON(writer, status)
	} else {
		err = WriteResponseHTML(writer, status)
	}

	if err != nil {
		return err
	}

	return nil
}
