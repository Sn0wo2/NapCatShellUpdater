package main

import (
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat/login"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
	flags.InitFlag()

	err := log.InitLogger("", log.DefaultFormatter(), logrus.TraceLevel)
	if err != nil {
		panic(err)
	}

	fmt.Println(`$$\   $$\                   $$$$$$\           $$\     $$$$$$\ $$\               $$\$$\$$\   $$\               $$\          $$\                       
$$$\  $$ |                 $$  __$$\          $$ |   $$  __$$\$$ |              $$ $$ $$ |  $$ |              $$ |         $$ |                      
$$$$\ $$ |$$$$$$\  $$$$$$\ $$ /  \__|$$$$$$\$$$$$$\  $$ /  \__$$$$$$$\  $$$$$$\ $$ $$ $$ |  $$ |$$$$$$\  $$$$$$$ |$$$$$$\$$$$$$\   $$$$$$\  $$$$$$\  
$$ $$\$$ |\____$$\$$  __$$\$$ |      \____$$\_$$  _| \$$$$$$\ $$  __$$\$$  __$$\$$ $$ $$ |  $$ $$  __$$\$$  __$$ |\____$$\_$$  _| $$  __$$\$$  __$$\ 
$$ \$$$$ |$$$$$$$ $$ /  $$ $$ |      $$$$$$$ |$$ |    \____$$\$$ |  $$ $$$$$$$$ $$ $$ $$ |  $$ $$ /  $$ $$ /  $$ |$$$$$$$ |$$ |   $$$$$$$$ $$ |  \__|
$$ |\$$$ $$  __$$ $$ |  $$ $$ |  $$\$$  __$$ |$$ |$$\$$\   $$ $$ |  $$ $$   ____$$ $$ $$ |  $$ $$ |  $$ $$ |  $$ $$  __$$ |$$ |$$\$$   ____$$ |      
$$ | \$$ \$$$$$$$ $$$$$$$  \$$$$$$  \$$$$$$$ |\$$$$  \$$$$$$  $$ |  $$ \$$$$$$$\$$ $$ \$$$$$$  $$$$$$$  \$$$$$$$ \$$$$$$$ |\$$$$  \$$$$$$$\$$ |      
\__|  \__|\_______$$  ____/ \______/ \_______| \____/ \______/\__|  \__|\_______\__\__|\______/$$  ____/ \_______|\_______| \____/ \_______\__|      
                  $$ |                                                                         $$ |                                                  
                  $$ |                                                                         $$ |                                                  
                  \__|                                                                         \__|                                                  `)

	log.Info("NapCatShellUpdater", "Loading...")

	if runtime.GOOS != "windows" {
		log.Error("NapCatShellUpdater", "Unsupported system:", runtime.GOOS)
	}
}

func main() {
	if !flags.Config.SkipCheck {
		cv := flags.Config.Version
		if cv == "" {
			napcat.CheckNapCatUpdate()
		} else {
			napcat.ProcessVersionUpdate(cv)
		}
	}
	if flags.Config.Login {
		log.Info("NapCatShellUpdater", "Wating NapCat process to login...")
		ncProc, err := napcat.WaitForProcess(filepath.Join(flags.Config.Path, "NapCatWinBootMain.exe"))
		select {
		case p := <-ncProc:
			log.Debug("NapCatShellUpdater", "NapCat process found:", p.String())
		case e := <-err:
			panic(e)
		}
		log.Info("NapCatShellUpdater", fmt.Sprintf("Waiting %s to full load NapCat", flags.Config.Sleep.String()))
		time.Sleep(flags.Config.Sleep)
		log.Info("NapCatShellUpdater", "Login to NapCat Panel...")
		if flags.Config.NapCatPanelURL == "" || flags.Config.NapCatToken == "" {
			log.Error("NapCatShellUpdater", "NapCatPanelURL or NapCatToken is empty, trying find NapCat Panel url and token in logs...")
			url, token, err := napcat.GetNapCatPanelURLInLogs(filepath.Join(flags.Config.Path, "logs"))
			if err != nil {
				panic(err)
			}
			flags.Config.NapCatPanelURL = url
			flags.Config.NapCatToken = token
		}
		log.Debug("NapCatShellUpdater", "Panel URL: ", flags.Config.NapCatPanelURL)
		log.Debug("NapCatShellUpdater", "Panel Token: ", flags.Config.NapCatToken)
		login.NapCatLogin()
	}
}
