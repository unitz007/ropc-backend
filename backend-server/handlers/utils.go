package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const (
	jsonDecodeError = "Couldn't decode JSON: "
)

func JsonToStruct[T any](r io.ReadCloser, t T) error {
	err := json.NewDecoder(r).Decode(t)
	if err != nil {
		log.Println("error = ", err)
		return errors.New(jsonDecodeError + err.Error())
	}

	return nil
}

func PrintResponse[T any](statusCode int, res http.ResponseWriter, payload T) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	err := json.NewEncoder(res).Encode(payload)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		_, err = res.Write([]byte("Invalid response"))
		if err != nil {
			return err
		}
	}

	return nil
}
