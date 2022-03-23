package main

import (
	"fmt"
	"wiki_search_go/modules"
	"wiki_search_go/modules/suffix_tree"
	"sort"
	"strings"
	"log"
	"encoding/json"
	badger "github.com/dgraph-io/badger/v3"
)

type Webpage struct{
	Text [] string
	Index []int
}

func (w Webpage) encodeWebpage() []byte {
	data, err := json.Marshal(w)
	if err != nil{
		panic(err)
	}
	return data
}

func deccodeWebpage(data []byte) (Webpage, error) {
	var w Webpage
	err := json.Unmarshal(data, &w)
	return w,err
}

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

	fmt.Println(len(raw_text),len(text))
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

	fmt.Println(output)

	// ################################################
	// ########  Storing Data in BadgerDB  ############
	// ################################################

	// url : {text []string, index []int}

	db, err := badger.Open(badger.DefaultOptions("/data/badger"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// ##############  inserting webpage Key #################
	// err = db.Update(func(txn *badger.Txn) error {
	// 	err := txn.Set([]byte(url), Webpage{
	// 		Text:text,
	// 		Index:index,
	// 	}.encodeWebpage())
	// 	return err
	//   })

	// if err != nil{
	// 	panic(err)
	// }
	
	// ##############  Iterating over all Keys #################
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
		  item := it.Item()
		  k := item.Key()
		  err := item.Value(func(v []byte) error {
			webpage , _ := deccodeWebpage(v)
			fmt.Printf("key=%s, value=%+v\n", k,webpage)
			
			return nil
		  })
		  if err != nil {
			return err
		  }
		}
		return nil
	  })

	if err != nil{
		panic(err)
	}
	  
}