package log

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"testing"
)

// TestEnumerateLogFiles checks log files enumeration
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
		t.Fatal("At least one log file should exist")
	}

	for _, fileName := range logFileNames {
		match, err := regexp.MatchString(logFileNameRegexp, fileName)
		if err != nil {
			t.Fatalf("Something wrong with regexp: %v", err)
		}
		if !match {
			t.Fatalf("Invalid name in list: %s", fileName)
		}
	}
}

// TestGetLogContentIfError checks get log content if error
func TestGetLogContentIfError(t *testing.T) {
	_, err := GetLogContent("filenotexist.log")
	if err == nil {
		t.Fatal("Expected error")
	}
}

// TestLogStack is a convenience test, it tests almost nothing
// launch it to see, does it work or not.
// Command line example:
//   rm /home/felian/go_code/src/github.com/mtfelian/log/logs/2017-04-10.log;
//   go test --run TestLogStackTrace;
//   cat /home/felian/go_code/src/github.com/mtfelian/log/logs/2017-04-10.log
func TestLogStackTrace(t *testing.T) {
	logger, err := InitLog()
	if err != nil {
		t.Fatal(err)
	}

	logger.InfofStack("{YTest message 1 (with stack): {G%v{0", "OK")
	logger.Errorf("{YTest message 2: {G%v{0", "OK")
	logger.Infof("{YTest message 3: {G%v{0", "OK")
	logger.FatalfStack("{YTest message 4: {G%v{0", "OK")
	logger.Infof("{YTest message 5: {G%v{0", "OK")
}
