package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	lorem "github.com/sbroekhoven/golorem"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Listen       string `json:"listen"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	Lines        int    `json:"lines"`
	WordFile     string `json:"word_file"`
	Words        []string
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	conf, err := LoadConfiguration("config.json")
	if err != nil {
		println(err)
	}

	lines, err := readLines(conf.WordFile)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	conf.Words = Shuffle(lines)
	println(len(conf.Words))

	// Template
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	// e.Pre(middleware.AddTrailingSlash())
	e.Renderer = t

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set server in header
	e.Use(ServerHeader)

	e.Static("/assets", "public/assets")
	e.File("/robots.txt", "public/robots.txt")
	e.File("/favicon.ico", "public/favicon.ico")

	e.GET("/", conf.Home)
	e.GET("/internal", conf.Internal)
	e.GET("/internal/", conf.Internal)
	e.GET("/internal/:next", conf.InternalNext)
	e.GET("/internal/:next/", conf.InternalNext)

	e.GET("/documents", conf.Internal)
	e.GET("/documents/", conf.Internal)
	e.GET("/documents/:next", conf.InternalNext)
	e.GET("/documents/:next/", conf.InternalNext)

	e.GET("/login", conf.Internal)
	e.GET("/login/", conf.Internal)
	e.GET("/login/:next", conf.InternalNext)
	e.GET("/login/:next/", conf.InternalNext)

	// Private pages
	private := e.Group("/login")
	private.Use(middleware.BasicAuth(conf.ValidateUser))
	private.GET("/", conf.Internal)
	private.GET("/:next", conf.InternalNext)

	e.Logger.Fatal(e.Start(conf.Listen))
}

func (conf Config) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"Title": "Home",
	})
}

func (conf Config) Internal(c echo.Context) error {

	real_ip := c.RealIP()
	user_agent := c.Request().UserAgent()
	request_url := c.Request().URL.String()
	logrus.WithFields(logrus.Fields{
		"authentication": "none",
		"real_ip":        real_ip,
		"user_agent":     user_agent,
		"request_url":    request_url,
	}).Warn("Ignored robots.txt")

	paragraph := lorem.Paragraph(5, 8)
	lines := getWords(conf.Words, conf.Lines)
	return c.Render(http.StatusOK, "template.html", map[string]interface{}{
		"Title":     "Internal",
		"Paragraph": paragraph,
		"Lines":     lines,
	})
}

func (conf Config) InternalNext(c echo.Context) error {

	real_ip := c.RealIP()
	user_agent := c.Request().UserAgent()
	request_url := c.Request().URL.String()
	logrus.WithFields(logrus.Fields{
		"authentication": "none",
		"real_ip":        real_ip,
		"user_agent":     user_agent,
		"request_url":    request_url,
	}).Warn("Ignored robots.txt and crawling")

	if !findWord(conf.Words, c.Param("next")) {
		return c.String(http.StatusNotFound, "Not found...")
	}
	lines := getWords(conf.Words, conf.Lines)

	paragraph := lorem.Paragraph(5, 8)

	return c.Render(http.StatusOK, "template.html", map[string]interface{}{
		"Title":     strings.Title(c.Param("next")),
		"Paragraph": paragraph,
		"Lines":     lines,
	})
}
