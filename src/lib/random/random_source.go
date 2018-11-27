package random

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

func GetRandomBytes(num int) ([][]byte, error) {
	intUrl := getUrl("/cgi-bin/randbyte")
	q := intUrl.Query()
	q.Add("nbytes", strconv.Itoa(num))
	q.Add("format", "h")
	intUrl.RawQuery = q.Encode()

	resp, err := http.Get(intUrl.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, errors.Errorf("Randomize return non-200 code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	listWithColumns := bytes.Split(body, []byte("\n"))
	out := make([][]byte, 0)
	i := 0

	for _, listData := range listWithColumns {
		listData = bytes.TrimSpace(listData)
		if len(listData) == 0 {
			continue
		}

		list := bytes.Split(listData, []byte(" "))
		for _, elem := range list {
			if len(out) != i+1 {
				out = append(out, make([]byte, 0))
			}

			elem = bytes.TrimSpace(elem)
			if len(elem) == 0 {
				continue
			}

			b := make([]byte, 1)
			_, err := hex.Decode(b, elem)
			if err != nil {
				return nil, err
			}

			out[i] = append(out[i], b[0])

			if len(out[i]) == 8 {
				i++
			}
		}
	}

	return out, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func getUrl(path string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   "www.random.org",
		Path:   path,
	}
}
