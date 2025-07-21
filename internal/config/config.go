package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/eiannone/keyboard"
	"github.com/fvmoraes/kubelist/internal/namespace"
)

const kubeListFile = "/.kube/kubelist.json"

func ReadNamespacesFromJSON() ([]namespace.Namespace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filePath := homeDir + kubeListFile
	data, err := os.ReadFile(filePath)
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
	// Ativa a tela alternativa (modo temporário)
	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l") // Sai da tela alternativa ao sair da função

	selectedIndex := 0

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return ""
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		// Limpa a tela da alternativa antes de desenhar
		fmt.Print("\033[H\033[2J")

		fmt.Println("Select the desired namespace using the arrow keys (↑ ↓) and press Enter to confirm:\n")
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
		default:
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
			fmt.Printf("\033[33m->\033[0m \033[1m%s\033[0m\n", ns.Name)
		} else {
			fmt.Printf("   %s\n", ns.Name)
		}
	}
}

func SetKubectlContext(namespace string) error {
	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace", namespace)
	return cmd.Run()
}
