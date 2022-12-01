package services

type ShorterService interface {
	PostURL(url string) string
	GetURLByID(urlID string) (string, error)
}

type Service struct {
	ShorterService
}

func NewService(shorter *Shorter) *Service {
	return &Service{
		ShorterService: shorter,
	}
}
