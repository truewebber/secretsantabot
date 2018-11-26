package random

import (
	"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type (
	IntRandomStuck struct {
		data []int
	}
)

func IntInit(num int, min int, max int) (*IntRandomStuck, error) {
	intUrl := getUrl("/integers")
	q := intUrl.Query()
	q.Add("num", strconv.Itoa(num))
	q.Add("min", strconv.Itoa(min))
	q.Add("max", strconv.Itoa(max))
	q.Add("col", "1")
	q.Add("base", "10")
	q.Add("format", "plain")
	q.Add("rnd", "new")
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

	list := bytes.Split(body, []byte("\n"))
	out := new(IntRandomStuck)

	for _, elem := range list {
		if len(elem) == 0 {
			continue
		}

		i, err := strconv.Atoi(string(elem))
		if err != nil {
			return nil, err
		}

		out.data = append(out.data, i)
	}

	return out, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (s *IntRandomStuck) Get() int {
	if len(s.data) == 0 {
		return -999999
	}

	var r int
	r, s.data = s.data[0], s.data[1:]

	return r
}

// ---------------------------------------------------------------------------------------------------------------------

func getUrl(path string) url.URL {
	return url.URL{
		Scheme: "https",
		Host:   "www.random.org",
		Path:   path,
	}
}
