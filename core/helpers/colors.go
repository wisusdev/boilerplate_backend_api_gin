package helpers

func ColorGreen(text string) string {
	return "\033[32m" + text + "\033[0m"
}

func ColorRed(text string) string {
	return "\033[31m" + text + "\033[0m"
}

func ColorYellow(text string) string {
	return "\033[33m" + text + "\033[0m"
}

func ColorBlue(text string) string {
	return "\033[34m" + text + "\033[0m"
}

func ColorMagenta(text string) string {
	return "\033[35m" + text + "\033[0m"
}

func ColorCyan(text string) string {
	return "\033[36m" + text + "\033[0m"
}

func ColorWhite(text string) string {
	return "\033[37m" + text + "\033[0m"
}
