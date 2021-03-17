package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

var PROJECT_NAME = ""
var MAIL_FROM = ""
var MAIL_TO = ""
var SMTP_USER = ""
var SMTP_PASSWORD = ""
var SMTP_HOST = ""
var SMTP_PORT = ""
var SECRET_TOKEN = ""

func sendMail(title string, body string) {
	auth := smtp.PlainAuth("", SMTP_USER, SMTP_PASSWORD, SMTP_HOST)
	msg := []byte("To: " + MAIL_FROM + "\r\nSubject: " + title + "\r\n\r\n" + body)
	err := smtp.SendMail(SMTP_HOST + ":" + SMTP_PORT, auth, MAIL_FROM, []string{MAIL_TO}, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func isValidSecret(c *gin.Context) bool {
	if val, ok := c.Request.Header["X-Gitlab-Token"]; ok {
		if val[0] == SECRET_TOKEN {
			return true
		}
	}
	return false
}

func mergeRequestAttrs(attrs map[string]interface{}) {
	var title = ""
	var description = ""
	var url = ""
	var action = ""

	for key, value := range attrs {
		if key == "title" {
			title = fmt.Sprintf("%v", value)
		} else if key == "description" {
			description = fmt.Sprintf("%v", value)
		} else if key == "url" {
			url = fmt.Sprintf("%v", value)
		} else if key == "action" {
			action = fmt.Sprintf("%v", value)
		}
	}

	log.Println("title", title)
	log.Println("description", description)
	log.Println("url", url)
	log.Println("action", action)

	if title != "" && url != "" {
		if action == "open" {
			var msgTitle = "[" + PROJECT_NAME + "] MR opened: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		} else if action == "close" {
			var msgTitle = "[" + PROJECT_NAME + "] MR closed: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		} else if action == "reopen" {
			var msgTitle = "[" + PROJECT_NAME + "] MR reopened: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		} else if action == "merge" {
			var msgTitle = "[" + PROJECT_NAME + "] MR merged: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		}
	}
}

func mergeRequest(c *gin.Context) {
	if !isValidSecret(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result map[string]interface{}
	json.Unmarshal(buf, &result)

	for key, value := range result {
		if key == "object_kind" {
			if value != "merge_request" {
				log.Println("not a merge_request!")
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
		} else if key == "object_attributes" {
			mergeRequestAttrs(value.(map[string]interface{}))
		}
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf))

	c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
}

func issueAttrs(attrs map[string]interface{}) {
	var title = ""
	var description = ""
	var url = ""
	var action = ""

	for key, value := range attrs {
		if key == "title" {
			title = fmt.Sprintf("%v", value)
		} else if key == "description" {
			description = fmt.Sprintf("%v", value)
		} else if key == "url" {
			url = fmt.Sprintf("%v", value)
		} else if key == "action" {
			action = fmt.Sprintf("%v", value)
		}
	}

	log.Println("title", title)
	log.Println("description", description)
	log.Println("url", url)
	log.Println("action", action)

	if title != "" && url != "" {
		if action == "open" {
			var msgTitle = "[" + PROJECT_NAME + "] Issue opened: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		} else if action == "close" {
			var msgTitle = "[" + PROJECT_NAME + "] Issue closed: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		} else if action == "reopen" {
			var msgTitle = "[" + PROJECT_NAME + "] Issue reopened: " + title
			var msgBody = description + "\r\n---\r\n" + url
			sendMail(msgTitle, msgBody)
		}
	}
}

func issue(c *gin.Context) {
	if !isValidSecret(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result map[string]interface{}
	json.Unmarshal(buf, &result)

	for key, value := range result {
		if key == "object_kind" {
			if value != "issue" {
				log.Println("not an issue!")
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
				return
			}
		} else if key == "object_attributes" {
			issueAttrs(value.(map[string]interface{}))
		}
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(buf))

	c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	PROJECT_NAME = os.Getenv("PROJECT_NAME")
	if PROJECT_NAME == "" {
		log.Fatal("$PROJECT_NAME must be set")
	}
	MAIL_FROM = os.Getenv("MAIL_FROM")
	if PROJECT_NAME == "" {
		log.Fatal("$MAIL_FROM must be set")
	}
	MAIL_TO = os.Getenv("MAIL_TO")
	if MAIL_TO == "" {
		log.Fatal("$MAIL_TO must be set")
	}
	SMTP_USER = os.Getenv("SMTP_USER")
	if SMTP_USER == "" {
		log.Fatal("$SMTP_USER must be set")
	}
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	if SMTP_PASSWORD == "" {
		log.Fatal("$SMTP_PASSWORD must be set")
	}
	SMTP_HOST = os.Getenv("SMTP_HOST")
	if SMTP_HOST == "" {
		log.Fatal("$SMTP_HOST must be set")
	}
	SMTP_PORT = os.Getenv("SMTP_PORT")
	if SMTP_PORT == "" {
		log.Fatal("$SMTP_PORT must be set")
	}
	SECRET_TOKEN = os.Getenv("SECRET_TOKEN")
	if SECRET_TOKEN == "" {
		log.Fatal("$SECRET_TOKEN must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.POST("/issue", issue)
	router.POST("/mergerequest", mergeRequest)

	router.Run(":" + port)
}
