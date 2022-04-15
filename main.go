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
	genMain := NewGenerate(file,packageName)
	genMain.writeMainFile()

	//创建config.json文件
	cFile := createFile(projectAddr, "/config.json")
	genC := NewGenerate(cFile,packageName)
	genC.writeJsonFile()

	//创建启动器
	bootAddr :=createDir(projectAddr,"bootstrap")
	bootFile := createFile(bootAddr, "/bootstrap.go")
	genBoot := NewGenerate(bootFile,packageName)
	genBoot.writeBootFile()

	//创建配置文件
	configAddr := createDir(projectAddr,"config")
	configFile := createFile(configAddr,"/config.go")
	genConfig := NewGenerate(configFile,packageName)
	genConfig.writeConfigFile()

	//创建db文件（包含redis及mysql）
	dbAddr :=createDir(projectAddr,"db")
	dbFile :=createFile(dbAddr,"/db.go")
	dbConfig := NewGenerate(dbFile,packageName)
	dbConfig.writeDBFile()

	//创建route文件
	routeAddr :=createDir(projectAddr,"route")
	routeFile :=createFile(routeAddr,"/route.go")
	routeConfig := NewGenerate(routeFile,packageName)
	routeConfig.writeRouteFile()

	//创建一个默认的控制器
	controllerAddr :=createDir(projectAddr,"controller")
	//默认的base控制器
	baseControllerFile :=createFile(controllerAddr,"/base.go")
	baseControllerConfig := NewGenerate(baseControllerFile,packageName)
	baseControllerConfig.writeBaseControllerFile()
	//默认的index控制器
	indexControllerFile :=createFile(controllerAddr,"/index.go")
	indexControllerConfig := NewGenerate(indexControllerFile,packageName)
	indexControllerConfig.writeIndexControllerFile()

	//创建model
	modelAddr :=createDir(projectAddr,"model")
	modelFile :=createFile(modelAddr,"/baseModel.go")
	modelConfig := NewGenerate(modelFile,packageName)
	modelConfig.writeModelFile()

	//创建service
	serviceAddr :=createDir(projectAddr,"service")
	serviceFile :=createFile(serviceAddr,"/service.go")
	serviceConfig := NewGenerate(serviceFile,packageName)
	serviceConfig.writeServiceFile()

	//创建dao
	daoAddr :=createDir(projectAddr,"dao")
	daoFile :=createFile(daoAddr,"/dao.go")
	daoConfig := NewGenerate(daoFile,packageName)
	daoConfig.writeDaoFile()


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
	PName string
}

func NewGenerate(file *os.File,packageName string) *Generate {
	return &Generate{
		File: file,
		PName:packageName,
	}
}

func (g *Generate)p(str string)  {
	io.WriteString(g.File,str)
}

//bootstrap文件代码生成
func (g *Generate)writeBootFile()  {
	//生成包名
	g.p("package bootstrap")
	g.p("\n")
	g.p("\n")
	//生成引入文件
	g.p("import (")
	g.p("\n\t")
	g.p(fmt.Sprintf("\"%s/config\"",g.PName))
	g.p("\n\t")
	g.p(fmt.Sprintf("\"%s/db\"",g.PName))
	g.p("\n\t")
	g.p(fmt.Sprintf("\"%s/route\"",g.PName))
	g.p("\n\t")
	g.p(fmt.Sprintf("\"github.com/gin-gonic/gin\""))
	g.p("\n\t")
	g.p(fmt.Sprintf("\"github.com/spf13/viper\""))
	g.p("\n\t")
	g.p(fmt.Sprintf("\"sync\""))
	g.p("\n")
	g.p(")")
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
//main文件代码生成
func (g *Generate)writeMainFile()  {
	//生成包名
	g.p("package main")
	g.p("\n")
	g.p("\n")

	//生成main函数
	g.p(fmt.Sprintf("func main() { "))
	g.p("\n\t")
	g.p(fmt.Sprintf("bootstrap.AppInit()"))
	g.p("\n}")
}
//生成config.json文件
func (g *Generate)writeJsonFile()  {
	//{
	//	"server": {
	//		"port": ":9998"
	//	},
	//	"mysql": {
	//		"host": "127.0.0.1",
	//		"port": "3306",
	//		"database": "jjmatch",
	//		"username": "root",
	//		"password": "root"
	//	},
	//	"redis": {
	//		"port": "6379",
	//		"password":"Bdz170635",
	//		"host":"127.0.0.1"
	//	},
	//	"test": {
	//		"dev": "123",
	//		"pro": "456"
	//	}
	//}
	g.p("{\n  \"server\": {\n    \"port\": \":9998\"\n  },\n  \"mysql\": {\n    \"host\": \"127.0.0.1\",\n    \"port\": \"3306\",\n    \"database\": \"jjmatch\",\n    \"username\": \"root\",\n    \"password\": \"root\"\n  },\n  \"redis\": {\n    \"port\": \"6379\",\n    \"password\":\"Bdz170635\",\n    \"host\":\"127.0.0.1\"\n  },\n  \"test\": {\n    \"dev\": \"123\",\n    \"pro\": \"456\"\n  }\n}\n\n")
}


//config文件代码生成
func (g *Generate) writeConfigFile() {
	//生成包名
	g.p("package config")
	g.p("\n")
	g.p("\n")
	//生成引入文件
	g.p("import (")
	g.p("\n\t")
	g.p("\"github.com/spf13/viper\"")
	g.p("\n\t")
	g.p("\"os\"")
	g.p("\n")
	g.p(")")
	g.p("\n")
	//生成parseConfig函数
	g.p("func ParseConfig()  {")
	g.p("\n\t")
	g.p("pwd, _ := os.Getwd()")
	g.p("\n\t")
	g.p("viper.SetConfigName(\"config\")")
	g.p("\n\t")
	g.p("viper.SetConfigType(\"json\")")
	g.p("\n\t")
	g.p("viper.AddConfigPath(pwd)")
	g.p("\n\t")
	g.p("err := viper.ReadInConfig()")
	g.p("\n\t")
	g.p("if err != nil {")
	g.p("\n\t\t")
	g.p("panic(\"read config file failed\")")
	g.p("\n\t")
	g.p("}")
	g.p("\n")
	g.p("}")
}
//生成db文件（包含redis和mysql）
func (g *Generate) writeDBFile() {
	//生成包名
	g.p("package db")
	g.p("\n")
	g.p("\n")
	//生成引入文件
	g.p("import (")
	g.p("\n\t")
	g.p("\"fmt\"")
	g.p("\n\t")
	g.p("\"github.com/spf13/viper\"")
	g.p("\n\t")
	g.p("\"github.com/go-redis/redis\"")
	g.p("\n\t")
	g.p("\"gorm.io/gorm\"")
	g.p("\n\t")
	g.p("\"github.com/biandz/bw/mysql\"")
	g.p("\n\t")
	g.p("re \"github.com/biandz/bw/redis\"")
	g.p("\n")
	g.p(")")

	g.p("\n")
	g.p("\n")

	//生成全局变量
	g.p("var (")
	g.p("\n\t")
	g.p("Msq   *gorm.DB")
	g.p("\n\t")
	g.p("Redis *redis.Client")
	g.p("\n")
	g.p(")")
	g.p("\n")
	//生成InitDB函数
	g.p("func InitDB() {\n\tdsn := fmt.Sprintf(\n\t\t\"%s:%s@tcp(%s:%s)/%s?timeout=90s&parseTime=true&loc=Local&collation\"+\n\t\t\t\"=utf8mb4_general_ci\",\n\t\t\tviper.GetString(\"mysql.username\"),\n\t\t\tviper.GetString(\"mysql.password\"),\n\t\t\tviper.GetString(\"mysql.host\"),\n\t\t\tviper.GetString(\"mysql.port\"),\n\t\t\tviper.GetString(\"mysql.database\"),\n\t)\n\tMsq = mysql.InitMysql(dsn)\n\n\tRedis = re.InitRedis(\n\t\tviper.GetString(\"redis.host\"),\n\t\tviper.GetString(\"redis.port\"),\n\t\tviper.GetString(\"redis.password\"),\n\t)\n}")
}
//生成route文件
func (g *Generate) writeRouteFile() {
	g.p(
		fmt.Sprintf("package route\n\nimport (\n\t\"%s/controller\"\n\t\"github.com/biandz/bw/middwear\"\n\t\"github.com/gin-gonic/gin\"\n)\n\nfunc LoadRouter(r *gin.Engine)  {\n\t//系统日志中间件,过滤器中间件(全局),拦截器中间件（全局）,跨域中间件\n\tr.Use(\n\t\tmiddwear.LoggerForFrame(),\n\t\tmiddwear.Filter(),\n\t\tmiddwear.Intercept(),\n\t\tmiddwear.Cors(),\n\t)\n\n\t//不需要加入jwt验证\n\tindex :=&controller.IndexController{}\n\tr.POST(\"/login\",middwear.Filter(),index.Login)\n\n\t//index组\n\tindexRouter := r.Group(\"/index\")\n\t//加入jst验证的路由\n\tindexRouter.Use(middwear.JwtVerify)\n\t{\n\t\tindexRouter.GET(\"/index\",index.Index)\n\t}\n}",
		g.PName))
}
//生成base控制器
func (g *Generate)writeBaseControllerFile()  {
	g.p("package controller\n\n\ntype BaseController struct {\n\n}\n")
}
//生成index控制器
func (g *Generate)writeIndexControllerFile()  {
	g.p(fmt.Sprintf("package controller\n\nimport (\n\t\"%s/model\"\n\t\"github.com/gin-gonic/gin\"\n\t\"net/http\"\n)\n\ntype IndexController struct {\n\tBaseController\n}\n\nfunc (i *IndexController)Login(g *gin.Context)  {\n\tg.JSON(http.StatusOK,&model.Rsp{Code: 0,Msg: \"success\",Data: nil})\n}\n\nfunc (i *IndexController)Index(g *gin.Context)  {\n\tg.JSON(http.StatusOK,&model.Rsp{Code: 0,Msg: \"success\",Data: nil})\n}\t",
		g.PName))
}
//生成基础model
func (g *Generate)writeModelFile () {
	g.p("package model\n\ntype Rsp struct{\n\tCode int `json:\"code\"`\n\tMsg  string `json:\"msg\"`\n\tData interface{} `json:\"data\"`\n}")
}

//生成service
func (g *Generate)writeServiceFile()  {
	g.p("package service")
}

//生成dao
func (g *Generate)writeDaoFile()  {
	g.p("package dao")
}














