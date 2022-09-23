package fonts

import (
	_ "embed"
	"fmt"
	"os"

	"golang.org/x/image/font/opentype"
)

//go:embed Roboto-Regular.ttf
var robotoRegularBytes []byte

func RobotoRegular() *opentype.Font {
	tt, err := opentype.Parse(robotoRegularBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "roboto regular font: %s\n", err)
		os.Exit(1)
	}
	return tt
}
