package color

import (
	"fmt"
)

type Col string

//https://misc.flogisoft.com/bash/tip_colors_and_formatting
//http://jafrog.com/2013/11/23/colors-in-terminal.html
var (
	Reset       Col = "\x1B[0m"
	Red         Col = "\x1B[38;5;1m"
	Orange      Col = "\x1B[38;5;9m"
	DeepPink    Col = "\x1B[38;5;5m"
	Pink        Col = "\x1B[38;5;13m"
	Yellow      Col = "\x1B[38;5;3m"
	Green       Col = "\x1B[38;5;2m"
	SpringGreen Col = "\x1B[38;5;10m"
	Blue        Col = "\x1B[38;5;4m"
	DeepSkyBlue Col = "\x1B[38;5;6m"
	SkyBlue     Col = "\x1B[38;5;14m"
	Grey        Col = "\x1B[38;5;8m"
	Black       Col = "\x1B[38;5;232m"
	White       Col = "\x1B[38;5;15m"
)

func Addf(c Col, format string, a ...interface{}) string {
	return Add(c, fmt.Sprintf(format, a...))
}

func Add(c Col, str string) string {
	return string(c) + str + string(Reset)
}

func Check() {
	for i := 1; i < 256; i++ {
		fmt.Println(Add(Col(fmt.Sprintf("\x1B[38;5;%dm", i)), fmt.Sprintf("Number %d:", i)))
	}
}
