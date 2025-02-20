package login

import (
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/helper"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
	"time"
)

func NapCatLogin() {
	if flags.Config.NapCatPanelURL == "" || flags.Config.NapCatToken == "" {
		log.Error("NapCatShellUpdater", "NapCatPanelURL or NapCatToken is empty")
		return
	}
	token := loginNapCatPanel()
	if token != "" {
		loginList := getNapCatPanelLoginList(token)
		if len(loginList) <= 0 {
			return
		}
		setNapCatQuickLogin(token, loginList)
	}
}

func loginNapCatPanel() (token string) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/auth/login", flags.Config.NapCatPanelURL), strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, flags.Config.NapCatToken)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Trace("NapCatShellUpdater", "Login to NapCat Panel:", helper.BytesToString(body))
	return gjson.Parse(helper.BytesToString(body)).Get("data.Credential").String()
}

func getNapCatPanelLoginList(token string) []int64 {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/QQLogin/GetQuickLoginList", flags.Config.NapCatPanelURL), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var loginList []int64
	gjson.Parse(helper.BytesToString(body)).Get("data").ForEach(func(key, value gjson.Result) bool {
		loginList = append(loginList, value.Int())
		return true
	})
	log.Trace("NapCatShellUpdater", "Get NapCat Panel Login List:", loginList)
	return loginList
}

func setNapCatQuickLogin(token string, loginList []int64) {
	for _, uin := range loginList {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/QQLogin/SetQuickLogin", flags.Config.NapCatPanelURL), strings.NewReader(fmt.Sprintf(`{"uin":"%d"}`, uin)))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		log.Trace("NapCatShellUpdater", uin, " | Set NapCat Panel Quick Login:", helper.BytesToString(body))
		time.Sleep(222 * time.Millisecond)
	}
}
