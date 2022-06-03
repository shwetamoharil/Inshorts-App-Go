package controllers

import (
	"Inshorts/models"
	"Inshorts/utils"
	"encoding/json"
	"errors"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
	if article.Title == "" || article.SubTitle == "" || article.Content == "" {
		if isError := utils.HandleHttpRequestErrors(w, "all fields are compulsory", http.StatusBadRequest, errors.New("")); isError {
			return
		}
	}

	uuid, err := exec.Command("uuidgen").Output()
	if err != nil {
		return
	}

	article.Id = strings.Trim(string(uuid), "\n")
	article.CreationTimestamp = time.Now()

	err = utils.AddDataToJsonFile(article)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
	err = json.NewEncoder(w).Encode(article)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
}

func GetArticleById(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadFile(utils.FILENAME)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	var articles []models.Article
	err = json.Unmarshal(data, &articles)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	var response models.Article
	var found bool = false

	for _, article := range articles {
		if article.Id == id {
			response = article
			found = true
			break
		}
	}

	if !found {
		if isError := utils.HandleHttpRequestErrors(w, "", http.StatusNotFound, errors.New("no record found")); isError {
			return
		}
	}
	err = json.NewEncoder(w).Encode(response)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	data, err := utils.ReadFile(utils.FILENAME)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	var response []models.Article
	err = json.Unmarshal(data, &response)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	strLimit := r.URL.Query().Get("limit")
	limit := -1
	if strLimit != "" {
		limit, err = utils.ConvertStringToInt(strLimit)
		if err != nil || limit == -1 {
			if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, errors.New("invalid limit")); isError {
				return
			}
		}
	}

	strOffset := r.URL.Query().Get("offset")
	offset := -1
	if strOffset != "" {
		offset, err = utils.ConvertStringToInt(strOffset)
		if err != nil || offset == -1 {
			if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, errors.New("invalid offset")); isError {
				return
			}
		}
	}

	if limit != -1 && offset == -1 {
		end := limit
		if end > len(response) {
			end = len(response)
		}
		response = response[:end]
	} else if limit != -1 && offset != -1 {
		start := offset + 1
		if start > len(response) {
			start = len(response)
		}
		end := offset + limit
		if end > len(response) {
			end = len(response)
		}

		response = response[start:end]
	}

	err = json.NewEncoder(w).Encode(response)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
}

func SearchArticle(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("q")
	data, err := utils.ReadFile(utils.FILENAME)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	var articles []models.Article
	err = json.Unmarshal(data, &articles)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}

	var response models.Article
	for _, article := range articles {
		if strings.Contains(article.Title, term) {
			response = article
			break
		} else if strings.Contains(article.SubTitle, term) {
			response = article
			break
		} else if strings.Contains(article.Content, term) {
			response = article
			break
		}
	}

	err = json.NewEncoder(w).Encode(response)
	if isError := utils.HandleHttpRequestErrors(w, "", http.StatusInternalServerError, err); isError {
		return
	}
}
