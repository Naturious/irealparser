package songbook

import (
	"os"
	"testing"
)

func TestParseIRealInput(t *testing.T) {
	data, err := os.ReadFile("spec/Tester.html")
	if err != nil {
		t.Fatalf("could not read test file: %v", err)
	}

	book := ParseIRealInput(string(data))

	if len(book.Songs) != 14 {
		t.Errorf("expected 14 songs, got %d", len(book.Songs))
	}

	testSong := book.Songs[0]
	neverBeSame := book.Songs[1]
	paperMoon := book.Songs[2]
	theBreezeAndI := book.Songs[3]
	imagination := book.Songs[4]
	forJan := book.Songs[5]
	allOfMe := book.Songs[6]
	comeFly := book.Songs[7]
	daphne := book.Songs[8]
	inHerFam := book.Songs[9]
	mySong := book.Songs[10]
	noComp := book.Songs[11]
	masquerade := book.Songs[12]
	aTasteOfHoney := book.Songs[13]

	// Example of testing one song
	t.Run("Test Song basic metadata", func(t *testing.T) {
		if testSong.Title != "Test" {
			t.Errorf("got title %q, want 'Test'", testSong.Title)
		}
		if testSong.Composer != "Florin" {
			t.Errorf("got composer %q, want 'Florin'", testSong.Composer)
		}
		if testSong.Style != "Medium Swing" {
			t.Errorf("got style %q, want 'Medium Swing'", testSong.Style)
		}
		if testSong.Key != "C" {
			t.Errorf("got key %q, want 'C'", testSong.Key)
		}
		if testSong.BPM != 140 {
			t.Errorf("got BPM %d, want 140", testSong.BPM)
		}
		if testSong.Transpose != 2 {
			t.Errorf("got transpose %d, want 2", testSong.Transpose)
		}
		if testSong.CompStyle != "Latin-Brazil: Bossa Acoustic" {
			t.Errorf("got compStyle %q, want 'Latin-Brazil: Bossa Acoustic'", testSong.CompStyle)
		}
		if testSong.Repeats != 3 {
			t.Errorf("got repeats %d, want 3", testSong.Repeats)
		}
		if testSong.TimeSig != "34" {
			t.Errorf("got timesig %q, want '34'", testSong.TimeSig)
		}
	})

	t.Run("Test Song chords start with A-G or null", func(t *testing.T) {
		for i, measure := range testSong.Music {
			for j, chord := range measure {
				if chord != "" && !isChordLetter(chord) {
					t.Errorf("measure %d chord %d = %q is invalid", i+1, j+1, chord)
				}
			}
		}
	})

	// ... write similar t.Run blocks for neverBeSame, paperMoon, etc.
}
