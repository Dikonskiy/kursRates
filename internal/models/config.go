package models

var Config struct {
	ListenPort            string `json:"listenPort"`
	MysqlConnectionString string `json:"mysqlConnectionString"`
	APIURL                string `json:"apiURL"`
}
