# GoPlayerTUI 🎵

Un reproductor de música ligero y elegante para la terminal (TUI) escrito en Go. GoPlayerTUI te permite explorar y reproducir tu biblioteca local de MP3 directamente desde la consola, destacando por su renderizado de carátulas de álbum en alta resolución utilizando caracteres ASCII.

## Características

- **Interfaz TUI Moderna:** Interfaz limpia construida con las librerías Bubbletea y Lipgloss.
- **Arte de Álbum con Estilo Retro Pixelado:** Extrae la carátula incrustada en tus MP3 y la renderiza en la terminal utilizando verdaderos colores (TrueColor ANSI) y medios bloques para lograr un atractivo y nostálgico estilo pixel art.
- **Gestión de Biblioteca:** Agrupa automáticamente tus pistas por álbum/carpeta (pestañas navegables).
- **Controles Integrados:** Pausa, reanuda, adelanta o atrasa la música con atajos de teclado rápidos.
- **Barra de Progreso:** Visualización del progreso de la canción actual.

## Instalación

Asegúrate de tener Go instalado en tu sistema.

```bash
# Clona el repositorio
git clone https://github.com/ader999/GoPlayerTUI.git
cd GoPlayerTUI

# Compila el ejecutable
go build -o goplayertui ./cmd/goplayertui

# (Opcional) Instala en el sistema para usarlo desde cualquier parte
sudo mv goplayertui /usr/local/bin/
```

## Uso

Simplemente ejecuta el comando en tu terminal:

```bash
goplayertui
```

*(Nota: En esta versión inicial, el directorio de música apunta a `~/Music/` por defecto).*

## Atajos de Teclado

- `[↑ / ↓]` o `[j / k]`: Navegar por la lista
- `[Enter]`: Entrar al álbum seleccionado / Reproducir canción
- `[Tab]`: Cambiar entre la pestaña de "Todas las Canciones" y "Álbumes"
- `[Espacio]`: Pausar o reanudar la reproducción
- `[← / →]` o `[h / l]`: Atrasar o adelantar la canción 5 segundos
- `[Esc]`: Volver atrás (de los detalles del álbum a la lista de carpetas)
- `[q]` o `Ctrl+C`: Salir del reproductor

## Tecnologías Utilizadas

- [Bubbletea](https://github.com/charmbracelet/bubbletea) & [Lipgloss](https://github.com/charmbracelet/lipgloss) para la interfaz gráfica TUI.
- [Beep](https://github.com/gopxl/beep) para la decodificación y reproducción de audio.
- [Tag](https://github.com/dhowden/tag) para leer metadatos de los MP3.
- [Resize](https://github.com/nfnt/resize) para el escalado de imágenes ANSI.
