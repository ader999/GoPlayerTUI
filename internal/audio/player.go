package audio

import (
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type Player struct {
	ctrl     *beep.Ctrl
	streamer beep.StreamSeekCloser
	format   beep.Format
}

func NewPlayer() *Player {
	return &Player{}
}

func (p *Player) PlayFile(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}

	p.streamer = streamer
	p.format = format

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	p.ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}

	speaker.Clear()
	speaker.Play(p.ctrl)

	return nil
}

func (p *Player) TogglePause() bool {
	if p.ctrl != nil {
		speaker.Lock()
		p.ctrl.Paused = !p.ctrl.Paused
		speaker.Unlock()
		return p.ctrl.Paused
	}
	return false
}

func (p *Player) Progress() (time.Duration, time.Duration) {
	if p.streamer == nil {
		return 0, 0
	}
	speaker.Lock()
	defer speaker.Unlock()

	pos := p.format.SampleRate.D(p.streamer.Position())
	total := p.format.SampleRate.D(p.streamer.Len())
	return pos, total
}

// Avanza la reproducción por la duración especificada
func (p *Player) SeekForward(d time.Duration) {
	if p.streamer == nil {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()

	current := p.streamer.Position()
	offset := p.format.SampleRate.N(d)
	newPos := current + offset
	if newPos > p.streamer.Len() {
		newPos = p.streamer.Len()
	}
	p.streamer.Seek(newPos)
}

// Retrocede la reproducción por la duración especificada
func (p *Player) SeekBackward(d time.Duration) {
	if p.streamer == nil {
		return
	}
	speaker.Lock()
	defer speaker.Unlock()

	current := p.streamer.Position()
	offset := p.format.SampleRate.N(d)
	newPos := current - offset
	if newPos < 0 {
		newPos = 0
	}
	p.streamer.Seek(newPos)
}


// NUEVO: Comprueba si la canción llegó a su fin
func (p *Player) IsFinished() bool {
	if p.streamer == nil {
		return false
	}
	speaker.Lock()
	defer speaker.Unlock()
	// Si la posición actual alcanzó (o superó por unos bytes) a la longitud total
	return p.streamer.Position() > 0 && p.streamer.Position() >= p.streamer.Len()
}