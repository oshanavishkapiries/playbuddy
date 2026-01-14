package ui

import (
	"fmt"
	"github.com/oshanavishkapiries/playbuddy/internal/configs"
)

const asciiArt = `
██████  ██       █████  ██    ██ ██████  ██    ██ ██████  ██████  ██    ██
██   ██ ██      ██   ██  ██  ██  ██   ██ ██    ██ ██   ██ ██   ██  ██  ██ 
██████  ██      ███████   ████   ██████  ██    ██ ██   ██ ██   ██   ████  
██      ██      ██   ██    ██    ██   ██ ██    ██ ██   ██ ██   ██    ██   
██      ███████ ██   ██    ██    ██████   ██████  ██████  ██████     ██  
`

func AsciiArt() string {
	return fmt.Sprintf("%sVersion: %s", asciiArt, configs.Version)
}
