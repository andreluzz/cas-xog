package log

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andreluzz/cas-xog/util"
	"github.com/mattn/go-colorable"
)

var logger *log.Logger

//Info prints text on the screen and into log file
func Info(format string, args ...interface{}) {
	format = fmt.Sprintf(format, args...)
	clearLog("Info", format)
	r := strings.NewReplacer("[red[", "\033[91m", "[green[", "\033[92m", "[yellow[", "\033[93m", "[blue[", "\033[96m", "]]", "\033[0m")
	format = r.Replace(format)
	fmt.Fprintf(colorable.NewColorableStdout(), format)
}

//Debug only insert text into the log file
func Debug(msg string) {
	clearLog("Debug", msg)
}

func clearLog(mode, msg string) {
	r := strings.NewReplacer("[red[", "", "[green[", "", "[yellow[", "", "[blue[", "", "]]", "", "\n", "", "\r", "")
	logger.Println(mode + ": " + r.Replace(msg))
}

//InitLog initialize the io.Writer with the log file
func InitLog() {
	folder := "_logs"
	util.ValidateFolder(folder)
	file, err := os.OpenFile(folder+"/cas-xog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\n[cas-xog]Error: Failed to open log file\n")
	}
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.LstdFlags)
}
