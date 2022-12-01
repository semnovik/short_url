package repositories

type URLRepo interface {
	Add(url string) (urlID int)
	Get(urlID int) (url string, err error)
}

type Repository struct {
	URLRepo
}

var urlStorage []string

func NewRepository() *Repository {
	return &Repository{
		URLRepo: NewURLRepo(urlStorage),
	}
}
