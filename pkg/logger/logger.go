package logger

import (
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var (
	Info    *log.Logger
	Warring *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

func init() {
	godotenv.Load()

	debug := os.Getenv("DEBUG")
	if debug == "" {
		debug = "false"
	}

	Info = log.New(os.Stdout, color.CyanString("[INFO] "), log.Ldate|log.Ltime|log.Lshortfile)
	Warring = log.New(os.Stdout, color.YellowString("[WARNING] "), log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stdout, color.RedString("[ERROR] "), log.Ldate|log.Ltime|log.Lshortfile)

	if debug == "true" || debug == "1" {
		Debug = log.New(os.Stdout, color.MagentaString("[DEBUG] "), log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Debug = log.New(io.Discard, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}
