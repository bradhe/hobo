package parsing

import (
	"bufio"
	"strconv"
)

func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
func readLine(buf *bufio.Reader) (string, error) {
	var str string

	for {
		arr, prefix, err := buf.ReadLine()

		str += string(arr)

		if err != nil {
			return "", err
		}

		if !prefix {
			break
		}
	}

	return str, nil
}
