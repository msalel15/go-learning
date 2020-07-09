package user

import (
	"github.com/couchbase/gocb/v2"
	"time"
)

type userService struct {
	db     *gocb.Cluster
	bucket *gocb.Bucket
}

type UserService interface {
	CreateUser(user UserDTO) error
	IsExistUsernameAndPassword(username, password string) (bool, error)
	IsExistByUsername(username string) (bool, error)
}

func NewService(cb *gocb.Cluster) UserService {
	bucket := cb.Bucket("user")

	err := bucket.WaitUntilReady(15*time.Second, nil)
	if err != nil {
		panic(err)
	}

	return &userService{db: cb, bucket: bucket}
}

func (u *userService) CreateUser(user UserDTO) error {

	collection := u.bucket.DefaultCollection()

	_, err := collection.Upsert(user.ID, user, &gocb.UpsertOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (u *userService) IsExistUsernameAndPassword(username, password string) (bool, error) {

	results, err := u.db.Query("SELECT * FROM `user` WHERE `username` = $1",
		&gocb.QueryOptions{
			PositionalParameters: []interface{}{username},
		},
	)

	if err != nil {
		return false, err
	}

	return results.Next(), nil
}

func (u *userService) IsExistByUsername(username string) (bool, error) {

	results, err := u.db.Query("SELECT * FROM `user` WHERE `username` = $1",
		&gocb.QueryOptions{
			PositionalParameters: []interface{}{username},
		},
	)

	if err != nil {
		return false, err
	}

	return results.Next(), nil
}
