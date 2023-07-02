package strutil

import "strconv"

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt64(str string) (int64, error) {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func StringToInt(str string) (int, error) {
	return strconv.Atoi(str)
}
