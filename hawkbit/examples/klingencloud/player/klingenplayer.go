package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sdobz/go-mpg123"
	"golang.org/x/mobile/exp/audio"
	"io/ioutil"
	"os"
	"strings"
	"time"
	_ "time"
)

func main() {
	fmt.Println("Starting KlingenCloud Player 1.0-ALPHA")

	speechWelcome := &Song{
		Filename: "../player/assets/speech_welcome.mp3",
		Duration: "2s",
	}

	PlaySong(speechWelcome)

	// Load playlist
	playlist, err := LoadPlaylist("../player/playlist")
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		for _, s := range playlist.Songs {
			// Play song and wait for delay
			PlaySong(s)

		}
	}
}

func PlaySong(s *Song) {
	mpg123.Initialize()

	filename := s.Filename

	fmt.Println("Playing ", filename)

	mp3, err := mpg123.Open(filename)
	if err != nil {
		fmt.Println(err.Error())
	}

	rate, channels, encoding, format := mp3.Format()
	fmt.Printf("Rate: %i Channels: %i Encoding: %i Format: %s\n", rate, channels, encoding, format)

	p, err := audio.NewPlayer(mp3, audio.Format(format), rate)
	if err != nil {
		fmt.Println(err.Error())
	}

	p.Play()
	for p.State() == audio.Playing {
		time.Sleep(time.Second)
	}
	mpg123.Exit()
}

func LoadPlaylist(f string) (pl *Playlist, err error) {

	pl = &Playlist{
		Songs: []*Song{},
	}

	b, err := ioutil.ReadFile(f)
	fmt.Println(b)

	fmt.Println("LoadPlaylist ", f)
	inFile, _ := os.Open(f)
	defer inFile.Close()
	fmt.Println("LoadPlaylist 2", inFile)
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	fmt.Println("LoadPlaylist 3")
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("LoadPlaylist 4", text)
		switch {
		case strings.HasPrefix(text, "#Title"):
			pl.Title = text[7:]
			break

		case strings.HasPrefix(text, "#Artist"):
			pl.Artist = text[8:]
			break

		default:
			s, err := ParseSongEntry(text)
			if err != nil {
				break
			}
			pl.Songs = append(pl.Songs, s)
			break
		}
	}
	return
}

func ParseSongEntry(str string) (s *Song, err error) {
	splits := strings.Split(str, ",")

	if len(splits) < 3 {
		err = errors.New("Invalid entry length")
		return
	}

	s = &Song{
		Filename: "../player/music/" + splits[0],
		Name:     splits[1],
		Duration: splits[2],
	}

	return
}

type Playlist struct {
	Title  string
	Artist string
	Songs  []*Song
}

type Song struct {
	Filename string
	Name     string
	Duration string
}
