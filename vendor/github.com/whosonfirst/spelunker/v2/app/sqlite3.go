//go:build sqlite3

package app

import (
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/whosonfirst/spelunker/v2/sql"
)
