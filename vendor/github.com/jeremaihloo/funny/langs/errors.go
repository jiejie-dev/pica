package langs

import (
	"fmt"
)

// P panic
func P(keyword string, pos Position) string {
	return fmt.Sprintf("funny error [%s] at position %s", keyword, pos.String())
}
