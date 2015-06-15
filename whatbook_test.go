package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

func TestCanary(t *testing.T) {
}

func TestIndex(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	indexHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home Page did not return %v", http.StatusOK)
	}
}
