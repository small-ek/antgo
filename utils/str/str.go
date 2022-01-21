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

//Trim Remove all spaces from the end of a string
func Trim(str string) string {
	spaceNum := 0
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == ' ' {
			spaceNum++
		} else {
			break
		}
	}
	return str[:len(str)-spaceNum]
}
