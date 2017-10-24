package common

import (
	"os"
	"fmt"
	"log"
	"strings"
	"github.com/beevik/etree"
	"github.com/mattn/go-colorable"
)

var Log *log.Logger

func Info(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	clearLog(format)
	r := strings.NewReplacer("[red[", "\033[91m", "[green[", "\033[92m", "[yellow[", "\033[93m", "[blue[", "\033[96m", "]]", "\033[0m")
	format = r.Replace(format)
	fmt.Fprintf(colorable.NewColorableStdout(), format)
}

func clearLog(msg string) {
	r := strings.NewReplacer("[red[", "", "[green[", "", "[yellow[", "", "[blue[", "", "]]", "", "\n", "", "\r", "")
	Log.Println(r.Replace(msg))
}

func InitLog(){
	file, err := os.OpenFile("cas-xog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\n[cas-xog][red[Error]] Failed to open log file\n")
	}
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LstdFlags)
}

func DebugElement(e *etree.Element) {
	doc := etree.NewDocument()
	doc.SetRoot(e.Copy())
	xml, _ := doc.WriteToString()
	fmt.Println(xml)
}