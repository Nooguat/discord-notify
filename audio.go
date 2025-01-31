// +build !gen

package main

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"os"
	"strings"
	"time"
)

type NotifySound struct {
	buffer *beep.Buffer
	name   string
	format beep.Format
}

var currentSound NotifySound

func setSoundFromDisk(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error loading sound: %w", err)
	}
	defer f.Close()
	return setSound(f, path)
}
func setSoundBuiltin(name string) error {
	if name == "" || name == "none" {
		setSoundNone()
		return nil
	}
	f, err := BinFS.Open("assets/sounds/" + name + ".mp3")
	if err != nil {
		return err
	}
	return setSound(f, name)
}

func builtInSounds() (results []string) {
	for _, file := range BinFS.Files {
		if strings.HasSuffix(file.Filename, ".mp3") {
			results = append(results, strings.TrimSuffix(file.Filename, ".mp3"))
		}
	}
	return results
}

func setSoundNone() {
	currentSound = NotifySound{name: "none"}
}
func setSound(data io.ReadCloser, name string) error {
	streamer, format, err := mp3.Decode(data)
	currentSound = NotifySound{name: name, format: format, buffer: beep.NewBuffer(format)}
	speaker.Init(currentSound.format.SampleRate, currentSound.format.SampleRate.N(time.Second/10))
	if err != nil {
		return fmt.Errorf("error decoding sound: %w", err)
	}
	currentSound.buffer.Append(streamer)
	return streamer.Close()
}

func playSound() {

	sound := currentSound.buffer.Streamer(0, currentSound.buffer.Len())	
	speaker.Play(sound)
}

