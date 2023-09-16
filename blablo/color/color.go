package color

var (
	Reset = "\033[0m"

	/////////////
	// Special //
	/////////////

	Bold      = "\033[1m"
	Underline = "\033[4m"

	/////////////////
	// Text colors //
	/////////////////

	Black       = "\033[30m"
	Red         = "\033[0;31m"
	Green       = "\033[0;32m"
	Yellow      = "\033[38;5;228m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	GreenCyan50 = "\033[38;5;50m"
	GreenCyan49 = "\033[38;5;49m"
	Gray        = "\033[37m"
	Gray254     = "\033[38;5;254m"
	Gray250     = "\033[38;5;250m"
	Gray247     = "\033[38;5;247m"
	White       = "\033[97m"

	///////////////////////
	// Background colors //
	///////////////////////

	BlackBackground   = "\033[40m"
	RedBackground     = "\033[41m"
	GreenBackground   = "\033[42m"
	YellowBackground  = "\033[43m"
	BlueBackground    = "\033[44m"
	MagentaBackground = "\033[45m"
	CyanBackground    = "\033[46m"
	GrayBackground    = "\033[47m"
	WhiteBackground   = "\033[107m"
)

func With(color, text string) string {
	return color + text + Reset
}
