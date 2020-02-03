package handlers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/session"

	"redditclone/pkg/items"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ItemsHandler struct {
	Tmpl      *template.Template
	ItemsRepo *items.ItemsRepo
	Logger    *zap.SugaredLogger
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	elems, err := h.ItemsRepo.GetAll()
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(elems)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ItemsHandler) Add(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	item := items.NewItem()
	err := json.Unmarshal(data, item)
	if err != nil {
		http.Error(w, `Cant unmarshal JSON`, http.StatusBadRequest)
	}
	sess, _ := session.SessionFromContext(r.Context())
	lastID, err := h.ItemsRepo.Add(sess, item)

	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(item)
	w.Write(resp)
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
}

func (h *ItemsHandler) Read(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(id)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	if item == nil {
		http.Error(w, `no item`, http.StatusNotFound)
		return
	}
	item.Views++

	data, _ := json.Marshal(item)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	_, err := h.ItemsRepo.Delete(id)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]string{
		"message": "success",
	})
	w.Write(respJSON)
}

func (h *ItemsHandler) Category(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	catName, found := vars["catName"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	elems, err := h.ItemsRepo.GetByCategory(catName)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(elems)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h *ItemsHandler) UserItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, found := vars["username"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	elems, err := h.ItemsRepo.GetByUserID(username)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	data, _ := json.Marshal(elems)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
