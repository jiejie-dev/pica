package langs

import (
	"fmt"
)

func P(keyword string, pos Position) string {
	return fmt.Sprintf("funny error [%s] at position %s", keyword, pos.String())
}
