package models

import "log"

// DoesREADMEExist will check if readme file exists for a given task
func DoesREADMEExist(id string) bool {
	sqlStatement := `SELECT EXISTS(SELECT 1 FROM TASK_README WHERE TASK_ID=$1)`
	var exists bool
	err := DB.QueryRow(sqlStatement, id).Scan(&exists)
	if err != nil {
		log.Println(err)
		return false
	}
	if exists {
		return true
	}
	return false
}
