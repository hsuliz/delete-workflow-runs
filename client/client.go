package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func createUrl(owner string, repository string) string {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/workflows",
		owner, repository)
	return url
}

func ListWorkflows(owner string, repository string) string {
	url := createUrl(owner, repository)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		os.Exit(1)
	}

	gitHubToken := "$token"
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gitHubToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		os.Exit(1)
	}
	log.Println(resp.Status)

	return string(body)
}
