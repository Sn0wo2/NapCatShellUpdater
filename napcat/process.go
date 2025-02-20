package napcat

import (
	"github.com/Sn0wo2/NapCatShellUpdater/flags"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"github.com/shirou/gopsutil/process"

	"os"
	"path/filepath"
	"time"
)

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

func WaitForProcess(targetPath string) (<-chan *process.Process, <-chan error) {
	resultChan := make(chan *process.Process)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		initialProcs, err := process.Processes()
		if err != nil {
			errorChan <- err
			return
		}

		initialProcIDs := make(map[int32]struct{})
		for _, p := range initialProcs {
			initialProcIDs[p.Pid] = struct{}{}
		}
		for {
			time.Sleep(2 * time.Second)

			procs, err := process.Processes()
			if err != nil {
				errorChan <- err
				return
			}

			for _, proc := range procs {
				if _, existed := initialProcIDs[proc.Pid]; !existed {
					path, err := proc.Exe()
					if err != nil {
						continue
					}
					if path == targetPath {
						resultChan <- proc
						return
					}
				}
			}
		}
	}()

	return resultChan, errorChan
}

func WaitForAllProcessesEnd(target string, abs bool) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		for {
			allProcs, err := process.Processes()
			if err != nil {
				errorChan <- err
				return
			}

			activeProcs := 0
			for _, proc := range allProcs {
				exePath, err := proc.Exe()
				if err != nil {
					continue
				}

				if abs {
					if exePath == target {
						activeProcs++
					}
				} else {
					if name, err := proc.Name(); err != nil && name == target {
						activeProcs++
					}
				}
			}

			if activeProcs == 0 {
				return
			}

			time.Sleep(2 * time.Second)
		}
	}()

	return errorChan
}
