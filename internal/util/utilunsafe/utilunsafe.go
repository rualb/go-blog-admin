package utilunsafe

import (
	"errors"
	"go-blog-admin/internal/config/consts"
	"go-blog-admin/internal/util/utilstring"
)

// unsafeName check if var name chars are valid [a-z0-9_-]
func UnsafeName(name string) (err error) {

	if len(name) > consts.DefaultTextLength {
		return errors.New("name len > 100")
	}

	// Check each character in the string
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_' || char == '-') {
			return errors.New("name must contain only safe chars")
		}
	}

	return nil
}

// UnsafeLen
func UnsafeLen(s string) string {
	return utilstring.Left(s, consts.DefaultTextLength)
}
