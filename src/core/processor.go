package core

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/tomhjx/gogfw"
)

const (
	VERSION = "1.0.4"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

// https://github.com/Dreamacro/clash/wiki/configuration
type Proxy struct {
	Name           string
	Type           string
	Server         string
	Port           uint16
	Password       string `yaml:"password,omitempty"`
	Sni            string `yaml:"sni,omitempty"`
	AlterId        uint8  `yaml:"alterId"`
	Cipher         string `yaml:"cipher,omitempty"`
	Network        string `yaml:"network,omitempty"`
	Tls            bool
	SkipCertVerify bool   `yaml:"skip-cert-verify"`
	Uuid           string `yaml:"uuid,omitempty"`
	Obfs           string `yaml:"obfs,omitempty"`
	Protocol       string `yaml:"protocol,omitempty"`
	WsHeaders      struct {
		Host string `yaml:"Host,omitempty"`
	} `yaml:"ws-headers,omitempty"`
	WsPath     string `yaml:"ws-path,omitempty"`
	Plugin     string `yaml:"plugin,omitempty"`
	PluginOpts struct {
		Mode string `yaml:"mode,omitempty"`
		Host string `yaml:"host,omitempty"`
	} `yaml:"plugin-opts,omitempty"`
}

type Proxyg struct {
	Name     string
	Type     string
	Url      string
	Interval uint16
	Proxies  []string
}

type stringsFlag []string

func (i *stringsFlag) String() string {
	return strings.Join([]string(*i), ",")
}

func (i *stringsFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *stringsFlag) Get() []string {
	return []string(*i)
}

func addGFWRules(i *target) (res bool, err error) {
	gfwh, err := gogfw.OpenOnline("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	if err != nil {
		return false, err
	}
	gfwd, err := gfwh.ReadItems()
	if err != nil {
		return false, err
	}
	for _, v := range gfwd {
		rtype := "DOMAIN"
		rval := v.Value
		switch v.Type {
		case gogfw.ITEM_TYPE_DOMAIN, gogfw.ITEM_TYPE_IP:
			rtype = "DOMAIN"
		case gogfw.ITEM_TYPE_DOMAIN_SUFFIX:
			rtype = "DOMAIN-SUFFIX"
		case gogfw.ITEM_TYPE_DOMAIN_KEYWORD:
			rtype = "DOMAIN-KEYWORD"
		}
		i.prePersistRule([]string{rtype, rval, "PROXY"})
	}
	return true, nil
}

func addProxies(t *target, srcp string) (res bool, err error) {

	s, err := OpenOnlineSource(srcp)
	if err != nil {
		log.Println(err)
		return false, err
	}
	proxies, err := s.Proxies()
	if err != nil {
		log.Println(err)
		return false, err
	}
	for _, p := range proxies {
		t.prePersistProxy(p)
	}

	return true, nil
}

func (i *Processor) Run() {
	var (
		srcps    stringsFlag
		outp     string
		tplp     string
		showHelp bool
		wg       sync.WaitGroup
	)
	flag.Var(&srcps, "s", "source's clashx configuration yaml url.")
	flag.StringVar(&outp, "o", "/work/out/clashx.yaml", "output clashx configuration yaml file path.")
	flag.StringVar(&tplp, "tpl", "/work/resources/tpl.yaml", "templet clashx configuration yaml file path.")
	flag.BoolVar(&showHelp, "h", false, "this help")

	flag.Parse()

	if showHelp {
		fmt.Printf("version %s\n", VERSION)
		flag.Usage()
		return
	}

	t, err := newTarget(tplp)
	if err != nil {
		log.Panic(err)
	}

	for _, srcp := range srcps {
		wg.Add(1)
		go func(t *target, srcp string) {
			defer wg.Done()
			addProxies(t, srcp)
		}(t, srcp)
	}

	wg.Add(1)
	go func(t *target) {
		defer wg.Done()
		addGFWRules(t)
	}(t)

	go func() {
		wg.Wait()
		t.finishNotifyPersist()
	}()

	t.consumePersistQ(outp)
}
