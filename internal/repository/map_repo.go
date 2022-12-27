package repository

import "errors"

type mapRepo struct {
	URLs map[string]string
}

func newMapRepo() *mapRepo {
	return &mapRepo{URLs: make(map[string]string)}
}

func (r *mapRepo) Add(url string) (string, error) {
	for {
		uuid := genUUID()
		if _, ok := r.URLs[uuid]; !ok {
			r.URLs[uuid] = url
			return uuid, nil
		}
	}
}

func (r *mapRepo) Get(uuid string) (string, error) {
	if uuid == "" {
		return "", errors.New("id of url isn't set")
	}

	url := r.URLs[uuid]

	if url == "" {
		return "", errors.New("url with that id is not found")
	}

	return url, nil
}
