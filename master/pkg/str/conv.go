package str

import "strconv"

func StringToInt(buf string, dval int) int {
	res, err := strconv.ParseInt(buf, 10, 0)

	if err == nil {
		return int(res)
	}

	return dval
}
