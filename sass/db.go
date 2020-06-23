package cargo_db

import (
	"database/sql"

	_ "github.com/lib/pq"
	cargo_db "github.com/major2015/new-cargo/models"
)

type DB = *sql.Tx

// as a convience re-export models/func from models
// so that other users of Sass don't have to import models
type Tenant = cargo_db.Tenant

var Tenants = cargo_db.Tenants
var FindTenantP = cargo_db.FindTenantP

type User = cargo_db.User

var Users = cargo_db.Users
var FindUserP = cargo_db.FindUserP

type Subscription = cargo_db.Subscription

var Subscriptions = cargo_db.Subscriptions

func ConnectDB(c Configuration) *sql.DB {

}
