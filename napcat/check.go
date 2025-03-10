package napcat

import (
	"bufio"
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/helper"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func CheckNapCatUpdate() {
	newVersion := fetchLastNapCatDownloadURL()
	currentVersion := getCurrentNapCatVersion()
	if newVersion != currentVersion {
		log.Info("NapCatShellUpdater", fmt.Sprintf("Updating NapCat from %s to %s", currentVersion, newVersion))
		processAndUpdate(downloadFile(fmt.Sprintf("https://github.com/NapNeko/NapCatQQ/releases/download/%s/NapCat.Shell.zip", newVersion)))
	} else {
		log.Info("NapCatShellUpdater", "NapCat is up to date: ", currentVersion)
	}
}

func ProcessVersionUpdate(ver string) {
	currentVersion := getCurrentNapCatVersion()
	if ver == "" || currentVersion == "" {
		log.Error("NapCatShellUpdater", "Failed to fetch version info", ver, currentVersion)
		return
	}
	if ver != currentVersion {
		processAndUpdate(downloadFile(fmt.Sprintf("https://github.com/NapNeko/NapCatQQ/releases/download/%s/NapCat.Shell.zip", ver)))
	} else {
		log.Info("NapCatShellUpdater", "NapCat is up to date: ", currentVersion)
	}
}

func GetNapCatPanelURLInLogs(dirPath string) (string, string, error) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil || !fileInfo.IsDir() {
		return "", "", fmt.Errorf("invalid directory path: %s", dirPath)
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read directory: %v", err)
	}
	urlTokenRegex := regexp.MustCompile(`(https?://[^\s:/]+:\d+)/webui\?token=([^\s]+)`)
	var logFiles []struct {
		Path    string
		ModTime time.Time
	}
	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || strings.ToLower(filepath.Ext(entry.Name())) != ".log" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		logFiles = append(logFiles, struct {
			Path    string
			ModTime time.Time
		}{
			Path:    fullPath,
			ModTime: fileInfo.ModTime(),
		})
	}

	if len(logFiles) == 0 {
		return "", "", fmt.Errorf("no log files found in %s", dirPath)
	}

	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime.After(logFiles[j].ModTime)
	})

	for _, logFile := range logFiles {
		f, err := os.Open(logFile.Path)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			matches := urlTokenRegex.FindStringSubmatch(scanner.Text())
			if len(matches) >= 3 {
				return matches[1], matches[2], nil
			}
		}
	}

	return "", "", fmt.Errorf("no matching URL found in %s", dirPath)
}

func processAndUpdate(filename string) {
	log.Info("NapCatShellUpdater", "Waiting NapCatWinBootMain.exe process to end...")
	err := <-WaitForAllProcessesEnd(filepath.Join(flags.Config.Path, "NapCatWinBootMain.exe"), true)
	if err != nil {
		panic(err)
	}

	log.Info("NapCatShellUpdater", "Clean target directory...")
	err = cleanDirectory(flags.Config.Path, []string{"config", "logs", "quickLoginExample.bat", "update.bat", filepath.Base(os.Args[0]), filename})
	if err != nil {
		log.RPanic(err)
	}

	log.Info("NapCatShellUpdater", "Extracting new version...")
	err = unzipWithExclusion(filename, flags.Config.Path, []string{"quickLoginExample.bat"})
	if err != nil {
		panic(err)
	}

	err = os.Remove(filename)
	if err != nil {
		panic(err)
	}
}

func getCurrentNapCatVersion() (ver string) {
	packageFile, err := os.ReadFile(filepath.Join(flags.Config.Path, "package.json"))
	if err != nil {
		log.Error("NapCatShellUpdater", "failed to read package.json:", err)
		return "v0.0.0(Error)"
	}
	version := gjson.GetBytes(packageFile, "version").String()
	if version == "" {
		version = "0.0.0(Not Found)"
	}
	return "v" + version
}

func fetchLastNapCatDownloadURL() (ver string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/NapNeko/NapCatQQ/releases?per_page=1", nil)
	if err != nil {
		panic(err)
	}

	client := http.DefaultClient
	if flags.Config.Proxy != "" {
		p, err := url.Parse(flags.Config.Proxy)
		if err != nil {
			panic(err)
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(p)}}
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error("NapCatShellUpdater", fmt.Sprintf("Failed to fetch version info: %v, status: %d", err, resp.StatusCode))
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	release := gjson.Parse(helper.BytesToString(body)).Array()[0]
	version := release.Get("tag_name").Str
	if version == "" {
		log.Error("NapCatShellUpdater", "Failed to fetch version info\n", helper.BytesToString(body))
		return version
	}
	return version
}
