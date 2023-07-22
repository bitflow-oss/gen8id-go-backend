package util

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"runtime"
)

type Config struct {
	ServerPort string `yaml:"ServerPort"`

	OgUrl         string `yaml:"OgUrl"`
	OgSiteName    string `yaml:"OgSiteName"`
	OgType        string `yaml:"OgType"`
	OgTitle       string `yaml:"OgTitle"`
	OgDescription string `yaml:"OgDescription"`
	OgImage       string `yaml:"OgImage"`
	OgImageType   string `yaml:"OgImageType"`
	OgImageWidth  string `yaml:"OgImageWidth"`
	OgImageHeight string `yaml:"OgImageHeight"`
	OgLocale      string `yaml:"OgLocale"`

	TwitterCard string `yaml:"TwitterCard"`

	Title     string `yaml:"Title"`
	Keywords  string `yaml:"Keywords"`
	Author    string `yaml:"Author"`
	Copyright string `yaml:"Copyright"`

	ObjStrgEndpnt   string `yaml:"ObjStrgEndpnt"`
	ObjStrgRegion   string `yaml:"ObjStrgRegion"`
	ObjStrgAccKey   string `yaml:"ObjStrgAccKey"`
	ObjStrgScrtKey  string `yaml:"ObjStrgScrtKey"`
	ObjStrgBcktName string `yaml:"ObjStrgBcktName"`
	ObjStrgFoldPblc string `yaml:"ObjStrgFoldPblc"`
	ObjStrgFoldPrvt string `yaml:"ObjStrgFoldPrvt"`

	WtmkThmbPath  string `yaml:"WtmkThmbPath"`
	UpldRltvPath  string `yaml:"UpldRltvPath"`
	OrgImgFileNm  string `yaml:"OrgImgFileNm"`
	HashImgFileNm string `yaml:"HashImgFileNm"`
	DstImgFileNm  string `yaml:"DstImgFileNm"`
}

var AppConfig Config

func GetConfig() Config {
	//	mu.RLock() // readers lock
	//	defer mu.RUnlock()
	return AppConfig
}

func LoadConfig(filename string) Config {

	_, filedir, _, _ := runtime.Caller(0)
	absdir := path.Join(path.Dir(filedir), "..", "..")
	err := os.Chdir(absdir)
	if err != nil {
		panic(err)
	}
	var absFilePath = path.Join(absdir, filename)
	log.Println("[init] config at", absFilePath, "is loading")

	bytes, err := os.ReadFile(absFilePath)
	if err != nil {
		log.Println(err.Error())
	}
	var lAppConfig Config
	err = yaml.Unmarshal(bytes, &lAppConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	AppConfig = lAppConfig
	return lAppConfig
}

func GetHtmlTemplateDir(env string) string {
	var filename = ".env"
	if len(env) > 0 {
		filename = "." + env + ".env"
	}
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("HTML_TMPLT_DIR")
}
