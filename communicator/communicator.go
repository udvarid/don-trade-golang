package communicator

import (
	"log"
	"net/smtp"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/udvarid/don-trade-golang/model"
)

var (
	activeConfiguration = &model.Configuration{}
)

func Init(config *model.Configuration) {
	activeConfiguration = config
}

var verifier = emailverifier.NewVerifier()

func SendMessage(toAddress string, message string) {
	ret, err := verifier.Verify(toAddress)
	if err == nil && ret.Syntax.Valid {
		msg := []byte("To: " + toAddress + "\r\n" +
			"Subject: Wellcome!!\r\n" +
			"\r\n" +
			message)
		sendMail(toAddress, msg)
	}
}

func sendMail(toAddress string, message []byte) {
	address, pw := getMailAndPw()
	auth := smtp.PlainAuth("", address, pw, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, address, []string{toAddress}, message)
	if err != nil {
		log.Print(err)
	}
}

func getMailAndPw() (string, string) {
	var address string
	var pw string
	if activeConfiguration.Environment == "local" {
		address = activeConfiguration.Mail_local_from
		pw = activeConfiguration.Mail_local_psw
	} else {
		address = activeConfiguration.Mail_from
		pw = activeConfiguration.Mail_psw
	}
	return address, pw
}
