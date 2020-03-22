package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"imooc.com/ccmouse/learngo/mockserver/generator/city"
	"imooc.com/ccmouse/learngo/mockserver/generator/profile"
	"imooc.com/ccmouse/learngo/mockserver/path"
	"imooc.com/ccmouse/learngo/mockserver/recommendation"
)

const (
	templateSuggestion = "Please make sure working directory is the root of the repository, where we have go.mod/go.sum. Suggested command line: go run mockserver/main.go"
	port               = 8080
)

func main() {
	profileTemplate, err := template.ParseFiles("mockserver/generator/profile/profile_tmpl.html")
	if err != nil {
		log.Fatalf("Cannot create profile template: %v. %s", err, templateSuggestion)
	}
	profileGen := &profile.Generator{
		Tmpl:           profileTemplate,
		Recommendation: recommendation.Client{},
	}

	cityTemplate, err := template.ParseFiles("mockServer/generator/city/city_tmpl.html")
	if err != nil {
		log.Fatalf("Cannot create city template: %v. %s", err, templateSuggestion)
	}
	cityGen := &city.Generator{
		Tmpl:       cityTemplate,
		ProfileGen: profileGen,
	}

	rand.Seed(time.Now().Unix())
	r := gin.Default()
	r.Use(path.Rewrite)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/static/index.html")
	})
	r.Static("/static", "mockserver/static")
	r.GET("mock/www.zhenai.com/zhenghun/:city/:page", cityGen.HandleRequest)
	r.GET("mock/www.zhenai.com/zhenghun/:city", cityGen.HandleRequest)
	r.GET("mock/album.zhenai.com/u/:id", profileGen.HandleRequest)

	log.Fatal(r.Run(fmt.Sprintf(":%d", port)))
}
