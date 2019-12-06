package models

import "log"

// GithubSHA represents relationship between Github repo and User
type GithubSHA struct {
	ID     int    `json:"id"`
	TaskID int    `json:"user_id"`
	SHA    string `json:"ssh"`
}

// AddNewSHA will add a new commit SHA key
func AddNewSHA(taskID int, SHA string) {
	sqlStatement := `INSERT INTO GITHUB_SHA(TASK_ID,SHA) VALUES($1,$2)`
	_, err := DB.Exec(sqlStatement, taskID, SHA)
	if err != nil {
		log.Println(taskID)
		log.Println(err)
	}
}
