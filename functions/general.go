package functions

func Includes(strs []string, args ...string) bool {
	for _, str := range strs {
		for _, arg := range args {
			if str == arg {
				return true
			}
		}
	}
	return false
}
