package ui

// Application states
type AppState int

const (
	StateWelcome AppState = iota
	StateMainMenu
	StateSearch
	StateResults
	StateTorrentDetails
)

// ASCII Art for PlayBuddy
const AsciiArt = `
██████╗ ██╗      █████╗ ██╗   ██╗██████╗ ██╗   ██╗██████╗ ██████╗ ██╗   ██╗
██╔══██╗██║     ██╔══██╗╚██╗ ██╔╝██╔══██╗██║   ██║██╔══██╗██╔══██╗╚██╗ ██╔╝
██████╔╝██║     ███████║ ╚████╔╝ ██████╔╝██║   ██║██║  ██║██║  ██║ ╚████╔╝ 
██╔═══╝ ██║     ██╔══██║  ╚██╔╝  ██╔══██╗██║   ██║██║  ██║██║  ██║  ╚██╔╝  
██║     ███████╗██║  ██║   ██║   ██████╔╝╚██████╔╝██████╔╝██████╔╝   ██║   
╚═╝     ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚═════╝  ╚═════╝ ╚═════╝ ╚═════╝    ╚═╝   
`

const Version = "V1.0.0"

// Navigation constants
const (
	KeyBack    = "left"
	KeyForward = "right"
	KeyUp      = "up"
	KeyDown    = "down"
	KeyEnter   = "enter"
	KeyEscape  = "esc"
	KeyQuit    = "q"
)
