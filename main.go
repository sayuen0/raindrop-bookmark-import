package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"net/http"
	"os"
	"time"
)

type Raindrop struct {
	Link  string   `json:"link"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

const collectionID = 0

func main() {
	f, err := os.Open("hoge.csv")
	panicErr(err)
	defer f.Close()
	csvFile := csv.NewReader(f)
	// ヘッダを捨てる
	{
		_, err := csvFile.Read()
		panicErr(err)
	}

	items := make([]*Raindrop, 0)
	for {
		record, err := csvFile.Read()
		if err == io.EOF {
			break
		}
		panicErr(err)
		link, title := record[5], record[6]
		items = append(items, &Raindrop{Link: link, Title: title, Tags: []string{"fromTwitter"}})
	}

	progressBar := pb.StartNew(len(items))
	for _, itm := range items {
		panicErr(postRaindrop(itm))
		// 秒間120リクエストまでなので、一応時間を取る
		time.Sleep(500 * time.Millisecond)
		progressBar.Increment()
	}
	progressBar.Finish()

}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func postRaindrop(rdp *Raindrop) error {
	body, _ := json.Marshal(rdp)
	req, _ := http.NewRequest(http.MethodPost, "https://api.raindrop.io/rest/v1/raindrop", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer API_TOKEN")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBody := string(bodyBytes)
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("status diff, %v", respBody))
	}
	return nil
}
