package main

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/kamva/csvtools"
)

func main() {
	content := `name,color,age,author.name,author.score
book a,white,1,ali,3
book a,,2,reza,4
book a,,10,,
book b,red,2,John,5
book b,,4,Jessy,6
book b,,6,,
`

	sc := bufio.NewScanner(bytes.NewReader([]byte(content)))
	sc.Split(new(csvtools.ScanGroupedCSVRecords).SplitFunc)
	for sc.Scan() {
		fmt.Println("--------------------")
		fmt.Println(string(sc.Bytes()))
	}
}
