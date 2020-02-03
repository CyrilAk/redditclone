package items

import (
	"fmt"
	"time"
)

type Comment struct {
	Author  Author `json:"author"`
	Body    string `json:"body"`
	Created string `json:"created"`
	ID      string `json:"id"`
}

func (item *Item) NewComment(username string, userID string, message string) (string, error) {
	item.commentLastID++
	comment := &Comment{
		ID:      fmt.Sprint(item.commentLastID),
		Created: time.Now().Format("2006-01-02 15:04:05"),
		Body:    message,
	}

	comment.Author.Username = username
	comment.Author.ID = userID

	item.Comments = append(item.Comments, comment)
	return comment.ID, nil
}

func (item *Item) DeleteComment(comID string) (bool, error) {
	for i, com := range item.Comments {
		if com.ID == comID {
			if i < len(item.Comments)-1 {
				copy(item.Comments[i:], item.Comments[i+1:])
			}
			item.Comments[len(item.Comments)-1] = nil
			item.Comments = item.Comments[:len(item.Comments)-1]
			return true, nil
		}
	}
	return false, nil
}
