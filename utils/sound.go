package utils

import (
    "fmt"
    "os"
    "time"

    "github.com/faiface/beep"
    "github.com/faiface/beep/speaker"
    "github.com/faiface/beep/wav"
)

func PlaySound(path string) error {
    if path == "" {
        return nil // no sound configured
    }
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("error opening sound file: %v", err)
    }
    defer f.Close()

    streamer, format, err := wav.Decode(f)
    if err != nil {
        return fmt.Errorf("error decoding wav: %v", err)
    }
    defer streamer.Close()

    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
    done := make(chan bool)
    speaker.Play(beep.Seq(streamer, beep.Callback(func() {
        done <- true
    })))
    <-done
    return nil
}
