package songbook

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/naturious/irealparser/internal/parser"
	"github.com/naturious/irealparser/internal/unscramble"
)

type Song struct {
	Title     string
	Composer  string
	Style     string
	Key       string
	Transpose *int
	Music     [][]string
	TimeSig   string
	CompStyle string
	BPM       *int
	Repeats   *int
}

type Book struct {
	Name  string
	Songs []Song
}

var protocolRegex = regexp.MustCompile(`irealb:\/\/([^"]*)`)

const musicPrefix = "1r34LbKcu7"

func ParseIRealInput(data string) Book {
	match := protocolRegex.FindStringSubmatch(data)
	if match == nil {
		return Book{}
	}
	raw := decode(match[1])
	parts := strings.Split(raw, "===")

	name := ""
	if len(parts) > 1 {
		name = strings.TrimSpace(parts[len(parts)-1])
		parts = parts[:len(parts)-1]
	}

	songs := []Song{}
	for _, p := range parts {
		songs = append(songs, makeSong(p))
	}

	return Book{Name: name, Songs: songs}
}

func decode(s string) string {
	decoded, _ := url.QueryUnescape(s)
	return decoded
}

func makeSong(data string) Song {
	fields := strings.Split(data, "=")
	// Remove empty entries
	nonEmpty := []string{}
	for _, f := range fields {
		if strings.TrimSpace(f) != "" {
			nonEmpty = append(nonEmpty, f)
		}
	}

	var title, composer, style, key, transpose, music, compStyle, bpm, repeats string

	switch len(nonEmpty) {
	case 7:
		title, composer, style, key, music, bpm, repeats = nonEmpty[0], nonEmpty[1], nonEmpty[2], nonEmpty[3], nonEmpty[4], nonEmpty[5], nonEmpty[6]
	case 8:
		if strings.HasPrefix(nonEmpty[4], musicPrefix) {
			title, composer, style, key, music, compStyle, bpm, repeats = nonEmpty[0], nonEmpty[1], nonEmpty[2], nonEmpty[3], nonEmpty[4], nonEmpty[5], nonEmpty[6], nonEmpty[7]
		} else {
			title, composer, style, key, transpose, music, bpm, repeats = nonEmpty[0], nonEmpty[1], nonEmpty[2], nonEmpty[3], nonEmpty[4], nonEmpty[5], nonEmpty[6], nonEmpty[7]
		}
	case 9:
		title, composer, style, key, transpose, music, compStyle, bpm, repeats = nonEmpty[0], nonEmpty[1], nonEmpty[2], nonEmpty[3], nonEmpty[4], nonEmpty[5], nonEmpty[6], nonEmpty[7], nonEmpty[8]
	}

	ts := 0
	if transpose != "" {
		ts, _ = strconv.Atoi(transpose)
	}
	bpmInt := 0
	if bpm != "" {
		bpmInt, _ = strconv.Atoi(bpm)
	}
	repeatInt := 0
	if repeats != "" {
		repeatInt, _ = strconv.Atoi(repeats)
	}

	decodedMusic := strings.Split(music, musicPrefix)[1]
	unscrambled := unscramble.IReal(decodedMusic)
	measures, timeSig := parser.ParseChart(unscrambled)

	return Song{
		Title:     title,
		Composer:  composer,
		Style:     style,
		Key:       key,
		Transpose: &ts,
		Music:     measures,
		TimeSig:   timeSig,
		CompStyle: compStyle,
		BPM:       &bpmInt,
		Repeats:   &repeatInt,
	}
}
