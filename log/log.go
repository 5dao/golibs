package log

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

//appDir/logs/appDir_20060102.log

var filePrefix string
var logFile *os.File

func init() {
	var err error
	filePrefix, err = getPrefix()
	if err != nil {
		panic(err)
	}
	//init
	makeDateFile()

	logCron := cron.New()
	logCron.AddFunc("0 0 * * *", func() {
		mkErr := makeDateFile()
		if mkErr != nil {
			std.Errorf("makeDateFile err: $v", mkErr)
		}
	})
	logCron.Start()
}

// makeDateFile appDir/logs/appNam_20060102.log
func makeDateFile() (err error) {
	defer func() {
		if rev := recover(); rev != nil {
			err = fmt.Errorf("makeDateFile recover,rev: %v", rev)
		}
	}()

	now := time.Now()

	var fileName string
	fileName = filePrefix + "_" + now.Format("20060102") + ".log"
	oldMask := syscall.Umask(0)
	newLogFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0640)
	syscall.Umask(oldMask)
	if err != nil {
		return fmt.Errorf("os.OpenFile err: %v", err)
	}
	std.SetOutput(newLogFile)

	//close old logfile
	if logFile != nil {
		err = logFile.Close()
		if err != nil {
			return fmt.Errorf("oldLogFile close err: %v", err)
		}
	}
	logFile = newLogFile

	return
}

// appDir/logs/appName
func getPrefix() (string, error) {
	instanceAbsPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", fmt.Errorf("filepath.Abs err: %v, os.Args[0]: %v", err, os.Args[0])
	}
	instanceName := filepath.Base(instanceAbsPath)

	dir := filepath.Dir(instanceAbsPath)

	oldMask := syscall.Umask(0)
	os.Mkdir(filepath.Join(dir, "logs"), 0750)
	syscall.Umask(oldMask)

	return filepath.Join(dir, "logs", instanceName), nil
}
