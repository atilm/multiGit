package utilities

import (
	"time"

	"github.com/gosuri/uilive"
)

type ConsolePrinter interface {
	PrintLines(lines []string)
}

type LiveConsolePrinter struct {
	writer *uilive.Writer
}

func NewLiveConsolePrinter() *LiveConsolePrinter {
	writer := uilive.New()
	writer.RefreshInterval = time.Hour
	return &LiveConsolePrinter{writer}
}

func (p *LiveConsolePrinter) PrintLines(lines []string) {
	if len(lines) == 0 {
		return
	}

	p.writer.Write([]byte(lines[0] + "\n"))

	for _, line := range lines[1:] {
		p.writer.Newline().Write([]byte(line + "\n"))
	}

	p.writer.Flush()
}
