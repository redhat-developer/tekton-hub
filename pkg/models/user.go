package models

import "log"

// User represents User model in database
type User struct {
	ID         int    `json:"id"`
	FirstName  string `json:"username"`
	SecondName string `json:"password"`
	EMAIL      string `json:"email"`
}

// UserTaskResponse represents all tasks uploaded by user
type UserTaskResponse struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Rating    float64 `json:"rating"`
	Downloads int     `json:"downloads"`
}

// GetAllTasksByUser will return all tasks uploaded by user
func GetAllTasksByUser(userID int) []UserTaskResponse {
	sqlStatement := `SELECT ID,NAME,DOWNLOADS,RATING FROM TASK T JOIN USER_TASK 
	U ON T.ID=U.TASK_ID WHERE U.USER_ID=$1`
	rows, err := DB.Query(sqlStatement, userID)
	if err != nil {
		log.Println(err)
	}
	tasks := []UserTaskResponse{}
	for rows.Next() {
		var id int
		var name string
		var rating float64
		var downloads int
		rows.Scan(&id, &name, &downloads, &rating)
		task := UserTaskResponse{id, name, rating, downloads}
		tasks = append(tasks, task)
	}
	return tasks
}
