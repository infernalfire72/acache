package config

import (
	"database/sql"
	"fmt"
)

type sqlConf struct {
	Host     string
	Username string
	Password string
	Database string
}

func (c sqlConf) String() string {
	return fmt.Sprintf("%s:%s@%s/%s", c.Username, c.Password, c.Host, c.Database)
}

var DB *sql.DB
