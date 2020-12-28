package parsing

import (
	"bufio"
	"io"
	"strings"

	"github.com/bradhe/hobo/pkg/content"
)

type AdminCode struct {
	Code      string
	Name      string
	ASCIIName string
	GeonameId int
}

func ParseAdminCode(row []string, ac *AdminCode) error {
	ac.Code = row[0]
	ac.Name = row[1]
	ac.ASCIIName = row[2]
	ac.GeonameId = ParseInt(row[3])
	return nil
}

func GetAdminCodes() map[string]AdminCode {
	codes := make(map[string]AdminCode)

	buf := bufio.NewReader(content.AssetReader("assets/admin1CodesASCII.txt"))

	for {
		li, err := readLine(buf)

		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		var adminCode AdminCode
		ParseAdminCode(strings.Split(li, "\t"), &adminCode)
		codes[adminCode.Code] = adminCode
	}

	return codes
}
