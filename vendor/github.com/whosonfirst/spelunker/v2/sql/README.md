# sql

The `sql` package implements the `Spelunker` interface for Who's On First data indexed in SQL databases with [database/sql](https://pkg.go.dev/database/sql) drivers.

Currently supported drivers are: `sqlite3`, `mysql` and `postgres` though in practice only the SQLite implementation has been thoroughly tested.

| Target | Driver | Build tags | Provider | Notes |
| --- | --- | --- | --- | --- |
| MySQL | `mysql` | `mysql` | [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) | Support for MySQL should probably still be considered "alpha" at best. |
| Postgres | `postgres` | `postgres` | [lib/pq](https://github.com/lib/pq) | Support for Postgres should probably still be considered "alpha" at best. |
| SQLite | `sqlite3` | `sqlite3,icu,json1,fts5` | [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) | |

## Example

New `database/sql`-backed Spelunker instances are created by passing a URI to the `NewSpelunker` method in the form of:

```
sql://{DATABASE_ENGINE}?dsn={DATABASE_ENGINE_DSN}
```

Where `{DATABASE_ENGINE}` is a registered (imported) `database/sql.Driver` name and `{DATABASE_ENGINE_DSN}` is that driver's specific DSN string for connecting to the database.

For example:

```
import (
       "context"

       "github.com/whosonfirst/spelunker"
       _ "github.com/whosonfirst/spelunker/sql"       
)

sp, _ := spelunker.NewSpelunker(context.Background(), "sql://sqlite3?dsn=example.db")
```

_Note how the code does NOT import any specific `database/sql` implementation. That is expected to be handled by build tags (described above)._

## Things the `database/sql` Spelunker implementation does NOT do yet

* The `database/sql` Spelunker does not implement any of the tag-related methods (`GetTags`, `HasTag`, `HasTagFaceted`) yet.

## Database schema(s)

Database table schemas used by the `SQLSpelunker` implementation are defined in the [whosonfirst/go-whosonfirst-database/sql/tables](https://github.com/whosonfirst/go-whosonfirst-database/tree/main/sql/tables) package.