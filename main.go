package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type worldToTr struct {
	Worlds []string
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

func respHTTP(url, methodReq string, metadataHTTP map[string]string, dataHTTP []byte) []byte {

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

	return body
}

func translate(worlds []string) map[string]string {

	var Translate translateYandex

	in := make(chan string, 1)
	out := make(chan string, 1)

	dictionary := make(map[string]string)

	go func(in, out chan string) {
		for world := range in {

			url := fmt.Sprintf(
				"https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=en-ru&text=%s&flags=4",
				os.Getenv(tokenYA),
				world)

			body := respHTTP(url, "GET", nil, nil)
			err := json.Unmarshal(body, &Translate)
			if err != nil {
				fmt.Println(err)
			}
			out <- Translate.Def[0].Tr[0].Text
		}
		close(out)
	}(in, out)

	for _, world := range worlds {
		in <- world
		transl := <-out

		dictionary[world] = transl
	}
	close(in)

	return dictionary
}

func main() {
	http.HandleFunc("/", parsePost)
	http.ListenAndServe(":8080", nil)

}

func parsePost(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var worlds worldToTr
	err := decoder.Decode(&worlds)

	if err != nil {
		panic(err)
	}

	worldsJS, err := json.Marshal(translate(worlds.Worlds))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Server", "A Go Web Server")
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(worldsJS)
}
