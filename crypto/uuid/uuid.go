package encoding

import (
	"github.com/satori/go.uuid"
)

//创建uuid
//Create uuid
func New() string {
	uuid := uuid.NewV4().String()
	return uuid
}
