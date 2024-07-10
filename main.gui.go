package main

import (
	"image/color"
	"log"
	"strings"

	"yubifyne/otp"

	"fyne.io/fyne/v2/canvas"
	"github.com/GeertJohan/yubigo"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type gui struct {
	win fyne.Window

	user fyne.Focusable
}

func newGUI(w fyne.Window) *gui {
	return &gui{win: w}
}

func (g *gui) makeUI() fyne.CanvasObject {
	user := widget.NewEntry()
	pass := widget.NewPasswordEntry()
	g.user = user

	key := otp.NewInput(func(key string) {
		if user.Text == "" || pass.Text == "" {
			dialog.ShowInformation("Login", "Missing username or password", g.win)
			return
		}
		if user.Text != pass.Text {
			dialog.ShowInformation("Login", "Failed login, wrong password", g.win)
			return
		}

		wait := widget.NewActivity()
		prop := canvas.NewRectangle(color.Transparent)
		prop.SetMinSize(fyne.NewSquareSize(64))
		content := container.NewStack(prop, wait)
		d := dialog.NewCustomWithoutButtons("Checking...", content, g.win)
		wait.Start()
		d.Show()

		go func() {
			id, err := auth(key)
			d.Hide()
			wait.Stop()

			if err == nil {
				dialog.ShowInformation("Login", "Logged in: "+user.Text+"\nKeyID: "+id, g.win)
			} else {
				dialog.ShowInformation("Login", "Failed key check!\n"+err.Error(), g.win)
			}
		}()
	})

	return container.NewVBox(
		widget.NewLabel("Yubikey Example"),
		widget.NewForm(
			widget.NewFormItem("Username", user),
			widget.NewFormItem("Password", pass),
			widget.NewFormItem("Key", key)))
}

func auth(key string) (string, error) {
	// get the API key from https://upgrade.yubico.com/getapikey/
	yubiAuth, err := yubigo.NewYubiAuth("<id>", "<secret key>")
	if err != nil {
		return "", err
	}

	// verify an OTP string
	result, ok, err := yubiAuth.Verify(key)
	if err != nil {
		return "", err
	}

	if ok && result.IsValidOTP() {
		log.Println("RES", result.GetRequestQuery())

		return strings.ToLower(key[0 : len(key)-32]), nil
	}
	return "", err
}
