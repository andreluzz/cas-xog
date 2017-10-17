package common

import (
	"fmt"
	"strings"
	"github.com/mattn/go-colorable"
	"github.com/beevik/etree"
)

func Debug(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	r := strings.NewReplacer("[red[", "\033[91m", "[green[", "\033[92m", "[yellow[", "\033[93m", "[blue[", "\033[96m", "]]", "\033[0m")
	format = r.Replace(format)
	fmt.Fprintf(colorable.NewColorableStdout(), format)
}

func DebugElement(e *etree.Element) {
	doc := etree.NewDocument()
	doc.SetRoot(e.Copy())
	xml, _ := doc.WriteToString()
	fmt.Println(xml)
}