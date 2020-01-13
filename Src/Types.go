package src

import "github.com/jinzhu/gorm"

var Db *gorm.DB

const DbLogMode = true

var Err error
var Port string

//var ScreenStatUrl string

const Version = "0.0.2"

var GetMaxTvStatInterval int64
