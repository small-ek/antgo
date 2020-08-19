package encoding

import (
	"github.com/satori/go.uuid"
)

//创建uuid
//Create uuid
func NewUUID() string {
	uuid := uuid.NewV4().String()
	return uuid
}
