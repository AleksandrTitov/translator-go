package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type wordToTr struct {
	Words []string
}

type translateYandex struct {
	Head struct {
	}
	Def []struct {
		Text string
		Pos  string
		Ts   string
		Tr   []struct {
			Text string
			Pos  string
			Syn  []struct {
				Text string
				Pos  string
				Gen  string
			}
			Mean []struct {
				Text string
			}
			Ex []struct {
				Text string
				Tr   []struct {
					Text string
				}
			}
			Gen string
		}
	}
}

const tokenYA = "TOKEN_YA"

func respHTTP(url, methodReq string, metadataHTTP map[string]string, dataHTTP []byte)([]byte, int64, string, int) {

	client := &http.Client{}
	httpReq, _ := http.NewRequest(methodReq, url, bytes.NewBuffer(dataHTTP))
	for key, value := range metadataHTTP {
		httpReq.Header.Add(key, value)
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println(err)
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer httpResp.Body.Close()

	return body, httpResp.ContentLength, httpResp.Header.Get("Date"), httpResp.StatusCode
}

func translate(words []string) map[string]string {

	var Translate translateYandex
	var translation string

	in := make(chan string, 1)
	out := make(chan string, 1)

	dictionary := make(map[string]string)

	go func(in, out chan string) {
		for word := range in {

			url := fmt.Sprintf(
				"https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-ru&text=%s&flags=4",
				os.Getenv(tokenYA),
				word)

			body, contentLen, date, respCode := respHTTP(url, "GET", nil, nil)
			err := json.Unmarshal(body, &Translate)
			if err != nil {
				fmt.Println(err)
			}
			if respCode == 200 {
				if contentLen != 20 {
					translation = Translate.Def[0].Tr[0].Text
				} else {
					translation = "none"
				}
			} else {
				translation = "none"
			}
			out <- translation

			fmt.Println(date,", status: [",respCode ,"], word: [",word, "], translate: [", translation, "]")
		}
		close(out)
	}(in, out)

	for _, word := range words {
		in <- word
		transl := <-out

		dictionary[word] = transl
	}
	close(in)

	return dictionary
}

func parsePost(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var words wordToTr
	err := decoder.Decode(&words)

	if err != nil {
		panic(err)
	}

	wordsJS, err := json.Marshal(translate(words.Words))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Server", "A Go Web Server")
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(wordsJS)
}

func main() {
	http.HandleFunc("/", parsePost)
	http.ListenAndServe(":8080", nil)
}