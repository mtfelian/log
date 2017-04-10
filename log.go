package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"sync"

	"github.com/ccding/go-logging/logging"
	"github.com/mtfelian/cli"
	"runtime/debug"
	"strings"
)

var (
	Log *Logger
)

type Logger struct {
	*logging.Logger

	logSync   sync.Mutex // мьютекс для показа стека
	showStack bool       // показать стек-трейс
}

const logFileExtension = ".log"
const logFileNameRegexp = `^\d{4}-\d{2}-\d{2}\.log$`
const logDirName = "logs"

// getDir возвращает путь к файлам логов
func getDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dirName := filepath.Join(wd, logDirName)
	if err := os.MkdirAll(dirName, 0700); err != nil {
		return "", err
	}

	return dirName, nil
}

// GetTodayLogName возвращает имя лога на сегодня
func GetTodayLogName() string {
	return time.Now().Format("2006-01-02") + logFileExtension
}

// getTodayPath возвращает путь к текущему файлу лога (на сегодня)
func getTodayPath() (string, error) {
	fileName := GetTodayLogName()
	return getPath(fileName)
}

// getPath возвращает путь к заданному файлу лога с именем logName
func getPath(logName string) (string, error) {
	dirName, err := getDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirName, logName), nil
}

// EnumerateLogFiles возвращает срез имён файлов лога
// postfix задаёт то что нужно добавить к директории с логами (для тестов)
func EnumerateLogFiles(postfix string) ([]string, error) {
	logDirName, err := getDir()
	if err != nil {
		return []string{}, fmt.Errorf("Error getting logdir: %v", err)
	}
	logDirName = filepath.Join(logDirName, postfix)

	files, err := ioutil.ReadDir(logDirName)
	if err != nil {
		return []string{}, fmt.Errorf("Error reading logdir: %v", err)
	}

	fileNames := []string{}
	for _, file := range files {
		fileName := file.Name()

		matches, err := regexp.MatchString(logFileNameRegexp, fileName)
		if err != nil {
			return []string{}, fmt.Errorf("Failed to check filename regexp: %v", err)
		}
		if !matches {
			// судя по имени, это не файл лога
			continue
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}

// GetLog возвращает содержимое файлов лога с именем name
func GetLogContent(name string) (string, error) {
	logPath, err := getPath(name)
	if err != nil {
		return "", fmt.Errorf("Error getting path to log file %s: %v", name, err)
	}

	b, err := ioutil.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("Error reading log file %s: %v", name, err)
	}
	return string(b), nil
}

// InitLog инициализирует лог-файл согласно текущей дате
func InitLog() (*Logger, error) {
	logPath, err := getTodayPath()
	if err != nil {
		return nil, err
	}

	logger, err := logging.FileLogger("log", // имя лога, нигде не используется пока
		logging.INFO, "[%6s] [%s] %s():%d -> %s\n levelname,time,funcname,lineno,message", "02.01.2006 15:04:05", logPath, true)
	if err != nil {
		return nil, err
	}
	return &Logger{logger, sync.Mutex{}, false}, nil
}

// NextWithStack() включает показ стек-трейса при следующем вызове одной из функций записи в лог logger.*f()
func (logger *Logger) NextWithStack() {
	logger.showStack = true
}

func (logger *Logger) stripFromStackTrace(depth int, stackTrace string) string {
	parts := strings.Split(stackTrace, "\t")
	truncatedStackTrace := ""
	for i, part := range parts {
		if i < depth {
			continue
		}
		truncatedStackTrace += fmt.Sprintf("%s\t", part)
	}
	return truncatedStackTrace
}

// logf() выводит в лог сообщение с уровнем level, заданной строкой format с параметрами спецификаторов v
func (logger *Logger) logf(level logging.Level, format string, v ...interface{}) {
	logger.logSync.Lock()
	defer logger.logSync.Unlock()

	s := cli.Sprintf(format, v...)
	if logger.showStack {
		const stacktraceText = "Stacktrace follows: "
		// 3 это кол-во вырезаемых записей трассировки
		truncatedStackTrace := logger.stripFromStackTrace(3, string(debug.Stack()))
		s += cli.Sprintf("\n{R%s{0\n{A%s{0\n", stacktraceText, truncatedStackTrace)
		logger.showStack = false
	}

	logger.Logger.Logf(level, s)
}

// Criticalf добавляет в лог запись с уровнем CRITICAL
func (logger *Logger) Criticalf(format string, v ...interface{}) {
	logger.logf(logging.CRITICAL, format, v...)
}

// Fatalf добавляет в лог запись с уровнем FATAL
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.logf(logging.FATAL, format, v...)
}

// Errorf добавляет в лог запись с уровнем ERROR
func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.logf(logging.ERROR, format, v...)
}

// Warnf добавляет в лог запись с уровнем WARN
func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.logf(logging.WARN, format, v...)
}

// Warningf добавляет в лог запись с уровнем WARNING
func (logger *Logger) Warningf(format string, v ...interface{}) {
	logger.Logger.Warningf(cli.Sprintf(format, v...))
}

// Infof добавляет в лог запись с уровнем INFO
func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.logf(logging.INFO, format, v...)
}

// Debugf добавляет в лог запись с уровнем DEBUG
func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.logf(logging.DEBUG, format, v...)
}

// Notsetf добавляет в лог запись с неустановленным уровнем
func (logger *Logger) Notsetf(format string, v ...interface{}) {
	logger.logf(logging.NOTSET, format, v...)
}

// LogPrefixedError записывает ошибку с заданным префиксом prefix и сообщением msg
func (logger *Logger) LogPrefixedError(prefix string, msg string) {
	logger.Errorf("[%s ERROR] %s", prefix, msg)
}

// LogPrefixedSuccess записывает успех с заданным префиксом prefix и сообщением msg
func (logger *Logger) LogPrefixedSuccess(prefix string, msg string) {
	logger.Infof("[%s SUCCESS] %s", prefix, msg)
}
