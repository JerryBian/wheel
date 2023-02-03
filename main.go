package main

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//C "github.com/kelseyhightower/envconfig"
)

//go:embed static/style.min.css static/script.min.js static/fonts/* internal/html/*
var f embed.FS
var AppVer = "1.0"
var GitHash = "1234567"
var BuildTime = "2022-12-30"

func main() {
	log.SetOutput(os.Stdout)
	c, err := getConfig()
	if err != nil {
		log.Printf("Load config file failed, will use in memory config only. Error = %s", err)
	}



	log.Println("Hello World!")
	runGin(c)
}

func runGin(c Config) {
	r := gin.Default()
	uid := uuid.New()
	sessionSecret := []byte(uid.String())
	r.Use(sessions.Sessions("wheel__", sessions.NewCookieStore(sessionSecret)))

	templ := template.Must(template.New("").ParseFS(f, "internal/html/*.html"))
	r.SetHTMLTemplate(templ)

	staticFs, _ := fs.Sub(f, "static")
	r.StaticFS("/static", http.FS(staticFs))

	/* r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": "Login",
			"Config": h.Conf,
		})
	})

	r.POST("/login", h.loginHander)
	r.GET("/logout", logoutHandler)

	authRoute := r.Group("/")
	authRoute.Use(AuthRequired)
	authRoute.GET("/", h.indexHandler)
	authRoute.GET("/add", h.addDiaryGetHandler)
	authRoute.POST("/api/add", h.addWordHandler)
	authRoute.GET("/:year/:month/:day", h.getDiariesHandler)
	authRoute.GET("/edit/:id", h.editDiaryGetHandler)
	authRoute.GET("/revision/:id", h.revisionHandler) */

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title" : "test",
		})
	})

	r.Run()
}

func getConfig() (Config, error) {
	var c Config
	configFile, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		if len(os.Args) > 1 {
			configFile = os.Args[1]
		}
	}

	if len(configFile) <= 0 {
		return c, errors.New("missing config file, specify JSON file either in first command line argument or environment variable CONFIG_FILE")
	}

	configFile, err := filepath.Abs(configFile)
	if err != nil {
		return c, err
	}

	_, err = os.Stat(configFile)
	if err != nil {
		return c, err
	}

	f, err := ioutil.ReadFile(configFile)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal([]byte(f), &c)
	if err != nil {
		return c, err
	}

	log.Println(fmt.Sprintf("Load config file from %s succeed.", configFile))
	return c, nil
}

type Site struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Protocol string `json:"protocol"`
	Host string `json:"host"`
	Port int `json:"port"`
	PathName string `json:"path_name"`
	OkCode string `json:"ok_code"`
	StartScript string `json:"start_script"`
	StopScript string `json:"stop_script"`
	RestartScript string `json:"restart_script"`
}

type Config struct {
	Shell string `json:"shell"`
	ShellArgument string `json:"shell_arg"`
	Sites []Site `json:"sites"`
}