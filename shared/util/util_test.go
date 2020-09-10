package util

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
)

func TestRandomString(t *testing.T) {
	randomTypes := []string{
		"UPPERCASE",
		"LOWERCASE",
		"UPPERCASE_ALPHANUMERIC",
		"LOWERCASE_ALPHANUMERIC",
		"DEFAULT",
	}

	for _, randomType := range randomTypes {
		t.Run(fmt.Sprintf("Generate random string %s", randomType), func(t *testing.T) {
			resp := RandomString(3, "UPPERCASE")
			if resp == "" {
				t.Error("failed return generate random string")
			}
		})
	}

}

func TestRandomStringEngine(t *testing.T) {
	t.Run("test random string engine", func(t *testing.T) {
		var letter = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		resp := randomStringEngine(letter, 3)
		if resp == "" {
			t.Error("failed return random string engine")
		}
	})
}

func TestCustomPaginationQuery(t *testing.T) {
	t.Run("test custom pagination gorm query builder", func(t *testing.T) {

		trx := &gorm.DB{}
		_, err := CustomGormPaginationQuery(trx, 10, 1, "orderBy", "asc")
		if err != nil {
			t.Error("failed test custom pagination gorm query builder, err: ", err.Error())
		}

	})
}
