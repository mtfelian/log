package cli

import (
	"fmt"
	"strings"
)

// список возможных кодов цветов
const (
	Black = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// список возможных SGR-кодов
const (
	Reset = iota
	Bold
	Faint
	Italic
	Underline
)

const (
	// ClearScreen это escape-код для очистки экрана
	ClearScreen = "\x1b[2J"
	// Move это escape-код для перемещения курсора на x, y
	Move = "\x1b[%d;%dH"
	// LineReset это escape-код для возврата курсора в начало строки и очистки её
	LineReset = "\r\x1b[K"
)

// getColor возвращает escape-код для цвета с кодом code
func getColor(code int) string {
	return getParam(30 + code)
}

// getBgColor возвращает escape-код для фона цвета с кодом code
func getBgColor(code int) string {
	return getParam(40 + code)
}

// getParam возвращает escape-код для параметров текста
func getParam(code int) string {
	return fmt.Sprintf("\x1b[%dm", code)
}

// Println печатает строку str и добавляет в конец перенос строки
func Println(str string, a ...interface{}) {
	Printf(str+"\n", a...)
}

// Printf печатает строку str и подставляет туда параметры
func Printf(str string, a ...interface{}) {
	fmt.Printf(Colorize(str, a...))
}

// Sprintf возвращает строку str и подставляет туда параметры
func Sprintf(str string, a ...interface{}) string {
	return fmt.Sprintf(Colorize(str, a...))
}

// Colorize возвращает цветную строку преобразуя {-теги
// Пример: cli.Colorize("{Rred string{0 and {Bblue part{0")
func Colorize(str string, a ...interface{}) string {
	const (
		prefix  = "{"
		postfix = "|"
	)

	str = fmt.Sprintf(str, a...)
	/*
		Предыдущий формат: {G
		Новый формат: {G|
		Переход:
		замена
		  \{((?:(?:_)?(?:[wargybmcWARGYBMC]))|(?:[ius0]))
		на
		  {$1|
		в строковых литералах
	*/
	changeMap := map[string]string{
		prefix + "w" + postfix:  getColor(White),
		prefix + "a" + postfix:  getColor(Black),
		prefix + "r" + postfix:  getColor(Red),
		prefix + "g" + postfix:  getColor(Green),
		prefix + "y" + postfix:  getColor(Yellow),
		prefix + "b" + postfix:  getColor(Blue),
		prefix + "m" + postfix:  getColor(Magenta),
		prefix + "c" + postfix:  getColor(Cyan),
		prefix + "W" + postfix:  getParam(Bold) + getColor(White),
		prefix + "A" + postfix:  getParam(Bold) + getColor(Black),
		prefix + "R" + postfix:  getParam(Bold) + getColor(Red),
		prefix + "G" + postfix:  getParam(Bold) + getColor(Green),
		prefix + "Y" + postfix:  getParam(Bold) + getColor(Yellow),
		prefix + "B" + postfix:  getParam(Bold) + getColor(Blue),
		prefix + "M" + postfix:  getParam(Bold) + getColor(Magenta),
		prefix + "C" + postfix:  getParam(Bold) + getColor(Cyan),
		prefix + "_w" + postfix: getBgColor(White),
		prefix + "_a" + postfix: getBgColor(Black),
		prefix + "_r" + postfix: getBgColor(Red),
		prefix + "_g" + postfix: getBgColor(Green),
		prefix + "_y" + postfix: getBgColor(Yellow),
		prefix + "_b" + postfix: getBgColor(Blue),
		prefix + "_m" + postfix: getBgColor(Magenta),
		prefix + "_c" + postfix: getBgColor(Cyan),
		prefix + "_W" + postfix: getParam(Bold) + getBgColor(White),
		prefix + "_A" + postfix: getParam(Bold) + getBgColor(Black),
		prefix + "_R" + postfix: getParam(Bold) + getBgColor(Red),
		prefix + "_G" + postfix: getParam(Bold) + getBgColor(Green),
		prefix + "_Y" + postfix: getParam(Bold) + getBgColor(Yellow),
		prefix + "_B" + postfix: getParam(Bold) + getBgColor(Blue),
		prefix + "_M" + postfix: getParam(Bold) + getBgColor(Magenta),
		prefix + "_C" + postfix: getParam(Bold) + getBgColor(Cyan),
		prefix + "i" + postfix:  getParam(Italic),
		prefix + "u" + postfix:  getParam(Underline),
		prefix + "s" + postfix:  ClearScreen,
	}
	for key, value := range changeMap {
		str = strings.Replace(str, key, getParam(Reset)+value, -1)
	}
	str = strings.Replace(str, prefix+"0"+postfix, getParam(Reset), -1)
	return str
}
