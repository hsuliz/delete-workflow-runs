package main

import (
	"fmt"
	"hsuliz/dwr/client"
)

func main() {
	workflows := client.ListWorkflows("hsuliz", "terraform-ansible-sample")
	fmt.Println(workflows)
}
