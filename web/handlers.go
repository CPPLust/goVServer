package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomePage struct {
	Name string
}

type UserPage struct {
	Name string
}

//const TEMP_DIR = "d:/gowork/code/3.code/goVServer/web/"

func homeHandler(c *gin.Context) {
	_, err1 := c.Cookie("username")
	_, err2 := c.Cookie("session")

	if err1 != nil || err2 != nil {
		fmt.Println(err1, err2)
		p := &HomePage{Name: "zz"}
		t, e := template.ParseFiles("./templates/home.html")
		if e != nil {
			log.Printf("Parsing template home.html error: %s", e)
			return
		}
		t.Execute(c.Writer, p)
		return
	}

	c.Redirect(http.StatusFound, "/userhome")
}

func userHomeHandler(c *gin.Context) {
	cname, err1 := c.Cookie("username")
	_, err2 := c.Cookie("session")

	if err1 != nil || err2 != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	fname := c.Request.FormValue("username")

	var p *UserPage
	if len(cname) != 0 {
		p = &UserPage{Name: cname}
	} else if len(fname) != 0 {
		p = &UserPage{Name: fname}
	}

	t, e := template.ParseFiles("./templates/userhome.html")
	if e != nil {
		log.Printf("Parsing userhome.html error: %s", e)
		return
	}

	t.Execute(c.Writer, p)
}

func apiHomeHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		re, _ := json.Marshal(ErrorRequestNotRecognized)
		c.String(http.StatusOK, string(re))
		return
	}

	res, err := io.ReadAll(c.Request.Body)
	if err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		c.String(http.StatusOK, string(re))
		return
	}

	apibody := &ApiBody{}
	if err := json.Unmarshal(res, apibody); err != nil {
		re, _ := json.Marshal(ErrorRequestBodyParseFailed)
		c.String(http.StatusOK, string(re))
		return
	}

	request(apibody, c)
}
