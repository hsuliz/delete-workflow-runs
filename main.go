package main

import (
    "flag"
    "log"
)

func main() {
    owner := flag.String("o", "", "")
    repository := flag.String("r", "", "")
    bearerToken := flag.String("bt", "", "")
    workflowName := flag.String("wn", "", "Optional")
    flag.Parse()

    if *owner == "" || *repository == "" || *bearerToken == "" {
        log.Println("Please provide arguments.")
        return
    }

    gitHubClient := GitHub{
        Owner:       *owner,
        Repository:  *repository,
        BearerToken: *bearerToken,
    }

    workflowRuns := gitHubClient.ListWorkflows(*workflowName)
    log.Println("Found", len(workflowRuns), "workflow runs")
    for i, workflow := range workflowRuns {
        gitHubClient.DeleteWorkflowRun(workflow.ID)
        log.Println("Deleted:", i+1, "/", len(workflowRuns))
    }
}
