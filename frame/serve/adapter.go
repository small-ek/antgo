package serve

import "gorm.io/gorm"

// WebFrameWork is an interface which is used as an adapter of
// framework and goAdmin. It must implement two methods. Use registers
// the routes and the corresponding handlers. Content writes the
// response to the corresponding context of framework.
type WebFrameWork interface {
	Name() string
	SetApp(app interface{}) error
	SetConnection()
}

// BaseAdapter is a base adapter contains some helper functions.
type BaseAdapter struct {
	Db *gorm.DB
}

// SetConnection set the db connection.
func (base *BaseAdapter) SetConnection() {

}
