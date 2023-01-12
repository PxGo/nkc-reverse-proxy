package modules

import "net/http"

func WriteResponse(writer http.ResponseWriter, status int, content []byte) error {
	writer.WriteHeader(status)
	_, err := writer.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func WriteResponsePage(writer http.ResponseWriter, status int) error {
	pageContent, err := GetPageByStatus(status)
	if err != nil {
		return err
	}
	err = WriteResponse(writer, status, pageContent)
	if err != nil {
		return err
	}
	return nil
}
