package main

import (
	"net/http"
	"strings"
	"regexp"
	"sort"
	"flag"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/fatih/color"
)

func main() {
	parseflags()
}

// Parse the 3 flags: "-topic", "-filter", and "-sentences"
func parseflags() {
	topicPtr := flag.String("topic", "fart", "Wikipedia topic to summarize")
	filterPtr := flag.Bool("filter", false, "Use filter")
	sentencesPtr := flag.Int("sentences", 5, "Number of sentences to summarize in")
	flag.Parse()

	red := color.New(color.FgRed, color.Bold)
	if *topicPtr == "" {
		red.Println("Error. Page cannot be blank.") 
	} else if *sentencesPtr <= 0 {
		red.Println("Error. Must summarize in at least 1 sentence.")
	} else {
		wikisum(*topicPtr, *filterPtr, *sentencesPtr)
	}
}

// Combines it all under one function
func wikisum(topic string, filterpage bool, x int) {
	// Gotta get those colors
	blue := color.New(color.FgBlue, color.Bold)
	cyan := color.New(color.FgCyan, color.Bold)	
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	// Scrape the wikipage and remove citations
	content := removecitations(scrapewiki("https://en.wikipedia.org/wiki/"+topic))
	// Generate a word map
	wordmap := genwordmap(content, filterpage)
	// Generate a sentence map
	sentmap := gensentmap(content, wordmap)
	// Get top X sentences
	topsents := gettopsents(sentmap, x)
	// Range through and print them, alternating between cyan and blue
	for i, topsent := range topsents {
		yellow.Println()
		if i%2 == 0 {
			blue.Println(topsent)
		} else {
			cyan.Println(topsent)
		}
	}
	// Print size reduction
	yellow.Println()
	yellow.Println("|-------------------------------------------|")
	yellow.Printf("| Wikipedia page size reduced by ")
	green.Printf("%d",int((float64(1)-float64(len(strings.Join(topsents, "")))/float64(len(content)))*100))
	yellow.Printf(" percent |")
	yellow.Println()
	yellow.Println("|-------------------------------------------|")	
	yellow.Println()
}

// Gets the top X sentences from the sentmap
func gettopsents(sentmap []sm, x int) []string {
	// Create an arrray of all the frequencies of the given sentmap
	var frequencies []int 
	for _, s := range sentmap {
		frequencies = append(frequencies, s.Freq)
	}
	// Sort the array by greatest to least
	sort.Sort(sort.Reverse(sort.IntSlice(frequencies)))
	
	// Only keep first X frequencies
	frequencies = frequencies[:x]
	
	// Order the frequencies by their original chronological order of appearance
	var orderedfrequencies []int
	for _, s := range sentmap {
		for _, frequency := range frequencies {
			if s.Freq == frequency {
				orderedfrequencies = append(orderedfrequencies, frequency)
			}
		}
	}
	// Array of top sentences
	var topsents []string
	// Range through first X orderedfrequencies
	for i:=0; i<x; i++ {
		// Match the frequency to it's orignal sentence
		for _, s := range sentmap {
			if s.Freq == orderedfrequencies[i] {
				topsents = append(topsents, strings.TrimSpace(s.Sent+"."))
			}
		}
	}
	return topsents

}

// Struct for storing sentence maps 
type sm struct {
	Sent string 
	Freq int
}

// Generates a map of each sentence and the total frequency of the words that make it up
func gensentmap(page string, wordmap map[string]int) []sm {
	// Get the total frequency of the words that make up given sentence
	getsentvalue := func(sent string, wordmap map[string]int) int {
		totalpoints := 0
		// Split sentence into words
		wordarray := strings.Fields(sent)
		// Range through the words
		for _, word := range wordarray {
			// Range through all the words in the wordmap
			for wordm := range wordmap {
				// If the word from the sentence matches the word from the wordmap
				if word == wordm {
					// Increase totalpoints for that sentence by the word's frequency from the wormdpa
					totalpoints += wordmap[wordm]
				}
			}
		}
		return totalpoints
	}
	// Create map of sentences
	var sentmap []sm
	// Split the page into individual sentences
	sentarray := strings.Split(page, ".")
	// Range through the individual sentences
	for _, sent := range sentarray {
		sentmap = append(sentmap, sm{Sent: sent, Freq: getsentvalue(sent, wordmap)})
	}
	return sentmap
}


// Scrapes the wiki page and returns all the text it contains
func scrapewiki(url string) string {
	resp, _ := http.Get(url)
	root, _ := html.Parse(resp.Body)
	// Will only return text inside <p> elements
	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.P && n.Parent != nil && n.Parent.Parent != nil {
			return true
		}
		return false
	}
	var content []string
	pelements := scrape.FindAll(root, matcher)
	for _, pelement := range pelements {
		content = append(content, scrape.Text(pelement))
	}
	return strings.Join(content, "")
}

// Generates a map of each word and it's frequency
func genwordmap(page string, filterpage bool) map[string]int {
	// Convert page into array of words
	wordarray := strings.Fields(page)
	if filterpage == true {
		// Filter the wordarray
		wordarray = filter(wordarray)
	}
	// Checks if word is already in map. Returns bool.
	isalreadyinmap := func(w string, m map[string]int) bool {
		for mw, _ := range m {
			if w == mw {
				return true
			}
		}
		return false
	}
	// Range through  word array, and add the word to map with a frequency of 1. If it already in the map, add +1 to its frequency.
	wordmap := make(map[string]int)
	for _, word := range wordarray {
		if isalreadyinmap(word, wordmap) == true {
			wordmap[word] += 1
		} else {
			wordmap[word] = 1
		}
	
	}
	return wordmap
}

// Filters through given word array, removing common words, words shorter than 2 characters, or containing non-letter characters. Converts them to lowercase.
func filter(wordarray []string) []string {
	var filteredwordarray []string
 	// Checks if word is common, returns bool
	iscommonword := func(word string) bool {
		// Array of common words
 		commonwords := []string{"the","of","and","to","a","in","for","is","on","that","by","this","with","i","you","it","not","or","be","are","from","at","as","your","all","have","new","more","an","was","we","will","can","us","if","my","has","but","our","one","do","no","they","he","up","may","what","which","their","out","use","any","there","see","only","so","his","when","here","who","also","now","help","get","view","am","been","would","how","were","me","some","these","its","like","than","find"}
		for _, commonword := range commonwords {
			if word == commonword {
				return true
			}
		}
		return false
	}
	// Checks if character is letter, returns bool
	isletter := func(r rune) bool {
		return r < 'A' || r > 'z'
	}
	// Range through word array
	for _, word := range wordarray {
		// Check if word is at least longer than 1 character
		if len(word) > 1 {
			// Check if word contains any non-letter characters
			if strings.IndexFunc(word, isletter) == -1 {
				// Check if word is common
				if iscommonword(strings.ToLower(word)) == false {
					// Add to filteredwordarray in its lowercase form
					filteredwordarray = append(filteredwordarray, strings.ToLower(word))
				}
			}
		}
	}
	return filteredwordarray
}

// Removes any citations from given string
func removecitations(page string) string {
	rgx := regexp.MustCompile(`\[(.*?)\]`)
	return rgx.ReplaceAllString(page, "")
}