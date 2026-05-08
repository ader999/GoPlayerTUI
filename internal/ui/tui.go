package ui

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"goplayertui/internal/audio"
)

const (
	modeAllSongs = iota
	modeAlbums
	modeAlbumDetail
)

type Song struct {
	Title    string
	Artist   string
	FilePath string
}

type Model struct {
	allSongs   []Song
	albums     map[string][]Song
	albumNames[]string

	currentMode   int
	selectedAlbum string
	cursor        int

	playingList[]Song
	playingIndex int
	playingFile  string
	isPlaying    bool

	status     string
	player     *audio.Player
	coverCache map[string]string
}

var (
	headerStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF")).Padding(0, 1).MarginBottom(1)
	activeTabStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#00FFFF")).Padding(0, 1).MarginRight(1)
	inactiveTabStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Padding(0, 1).MarginRight(1)

	playingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	cursorStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#333333"))
	normalStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
	folderStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")).Bold(true)

	// Ajustamos el ancho de la lista para dar más espacio
	listColumnStyle  = lipgloss.NewStyle().Width(55).MarginRight(4)
	
	// Estilo de la carátula
	coverColumnStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true).BorderForeground(lipgloss.Color("#444444")).Padding(0, 1)

	// Ajustamos el footer al nuevo ancho total (Lista + Espacio + Carátula HD)
	footerStyle     = lipgloss.NewStyle().MarginTop(1).Border(lipgloss.NormalBorder(), true, false, false, false).PaddingTop(1).Width(135)
	playStateStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	pauseStateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
)

type tickMsg time.Time

func loadMusic(dir string) ([]Song, map[string][]Song, []string) {
	var all[]Song
	albums := make(map[string][]Song)
	var albumNames[]string

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil { return nil }
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".mp3") {
			title := strings.TrimSuffix(d.Name(), ".mp3")
			artist := filepath.Base(filepath.Dir(path))

			song := Song{Title: title, Artist: artist, FilePath: path}
			all = append(all, song)

			if len(albums[artist]) == 0 {
				albumNames = append(albumNames, artist)
			}
			albums[artist] = append(albums[artist], song)
		}
		return nil
	})

	sort.Strings(albumNames)
	return all, albums, albumNames
}

func InitialModel() Model {
	musicDir := "/home/overader/Music/"
	allSongs, albums, albumNames := loadMusic(musicDir)

	if len(allSongs) == 0 {
		allSongs = append(allSongs, Song{Title: "Sin MP3", Artist: "Revisa la ruta", FilePath: ""})
	}

	return Model{
		allSongs:    allSongs,
		albums:      albums,
		albumNames:  albumNames,
		currentMode: modeAllSongs,
		cursor:      0,
		status:      fmt.Sprintf("Cargados %d tracks y %d álbumes", len(allSongs), len(albumNames)),
		player:      audio.NewPlayer(),
		coverCache:  make(map[string]string),
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) Init() tea.Cmd { return tickCmd() }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.isPlaying && m.player.IsFinished() {
			if m.playingIndex < len(m.playingList)-1 {
				m.playingIndex++
				nextSong := m.playingList[m.playingIndex]
				m.player.PlayFile(nextSong.FilePath)
				m.playingFile = nextSong.FilePath
				m.status = "Reproduciendo siguiente: " + nextSong.Title
			} else {
				m.isPlaying = false
				m.status = "Fin de la lista."
			}
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.currentMode == modeAllSongs { m.currentMode = modeAlbums } else { m.currentMode = modeAllSongs }
			m.cursor = 0
		case "esc", "backspace":
			if m.currentMode == modeAlbumDetail {
				m.currentMode = modeAlbums
				m.cursor = 0
			}
		case "up", "k":
			if m.cursor > 0 { m.cursor-- }
		case "down", "j":
			listLen := len(m.allSongs)
			if m.currentMode == modeAlbums { listLen = len(m.albumNames) }
			if m.currentMode == modeAlbumDetail { listLen = len(m.albums[m.selectedAlbum]) }
			if m.cursor < listLen-1 { m.cursor++ }
		case "enter":
			if m.currentMode == modeAlbums {
				m.selectedAlbum = m.albumNames[m.cursor]
				m.currentMode = modeAlbumDetail
				m.cursor = 0
			} else {
				if m.currentMode == modeAllSongs { m.playingList = m.allSongs } else if m.currentMode == modeAlbumDetail { m.playingList = m.albums[m.selectedAlbum] }
				m.playingIndex = m.cursor
				song := m.playingList[m.playingIndex]

				if song.FilePath != "" {
					err := m.player.PlayFile(song.FilePath)
					if err == nil {
						m.playingFile = song.FilePath
						m.isPlaying = true
						m.status = "Reproduciendo: " + song.Title
					}
				}
			}
		case " ":
			if m.playingFile != "" {
				isPaused := m.player.TogglePause()
				m.isPlaying = !isPaused
				if m.isPlaying { m.status = "Reproducción reanudada." } else { m.status = "Reproducción pausada." }
			}
		case "right", "l":
			if m.isPlaying || m.playingFile != "" {
				m.player.SeekForward(5 * time.Second)
				m.status = "Adelantado 5s"
			}
		case "left", "h":
			if m.isPlaying || m.playingFile != "" {
				m.player.SeekBackward(5 * time.Second)
				m.status = "Atrasado 5s"
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	header := headerStyle.Render("GoPlayer TUI 🎵")

	var tabs string
	if m.currentMode == modeAllSongs {
		tabs = activeTabStyle.Render("Todas las Canciones") + inactiveTabStyle.Render("Álbumes (Carpetas)")
	} else if m.currentMode == modeAlbums {
		tabs = inactiveTabStyle.Render("Todas las Canciones") + activeTabStyle.Render("Álbumes (Carpetas)")
	} else {
		tabs = inactiveTabStyle.Render("Todas las Canciones") + activeTabStyle.Render("Álbum: "+m.selectedAlbum)
	}
	headerAndTabs := lipgloss.JoinVertical(lipgloss.Left, header, tabs, "")

	var listLen int
	var highlightedFilePath string

	if m.currentMode == modeAllSongs {
		listLen = len(m.allSongs)
		if listLen > 0 { highlightedFilePath = m.allSongs[m.cursor].FilePath }
	} else if m.currentMode == modeAlbums {
		listLen = len(m.albumNames)
		if listLen > 0 { highlightedFilePath = m.albums[m.albumNames[m.cursor]][0].FilePath }
	} else {
		listLen = len(m.albums[m.selectedAlbum])
		if listLen > 0 { highlightedFilePath = m.albums[m.selectedAlbum][m.cursor].FilePath }
	}

	// AUMENTADO: Mostraremos 32 elementos en la lista para emparejar la altura de la nueva imagen
	const maxVisible = 32
	start, end := 0, listLen
	if end > maxVisible {
		start = m.cursor - maxVisible/2
		if start < 0 { start = 0 }
		end = start + maxVisible
		if end > listLen { end = listLen; start = end - maxVisible }
	}

	var body strings.Builder
	for i := start; i < end; i++ {
		var row string
		var isPlayingRow bool

		if m.currentMode == modeAlbums {
			albumName := m.albumNames[i]
			count := len(m.albums[albumName])
			nameDisplay := albumName
			if len(nameDisplay) > 30 { nameDisplay = nameDisplay[:27] + "..." }

			row = fmt.Sprintf("  📁 %-30s | %02d trks", nameDisplay, count)
			if i == m.cursor { row = cursorStyle.Render(folderStyle.Render(row)) } else { row = normalStyle.Render(row) }
		} else {
			var song Song
			if m.currentMode == modeAllSongs { song = m.allSongs[i] } else { song = m.albums[m.selectedAlbum][i] }

			prefix := "  "
			if song.FilePath == m.playingFile { prefix = "▶ "; isPlayingRow = true }

			titleDisplay := song.Title
			if len(titleDisplay) > 25 { titleDisplay = titleDisplay[:22] + "..." }
			artistDisplay := song.Artist
			if len(artistDisplay) > 15 { artistDisplay = artistDisplay[:12] + "..." }

			rawStr := fmt.Sprintf("%s %-25s | %-15s", prefix, titleDisplay, artistDisplay)

			if isPlayingRow { row = playingStyle.Render(rawStr) } else { row = normalStyle.Render(rawStr) }
			if i == m.cursor { row = cursorStyle.Render(row) }
		}
		body.WriteString(row + "\n")
	}

	// CARÁTULA HD Y PROPORCIÓN CORREGIDA
	coverArt := ""
	if highlightedFilePath != "" {
		if cachedCover, exists := m.coverCache[highlightedFilePath]; exists {
			coverArt = cachedCover
		} else {
			// Antes era (50, 25).
			// Ahora es (70, 35).
			// 70 caracteres de ancho x 35 de alto.
			// Como usamos medio bloque (▀), 35 de alto son en realidad 70 "píxeles" verticales.
			// Esto nos da una imagen de 70x70 píxeles reales (cuadrado perfecto 1:1) y mucha más definición.
			newCover := generateAsciiCover(highlightedFilePath, 70, 35) 
			m.coverCache[highlightedFilePath] = newCover
			coverArt = newCover
		}
	}
	coverBox := coverColumnStyle.Render(coverArt)

	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, listColumnStyle.Render(body.String()), coverBox)

	progressBar := ""
	if m.playingFile != "" {
		pos, total := m.player.Progress()
		if total > 0 {
			pct := float64(pos) / float64(total)
			barWidth := 50
			filled := int(pct * float64(barWidth))
			if filled > barWidth { filled = barWidth }
			if filled < 0 { filled = 0 }

			progressBar = fmt.Sprintf("\nProgreso: %02d:%02d [%s%s] %02d:%02d",
				int(pos.Minutes()), int(pos.Seconds())%60,
				strings.Repeat("█", filled), strings.Repeat("░", barWidth-filled),
				int(total.Minutes()), int(total.Seconds())%60)
		}
	}

	stateVisual := pauseStateStyle.Render("[ PAUSE ]")
	if m.isPlaying { stateVisual = playStateStyle.Render("[ PLAY  ]") }

	controls := "Atajos: [↑/↓] Navegar | [←/→] Adelantar/Atrasar | [enter] Entrar/Play | [tab] Pestaña | [esc] Volver | [espacio] Pausa"
	footerContent := fmt.Sprintf("%s  %s\n%s\n\n%s", stateVisual, normalStyle.Render(m.status), playStateStyle.Render(progressBar), controls)
	footer := footerStyle.Render(footerContent)

	app := lipgloss.JoinVertical(lipgloss.Left, headerAndTabs, mainContent, footer)
	return lipgloss.NewStyle().Margin(1, 2).Render(app)
}