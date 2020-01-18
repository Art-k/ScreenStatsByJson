package src

import (
	"github.com/allegro/bigcache"
	"github.com/jinzhu/gorm"
)

var AuthorizationHash string

var Db *gorm.DB

var Cache *bigcache.BigCache

const DbLogMode = false

var Err error
var Port string
var ScreenStatUrl string

const Version = "0.0.2"

var GetMaxTvStatInterval int64
