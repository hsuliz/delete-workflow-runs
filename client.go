package main

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
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	NodeID string `json:"node_id"`
}

type GitHub struct {
	Owner       string
	Repository  string
	BearerToken string
}

func (gitHub GitHub) ListWorkflows(workflowName string) []WorkflowRun {
	pageSize := 1
	body := gitHub.getWorkflowRuns(gitHub.Owner,
		gitHub.Repository,
		pageSize,
	)
	var result ActionsRuns
	err := json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Could not unmarshal json:", err)
		os.Exit(1)
	}
	workflowRuns := result.WorkflowRuns

	for {
		result = ActionsRuns{}
		pageSize++
		body := gitHub.getWorkflowRuns(gitHub.Owner,
			gitHub.Repository,
			pageSize,
		)

		err := json.Unmarshal(body, &result)
		if err != nil {
			log.Println("Could not unmarshal json:", err)
			os.Exit(1)
		}

		if len(result.WorkflowRuns) == 0 {
			break
		}

		workflowRuns = append(workflowRuns, result.WorkflowRuns...)
	}
	if workflowName != "" {
		var filteredWorkflowRuns []WorkflowRun
		for _, workflowRun := range workflowRuns {
			if (func(s string) bool {
				return s == workflowName
			})(workflowRun.Name) {
				filteredWorkflowRuns = append(filteredWorkflowRuns, workflowRun)
			}
		}
		return filteredWorkflowRuns
	}

	return workflowRuns
}

func (gitHub GitHub) DeleteWorkflowRun(runID int64) {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs/%d",
		gitHub.Owner, gitHub.Repository, runID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		os.Exit(1)
	}

	gitHub.pupulateRequest(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	if resp.StatusCode != 204 {
		log.Fatal("Error making request: ", resp.StatusCode)
	}
}

func (gitHub GitHub) getWorkflowRuns(owner string, repository string, pageSize int) []byte {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/runs?page=%d",
		owner, repository, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	gitHub.pupulateRequest(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body:", err)
	}

	return body
}

func (gitHub GitHub) pupulateRequest(req *http.Request) {
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gitHub.BearerToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
}
