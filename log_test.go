package log

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"
)

func TestEnumerateLogFiles(t *testing.T) {
	logDir, err := getDir()
	if err != nil {
		t.Fatal(err)
	}

	fp := filepath.Join(logDir, "2016-09-11.log")
	fmt.Println(fp)
	if err = ioutil.WriteFile(fp, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	logFileNames, err := EnumerateLogFiles("")
	if err != nil {
		t.Fatal(err)
	}

	if len(logFileNames) == 0 {
		t.Fatal("Должен быть хоть один файл лога")
	}

	for _, fileName := range logFileNames {
		match, err := regexp.MatchString(logFileNameRegexp, fileName)
		if err != nil {
			t.Fatalf("Что-то не то с regexp: %v", err)
		}
		if !match {
			t.Fatal("Неудовлетворительное имя в списке: %s", fileName)
		}
	}
}

func TestGetLogContentIfError(t *testing.T) {
	_, err := GetLogContent("filenotexist.log")
	if err == nil {
		t.Fatal("Ожидалась ошибка")
	}
}
