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
	"runtime"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//C "github.com/kelseyhightower/envconfig"
)

//go:embed static/style.min.css static/script.min.js static/fonts/* internal/html/*
var f embed.FS
var AppVer = "1.0"
var GitHash = "1234567"
var BuildTime = "2023-02-01"

func main() {
	log.SetOutput(os.Stdout)
	c, err := getConfig()
	if err != nil {
		log.Printf("Load config file failed, will use in memory config only. Error = %s", err)
	}

	log.Println("Hello World!")
	runGin(c)
}

func runGin(conf Config) {
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

	normalizedConfig, err := normalizeConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":     "test",
			"AppVer":    AppVer,
			"GitHash":   GitHash,
			"BuildTime": BuildTime,
			"GoVer":     runtime.Version(),
			"Items":     normalizedConfig,
		})
	})

	r.Run()
}

func normalizeConfig(c Config) (Config, error) {
	conf := c
	if len(conf.Shell) <= 0 {
		if runtime.GOOS == "windows" {
			conf.Shell = "cmd.exe"
			conf.ShellArgument = "/q /c"
		} else {
			conf.Shell = "/bin/sh"
			conf.ShellArgument = ""
		}
	}

	for index, ele := range conf.Sites {
		if len(ele.Name) <= 0 {
			return conf, fmt.Errorf("invalid name for site %d", index+1)
		}

		if len(ele.Protocol) <= 0 {
			conf.Sites[index].Protocol = "http"
		}

		if len(ele.Host) <= 0 {
			return conf, fmt.Errorf("invalid host for site %d", index+1)
		}

		if len(ele.Port) <= 0 {
			conf.Sites[index].Port = "80"
		}

		if len(ele.OkCode) <= 0 {
			conf.Sites[index].OkCode = "200"
		}

		if len(ele.StartScript) > 0 {
			_, err := os.Stat(ele.StartScript)
			if err != nil {
				return c, err
			}
		}

		if len(ele.StopScript) > 0 {
			_, err := os.Stat(ele.StopScript)
			if err != nil {
				return c, err
			}
		}

		if len(ele.RestartScript) > 0 {
			_, err := os.Stat(ele.RestartScript)
			if err != nil {
				return c, err
			}
		}
	}

	return conf, nil
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

	for index := range c.Sites {
		c.Sites[index].Id = index + 1
	}

	log.Printf("Load config file from %s succeed.\n", configFile)
	return c, nil
}

type Site struct {
	Id            int    `json:"-"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Protocol      string `json:"protocol"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	PathName      string `json:"path_name"`
	OkCode        string `json:"ok_code"`
	StartScript   string `json:"start_script"`
	StopScript    string `json:"stop_script"`
	RestartScript string `json:"restart_script"`
	IsRunning     bool   `json:"-"`
}

type Config struct {
	Shell         string `json:"shell"`
	ShellArgument string `json:"shell_arg"`
	Sites         []Site `json:"sites"`
}
