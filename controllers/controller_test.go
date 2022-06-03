package controllers

import (
	"Inshorts/models"
	"Inshorts/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetAllArticles(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllArticles)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	data, err := utils.ReadFile(utils.FILENAME)
	if err != nil {
		t.Fatal(err)
	}

	var expectedArtilces []models.Article
	err = json.Unmarshal(data, &expectedArtilces)
	if err != nil {
		t.Fatal(err)
	}

	var gotArticles []models.Article
	err = json.NewDecoder(rr.Body).Decode(&gotArticles)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedArtilces, gotArticles) {
		t.Errorf("handler returned unexpected body: got %v want %v", gotArticles, expectedArtilces)
	}
}

func TestCreateArticle(t *testing.T) {
	article := models.Article{
		Title:    "Test 1",
		SubTitle: "Subtitile 1",
		Content:  "Content for article test 1",
	}

	jsonStr, err := json.Marshal(article)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateArticle)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var gotArticle models.Article
	err = json.NewDecoder(rr.Body).Decode(&gotArticle)
	if err != nil {
		t.Fatal(err)
	}

	if article.Content != gotArticle.Content {
		t.Errorf("handler returned unexpected body: got %v want %v", gotArticle, article)
	}
}

func TestGetArticleById(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles/tttt", nil)
	if err != nil {
		t.Fatal(err)
	}

	vars := map[string]string{
		"id": "tttt",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetArticleById)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var gotArticle models.Article
	err = json.NewDecoder(rr.Body).Decode(&gotArticle)
	if err != nil {
		t.Fatal(err)
	}

	expected := models.Article{
		Id:                "tttt",
		Title:             "RIP KK",
		SubTitle:          "Legendary singer passes away at 55",
		Content:           "Singer KK passed away ",
		CreationTimestamp: gotArticle.CreationTimestamp,
	}

	if !reflect.DeepEqual(gotArticle, expected) {
		t.Errorf("handler returned unexpected output: got %v want %v", gotArticle, expected)
	}

}
