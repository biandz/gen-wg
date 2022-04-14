package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)
func main()  {
	defer func() {
		if re := recover();re!=nil{
			fmt.Println("记录错误日志：",re)
		}
	}()
	list := os.Args
	//项目路径
	projectPwd := strings.Replace(list[1],"--","",1)
	packageName := strings.Replace(list[2],"--","",1)

	//创建项目
	projectAddr :=createDir(projectPwd,packageName)

	//创建项目main.go文件
	file := createFile(projectAddr, "/main.go")
	genMain := NewGenerate(file)
	genMain.writeMainFile()

	//创建启动器
	bootAddr :=createDir(projectAddr,"bootstrap")
	bootFile := createFile(bootAddr, "/bootstrap.go")
	genBoot := NewGenerate(bootFile)
	genBoot.writeBootFile()
}

//根据所给文件夹路径创建文件夹
func createDir(projectPwd ,packageName string) string {
	dir := projectPwd+"/"+packageName
	//查找系统是否存在当前路径，不存在直接创建
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir,os.ModePerm); err != nil {
			panic(err.Error())
		}
	}
	return dir
}

//创建文件
func createFile(proAddr,fileName string) *os.File {
	dir:=proAddr+fileName
	if _, err := os.Stat(dir); err != nil {
		if _,err := os.Create(dir); err != nil {
			panic(err.Error())
		}
	}
	file, err := os.OpenFile(dir, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file
}

































type Generate struct {
	File *os.File
}

func NewGenerate(file *os.File) *Generate {
	return &Generate{
		File: file,
	}
}

func (g *Generate)p(str string)  {
	io.WriteString(g.File,str)
}

func (g *Generate)writeBootFile()  {
	//生成包名
	g.p("package bootstrap")
	g.p("\n")
	g.p("\n")
	//生成全局变量once
	g.p("var once sync.Once")
	g.p("\n")
	g.p("\n")
	//生成AppInit函数
	g.p(fmt.Sprintf("func AppInit() { "))
	g.p("\n\t")
	g.p(fmt.Sprintf("once.Do(func() {"))
	g.p("\n\t\t")

	g.p("//设置模式")
	g.p("\n\t\t")
	g.p("gin.SetMode(gin.ReleaseMode)")
	g.p("\n\t\t")
	g.p("r := gin.Default()")

	g.p("\n\t\t")
	g.p("//加载路由")
	g.p("\n\t\t")
	g.p("route.LoadRouter(r)")

	g.p("\n\t\t")
	g.p("//加载配置文件")
	g.p("\n\t\t")
	g.p("config.ParseConfig()")

	g.p("\n\t\t")
	g.p("//加载mysql和redis")
	g.p("\n\t\t")
	g.p("db.InitDB()")

	g.p("\n\t\t")
	g.p("//启动服务")
	g.p("\n\t\t")
	g.p("_ = r.Run(viper.GetString(\"server.port\"))")

	g.p("\n\t")
	g.p("})")
	g.p("\n")
	g.p("}")
}

func (g *Generate)writeMainFile()  {
	//生成包名
	g.p("package main")
	g.p("\n")
	g.p("\n")

	//生成main函数
	g.p(fmt.Sprintf("func main() { "))
	g.p("\n\t")
	g.p(fmt.Sprintf("bootstrap.Appinit()"))
	g.p("\n}")
}