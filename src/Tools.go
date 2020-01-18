package src

import (
	"fmt"
	guuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DoEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func GetHash() string {
	id, _ := guuid.NewV4()
	return id.String()
}

// OptionsAnswer create options answer for browser
func OptionsAnswer(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
}

// FillAnswerHeader add some important headers to answer
func FillAnswerHeader(w http.ResponseWriter) {
	w.Header().Set("content-type", "application/json")
}

var Counter int
var BadCounter int
var GoodCounter int
var AllFiles int
var PaxHeaders int

func visit(path string, f os.FileInfo, err error) error {
	if !strings.Contains(path, "PaxHeaders") {
		AllFiles++
		if strings.Contains(path, "result") {
			Counter++
			buf, _ := ioutil.ReadFile(path)
			if string(buf) != "[]" {
				//fmt.Println(string(buf))
				GoodCounter++
			} else {
				BadCounter++
			}
		}

	} else {
		PaxHeaders++
		//os.Remove(path)
	}
	return nil
}

func RunScan(w http.ResponseWriter, r *http.Request) {
	//flag.Parse()

	Counter = 0
	GoodCounter = 0
	BadCounter = 0
	AllFiles = 0
	PaxHeaders = 0

	st := time.Now().Unix()
	fmt.Println(time.Now())
	root := r.URL.Query().Get("path")
	//go func() {
	filepath.Walk(root, visit)
	//}()
	//fmt.Printf("filepath.Walk() returned %v\n", err)
	fmt.Println("Total Count ", Counter)
	fmt.Println("Good One ", GoodCounter)
	fmt.Println("Bad ", BadCounter)
	fmt.Println("PaxHeaders ", PaxHeaders)

	fmt.Println(time.Now().Unix() - st)
}
