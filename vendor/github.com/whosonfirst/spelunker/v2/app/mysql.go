//go:build mysql

package app

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/whosonfirst/spelunker/v2/sql"
)
