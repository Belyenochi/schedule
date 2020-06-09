package store

import (
	"django-go/pkg/types"
	"io/ioutil"
	"path/filepath"
	"django-go/pkg/constants"
	"encoding/json"
	"os"
	"fmt"
	"runtime"
	"strings"
)

func StoreSchedule(results []types.ScheduleResult, directory string) {

	dir, err := mkdirAndReturn(directory)

	if err != nil {
		fmt.Println("open file err, msg:" + err.Error())
		return
	}

	path := filepath.Join(dir, constants.SCHEDULE_RESULT)

	fmt.Println(fmt.Sprintf("schedule file : %v ,result count : %v", path, len(results)))

	bytes, err := json.Marshal(results)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)

	if err != nil {
		fmt.Println("write file err, msg:" + err.Error())
	}

}

func StoreReschedule(results []types.RescheduleResult, directory string) {

	dir, err := mkdirAndReturn(directory)

	if err != nil {
		fmt.Println("open file err, msg:" + err.Error())
		return
	}

	path := filepath.Join(dir, constants.RESCHEDULE_RESULT)

	fmt.Println(fmt.Sprintf("reschedule file : %v ,result count : %v", path, len(results)))

	bytes, err := json.Marshal(results)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)

	if err != nil {
		fmt.Println("write file err, msg:" + err.Error())
	}
}

func LoadScheduleResults(directory string) ([]types.ScheduleResult, error) {

	fmt.Println(fmt.Sprintf("read file from directory:%v , sub dir:%v , fileName:%v", constants.RESULT_DIRECTORY, directory, constants.SCHEDULE_RESULT))

	scheduleResults := make([]types.ScheduleResult, 0)

	dir, err := mkdirAndReturn(directory)

	if err != nil {
		fmt.Println("open file err, msg:" + err.Error())
		return nil, err
	}

	path := filepath.Join(dir, constants.SCHEDULE_RESULT)

	data, err := ioutil.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(data, &scheduleResults)
		return scheduleResults, err
	}
	return nil, err
}

func LoadRescheduleResults(directory string) ([]types.RescheduleResult, error) {

	fmt.Println(fmt.Sprintf("read file from directory:%v , sub dir:%v , fileName:%v", constants.RESULT_DIRECTORY, directory, constants.RESCHEDULE_RESULT))

	rescheduleResults := make([]types.RescheduleResult, 0)

	dir, err := mkdirAndReturn(directory)

	if err != nil {
		fmt.Println("open file err, msg:" + err.Error())
		return nil, err
	}

	path := filepath.Join(dir, constants.RESCHEDULE_RESULT)

	data, err := ioutil.ReadFile(path)

	if err == nil {
		err = json.Unmarshal(data, &rescheduleResults)
		return rescheduleResults, err
	}
	return nil, err
}

func mkdirAndReturn(directory string) (string, error) {

	_, currentFilePath, _, _ := runtime.Caller(0)

	dataBaseDir := strings.Replace(currentFilePath, "pkg/store/result_store.go", "", 1)

	dir := filepath.Join(dataBaseDir, constants.RESULT_DIRECTORY, directory)

	err := os.MkdirAll(dir, os.ModePerm)

	return dir, err
}
