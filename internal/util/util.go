package util

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func NewID() string {
	return strings.Join(strings.Split(uuid.New().String(), "-"), "")
}

func UserIdFromContext(ctx context.Context) (string, bool) {
	session, exist := ctx.Value("sessionKey").(*Session)

	if session == nil {
		return "", false
	}

	return session.UserId, exist
}

func Md5String(str string) string {
	hashByte := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hashByte)
}
