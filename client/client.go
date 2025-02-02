package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ActionsRuns struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

type WorkflowRun struct {
	ID   int64 `json:"id"`
	Name string  `json:"name"`
	NodeID string `json:"node_id"`
}

func ListWorkflows(owner string, repository string, workflowName string) []WorkflowRun {
	pageSize := 1
	body := getWorkflowRuns(owner, repository, pageSize)
	var result ActionsRuns
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("Could not unmarshal json: %s\n", err)
		os.Exit(1)
	}
	workflowRuns := result.WorkflowRuns

	for {
		result = ActionsRuns{}
		pageSize++
		body := getWorkflowRuns(owner, repository, pageSize)

		err := json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("Could not unmarshal json: %s\n", err)
			os.Exit(1)
		}

		if len(result.WorkflowRuns) == 0 {
			break
		}

		workflowRuns = append(workflowRuns, result.WorkflowRuns...)
	}


	var filteredWorkflowRuns []WorkflowRun
	for i, workflowRun := range workflowRuns {
		log.Println(i, workflowRun.Name, workflowRun.NodeID)
		if (func(s string) bool {
			log.Println(s, workflowName, workflowRun.ID, s == workflowName)
			return s == workflowName
		})(workflowRun.Name) {
			filteredWorkflowRuns = append(filteredWorkflowRuns, workflowRun)
		}
	}

	return filteredWorkflowRuns
}

func DeleteWorkflowRun(owner string, repository string, runID int64) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs/%d",
		owner, repository, runID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		os.Exit(1)
	}

	pupulateRequest(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		os.Exit(1)
	}
	if resp.StatusCode != 204 {
		log.Fatal("Error making request: ", resp.StatusCode)
	}
}

func getWorkflowRuns(owner string, repository string, pageSize int) []byte {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs?page=%d",
		owner, repository, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		os.Exit(1)
	}

	pupulateRequest(req)

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

	return body
}

func pupulateRequest(req *http.Request) {
	gitHubToken := "$token"
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gitHubToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
}
