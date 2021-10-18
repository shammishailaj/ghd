package utils

import (
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

// InsertEmailPool - insert a row into email_pool
func (u *Utils) InsertEmailPool(email, subject, content string, dbProd *RDBMS) {
	t := time.Now()

	alertEmail := make(map[string]string)
	alertEmail["email"] = email
	alertEmail["subject"] = subject
	alertEmail["content"] = content
	alertEmail["sendtime"] = t.Format("2006-01-02 15:04:05")
	alertEmail["addtime"] = t.Format("2006-01-02 15:04:05")
	alertEmail["status"] = "1"
	alertEmail["isactive"] = "1"
	insertEmail := dbProd.InsertRow("email_pool", alertEmail)

	log.Println("InsertEmailPool = ", insertEmail)
}

func (u *Utils) SendEmailAlert(emails []string, subject, emailContent string) {
	var dbM2 RDBMS

	emailPoolTable := os.Getenv("EMAIL_POOL_TABLE_NAME")

	err := dbM2.ConnectM2()

	if err != nil {
		log.Errorf("Error connecting to database with provided dsn %s", dbM2.Conf.DSN())
		log.Errorf("Error details %s", err.Error())
	} else {
		log.Infof("Alert DB connected Successfully...")
		log.Infof("AlertDBTable = %s:%d/%s", dbM2.Conf.HostName, dbM2.Conf.Port, emailPoolTable)
		log.Infof("Writing error email to queue")
		for _, email := range emails {
			u.InsertEmailPool(email, subject, emailContent, &dbM2)
		}
		dbM2CloseErr := dbM2.Close() // close connection
		if dbM2CloseErr != nil {
			log.Errorf("Error closing DB connection. %s", dbM2CloseErr.Error())
		}
	}
}
