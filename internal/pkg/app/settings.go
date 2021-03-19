package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/labstack/echo/v4"
	"github.com/pkg/browser"
	"go.bmvs.io/ynab"
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
)

var (
	accessToken string
	expiresIn   string
)

func login(s GlobalState) {
	urlToOpen := fmt.Sprintf("https://app.youneedabudget.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=token", oauthClientID, redirectURL)
	ctx, cancel := context.WithCancel(context.Background())
	e := echo.New()
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
		// get setting
		var setting Setting
		result := s.DB.First(&setting)
		now := time.Now().UTC()
		if result.Error != nil {
			newSetting := Setting{
				AccessToken: &accessToken,
				ExpiresOn:   now.Add(time.Second * time.Duration(expiresInInt)),
				Model: Model{
					CreatedOn: now,
					UpdatedOn: now,
				},
			}
			s.YNABClient = ynab.NewClient(accessToken)
			result := s.DB.Create(&newSetting)
			if result.Error != nil {
				log.Println(result.Error)
			}
		} else {
			// update
			setting.AccessToken = &accessToken
			setting.Model.UpdatedOn = now
			setting.ExpiresOn = now.UTC().Add(time.Second * time.Duration(expiresInInt))
			s.YNABClient = ynab.NewClient(accessToken)
			s.DB.Save(&setting)
		}
		var accountsToInsert []Account
		accounts, err := s.YNABClient.Account().GetAccounts("default")
		if err != nil {
			log.Println(err)
		}
		for _, account := range accounts {
			accountsToInsert = append(accountsToInsert, Account{
				YNABID:  account.ID,
				Name:    account.Name,
				Type:    string(account.Type),
				Closed:  account.Closed,
				Deleted: account.Deleted,
				Model: Model{
					CreatedOn: now,
					UpdatedOn: now,
				},
			})
		}
		s.DB.Create(&accountsToInsert)
		categories, err := s.YNABClient.Category().GetCategories("default")
		if err != nil {
			log.Println(err)
		}
		var categoriesToCreate []Category
		for _, group := range categories {
			for _, c := range group.Categories {
				categoriesToCreate = append(categoriesToCreate, Category{
					YNABGroupID:    group.ID,
					GroupName:      group.Name,
					GroupHidden:    group.Hidden,
					GroupDeleted:   group.Deleted,
					YNABCategoryID: c.ID,
					Name:           c.Name,
					Hidden:         c.Hidden,
					Deleted:        c.Deleted,
					Model: Model{
						CreatedOn: now,
						UpdatedOn: now,
					},
				})
			}
		}
		result = s.DB.Create(&categoriesToCreate)
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
}

func (s GlobalState) NewSettingsScreen() *Component {
	/*
		loginButton := widget.NewButton("Login to YNAB", func() {
			login(s)
		})
		logoutButton := widget.NewButton("Logout of YNAB", func() {
			log.Println("tapped")
		})
		refreshCategoriesButton := widget.NewButton("Refresh YNAB Categories", func() {
			log.Println("tapped")
		})
		log.Println(loginButton, logoutButton)
		return NewComponent(
			container.NewVBox(loginButton, logoutButton, refreshCategoriesButton),
		)
	*/
	return NewComponent(
		container.NewMax(container.NewVBox(widget.NewButton("beep boop bop", func() {
			log.Println("tapped")
		}), widget.NewButton("beep boop bop", func() {
			log.Println("tapped")
		}))),
	)
}
