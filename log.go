package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ccding/go-logging/logging"
)

var (
	Log *Logger
)

type Logger struct {
	*logging.Logger
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
		return []string{}, fmt.Errorf("Ошибка получения директории с логами: %v", err)
	}
	logDirName = filepath.Join(logDirName, postfix)

	files, err := ioutil.ReadDir(logDirName)
	if err != nil {
		return []string{}, fmt.Errorf("Ошибка чтения директории с логами: %v", err)
	}

	fileNames := []string{}
	for _, file := range files {
		fileName := file.Name()

		matches, err := regexp.MatchString(logFileNameRegexp, fileName)
		if err != nil {
			return []string{}, fmt.Errorf("Ошибка проверки regexp имени файла: %v", err)
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
		return "", fmt.Errorf("Ошибка получения пути к файлу лога %s: %v", name, err)
	}

	b, err := ioutil.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("Ошибка чтения файла лога %s: %v", name, err)
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
		logging.WARN, "[%6s] [%s] %s():%d -> %s\n levelname,time,funcname,lineno,message", "02.01.2006 15:04:05", logPath, true)
	if err != nil {
		return nil, err
	}
	return &Logger{logger}, nil
}
