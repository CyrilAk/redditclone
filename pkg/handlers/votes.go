package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"redditclone/pkg/session"

	"github.com/gorilla/mux"
)

func (h *ItemsHandler) Upvote(w http.ResponseWriter, r *http.Request) {
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

	sess, _ := session.SessionFromContext(r.Context())
	err = h.ItemsRepo.Upvote(sess, item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(item)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	h.Logger.Infof("Update for item with id : %v", id)
}

func (h *ItemsHandler) Downvote(w http.ResponseWriter, r *http.Request) {
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

	sess, _ := session.SessionFromContext(r.Context())
	fmt.Println("FLAG:DV(1);\tsess =", sess)
	err = h.ItemsRepo.Downvote(sess, item)
	fmt.Println("FLAG:DV(2)")
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(item)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	h.Logger.Infof("Update for item with id : %v", id)
}

func (h *ItemsHandler) Unvote(w http.ResponseWriter, r *http.Request) {
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

	sess, _ := session.SessionFromContext(r.Context())
	err = h.ItemsRepo.Unvote(sess, item)
	if err != nil {
		http.Error(w, `DB err`, http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(item)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	h.Logger.Infof("Update for item with id : %v", id)
}
