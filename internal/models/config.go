package models

type Config struct {
	ListenPort            string `json:"listenPort"`
	MysqlConnectionString string `json:"mysqlConnectionString"`
	APIURL                string `json:"apiURL"`
	IsProd                bool   `json:"isProd"`
}
