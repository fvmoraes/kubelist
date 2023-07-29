package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/fvmoraes/kubelist/internal/namespace"

	"github.com/eiannone/keyboard"
)

const kubeListFile = "~/.kube/kubelist.json"

func ReadNamespacesFromJSON() ([]namespace.Namespace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filePath := homeDir + "/.kube/kubelist.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var namespaceList namespace.NamespaceList
	if err := json.Unmarshal(data, &namespaceList); err != nil {
		return nil, err
	}

	return namespaceList.Namespaces, nil
}

func SelectNamespace(namespaces []namespace.Namespace) string {
	fmt.Println("Select the desired namespace using the arrow keys (↑ ↓) and press Enter to confirm:")

	selectedIndex := 0
	for {
		clearScreen()
		listNamespacesWithHighlight(namespaces, selectedIndex)

		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			return ""
		}

		switch key {
		case keyboard.KeyArrowUp:
			selectedIndex = (selectedIndex - 1 + len(namespaces)) % len(namespaces)
		case keyboard.KeyArrowDown:
			selectedIndex = (selectedIndex + 1) % len(namespaces)
		case keyboard.KeyEnter:
			return namespaces[selectedIndex].Name
		case keyboard.KeyEsc:
			return "default"
		}

		if char != 0 {
			continue
		}
	}
}

func listNamespacesWithHighlight(namespaces []namespace.Namespace, selectedIndex int) {
	fmt.Println("Available namespaces:")
	for i, ns := range namespaces {
		if i == selectedIndex {
			fmt.Printf("-> %s\n", ns.Name)
		} else {
			fmt.Printf("   %s\n", ns.Name)
		}
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func SetKubectlContext(namespace string) error {
	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace", namespace)
	return cmd.Run()
}
