package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var path, exin, exout, prefix string

func init() {
	rand.Seed(time.Now().UnixNano())
	prefix = strconv.Itoa(rand.Int())
}

func main() {
	if len(os.Args) != 4 {
		// 参数数量不对
		fmt.Println("参数个数不对")
		return
	}
	path = os.Args[1]
	exin = os.Args[2]
	exout = os.Args[3]

	fmt.Println("程序开始")
	fmt.Println("--安全检查")
	if strings.Count(path, "/") < 3 {
		fmt.Println("----不允许处理3层以内的目录")
		return
	}
	fmt.Println("--安全检查结束")
	fmt.Println("--解开所有目录结构")
	err := fetchFileAndRemoveFolder(path)
	if err != nil {
		fmt.Println(err)
		fmt.Println("----有异常，中断")
		return
	}

	// 文件正则替换
	fmt.Println("--文件名处理")
	cmd0 := exec.Command("/bin/sh", "-c", "ls -l "+path+"|grep \"^\\-\"|awk '{print $9}'")
	output, err := cmd0.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	}

	files := strings.Split(string(output), "\n")

	if len(files) < 2 {
		// 没有文件
		return
	}

	files = files[:len(files)-1]

	for _, file := range files {
		// 重命名
		re := regexp.MustCompile(exin)
		newName := re.ReplaceAllString(file, exout)

		err := os.Rename(path+"/"+file, path+"/"+newName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fmt.Println("--文件名处理结束")
}

func fetchFileAndRemoveFolder(folderPath string) error {
	fmt.Println("----解开 " + folderPath)
	cmd0 := exec.Command("/bin/sh", "-c", "ls -l "+folderPath+"|grep \"^d\"|awk '{print $9}'")
	output, err := cmd0.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return err
	}

	folders := strings.Split(string(output), "\n")

	if len(folders) < 2 {
		// 没有目录
		return nil
	}

	folders = folders[:len(folders)-1]
	for i, folder := range folders {
		// 遍历子目录
		// 改名，避免重名
		newFolderName := "R" + prefix + "Tmp00" + strconv.Itoa(i)
		newPath := folderPath + "/" + newFolderName
		err := os.Rename(folderPath+"/"+folder, newPath)
		if err != nil {
			return err
		}

		// 解开子目录
		err = fetchFileAndRemoveFolder(newPath)
		if err != nil {
			return err
		}

		// 解开当前目录
		err = fetchFile(newPath, folderPath)
		if err != nil {
			return err
		}

		// 删除当前目录
		os.RemoveAll(newPath)
	}
	return nil
}

func fetchFile(folderPath, distPath string) error {
	fmt.Println("----提取 " + folderPath)
	cmd0 := exec.Command("/bin/sh", "-c", "ls -l "+folderPath+"|grep \"^\\-\"|awk '{print $9}'")
	output, err := cmd0.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return err
	}

	files := strings.Split(string(output), "\n")

	if len(files) < 2 {
		// 没有文件
		return nil
	}

	files = files[:len(files)-1]

	for _, file := range files {
		// 把文件挪到上级目录
		err := os.Rename(folderPath+"/"+file, distPath+"/"+file)
		if err != nil {
			return err
		}
	}
	return nil
}
