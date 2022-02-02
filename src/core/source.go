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
	resp, err := http.Get(url)
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
