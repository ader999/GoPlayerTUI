package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// Asegúrate de que el módulo coincida con lo que pusiste en 'go mod init'
	"goplayertui/internal/ui"
)

func main() {
	// Inicializamos el modelo desde nuestro paquete ui
	initialModel := ui.InitialModel()

	// Creamos el programa en modo "Alt Screen" para que actúe como una app nativa
	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	// Ejecutamos el motor
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error fatal al iniciar GoPlayerTUI: %v", err)
		os.Exit(1)
	}
}