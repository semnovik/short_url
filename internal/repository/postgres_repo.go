package repository

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"short_url/configs"
)

type PostgresRepo struct {
	Conn *sql.DB
}

func NewPostgresRepo() (*PostgresRepo, error) {
	var err error
	var db *sql.DB

	if configs.Config.DatabaseDSN != "" {
		db, err = sql.Open("pgx", configs.Config.DatabaseDSN)
	}

	err = goose.Up(db, configs.Config.MigrationsDir, goose.WithAllowMissing())

	return &PostgresRepo{Conn: db}, err
}

//func (r *PostgresRepo) Add(url string) (string, error) {
//	for {
//		uuid := GenUUID()
//		if _, ok := r.mapRepo.URLs[uuid]; !ok {
//
//			// Если есть интеграция с файлом, то пишем еще и в файл
//			if r.File != nil {
//				err := writeURLInFile(r.File, uuid, url)
//				if err != nil {
//					return "", err
//				}
//			}
//
//			r.mapRepo.URLs[uuid] = url
//			return uuid, nil
//		}
//	}
//}
//
//func (r *PostgresRepo) Get(uuid string) (string, error) {
//	return r.mapRepo.Get(uuid)
//}
//
//func (r *PostgresRepo) AddByUser(userID, originalURL, shortURL string) {
//	r.mapRepo.UserUrls[userID] = append(r.mapRepo.UserUrls[userID], URLObj{OriginalURL: originalURL, ShortURL: shortURL})
//}
//
//func (r *PostgresRepo) AllUsersURLS(userID string) []URLObj {
//	return r.mapRepo.UserUrls[userID]
//}
//
//func (r *PostgresRepo) IsUserExist(userID string) bool {
//	_, ok := r.mapRepo.UserUrls[userID]
//	return ok
//}
//
//func (r *PostgresRepo) Ping() error {
//	ctx := context.Background()
//	if r.PostgresDB == nil {
//		return errors.New("something wrong with DB-connection")
//	}
//	return r.PostgresDB.Ping(ctx)
//}
