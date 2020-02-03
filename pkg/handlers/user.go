package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"go.uber.org/zap"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo *user.UserRepo
	Sessions *session.SessionsManager
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, `Template errror`, http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	username, password, err := getUserParams(r)
	if err != nil {
		http.Error(w, `Incorrect JSON`, http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.Authorize(username, password)
	if err == user.ErrBadPass {
		w.WriteHeader(http.StatusUnauthorized)
		resp, _ := json.Marshal(map[string]interface{}{"message": "invalid password"})
		w.Write(resp)
		return
	} else if err == user.ErrNoUser {
		w.WriteHeader(http.StatusUnauthorized)
		resp, _ := json.Marshal(map[string]interface{}{"message": "user not found"})
		w.Write(resp)
		return
	}

	secret := h.Sessions.GetSessSecret("")
	sess, err := h.Sessions.Create(string(u.ID), u.Username, secret)
	if err != nil {
		http.Error(w, `Cant create session`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(map[string]interface{}{"token": sess.Token})
	w.Write(resp)
	h.Logger.Infof("created session for %v", sess.UserID)
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	username, password, err := getUserParams(r)
	if err != nil {
		http.Error(w, `Incorrect JSON`, http.StatusBadRequest)
		return
	}

	u, err := h.UserRepo.Registration(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp, _ := json.Marshal(map[string]interface{}{
			"errors": []map[string]interface{}{{
				"location": "body",
				"param":    "username",
				"value":    username,
				"msg":      "already exist",
			}}})
		w.Write(resp)
		return
	}

	secret := h.Sessions.GetSessSecret("")
	sess, err := h.Sessions.Create(fmt.Sprint(u.ID), u.Username, secret)
	if err != nil {
		http.Error(w, `Cant create session`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(map[string]interface{}{"token": sess.Token})
	w.Write(resp)
	h.Logger.Infof("created session for %x", sess.UserID)
}

type UserJSON struct {
	Username string
	Password string
}

func getUserParams(r *http.Request) (string, string, error) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	userJSON := &UserJSON{}
	err = json.Unmarshal(body, userJSON)
	if err != nil {
		return "", "", fmt.Errorf("Incorrect JSON")
	}
	return userJSON.Username, userJSON.Password, nil
}
