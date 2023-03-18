package utilities

type ConsolePrinter interface {
	PrintLines(lines []string)
}

type LiveConsolePrinter struct {
}

func (p *LiveConsolePrinter) PrintLines(lines []string) {

}
