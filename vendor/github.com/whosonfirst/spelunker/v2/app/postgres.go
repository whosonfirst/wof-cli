//go:build postgres

package app

import (
	_ "github.com/lib/pq"
	_ "github.com/whosonfirst/spelunker/v2/sql"
)
