package ynab

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pkg/browser"
	"github.com/thebenwaters/ynaboosh-desktop/pkg/internal/ynaboosh/models"
	externalYnab "go.bmvs.io/ynab"
)

const (
	oauthClientID = "1016b552249ee6006b6524144378da195efea9f7e9bb156ce63c01d04d593dea"
	redirectURL   = "http://localhost:8080"
	htmlResponse  = `
	<html>
	<head>
	<script type="text/javascript">
	let fragment = new URLSearchParams(window.location.hash.substr(1));
	let data = {
		token: fragment.get("access_token"),
		expires: fragment.get("expires_in")
	};
	
	fetch('http://localhost:8080/access', {
		method: 'POST', // or 'PUT'
  		headers: {
    		'Content-Type': 'application/json',
  		},
  		body: JSON.stringify(data)
	})
	window.close()
	</script>
	</head>
	<body>
	<h1> Success, you may now close the browser </h1>
	</body>
	</html>
	`
	authorizeURL = "https://app.youneedabudget.com/oauth/authorize?"
)

var (
	accessToken    string
	expiresIn      string
	externalClient externalYnab.ClientServicer
)

func Login(manager models.DBManager) externalYnab.ClientServicer {
	urlToOpen := fmt.Sprintf("%sclient_id=%s&redirect_uri=%s&response_type=token", authorizeURL, oauthClientID, redirectURL)
	ctx, cancel := context.WithCancel(context.Background())
	e := echo.New()
	e.HideBanner = true
	e.Debug = false
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, htmlResponse)
	})
	e.POST("/access", func(c echo.Context) error {
		m := echo.Map{}
		if err := c.Bind(&m); err != nil {
			return err
		}
		accessToken = m["token"].(string)
		expiresIn = m["expires"].(string)
		expiresInInt, _ := strconv.Atoi(expiresIn)
		externalClient = externalYnab.NewClient(accessToken)
		log.Println(expiresInInt, externalClient)
		var accountsToInsert []models.Account
		accounts, err := externalClient.Account().GetAccounts("default")
		if err != nil {
			log.Println(err)
		}
		for _, account := range accounts {
			accountsToInsert = append(accountsToInsert, models.Account{
				YNABID:  account.ID,
				Name:    account.Name,
				Type:    string(account.Type),
				Closed:  account.Closed,
				Deleted: account.Deleted,
			})
		}
		manager.CreateAccounts(accountsToInsert)
		categories, err := externalClient.Category().GetCategories("default")
		if err != nil {
			log.Println(err)
		}
		var categoriesToCreate []models.Category
		for _, group := range categories {
			for _, c := range group.Categories {
				categoriesToCreate = append(categoriesToCreate, models.Category{
					YNABGroupID:    group.ID,
					GroupName:      group.Name,
					GroupHidden:    group.Hidden,
					GroupDeleted:   group.Deleted,
					YNABCategoryID: c.ID,
					Name:           c.Name,
					Hidden:         c.Hidden,
					Deleted:        c.Deleted,
				})
			}
		}
		manager.CreateCategories(categoriesToCreate)
		cancel()
		return c.String(200, "success")
	})
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	_ = browser.OpenURL(urlToOpen)
	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Info(err)
	}
	return externalClient
}
