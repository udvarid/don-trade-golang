package userRepository

import (
	"encoding/json"
	"errors"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	bolt "go.etcd.io/bbolt"
)

func GetAllUsers() []model.User {
	db := repoUtil.OpenDb()
	var result []model.User
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))

		b.ForEach(func(k, v []byte) error {
			var user model.User
			json.Unmarshal([]byte(v), &user)
			result = append(result, user)
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func DeleteUser(id string) {
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		err := b.Delete([]byte(id))
		return err
	})
	defer db.Close()
}

func FindUser(id string) (model.User, error) {
	db := repoUtil.OpenDb()
	var result model.User
	foundUser := false
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		b.ForEach(func(k, v []byte) error {
			userId := string(k[:])
			if userId == id {
				var user model.User
				json.Unmarshal([]byte(v), &user)
				result = user
				foundUser = true
			}
			return nil
		})
		return nil
	})
	defer db.Close()

	if foundUser {
		return result, nil
	} else {
		return result, errors.New("not found in db")
	}
}

func AddUser(user model.User) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.ID), buf)
	})

	defer db.Close()
}

func UpdateUser(user model.User) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("User"))
		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return b.Put([]byte(user.ID), buf)
	})
	defer db.Close()
}
