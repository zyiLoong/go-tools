package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// projectPath 本地java项目的路径
const projectPath = "/path/to/project"

// commonPackagePrefix java工程中所有类的公共包前缀
const commonPackagePrefix = "com.example.demo"

// clearUnused 是否要删除掉没有被引用的java文件
const clearUnused = false
const printDetail = true

// main 查询一个java项目中，指定类所依赖的其他本项目类。
// 删除没有被引用的类
func main() {

	// root，从这些类开始向下查询
	rootClassNameList := []string{
		"com.example.demo.Test",
		"com.example.demo.TestUtils",
	}
	javaFilePathList := listJavaFiles()
	allClassMap := readClassMap(javaFilePathList)

	usedClassMap := make(map[string]bool)

	for _, rootClass := range rootClassNameList {
		findUsedClass(rootClass, allClassMap, &usedClassMap)
	}
	println("used class size: ", len(usedClassMap))
	if printDetail {
		fmt.Println(usedClassMap)
	}
	println("all class size: ", len(allClassMap))
	if printDetail {
		fmt.Println(allClassMap)
	}

	unUsedMap := CopyMap(allClassMap)

	for k, _ := range usedClassMap {
		delete(unUsedMap, k)
	}
	println("unused class size:", len(unUsedMap))
	if printDetail {
		fmt.Println(unUsedMap)
	}

	if clearUnused {
		for k, path := range unUsedMap {
			println(k)
			println(path)
			os.Remove(path)
		}
	}
}

// findUsedClass 找出被root Class依赖的class
// todo 会漏掉引用的同包类
func findUsedClass(class string, classMap map[string]string, resultMap *map[string]bool) {
	path := classMap[class]
	f, _ := os.Open(path)
	scanner := bufio.NewScanner(f)
	(*resultMap)[class] = true

	var classList []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "import") {
			class := getClassFromImport(line)
			// 确认是工程内定义的class
			// 确认这个类没有被添加过，避免循环导入
			if strings.Contains(class, commonPackagePrefix) && (*resultMap)[class] != true {
				classList = append(classList, class)
			}
		} else if strings.Contains(line, " class ") {
			break
		}
	}
	for _, class := range classList {
		findUsedClass(class, classMap, resultMap)
	}
}

func getClassFromImport(importStr string) string {
	str := strings.TrimPrefix(importStr, "import ")
	str = strings.TrimSuffix(str, ";")
	if strings.Contains(str, " ") {
		panic(str)
	}
	return str
}

// readClassMap 读取指定路径的java文件，然后以class全路径作为key，java文件路径作为value，放入map中返回
// key is the full path of class, eg: java.lang.String
// value is the full path of file.
func readClassMap(pathList []string) map[string]string {
	m := make(map[string]string)
	for _, path := range pathList {
		file, _ := os.Open(path)

		sc := bufio.NewScanner(file)
		sc.Scan()
		pkg := strings.TrimPrefix(sc.Text(), "package ")
		pkg = strings.TrimSuffix(pkg, ";")

		stat, _ := file.Stat()
		classFileName := stat.Name()
		className := strings.TrimSuffix(classFileName, ".java")

		class := pkg + "." + className
		m[class] = path
		file.Close()
	}
	return m
}

// listJavaFiles 获取项目中的java文件列表
func listJavaFiles() []string {
	var fileList []string
	filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".java") {
			fileList = append(fileList, path)
		}
		return nil
	})
	return fileList
}
