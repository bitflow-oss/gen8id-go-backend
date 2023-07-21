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
	ServerPort string `yaml:"server_port"`

	OgUrl         string `yaml:"og_url"`
	OgSiteName    string `yaml:"og_site_name"`
	OgType        string `yaml:"og_type"`
	OgTitle       string `yaml:"og_title"`
	OgDescription string `yaml:"og_description"`
	OgImage       string `yaml:"og_image"`
	OgImageType   string `yaml:"og_image_type"`
	OgImageWidth  string `yaml:"og_image_width"`
	OgImageHeight string `yaml:"og_image_height"`
	OgLocale      string `yaml:"og_locale"`

	Title     string `yaml:"title"`
	Keywords  string `yaml:"keywords"`
	Author    string `yaml:"author"`
	Copyright string `yaml:"copyright"`
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
	// confFilepath, _ := filepath.Abs(filepath)
	bytes, err := os.ReadFile(path.Join(absdir, filename))
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
