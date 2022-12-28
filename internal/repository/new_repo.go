package repository

import "errors"

type URLObj struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type SomeRepo struct {
	URLs     map[string]string
	UserUrls map[string][]URLObj
}

func NewSomeRepo() *SomeRepo {
	return &SomeRepo{
		URLs:     make(map[string]string),
		UserUrls: make(map[string][]URLObj),
	}
}

func (r *SomeRepo) Add(url string) (string, error) {
	for {
		uuid := GenUUID()
		if _, ok := r.URLs[uuid]; !ok {
			r.URLs[uuid] = url
			return uuid, nil
		}
	}
}

func (r *SomeRepo) Get(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
}

func (r *SomeRepo) AddByUser(userId, originalUrl, shortUrl string) {

	r.UserUrls[userId] = append(r.UserUrls[userId], URLObj{OriginalURL: originalUrl, ShortURL: shortUrl})

}

func (r *SomeRepo) AllUsersURLS(userId string) []URLObj {
	return r.UserUrls[userId]
}

func (r *SomeRepo) IsUserExist(userId string) bool {
	_, ok := r.UserUrls[userId]
	return ok
}
