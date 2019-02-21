package main

import (
	"honnef.co/go/js/dom"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Sanitize converts the input string to uppercase and removes all characters
// that don't match [A-Z].
func sanitize(s string) (string, error) {
	re, err := regexp.Compile("[^A-Z]")
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(strings.ToUpper(s), ""), nil
}

// sortWords reads a slice of strings and sorts each line by word.
func sortWords(lines []string) {
	for idx, line := range lines {
		w := strings.Split(line, " ")
		sort.Strings(w)
		lines[idx] = strings.Join(w, " ")
	}
}

func main() {
	doc := dom.GetWindow().Document()
	resultsDiv := doc.GetElementByID("results").(*dom.HTMLDivElement)
	waitDiv := doc.GetElementByID("wait").(*dom.HTMLDivElement)
	button := doc.GetElementByID("button1").(*dom.HTMLButtonElement)

	button.AddEventListener("click", false, func(event dom.Event) {
		event.PreventDefault()
		go func(d *dom.HTMLDivElement) {
			d.SetAttribute("style", "background-color: #ffff00")
			d.SetInnerHTML("<h2>Please wait...</h2>")
			d.Style().SetProperty("display", "block", "")
		}(waitDiv)

		go func(d, w *dom.HTMLDivElement) {
			word := doc.GetElementByID("word").(*dom.HTMLInputElement).Value
			phrase, err := sanitize(word)
			if err != nil {
				d.SetInnerHTML("<h2>Error: " + err.Error() + "</h2>")
				return
			}

			minWordLen, _ := strconv.Atoi(doc.GetElementByID("minWordLen").(*dom.HTMLInputElement).Value)
			maxWordLen, _ := strconv.Atoi(doc.GetElementByID("maxWordLen").(*dom.HTMLInputElement).Value)
			maxWords, _ := strconv.Atoi(doc.GetElementByID("maxWords").(*dom.HTMLInputElement).Value)

			// Generate list of candidate and alternate words.
			cand := candidates(dictWords, phrase, minWordLen, maxWordLen)

			// Anagram & Print sorted by word (and optionally, by line.)
			var an []string
			an = anagrams(freqmap(&phrase), cand, an, 0, maxWords)

			sortWords(an)
			sort.Strings(an)
			d.SetInnerHTML("<pre>" + strings.Join(an, "\n") + "</pre>")

			//w.Style().SetProperty("display", "none", "")

		}(resultsDiv, waitDiv)
		println("I'm out of here...")
	})
}
