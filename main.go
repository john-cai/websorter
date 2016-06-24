package main

import(
	"net/http"

	"github.com/john-cai/websorter/sorter"

)

func main() {
	http.HandleFunc("/sort", sorter.SortArray)
	http.ListenAndServe(":8080", nil)
}
