package mysql

import (
	"database/sql"
	"errors"
)

type Message struct {
	Mobile string `db:"mobile"`
	Code   string `db:"code"`
}

type SMService struct {
	Provider     string `db:"provider"`
	AppCode      string `db:"appcode"`
	Balance      string `db:"balance"`
	SignName     string `db:"signname"`
	TemplateCode string `db:"template"`
	AccessKeyId  string `db:"accesskeyid"`
	AccessSecret string `db:"accesssecret"`
}

const (
	mysqlCreateDatabase = iota
)

const (
	mysqlMessageCreateTable = iota
	mysqlMessageInsert
	mysqlMessageGetMobile
	mysqlMessageGetCode
	mysqlMessageDelete
)

const (
	mysqlSMServiceCreateTable = iota
	mysqlSMServiceALiYunInsert
	mysqlSMServiceTianYanInsert
	mysqlSMServiceDelete
	mysqlSMServiceUpDateBalance
	mysqlSMServiceGetProvider
	mysqlSMServiceGetAccessKeyId
	mysqlSMServiceGetAccessSecret
	mysqlSMServiceGetAppCode
	mysqlSMServiceGetTemplateCode
	mysqlSMServiceGetSignName
	mysqlSMServiceGetBalance
	mysqlSMServiceGetId
)

var (
	errInvalidInsert = errors.New("errInvalidInsert")
	errInvalidUpDate = errors.New("errInvalidUpDate")
	errInvalidDelete = errors.New("errInvalidDelete")
)

var SQLString = []string{
	`CREATE DATABASE IF NOT EXISTS detail`,
}

var (
	messageSQLString = []string{
		`CREATE TABLE IF NOT EXISTS detail.message(
		id INT NOT NULL AUTO_INCREMENT,
		mobile VARCHAR(32) NOT NULL,
		code VARCHAR(32) NOT NULL,
		PRIMARY KEY(id)
	   )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO detail.message (id, mobile, code) VALUES (?, ?, ?)`,
		`SELECT mobile FROM detail.message WHERE id = ?`,
		`SELECT code FROM detail.message WHERE id = ?`,
		`DELETE FROM detail.message WHERE id = ? LIMIT 1`,
	}

	smserviceSQLString = []string{
		`CREATE TABLE IF NOT EXISTS detail.smservice(
		   id INT NOT NULL AUTO_INCREMENT,
		   provider VARCHAR(32) NOT NULL,
		   accesskeyid VARCHAR(32),
		   accesssecret VARCHAR(32) ,
		   templatecode VARCHAR(32) NOT NULL,
		   signname VARCHAR(32),
		   appcode VARCHAR(32),
		   balance VARCHAR(32),
		   PRIMARY KEY(id) 
		   )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;`,
		`INSERT INTO detail.smservice (provider,accesskeyid, accesssecret, templatecode, signname, balance) VALUES (?, ?, ?, ?, ?, ?)`,
		`INSERT INTO detail.smservice (provider,appcode, templatecode, balance) VALUES (?, ?, ?, ?)`,
		`DELETE FROM detail.smservice WHERE id = ? LIMIT 1`,
		`UPDATE detail.smservice SET balance = ? WHERE id = ?`,
		`SELECT provider FROM detail.smservice WHERE id = ?`,
		`SELECT accesskeyid FROM detail.smservice WHERE id = ?`,
		`SELECT accesssecret FROM detail.smservice WHERE id = ?`,
		`SELECT appcode FROM detail.smservice WHERE id = ?`,
		`SELECT templatecode FROM detail.smservice WHERE id = ?`,
		`SELECT signname FROM detail.smservice WHERE id = ?`,
		`SELECT balance FROM detail.smservice WHERE id = ?`,
		`SELECT id FROM detail.smservice WHERE balance > 0`,
	}
)

func CreateDatabase(db *sql.DB) error {
	_, err := db.Exec(SQLString[mysqlCreateDatabase])
	if err != nil {
		return err
	}

	return nil
}

func CreateTableMessage(db *sql.DB) error {
	_, err := db.Exec(messageSQLString[mysqlMessageCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertMessage(db *sql.DB, id, mobile, code string) error {
	result, err := db.Exec(messageSQLString[mysqlMessageInsert], id, mobile, code)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}

func DeleteMessage(db *sql.DB, id string) error {
	result, err := db.Exec(messageSQLString[mysqlMessageDelete], id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidDelete
	}

	return nil
}

func GetMessageMobile(db *sql.DB, id string) (string, error) {
	var mobile string
	err := db.QueryRow(messageSQLString[mysqlMessageGetMobile], id).Scan(&mobile)
	if err != nil {
		return mobile, err
	}

	return mobile, nil
}

func GetMessageCode(db *sql.DB, id string) (string, error) {
	var code string
	err := db.QueryRow(messageSQLString[mysqlMessageGetCode], id).Scan(&code)
	if err != nil {
		return code, err
	}

	return code, nil
}

func GetMessage(db *sql.DB, id string) (*Message, error) {
	var msg Message
	msg.Mobile, _ = GetMessageMobile(db, id)
	msg.Code, _ = GetMessageCode(db, id)

	return &msg, nil
}

func CreateTableSMService(db *sql.DB) error {
	_, err := db.Exec(smserviceSQLString[mysqlSMServiceCreateTable])
	if err != nil {
		return err
	}

	return nil
}

func InsertALiYunSMService(db *sql.DB, provider, accesskeyid, accesssecret, templatecode, signName string, balance string) error {
	result, err := db.Exec(smserviceSQLString[mysqlSMServiceALiYunInsert], provider, accesskeyid, accesssecret, templatecode, signName, balance)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}

func InsertTianYanSMService(db *sql.DB, provider, appCode, templatecode string, balance string) error {
	result, err := db.Exec(smserviceSQLString[mysqlSMServiceTianYanInsert], provider, appCode, templatecode, balance)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidInsert
	}

	return nil
}

func DeleteSMService(db *sql.DB, id string) error {
	result, err := db.Exec(smserviceSQLString[mysqlSMServiceDelete], id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidDelete
	}

	return nil
}

func UpDateBalance(db *sql.DB, balance string, id string) error {
	result, err := db.Exec(smserviceSQLString[mysqlSMServiceUpDateBalance], balance, id)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return errInvalidUpDate
	}

	return nil
}

func GetSMServiceId(db *sql.DB) (string, error) {
	var id string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetId]).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func GetSMServiceAccessKeyId(db *sql.DB, id string) (string, error) {
	var accesskeyid string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetAccessKeyId], id).Scan(&accesskeyid)
	if err != nil {
		return accesskeyid, err
	}

	return accesskeyid, nil
}

func GetSMServiceAccessSecret(db *sql.DB, id string) (string, error) {
	var accessSecret string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetAccessSecret], id).Scan(&accessSecret)
	if err != nil {
		return accessSecret, err
	}

	return accessSecret, nil
}

func GetSMServiceTemplateCode(db *sql.DB, id string) (string, error) {
	var templatecode string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetTemplateCode], id).Scan(&templatecode)
	if err != nil {
		return templatecode, err
	}

	return templatecode, nil
}

func GetSMServiceSignName(db *sql.DB, id string) (string, error) {
	var signname string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetSignName], id).Scan(&signname)
	if err != nil {
		return signname, err
	}

	return signname, nil
}

func GetSMServiceAppCode(db *sql.DB, id string) (string, error) {
	var appcode string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetAppCode], id).Scan(&appcode)
	if err != nil {
		return appcode, err
	}

	return appcode, nil
}

func GetSMServiceProvider(db *sql.DB, id string) (string, error) {
	var provider string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetProvider], id).Scan(&provider)
	if err != nil {
		return provider, err
	}

	return provider, nil
}

func GetSMServiceBalance(db *sql.DB, id string) (string, error) {
	var balance string
	err := db.QueryRow(smserviceSQLString[mysqlSMServiceGetBalance], id).Scan(&balance)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func GetSMService(db *sql.DB, id string) (*SMService, error) {
	var smservice SMService

	smservice.AccessKeyId, _ = GetSMServiceAccessKeyId(db, id)
	smservice.AccessSecret, _ = GetSMServiceAccessSecret(db, id)
	smservice.AppCode, _ = GetSMServiceAppCode(db, id)
	smservice.TemplateCode, _ = GetSMServiceTemplateCode(db, id)
	smservice.SignName, _ = GetSMServiceSignName(db, id)
	smservice.Balance, _ = GetSMServiceBalance(db, id)
	smservice.Provider, _ = GetSMServiceProvider(db, id)

	return &smservice, nil
}
