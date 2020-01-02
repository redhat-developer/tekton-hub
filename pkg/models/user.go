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

// ResourceGithubResponse represents response for GetResourceGithubDetails query
type ResourceGithubResponse struct {
	ResourceID     int
	Owner          string
	RepositoryName string
	Path           string
	ReadmePath     string
}

// GetAllResourcesByUser will return all tasks uploaded by user
func GetAllResourcesByUser(userID int) []UserTaskResponse {
	sqlStatement := `SELECT ID,NAME,DOWNLOADS,RATING FROM RESOURCE T JOIN USER_RESOURCE 
	U ON T.ID=U.RESOURCE_ID WHERE U.USER_ID=$1`
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

// GetGithubToken will return github token by ID
func GetGithubToken(userID int) string {
	var token string
	sqlStatement := `SELECT TOKEN FROM USER_CREDENTIAL WHERE ID=$1`
	DB.QueryRow(sqlStatement, userID).Scan(&token)
	return token
}

// AddResourceRawPath will add a raw path for resource
func AddResourceRawPath(resourcePath string, resourceID int, resourceType string) {
	sqlStatement := `INSERT INTO RESOURCE_RAW_PATH(RESOURCE_ID,RAW_PATH,TYPE) VALUES($1,$2,$3)`
	_, err := DB.Exec(sqlStatement, resourceID, resourcePath, resourceType)
	if err != nil {
		log.Println(err)
	}
}

// GetResourceGithubDetails will return resource path and github details
func GetResourceGithubDetails(resourceID int) ResourceGithubResponse {
	sqlStatement := `SELECT * FROM GITHUB_DETAIL WHERE RESOURCE_ID=$1`
	githubDetails := ResourceGithubResponse{}
	DB.QueryRow(sqlStatement, resourceID).Scan(&githubDetails.ResourceID, &githubDetails.Owner, &githubDetails.RepositoryName, &githubDetails.Path, &githubDetails.ReadmePath)
	return githubDetails
}

// GetResourceRawLinks will return raw github links by ID
func GetResourceRawLinks(resourceID int) RawLinksResponse {
	sqlStatement := `SELECT * FROM RESOURCE_RAW_PATH WHERE RESOURCE_ID=$1`
	rows, err := DB.Query(sqlStatement, resourceID)
	if err != nil {
		log.Println(err)
	}
	links := RawLinksResponse{}
	for rows.Next() {
		var link string
		var rawResourceType string
		var id int
		rows.Scan(&id, &link, &rawResourceType)
		if rawResourceType == "Task" {
			links.Tasks = append(links.Tasks, link)
		} else if rawResourceType == "Pipeline" {
			links.Pipelines = append(links.Pipelines, link)
		}
	}
	return links
}
