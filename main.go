package main

import (
	"fmt"
	"wiki_search_go/modules"
	"wiki_search_go/modules/suffix_tree"
	"sort"
	"strings"
)

func main(){
	fmt.Println("\nEnter Wikipedia URL to index")
	fmt.Println("\nPress 1 if url = https://en.wikipedia.org/wiki/Synthetic_diamond;\nelse: Enter Url ")
	
	var url string
	fmt.Scan(&url)
	if url == "1"{
			url = "https://en.wikipedia.org/wiki/Synthetic_diamond"
	}

	fmt.Println("\nFetching Data....")
	raw_text := helper.FetchData(url)
	raw_text = strings.ReplaceAll(raw_text, "\n", "")
 	text := strings.Split(raw_text,".")

	fmt.Println(raw_text,text)
	fmt.Println("Data Fetched \nBuilding Index...")
	var index [] int

	for _,line := range(text){
			padding := 0
			if len(index) > 0{
				padding = index[len(index)-1]
			}
			index = append(index,len(line)+int(padding))
		}
	fmt.Println(index)

	fmt.Println("Building Suffix Tree...")
	st := suffixtree.NewGeneralizedSuffixTree()
	for k, word := range text {
		st.Put(word, k)
	}

	fmt.Println("\nTree construction completed")
	fmt.Println("Enter pattern to search")
	
	var pattern string
	fmt.Scan(&pattern)
	// pattern = "diamonds"

	var occurences [] int
	occurences = st.Search(pattern, -1)
	sort.Ints(occurences)

	fmt.Println("Index of Occurences", occurences, len(occurences))

	var output [] string
	for _,occurence := range(occurences){

			pos := helper.LowerBound(index, occurence)
			if pos >= 0 && pos < len(text){
					output = append(output,text[pos])
				}
		}

	fmt.Println(len(output))
	fmt.Println(output)

}