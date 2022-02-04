# c(Configuration) clashc(ClashX)
Aggregate configuration for ClashX

[![GitHub release](https://img.shields.io/github/release/tomhjx/cclashx.svg)](https://github.com/tomhjx/cclashx/releases)
[![GitHub license](https://img.shields.io/github/license/tomhjx/cclashx.svg)](https://github.com/tomhjx/cclashx/blob/master/LICENSE)


## Example

all proxies source:

http://example01.com/clashx.yaml

https://example02.com/clashx.yaml

my configuration template:

/your/template/clashx.yaml


aggregate a new ClashX configuration:

/your/output/dirpath/clashx.yaml


```bash

cclashx -s "http://example01.com/clashx.yaml" -s "https://example02.com/clashx.yaml" -o "/your/output/dirpath/clashx.yaml" -tpl "/your/template/clashx.yaml"

```


## Help

```bash

cclashx -h

```