package sessionRepository

import (
	"encoding/json"
	"errors"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	bolt "go.etcd.io/bbolt"
)

func GetAllSessions() []model.SessionWithTime {
	db := repoUtil.OpenDb()
	var result []model.SessionWithTime
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))

		b.ForEach(func(k, v []byte) error {
			var session model.SessionWithTime
			json.Unmarshal([]byte(v), &session)
			result = append(result, session)
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func DeleteSession(id string) {
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		err := b.Delete([]byte(id))
		return err
	})
	defer db.Close()
}

func FindSession(id string) (model.SessionWithTime, error) {
	db := repoUtil.OpenDb()
	var result model.SessionWithTime
	foundSession := false
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		b.ForEach(func(k, v []byte) error {
			sessionId := string(k[:])
			if sessionId == id {
				var session model.SessionWithTime
				json.Unmarshal([]byte(v), &session)
				result = session
				foundSession = true
			}
			return nil
		})
		return nil
	})
	defer db.Close()

	if foundSession {
		return result, nil
	} else {
		return result, errors.New("not found in db")
	}
}

func AddSession(id string, session model.SessionWithTime) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sessions"))
		buf, err := json.Marshal(session)
		if err != nil {
			return err
		}
		return b.Put([]byte(id), buf)
	})

	defer db.Close()
}
