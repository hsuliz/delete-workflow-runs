package main

import (

    "hsuliz/dwr/client"
    "log"
    "os"
)

func main() {
    workflowRuns := client.ListWorkflows("hsuliz", "terraform-ansible-sample", "Destroy")
    log.Println("Workflows founded:")
    for _, workflow := range workflowRuns {
        client.DeleteWorkflowRun("hsuliz", "terraform-ansible-sample", workflow.ID)
    }
    os.Exit(0)
}
