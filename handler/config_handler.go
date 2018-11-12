package handler

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//解析yaml 文件样例

type Conf struct {
	Address string     //属性首字母必须大写
	Filter []string
}

func (c *Conf)  getConf() *Conf{
	data, err := ioutil.ReadFile("config/cfg.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	yaml.Unmarshal(data, &c)
	d, _ := yaml.Marshal(&c)
	fmt.Println("看看 :", string(d))
	return c
}

func main(){
	var c Conf
	a:=c.getConf()
	for i:=0;i<len(a.Filter);i++{
		fmt.Println(a.Filter[i])
	}
}





