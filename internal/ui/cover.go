package ui

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Decodificador necesario para JPEG
	_ "image/png"  // Decodificador necesario para PNG
	"os"
	"strings"

	"github.com/dhowden/tag"
	"github.com/nfnt/resize"
)

// Genera un arte ANSI a partir de un archivo MP3
func generateAsciiCover(path string, width, height int) string {
	f, err := os.Open(path)
	if err != nil {
		return fallbackCover(width, height)
	}
	defer f.Close()

	// Leer los metadatos del MP3
	m, err := tag.ReadFrom(f)
	if err != nil || m.Picture() == nil {
		return fallbackCover(width, height)
	}

	// Decodificar la imagen incrustada
	img, _, err := image.Decode(bytes.NewReader(m.Picture().Data))
	if err != nil {
		return fallbackCover(width, height)
	}

	// Redimensionamos la imagen.
	// NOTA: Multiplicamos el alto x2 porque los caracteres de la terminal son rectangulares.
	// Usaremos un truco de bloque medio (▀) para duplicar la resolución vertical.
	resized := resize.Resize(uint(width), uint(height*2), img, resize.Lanczos3)

	var buf strings.Builder
	bounds := resized.Bounds()

	// Renderizar usando TrueColor (ANSI 24-bit)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += 2 {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Color del píxel superior (para el Foreground)
			r1, g1, b1, _ := resized.At(x, y).RGBA()
			
			// Color del píxel inferior (para el Background)
			r2, g2, b2, _ := uint32(0), uint32(0), uint32(0), uint32(0)
			if y+1 < bounds.Max.Y {
				r2, g2, b2, _ = resized.At(x, y+1).RGBA()
			}

			// Convertimos RGBA a 8-bit (0-255)
			r1, g1, b1 = r1>>8, g1>>8, b1>>8
			r2, g2, b2 = r2>>8, g2>>8, b2>>8

			// Código de escape ANSI: \x1b[38;2;R;G;Bm (Texto) \x1b[48;2;R;G;Bm (Fondo) y el bloque ▀
			buf.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", r1, g1, b1, r2, g2, b2))
		}
		buf.WriteString("\x1b[0m\n") // Resetear color al final de la línea
	}

	return buf.String()
}

// Devuelve un logo genérico si no hay carátula
func fallbackCover(width, height int) string {
	var buf strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Dibujamos un patrón simple en gris
			if (x+y)%2 == 0 {
				buf.WriteString("\x1b[38;2;50;50;50m\x1b[48;2;30;30;30m▀")
			} else {
				buf.WriteString("\x1b[38;2;30;30;30m\x1b[48;2;50;50;50m▀")
			}
		}
		buf.WriteString("\x1b[0m\n")
	}
	return buf.String()
}