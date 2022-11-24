package repository

type URLRepo interface {
	Add(url string) (urlID int)
	Get(urlID int) (url string, err error)
}

type Repository struct {
	URLRepo
}

func NewRepository(storage []string) *Repository {
	return &Repository{
		URLRepo: NewURLRepo(storage),
	}
}
