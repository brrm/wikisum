# Wikisum
*Summarize Wikipedia pages from the command line*
![wikisum gif](https://i.imgur.com/mtm1fVT.gif)

# About
Wikisum is a golang replication of [SMMRY](https://smmry.com/about). The Wikipedia page is summarized in 6 simple steps:

**1)** Calculates the occurrence of each word on the page.

**2)** Assigns each word with points depending on their occurence (if a word occurs 100 times, it will have a point value of 100).

**3)** Splits the text up into individual sentences the occurence of the words 

**4)** Assigns each sentence with points (the sum of their words' points).

**5)** Ranks the sentences by their points.

**6)** Returns X of the most highly ranked sentences in their chronological order.

# Usage
**1)** [Download the binary](https://github.com/brrm/wikisum/releases/download/v0.1/wikisum) 

**2)** Enter the directory where you placed the binary

**3)** Use `./wikisum -topic [your topic of choice]` to begin

`-filter` Can be used to stop common words like "the" from receiving points

`-sentences [sentence number]` Can be used to specify the length of the summary

All flags are optional. `-sentences` defaults to `5`, `-filter` is disabled by default (`false`), and `-topic` defaults to `fart`.  

# Todo

* Support for converting plurals to singular (ex: cats --> cat) in order to increase summary accuracy
* User Interface
* Handle invalid pages