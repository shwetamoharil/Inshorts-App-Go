package utils

import (
	"Inshorts/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	FILENAME = "/home/shweta/go/src/Inshorts/articles.json"
	LIMIT    = 10
)

func SetHeaders(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		f(w, r)
	}
}

func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)

	return os.IsNotExist(err)
}

func ReadFile(filename string) ([]byte, error) {
	ok := checkFileExists(filename)
	if !ok {
		file, err := os.Open(filename)
		if err != nil {
			return []byte{}, err
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return []byte{}, err
		}
		return data, nil
	}
	return []byte{}, errors.New("error reading file")
}

func AddDataToJsonFile(article models.Article) error {
	var result []models.Article
	data, err := ReadFile(FILENAME)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		result = append(result, article)
	} else {
		err = json.Unmarshal(data, &result)
		if err != nil {
			return nil
		}
		result = append(result, article)
	}

	dataBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(FILENAME, dataBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func HandleHttpRequestErrors(w http.ResponseWriter, customErrMsg string, statusCode int, err error) (isError bool) {
	if err != nil {
		errMsg := err.Error()
		if customErrMsg != "" {
			errMsg = customErrMsg
		}
		response := models.ErrorResponse{
			ErrorMessage: errMsg,
			StatusCode:   statusCode,
		}
		message, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
		w.WriteHeader(response.StatusCode)
		_, err = w.Write(message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return true
		}
		return true
	}
	return
}

func ConvertStringToInt(s string) (int, error) {
	number, err := strconv.Atoi(s)
	return number, err
}
