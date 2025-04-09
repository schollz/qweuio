package cli

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"asdfgh/src/tli"

	"github.com/fsnotify/fsnotify"
	log "github.com/schollz/logger"
)

type CLI struct {
	Filename string
	TLI      *tli.TLI
}

// Init initializes the CLI with the given filename and starts watching for changes.
func Init(filename string) (cli *CLI, err error) {
	cli = &CLI{
		Filename: filename,
	}

	if err := cli.load(); err != nil {
		return cli, err
	}

	go cli.watchFile()

	return cli, nil
}

// load reads and parses the TLI file.
func (cli *CLI) load() error {

	f, err := os.Open(cli.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	tliParsed, err := tli.Parse(string(data))
	if err != nil {
		return err
	}
	cli.TLI = &tliParsed

	log.Debug("Loaded TLI from file.")
	return nil
}
func (cli *CLI) watchFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	dir := filepath.Dir(cli.Filename)
	file := filepath.Base(cli.Filename)

	err = watcher.Add(dir)
	if err != nil {
		log.Error("Error adding directory to watcher:", err)
		return
	}

	var (
		debounceTimer *time.Timer
		timerMu       sync.Mutex
	)

	triggerReload := func() {
		time.Sleep(10 * time.Millisecond) // wait for file to sync
		b, err := os.ReadFile(cli.Filename)
		if err != nil {
			log.Error("Error reading file:", err)
			return
		}
		tliNew, err := tli.Parse(string(b))
		if err != nil {
			log.Error("Error parsing TLI:", err)
			return
		}
		log.Infof("Before reload: %+v", cli.TLI.Components[0].ChainSteps[0])
		cli.TLI.Copy(tliNew)
		log.Infof("After reload: %+v", cli.TLI.Components[0].ChainSteps[0])
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if filepath.Base(event.Name) == file &&
				(event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename)) != 0 {
				log.Infof("Detected change to '%s', scheduling reload...", file)

				timerMu.Lock()
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(10*time.Millisecond, triggerReload)
				timerMu.Unlock()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Watcher error:", err)
		}
	}
}
