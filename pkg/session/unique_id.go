package session

import (
	"time"

	"github.com/speps/go-hashids"
	"github.com/wexel-nath/wexel-auth/pkg/logger"
)

var (
	hashID *hashids.HashID
)

func Configure() {
	var err error
	hashID, err = hashids.New()
	if err != nil {
		logger.Error(err)
	}
}

func generateUniqueID(userID int64) (string, error) {
	n := []int64{
		userID,
		time.Now().Unix(),
	}
	return hashID.EncodeInt64(n)
}
