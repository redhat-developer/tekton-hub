package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/backend/pkg/polling"
)

func getJSON(URL string) string {
	resp, err := http.Get(URL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	j := ""
	json.Unmarshal([]byte(body), &j)
	log.Println(j)
	return string(j)
}

func main() {
	client, ctx := polling.Authenticate()

	repoContents, err := polling.GetDirContents(client, ctx, "tektoncd", "catalog", "", nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, c := range repoContents {
		log.Println(c.GetName())
		dirContent, err := polling.GetDirContents(client, ctx, "tektoncd", "catalog", c.GetName(), nil)
		if err != nil {
			log.Fatalln(err)
		}
		for _, k := range dirContent {
			// Handle YAML files
			if strings.HasSuffix(k.GetName(), ".yaml") {
				yamlContent, err := polling.GetFileContent(client, ctx, "tektoncd", "catalog", c.GetName()+"/"+k.GetName(), nil)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(yamlContent.GetContent())
			}
			// Handle README files
			if strings.HasSuffix(k.GetName(), ".md") {
				readmeContent, err := polling.GetFileContent(client, ctx, "tektoncd", "catalog", c.GetName()+"/"+k.GetName(), nil)
				if err != nil {
					log.Fatalln(err)
				}
				log.Println(readmeContent.GetContent())
			}
		}
	}
}
