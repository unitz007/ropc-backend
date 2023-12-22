package utils

import "fmt"

const (
	Blank   = ""
	UserKey = "user"

	WhereClientIdAndUserIdIs = "whereClientIdAndUserIdIs"
	WhereNameAndUserIdIs     = "whereNameAndUserIdIs"
	WhereUserIdIs            = "whereUserIdIs"
	WhereIdIs                = "whereIdIs"
	WhereUsernameOrEmailIs   = "whereUsernameOrEmailIs"
	WhereClientIdIs          = "whereClientIdIs"
)

var (
	Queries = map[string]func(args ...any) string{
		WhereClientIdAndUserIdIs: func(args ...any) string {
			return fmt.Sprintf("client_id = '%s' and user_id = %d", args[0], args[1])
		},

		WhereNameAndUserIdIs: func(args ...any) string {
			return fmt.Sprintf("name = '%s' and user_id = %d", args[0], args[1])
		},

		WhereUserIdIs: func(args ...any) string {
			return fmt.Sprintf("user_id = %d", args[0])
		},

		WhereIdIs: func(args ...any) string {
			return fmt.Sprintf("id = %d", args[0])
		},

		WhereUsernameOrEmailIs: func(args ...any) string {
			return fmt.Sprintf("username = '%s' OR email = '%s'", args[0], args[0])
		},

		WhereClientIdIs: func(args ...any) string {
			return fmt.Sprintf("client_id = '%s'", args[0])
		},
	}
)
