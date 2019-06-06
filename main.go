package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//目标目录的父目录
var Parents_dir string

//根据目标目录的路径，在目标目录的父目录下创建重复文件目录
var Same_file_dir	string

//指定要去重的目录
var source_path string

//遍历目标目录的所有文件，存到slice中
var file_list []string

//经过计算对比md5值，将重复文件目录放到slice中
var same_file_list []string

//根据参数初始化变量
//父目录及重复文件目录创建
func init()  {
	s := flag.String("f","","-f C:\\Users\\yyw\\go\\src\\github.com\\yangyouwei\\quchong\\test 访问url")
	flag.Parse()

	if *s == "" {
		flag.Usage()
		panic("process exsit!")
	}
	source_path, _ = filepath.Abs(*s)
	a := strings.LastIndex(source_path, "\\") //如果是linux系统使用 / 为分隔符,wondows 使用 \\
	rs := []rune(source_path)
	parents_dir := rs[:a]
	Parents_dir = string(parents_dir)
	dir_name := rs[a:]
	Same_file_dir = Parents_dir + string(dir_name) + "-samefile"
	err := os.Mkdir(Same_file_dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}


func main()  {
	//遍历目录
	fn := getfilelist(source_path)
	file_list = *fn

	//比较md5值，重复的放到重复的slice中，不重复的放到map中
	sn := diff_md5(file_list)
	same_file_list = *sn

	//读取重复文件的slice，移动文件到指定目录
	move_file(same_file_list,Parents_dir,Same_file_dir)
}

func getfilelist(source_path string) *[]string {
	var s []string
	s, err := GetAllFile(source_path, s)
	if err != nil {
		panic(err)
	}

	return &s
}

func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func diff_md5(fl []string) *[]string{
	f := make(map[string]string)
	var s []string

	for _,k := range fl {
		//计算md5值
		a := md5_sum(k)
		if _, ok := f[a]; ok {
			//fmt.Println("true")
			//重复的吸入重复slice中
			s = append(s,k)
		}else{
			//fmt.Println("false")
			//不重复的加入map中
			f[a] = k
		}
	}
	//fmt.Println("map : ", f)
	//fmt.Println("slice :",s)
	//fmt.Println("file_list : ",fl)
	//fmt.Println("same file list : ",s)
	return &s
}

func move_file(same_file []string, parents_dir string, same_file_path string) {
	for _, v := range same_file {
		filename := path.Base(v)
		//创建目录
		rs := []rune(parents_dir)
		n := len(rs)
		p1 := []rune(path.Dir(v))
		//取出目标目录文件的子目录
		sub_path := p1[n:]

		//需要在存放相同文件的目录里创建文件相应的层级目录
		d_patch := same_file_path + string(sub_path)

		//创建文件原来的目录层级结构
		err := os.MkdirAll(d_patch, os.ModePerm)
		if err != nil {
			log.Println(err)
			break
		}

		//移动文件
		err = os.Rename(v, d_patch+"\\"+filename)
		if err != nil {
			log.Println(err)
			break
		}
		//打印移动结果
		fmt.Println(v + "move to " + d_patch)
	}
}

func md5_sum(file_path string) string {
	f, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Copy", err)
		return ""
	}
	md5hash.Sum(nil)
	a := fmt.Sprint(md5hash.Sum(nil))
	return string(a)
}
