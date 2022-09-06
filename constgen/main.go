package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("input loader bin"))
	}
	var out []byte

	bin, e := os.ReadFile(os.Args[1])
	if e != nil {
		panic(e)
	}

	out = []byte("package main\n\n")

	out = append(out, []byte("var BofLdr = []byte{")...)
	for _, a := range bin {
		out = append(out, []byte(fmt.Sprintf("0x%02X,", a))...)
	}
	out = append(out, []byte("}\n\n")...)

	os.WriteFile("out\\const.go", out, 0644)

}
