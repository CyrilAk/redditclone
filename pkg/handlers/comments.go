package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
)

type CommentJSON struct {
	Message string `json:"comment"`
}

func (h *ItemsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, found := vars["id"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	item, err := h.ItemsRepo.GetByItemID(itemID)
	if err != nil {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	message := new(CommentJSON)
	err = json.Unmarshal(data, message)
	if err != nil {
		http.Error(w, `Cant unmarshal JSON`, http.StatusBadRequest)
	}

	sess, _ := session.SessionFromContext(r.Context())
	lastID, err := h.ItemsRepo.AddComment(sess, item, message.Message)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(item)
	w.Write(resp)
	h.Logger.Infof("Insert with id LastInsertId: %v", lastID)
}

func (h *ItemsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, found := vars["postID"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}
	comID, found := vars["comID"]
	if !found {
		http.Error(w, `{"error": "bad id"}`, http.StatusBadGateway)
		return
	}

	item, err := h.ItemsRepo.GetByItemID(postID)
	if (err != nil) || (item == nil) {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	err = h.ItemsRepo.DeleteComment(comID, item)
	if err != nil {
		http.Error(w, `{"error": "db error"}`, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(item)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
