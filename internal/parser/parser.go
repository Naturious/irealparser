package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	measures            [][]string
	startRepeatLocation int
	endRepeatLocation   *int
	lastChord           string
	codaLocation        int
	segnoLocation       int
	timeSignature       string
	thirdEndingImminent bool
	dcAlFineImminent    bool
	dcAlCodaImminent    bool
	dsAlCodaImminent    bool
	fineLocation        int
)

func resetVars() {
	measures = [][]string{}
	startRepeatLocation = 0
	endRepeatLocation = nil
	lastChord = ""
	codaLocation = 0
	segnoLocation = 0
	timeSignature = ""
	thirdEndingImminent = false
	dcAlFineImminent = false
	dcAlCodaImminent = false
	dsAlCodaImminent = false
	fineLocation = 0
}

type rule struct {
	token       interface{} // string or *regexp.Regexp
	description string
	operation   func([]string)
}

var rules = []rule{
	{token: "XyQ", description: "Empty space"},
	{token: regexp.MustCompile(`\*\w`), description: "Section marker"},
	{token: regexp.MustCompile(`<([^>]+)>`), description: "Comments", operation: checkForRepeats},
	{token: regexp.MustCompile(`T(\d+)`), description: "Time signature", operation: setTimeSignature},
	{token: "x", description: "Repeat last measure", operation: repeatLastMeasure},
	{token: "Kcl", description: "Repeat and add new", operation: repeatLastMeasureAndAddNew},
	{token: "r|XyQ", description: "Repeat last two", operation: repeatLastTwoMeasures},
	{token: regexp.MustCompile(`Y+`), description: "Spacers"},
	{token: "n", description: "No Chord", operation: pushNull},
	{token: "p", description: "Pause"},
	{token: "U", description: "Player ending"},
	{token: "S", description: "Segno", operation: setSegnoLocation},
	{token: "Q", description: "Coda", operation: setCodaLocation},
	{token: "{", description: "Start repeat", operation: setStartRepeatLocation},
	{token: "}", description: "End repeat", operation: repeatEverythingToEndRepeatLocation},
	{token: "LZ|", description: "Bar line", operation: createNewMeasure},
	{token: "|", description: "Bar line", operation: createNewMeasure},
	{token: "LZ", description: "Bar line", operation: createNewMeasure},
	{token: "[", description: "Double bar start", operation: createNewMeasure},
	{token: "]", description: "Double bar end", operation: repeatRemainingEndings},
	{token: regexp.MustCompile(`N(\d)`), description: "Numbered endings", operation: setEndRepeatLocation},
	{token: "Z", description: "Final bar", operation: repeatRemainingEndings},
	{token: regexp.MustCompile(`[A-GW]{1}[\+\-\^\dhob#suadlt]*(/[A-G][#b]?)?`), description: "Chord", operation: pushChordInMeasures},
}

func parse(input string) {
	input = strings.TrimSpace(input)
	for i := 0; i < len(rules); i++ {
		r := rules[i]
		switch t := r.token.(type) {
		case string:
			if strings.HasPrefix(input, t) {
				if r.operation != nil {
					r.operation(nil)
				}
				parse(strings.TrimSpace(input[len(t):]))
				return
			}
		case *regexp.Regexp:
			if match := t.FindStringSubmatchIndex(input); match != nil && match[0] == 0 {
				submatch := t.FindStringSubmatch(input)
				if r.operation != nil {
					r.operation(submatch)
				}
				parse(strings.TrimSpace(input[match[1]:]))
				return
			}
		}
	}
	if len(input) > 1 {
		parse(input[1:])
	}
}

func checkForRepeats(match []string) {
	switch strings.ToLower(match[1]) {
	case "d.c. al 3rd ending":
		thirdEndingImminent = true
	case "d.c. al fine":
		dcAlFineImminent = true
	case "d.c. al coda":
		dcAlCodaImminent = true
	case "d.s. al coda":
		dsAlCodaImminent = true
	case "fine":
		fineLocation = len(measures)
	}
}

func setTimeSignature(match []string) {
	timeSignature = match[1]
}

func pushNull(_ []string) {
	if len(measures) == 0 {
		measures = append(measures, []string{})
	}
	measures[len(measures)-1] = append(measures[len(measures)-1], "N.C.")
}

func repeatLastMeasure(_ []string) {
	measures[len(measures)-1] = measures[len(measures)-2]
}

func repeatLastMeasureAndAddNew(_ []string) {
	measures = append(measures, measures[len(measures)-1])
}

func repeatLastTwoMeasures(_ []string) {
	measures[len(measures)-1] = measures[len(measures)-3]
	measures = append(measures, measures[len(measures)-2])
}

func setStartRepeatLocation(_ []string) {
	createNewMeasure(nil)
	startRepeatLocation = len(measures) - 1
	endRepeatLocation = nil
}

func setEndRepeatLocation(match []string) {
	if n, _ := strconv.Atoi(match[1]); n == 1 {
		loc := len(measures) - 1
		endRepeatLocation = &loc
	}
}

func setSegnoLocation(_ []string) {
	segnoLocation = len(measures) - 1
}

func setCodaLocation(_ []string) {
	codaLocation = len(measures)
}

func createNewMeasure(_ []string) {
	if len(measures) == 0 || len(measures[len(measures)-1]) != 0 {
		measures = append(measures, []string{})
	}
}

func repeatEverythingToEndRepeatLocation(_ []string) {
	end := len(measures)
	if endRepeatLocation != nil {
		end = *endRepeatLocation
	}
	measures = append(measures, measures[startRepeatLocation:end]...)
	createNewMeasure(nil)
}

func repeatRemainingEndings(_ []string) {
	switch {
	case thirdEndingImminent:
		repeatEverythingToEndRepeatLocation(nil)
		thirdEndingImminent = false
	case dcAlFineImminent:
		measures = append(measures, measures[:fineLocation]...)
		dcAlFineImminent = false
	case dcAlCodaImminent:
		measures = append(measures, measures[:codaLocation]...)
		dcAlCodaImminent = false
	case dsAlCodaImminent:
		measures = append(measures, measures[segnoLocation:codaLocation]...)
		dsAlCodaImminent = false
	default:
		createNewMeasure(nil)
	}
}

func pushChordInMeasures(match []string) {
	if len(measures) == 0 {
		measures = append(measures, []string{})
	}
	chord := match[0]
	if strings.HasPrefix(chord, "W") && lastChord != "" {
		chord = strings.Replace(chord, "W", lastChord, 1)
	} else {
		lastChord = strings.Split(chord, "/")[0]
	}
	measures[len(measures)-1] = append(measures[len(measures)-1], chord)
}

func removeEmptyMeasures(input [][]string) [][]string {
	var output [][]string
	for _, m := range input {
		if len(m) > 0 {
			output = append(output, m)
		}
	}
	return output
}

func ParseChart(raw string) ([][]string, string) {
	resetVars()
	parse(raw)
	return removeEmptyMeasures(measures), timeSignature
}
