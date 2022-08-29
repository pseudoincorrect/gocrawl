package fetch

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Fetch(url string) string {
	u := url
	if !strings.Contains(url, "http") {
		u = "http://" + url
	}
	resp, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		fmt.Println("could not fetch on url :", url)
	}
	return buf.String()
}
