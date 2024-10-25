package authenticator

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/communicator"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/sessionRepository"
)

var sessions = make(map[string]model.SessionWithTime)
var sessionTime = 60.0 * 24.0 * 4.0

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ClearOldSessions() {
	sessions := sessionRepository.GetAllSessions()
	now := time.Now()
	clearedSessions := 0
	for _, session := range sessions {
		diff := now.Sub(session.SessDate)
		if diff.Minutes() > sessionTime {
			sessionRepository.DeleteSession(session.ID)
			clearedSessions++
		}
	}
	fmt.Println("Cleared old sessions: ", clearedSessions)
}

func IsValid(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	if !isPresent {
		sessionInDb, err := sessionRepository.FindSession(id)
		if err == nil {
			sessions[id] = sessionInDb
			sessionInMap = sessionInDb
		} else {
			return false
		}
	}
	now := time.Now()
	diff := now.Sub(sessionInMap.SessDate)
	if diff.Minutes() > sessionTime {
		delete(sessions, id)
		sessionRepository.DeleteSession(id)
		chart.DeleteSpecificHtml(session)
		return false
	}
	isValid := sessionInMap.Session == session
	if isValid {
		sessionInMap.SessDate = time.Now()
		sessions[id] = sessionInMap
		sessionRepository.AddSession(&sessionInMap)
	}
	return isValid
}

func Logout(id string, session string) {
	delete(sessions, id)
	chart.DeleteSpecificHtml(session)
	sessionRepository.DeleteSession(id)
}

func CheckIn(id string, session string) {
	sessionInMap, isPresent := sessions[id]
	if isPresent {
		sessionInMap.IsChecked = true
		sessions[id] = sessionInMap
		sessionRepository.AddSession(&sessionInMap)
	} else {
		sessionInDb, err := sessionRepository.FindSession(id)
		if err == nil {
			sessionInDb.IsChecked = true
			sessions[id] = sessionInDb
			sessionRepository.AddSession(&sessionInDb)
		}
	}
}

func GiveSession(id string) (string, error) {
	if len(id) == 0 {
		return "", errors.New("empty id")
	}
	sessionGenerated := randStringBytes(50)
	session := model.SessionWithTime{
		ID:        id,
		Session:   sessionGenerated,
		SessDate:  time.Now(),
		IsChecked: false,
	}
	sessions[id] = session
	sessionRepository.AddSession(&session)
	return sessionGenerated, nil
}

func Validate(activeConfiguration *model.Configuration, id string, session string, sendCommunicate bool) bool {
	isValidatedInTime := false
	if activeConfiguration.Environment == "local" {
		fmt.Println("Local environment, validation process skipped")
		isValidatedInTime = true
		CheckIn(id, session)
	} else {
		if sendCommunicate {
			linkToSend := activeConfiguration.RemoteAddress + "checkin/" + id + "/" + session
			communicator.SendMessageWithLink(id, linkToSend)
		}

		foundChecked := make(chan string)
		timer := time.NewTimer(90 * time.Second)
		go func() {
			for {
				time.Sleep(1 * time.Second)
				isCheckedAlready := isChecked(id, session)
				if isCheckedAlready {
					foundChecked <- "one"
				}

			}
		}()
		select {
		case <-foundChecked:
			fmt.Println("Id is validated")
			isValidatedInTime = true
		case <-timer.C:
			fmt.Println("Id is not validated in time")
		}
	}
	return isValidatedInTime
}

func isChecked(id string, session string) bool {
	sessionInMap, isPresent := sessions[id]
	if !isPresent {
		sessionInDb, err := sessionRepository.FindSession(id)
		if err == nil {
			sessions[id] = sessionInDb
			sessionInMap = sessionInDb
		} else {
			return false
		}

	}
	return sessionInMap.Session == session && sessionInMap.IsChecked
}

func randStringBytes(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng := rand.New(r)
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rng.Intn(len(letterBytes))]
	}
	return string(b)
}
