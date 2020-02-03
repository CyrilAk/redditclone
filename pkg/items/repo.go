package items

import (
	"fmt"
	"redditclone/pkg/session"
	"sync"
	"time"
)

type ItemsRepo struct {
	mu     *sync.RWMutex
	lastID uint32
	data   []*Item
}

func NewRepo() *ItemsRepo {
	return &ItemsRepo{
		mu:     &sync.RWMutex{},
		lastID: 0,
		data:   make([]*Item, 0, 10),
	}
}

func (repo *ItemsRepo) GetAll() ([]*Item, error) {
	return repo.data, nil
}

func (repo *ItemsRepo) GetByCategory(categoryName string) ([]*Item, error) {
	reqItems := make([]*Item, 0)
	for _, item := range repo.data {
		repo.mu.RLock()
		if item.Category == categoryName {
			reqItems = append(reqItems, item)
		}
		repo.mu.RUnlock()
	}
	return reqItems, nil
}

func (repo *ItemsRepo) GetByItemID(id string) (*Item, error) {
	for _, item := range repo.data {
		repo.mu.RLock()
		if item.ID == id {
			repo.mu.RUnlock()
			return item, nil
		}
		repo.mu.RUnlock()
	}
	return nil, nil
}

func (repo *ItemsRepo) GetByUserID(username string) ([]*Item, error) {
	reqItems := make([]*Item, 0)
	for _, item := range repo.data {
		repo.mu.RLock()
		if item.Author.Username == username {
			reqItems = append(reqItems, item)
		}
		repo.mu.RUnlock()
	}
	return reqItems, nil
}

func (repo *ItemsRepo) Add(sess *session.Session, item *Item) (uint32, error) {
	item.Author.Username = sess.Username
	item.Author.ID = sess.UserID
	item.Created = time.Now().Format("2006-01-02 15:04:05")

	defaultVote := +1
	item.NewVote(defaultVote, sess.UserID)

	repo.mu.Lock()
	repo.lastID++
	item.ID = fmt.Sprint(repo.lastID)
	repo.data = append(repo.data, item)
	repo.mu.Unlock()

	return repo.lastID, nil
}

func (repo *ItemsRepo) Delete(id string) (bool, error) {
	i := -1
	repo.mu.Lock()
	for idx, item := range repo.data {
		if item.ID == id {
			i = idx
			break
		}
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]
	repo.mu.Unlock()

	return true, nil
}

func (repo *ItemsRepo) Unvote(sess *session.Session, item *Item) error {
	item.muVt.Lock()
	_, err := item.DeleteVote(sess.UserID)
	item.muVt.Unlock()
	if err != nil {
		return fmt.Errorf("DB err")
	}
	return nil
}

func (repo *ItemsRepo) Upvote(sess *session.Session, item *Item) error {
	item.muVt.Lock()
	_, err := item.DeleteVote(sess.UserID)
	item.muVt.Unlock()
	if err != nil {
		return fmt.Errorf("DB err")
	}

	item.muVt.Lock()
	item.NewVote(1, sess.UserID)
	item.muVt.Unlock()
	return nil
}

func (repo *ItemsRepo) Downvote(sess *session.Session, item *Item) error {
	item.muVt.Lock()
	_, err := item.DeleteVote(sess.UserID)
	item.muVt.Unlock()
	if err != nil {
		return fmt.Errorf("DB err")
	}

	item.muVt.Lock()
	item.NewVote(-1, sess.UserID)
	item.muVt.Unlock()
	return nil
}

func (repo *ItemsRepo) AddComment(sess *session.Session, item *Item, message string) (string, error) {
	item.muCm.Lock()
	id, err := item.NewComment(sess.Username, sess.UserID, message)
	item.muCm.Unlock()
	return id, err
}

func (repo *ItemsRepo) DeleteComment(comID string, item *Item) error {
	item.muCm.Lock()
	defer item.muCm.Unlock()
	_, err := item.DeleteComment(comID)
	return err
}
