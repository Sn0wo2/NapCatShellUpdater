package main

import (
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat"
	"github.com/Sn0wo2/NapCatShellUpdater/napcat/login"
	"github.com/sirupsen/logrus"
	"os"
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
	cv := flags.Config.Version
	if cv == "" {
		napcat.CheckNapCatUpdate()
	} else {
		napcat.ProcessVersionUpdate(cv)
	}
	if flags.Config.Login {
		log.Info("NapCatShellUpdater", "Wating NapCat process to login...")
		ncProc, err := napcat.WaitForProcess(filepath.Join(flags.Config.Path, "NapCatWinBootMain.exe"))
		select {
		case p := <-ncProc:
			log.Debug("NapCatShellUpdater", "NapCat process found:", p.String())
		case e := <-err:
			log.Error("NapCatShellUpdater", "Failed to find NapCat process:", e)
			os.Exit(1)
		}
		log.Debug("NapCatShellUpdater", "Waiting 30s to full load NapCat")
		time.Sleep(30 * time.Second)
		log.Info("NapCatShellUpdater", "Login to NapCat Panel...")
		login.NapCatLogin()
	}
}
