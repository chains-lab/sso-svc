package sqldb

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS

type txKeyType struct{}

var txKey = txKeyType{}
