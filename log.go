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
	logSync   sync.Mutex         // мьютекс для логирования
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
		logging.INFO, "[%8s] [%s] :: %s\n levelname,time,message",
		"02.01.2006 15:04:05", logPath, true)
	if err != nil {
		return nil, err
	}
	return &Logger{logger, sync.Mutex{}}, nil
}

// stripFromStackTrace вырезает записи на глубину depth из stackTrace
func (logger *Logger) stripFromStackTrace(depth int, stackTrace string) string {
	parts := strings.Split(stackTrace, "\t")
	truncatedStackTrace := ""
	for i, part := range parts {
		if i < depth {
			continue
		}
		truncatedStackTrace += fmt.Sprintf("%s\t", part)
	}

	parts = strings.Split(truncatedStackTrace, "\n")
	if len(parts) == 0 {
		return truncatedStackTrace
	}
	for i, part := range parts {
		parts[i] = strings.TrimRight(part, "\t")
	}

	return strings.Join(parts[1:], "\n")
}

// logf() выводит в лог сообщение с уровнем level, заданной строкой format с параметрами спецификаторов v
func (logger *Logger) logf(showStackTrace bool, level logging.Level, format string, v ...interface{}) {
	logger.logSync.Lock()
	defer logger.logSync.Unlock()

	s := cli.Sprintf(format, v...)
	if showStackTrace {
		const stackTraceTextBegin = "Stacktrace follows: "
		const stackTraceTextEnd = "End of stacktrace. "
		// 3 это кол-во вырезаемых записей трассировки
		truncatedStackTrace := logger.stripFromStackTrace(3, string(debug.Stack()))
		s += cli.Sprintf("\n{R%s\n{A%s{R%s{0\n", stackTraceTextBegin, truncatedStackTrace, stackTraceTextEnd)
	}

	logger.Logger.Logf(level, s)
}

// Criticalf добавляет в лог запись с уровнем CRITICAL
func (logger *Logger) Criticalf(format string, v ...interface{}) {
	logger.logf(false, logging.CRITICAL, format, v...)
}

// CriticalfStack добавляет в лог запись с уровнем CRITICAL и трассировкой вызовов
func (logger *Logger) CriticalfStack(format string, v ...interface{}) {
	logger.logf(true, logging.CRITICAL, format, v...)
}

// Fatalf добавляет в лог запись с уровнем FATAL
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	logger.logf(false, logging.FATAL, format, v...)
}

// FatalfStack добавляет в лог запись с уровнем FATAL и трассировкой вызовов
func (logger *Logger) FatalfStack(format string, v ...interface{}) {
	logger.logf(true, logging.FATAL, format, v...)
}

// Errorf добавляет в лог запись с уровнем ERROR
func (logger *Logger) Errorf(format string, v ...interface{}) {
	logger.logf(false, logging.ERROR, format, v...)
}

// ErrorfStack добавляет в лог запись с уровнем ERROR и трассировкой вызовов
func (logger *Logger) ErrorfStack(format string, v ...interface{}) {
	logger.logf(true, logging.ERROR, format, v...)
}

// Warnf добавляет в лог запись с уровнем WARN
func (logger *Logger) Warnf(format string, v ...interface{}) {
	logger.logf(false, logging.WARN, format, v...)
}

// WarnfStack добавляет в лог запись с уровнем WARN и трассировкой вызовов
func (logger *Logger) WarnfStack(format string, v ...interface{}) {
	logger.logf(true, logging.WARN, format, v...)
}

// Warningf добавляет в лог запись с уровнем WARNING
func (logger *Logger) Warningf(format string, v ...interface{}) {
	logger.logf(false, logging.WARNING, format, v...)
}

// WarningfStack добавляет в лог запись с уровнем WARNING и трассировкой вызовов
func (logger *Logger) WarningfStack(format string, v ...interface{}) {
	logger.logf(true, logging.WARNING, format, v...)
}

// Infof добавляет в лог запись с уровнем INFO
func (logger *Logger) Infof(format string, v ...interface{}) {
	logger.logf(false, logging.INFO, format, v...)
}

// InfofStack добавляет в лог запись с уровнем INFO и трассировкой вызовов
func (logger *Logger) InfofStack(format string, v ...interface{}) {
	logger.logf(true, logging.INFO, format, v...)
}

// Debugf добавляет в лог запись с уровнем DEBUG
func (logger *Logger) Debugf(format string, v ...interface{}) {
	logger.logf(false, logging.DEBUG, format, v...)
}

// DebugfStack добавляет в лог запись с уровнем DEBUG и трассировкой вызовов
func (logger *Logger) DebugfStack(format string, v ...interface{}) {
	logger.logf(true, logging.DEBUG, format, v...)
}

// Notsetf добавляет в лог запись с неустановленным уровнем
func (logger *Logger) Notsetf(format string, v ...interface{}) {
	logger.logf(false, logging.NOTSET, format, v...)
}

// NotsetfStack добавляет в лог запись с неустановленным уровнем и трассировкой вызовов
func (logger *Logger) NotsetfStack(format string, v ...interface{}) {
	logger.logf(true, logging.NOTSET, format, v...)
}

// LogPrefixedError записывает ошибку с заданным префиксом prefix и сообщением msg
func (logger *Logger) LogPrefixedError(prefix string, msg string) {
	logger.Errorf("[%s ERROR] %s", prefix, msg)
}

// LogPrefixedSuccess записывает успех с заданным префиксом prefix и сообщением msg
func (logger *Logger) LogPrefixedSuccess(prefix string, msg string) {
	logger.Infof("[%s SUCCESS] %s", prefix, msg)
}
