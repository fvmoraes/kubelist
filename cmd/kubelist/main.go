package main

import (
	"fmt"
	"os"

	"github.com/fvmoraes/kubelist/internal/config"
)

func main() {
	namespaces, err := config.ReadNamespacesFromJSON()
	if err != nil {
		fmt.Println("Error reading namespaces from JSON file:", err)
		return
	}

	var selectedNamespace string

	if len(os.Args) > 1 {
		selectedNamespace = os.Args[1]
	} else {
		selectedNamespace = config.SelectNamespace(namespaces)
		if selectedNamespace == "" {
			fmt.Println("No namespace selected. Exiting...")
			return
		}
	}

	if err := config.SetKubectlContext(selectedNamespace); err != nil {
		fmt.Println("Error setting the namespace context:", err)
		return
	}

	fmt.Printf("Selected namespace: %s\n", selectedNamespace)
}
