package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var Parents_dir string
var Same_file_dir	string
var source_path string
var file_list []string
var same_file_list []string
//var file_md5 = make(map[string]string)

func init()  {
	s := "D:\\gocode\\src\\github.com\\yangyouwei\\quchong\\test"
	source_path, _ = filepath.Abs(s)
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
	fn := getfilelist(source_path)
	file_list = *fn
	fmt.Println(file_list)

	sn := diff_md5(file_list)
	same_file_list = *sn
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
		a := md5_sum(k)

		if v, ok := f[a]; ok {
			fmt.Println("true")
			s = append(s,v)
		}else{
			fmt.Println("false")
			f[a] = k
		}
	}
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
		fmt.Println(string(sub_path))

		//需要在存放相同文件的目录里创建文件相应的层级目录
		d_patch := same_file_path + string(sub_path)
		fmt.Println(d_patch)

		err := os.MkdirAll(d_patch, os.ModePerm)
		if err != nil {
			log.Println(err)
			break
		}

		//移动文件
		err = os.Rename(v, d_patch+"\\"+filename)
		if err != nil {
			fmt.Println("mv")
			log.Println(err)
			break
		}
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
