package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"redditclone/pkg/session"
)

type AuthURL struct {
	Url    string
	Method string
}

var authUrls = []AuthURL{
	AuthURL{
		Url:    "/api/posts/?$",
		Method: "POST",
	},
	AuthURL{
		Url:    `/api/post/[0-9a-zA-Z-]+/?$`,
		Method: "POST",
	},
	AuthURL{
		Url:    `/api/post/[0-9a-zA-Z-]+/?$`,
		Method: "DELETE",
	},
	AuthURL{
		Url:    `/api/post/[0-9a-zA-Z-]+/[0-9a-zA-Z-]+/?$`,
		Method: "DELETE",
	},
	AuthURL{
		Url:    `/api/post/[0-9a-zA-Z-]+/upvote/?$`,
		Method: "GET",
	},
	AuthURL{
		Url:    `/api/post/[0-9a-zA-Z-]+/downvote/?$`,
		Method: "GET",
	},
	AuthURL{
		Url:    `/api/post/[0-9a-zA-Z-]+/unvote/?$`,
		Method: "GET",
	},
}

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, item := range authUrls {
			matched, _ := regexp.MatchString(item.Url, r.URL.Path)
			if !matched {
				continue
			}
			if item.Method != r.Method {
				continue
			}
			sess, err := sm.Check(r)
			if err != nil {
				fmt.Println("no auth")
				http.Redirect(w, r, "/", 302)
				return
			}
			ctx := context.WithValue(r.Context(), session.SessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}
