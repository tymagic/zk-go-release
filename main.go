package main

import (
	"bufio"
	"io"
	"net/http"

	//"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"regexp"
	"strings"
	"time"
)
// 查询，输入节点路径
//if  exists{
//	if 最终节点，yml解析，一条条输出所有配置项
//		if  配置项==password    value = ****
//	else
//	   输出该节点所有子节点信息，并提示不是最终路径
// }
var infos map[string]interface{}
var block []string

//写入临时文件函数
func write_tmp(fileName,strs string){
	dstFile,_ := os.Create(fileName)
	defer dstFile.Close()
	dstFile.WriteString(strs + "\n")
}
func Serach_content(value string) ([]string,error){   //内容查询
	//字符串格式化空格变为换行。password后面=号后隐藏
	//写入临时文件
	var result_array []string
	write_tmp("tmp.txt", value)
	//逐行读取文件
	fi, _ := os.Open("tmp.txt")
	defer os.Remove("tmp.txt")
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		val := string(a)
		//遍历数组替换字符串,屏蔽password,屏蔽其他字符直接
		for _, blocks := range block {
			if strings.Contains(val, blocks) {
				rule := blocks + ".+"
				reg := regexp.MustCompile(rule)
				result := blocks + "=******"
				val = reg.ReplaceAllString(val, result)
			}
			result_array = append(result_array, val)
		}

	}
	//遍历数组替换字符串,屏蔽password,屏蔽其他字符直接
	return result_array,nil
}

func Search(uri string) (map[string]interface{},error,bool) {
	tmp_str:=string(uri[len(uri)-1])
	if tmp_str == "/"{
		uri=string(uri[0:len(uri)-1])
	}
	//需要屏蔽的key写在block数组中
	block = []string{"unipay_server_callback_url", "unipay_pay_url"}
	//map 需要在函数里初始化，如下
	infos := make(map[string]interface{})
	var hosts = []string{"192.168.1.34:9090"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return nil,err,false
	}
	defer conn.Close()
	var path = uri
	//子节点查询
	children, _, err := conn.Children(path)
	//如果存在子节点
	fmt.Println("子节点：：",len(children))
	if len(children) >0{
		var child_array []string
		if err != nil {
			fmt.Println(err)
			return nil,err,false
		}
		fmt.Printf("root_path[%s] child_count[%d]\n", path, len(children))
		for _, ch := range children {
			//
			// fmt.Printf("%d, %s\n", __, ch)
			child_array = append(child_array,ch )
		}
		infos["values"] = child_array
	} else{//如果不存在子节点
	fmt.Println("不存在子节点")
	value, _, err := conn.Get(path)
	//fmt.Println(value)
	if err != nil {
		fmt.Println(err)
		return nil,err,false
	}
	if value == nil {
		infos["values"] = "无数据"
	} else {
		value:=string(value)
		result,err:=Serach_content(value)//查询内容
		if err != nil {
			fmt.Println(err)
			return nil,err,false
		}
		infos["values"]=result
		}
	}
	return infos,nil,true

}
func main() {
	router := gin.Default()
	router.LoadHTMLFiles("templates/index.html")
	router.StaticFS("Stastic", http.Dir("stastic"))
	router.POST("/zkgo", func(c *gin.Context) {
		var uri string
		uri = c.Request.FormValue("location")
		fmt.Println("uri: ",uri)
		result,_,msg:=Search(uri)
		//fmt.Println(result["values"])
		c.JSON(200, gin.H{
			"OK":msg,
			"message": result ,
			"location": c.Request.FormValue("location"),
		})
	})
	router.GET("/index", func(c *gin.Context){
		c.HTML(http.StatusOK,"index.tmpl",gin.H{
		})
	})

	router.Run(":8001") // listen and serve on 0.0.0.0:8080
}