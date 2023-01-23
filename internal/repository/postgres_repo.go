package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"short_url/configs"
)

type PostgresRepo struct {
	Conn *pgx.Conn
}

func NewPostgresRepo() *PostgresRepo {
	ctx := context.Background()
	var dbConn *pgx.Conn
	var err error

	if configs.Config.DatabaseDSN != "" {

		dbConn, err = pgx.Connect(ctx, configs.Config.DatabaseDSN)

		switch {
		case err != nil:
			log.Print("DB not configured")
		default:

			err = dbConn.Ping(ctx)
			if err != nil {
				panic(err)
			}
			log.Print("DB successfully configured")
		}
	} else {
		log.Print("DB not configured")
	}

	return &PostgresRepo{Conn: dbConn}
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
