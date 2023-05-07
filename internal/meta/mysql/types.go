package mysql

import "database/sql"

const (
	DatabaseIsolation = sql.LevelSerializable
	LogAlias          = "[mysql_client]"
)
