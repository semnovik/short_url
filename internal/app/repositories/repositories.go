package repositories

type URLRepo interface {
	Add(url string) (uuid string)
	Get(urlID string) (url string, err error)
}

type Repository struct {
	URLRepo
}

var urlStorage = make(map[string]string)

func NewRepository() *Repository {
	return &Repository{
		URLRepo: NewURLRepo(urlStorage),
	}
}
