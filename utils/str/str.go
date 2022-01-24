package str

//ClearQuotes remove double quotes
func ClearQuotes(str string) string {
	if len(str) > 0 && (str[0] == '"') {
		str = str[1:]
	}
	if len(str) > 0 && str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}
	return str
}
