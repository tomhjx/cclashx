package core

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v3"
)

type source struct {
	Proxies []*Proxy
}

type Source struct {
	data []byte
}

func OpenOfflineSource(file string) (s *Source, err error) {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &Source{data: d}, nil
}

func OpenOnlineSource(url string) (s *Source, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36 cclashx"+VERSION)
	httpc := &http.Client{}
	resp, err := httpc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("online request returned code: %d", resp.StatusCode)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Source{data: d}, nil
}

func (i *Source) Proxies() (proxies []*Proxy, err error) {
	s := source{}
	err = yaml.Unmarshal(i.data, &s)
	if err != nil {
		return nil, err
	}
	return s.Proxies, nil
}
