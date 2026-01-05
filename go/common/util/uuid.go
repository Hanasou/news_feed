package util

import "github.com/google/uuid"

func NewUUID() string {
	myUuid := uuid.NewString()
	myUuid = FilterChar(myUuid, '-')
	return myUuid
}
