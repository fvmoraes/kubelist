package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fvmoraes/kubelist/internal/namespace"
)

const kubeListFile = "/.kube/kubelist.json"

// Lê a lista de namespaces do arquivo JSON
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

// Função para obter o namespace atual do contexto kubectl
func getCurrentNamespace() string {
	cmd := exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath={..namespace}")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// fallback para default em caso de erro
		return "default"
	}
	ns := strings.TrimSpace(out.String())
	if ns == "" {
		return "default"
	}
	return ns
}

// Seleciona namespace via menu interativo no terminal
func SelectNamespace(namespaces []namespace.Namespace) string {
	currentNamespace := getCurrentNamespace()

	fmt.Print("\033[?1049h")       // ativa tela alternativa
	defer fmt.Print("\033[?1049l") // desativa tela alternativa ao sair

	selectedIndex := 0

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error opening keyboard:", err)
		return currentNamespace
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		fmt.Print("\033[H\033[2J") // limpa tela
		fmt.Println("Select the desired namespace using the arrow keys (↑ ↓) and press Enter to confirm:\n")
		listNamespacesWithHighlight(namespaces, selectedIndex)

		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error reading keyboard input:", err)
			return currentNamespace
		}

		switch key {
		case keyboard.KeyArrowUp:
			selectedIndex = (selectedIndex - 1 + len(namespaces)) % len(namespaces)
		case keyboard.KeyArrowDown:
			selectedIndex = (selectedIndex + 1) % len(namespaces)
		case keyboard.KeyEnter:
			return namespaces[selectedIndex].Name
		case keyboard.KeyEsc:
			// Ao sair com ESC, mantém namespace atual
			return currentNamespace
		default:
			// Para qualquer outra tecla, mantém namespace atual
			return currentNamespace
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

// Define o contexto kubectl com o namespace informado
func SetKubectlContext(namespace string) error {
	cmd := exec.Command("kubectl", "config", "set-context", "--current", "--namespace", namespace)
	return cmd.Run()
}
