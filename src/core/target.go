package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type target struct {
	Port               uint16
	SocksPort          uint16 `yaml:"socks-port"`
	AllowLan           bool   `yaml:"allow-lan"`
	Mode               string
	LogLevel           string `yaml:"log-level"`
	ExternalController string `yaml:"external-controller"`
	Dns                struct {
		Enable bool
	}
	Proxies     []*Proxy
	ProxyGroups []*Proxyg `yaml:"proxy-groups"`
	Rules       []string
	proxyExists map[string]bool
}

func newTarget(def string) (t *target, err error) {
	tb, err := ioutil.ReadFile(def)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(tb, &t)
	if err != nil {
		return nil, err
	}

	t.proxyExists = make(map[string]bool)

	return t, nil
}

func (i *target) addProxy(p *Proxy) {
	_, err := net.ResolveIPAddr("ip", p.Server)
	if err == nil {
		return
	}

	if p == nil {
		return
	}
	pxykey := fmt.Sprintf("%s:%d", p.Server, p.Port)
	log.Printf("add proxy: %s, %s, %d", p.Name, p.Server, p.Port)
	if i.proxyExists[pxykey] {
		return
	}
	i.proxyExists[pxykey] = true
	i.Proxies = append(i.Proxies, p)
}

func (i *target) addRule(rule []string) {
	i.Rules = append(i.Rules, strings.Join(rule, ","))
}

func (i *target) persist(dstpath string) (res bool, err error) {
	if len(i.Proxies) == 0 {
		return false, errors.New("target.proxies is empty")
	}

	if len(i.Rules) == 0 {
		return false, errors.New("target.rules is empty")
	}

	// define proxy group
	autopg := &Proxyg{
		Name:     "Auto",
		Type:     "url-test",
		Url:      "http://www.gstatic.com/generate_204",
		Interval: 300,
	}
	cpg := &Proxyg{Name: "PROXY", Type: "select"}
	cpg.Proxies = append(cpg.Proxies, autopg.Name)

	hijackingpg := &Proxyg{
		Name:    "Hijacking",
		Type:    "select",
		Proxies: []string{"DIRECT", "REJECT", cpg.Name},
	}

	for _, proxy := range i.Proxies {
		autopg.Proxies = append(autopg.Proxies, proxy.Name)
		cpg.Proxies = append(cpg.Proxies, proxy.Name)
	}
	i.ProxyGroups = append(i.ProxyGroups, autopg, cpg, hijackingpg)

	i.addRule([]string{"MATCH", "DIRECT"})

	writer, err := os.OpenFile(dstpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return false, err
	}

	encoder := yaml.NewEncoder(writer)
	err = encoder.Encode(i)
	if err != nil {
		return false, err
	}
	encoder.Close()
	return true, nil
}
