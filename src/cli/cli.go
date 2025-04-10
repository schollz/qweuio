package cli

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"asdfgh/src/constants"
	"asdfgh/src/tli"

	"github.com/fsnotify/fsnotify"
	log "github.com/schollz/logger"
)

type CLI struct {
	Filename string
	TLI      []*tli.TLI
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

	// split the file when the line equals SYMBOL_BREAK
	var pieces []string
	var piece string
	for _, line := range strings.Split(string(data), "\n") {
		if line == constants.SYMBOL_BREAK {
			piece = strings.TrimSpace(piece)
			if piece != "" {
				pieces = append(pieces, piece)
			}
			piece = ""
		} else {
			piece += string(line) + "\n"
		}
	}
	if piece != "" {
		pieces = append(pieces, piece)
	}

	tlis := make([]*tli.TLI, len(pieces))
	for i, piece := range pieces {
		tliParsed, err := tli.Parse(piece)
		if err != nil {
			return err
		}
		tlis[i] = &tliParsed
	}

	if len(tlis) == len(cli.TLI) {
		// do copying
		log.Debugf("before copy: %+v", cli.TLI[0].Components[0])
		for i := range tlis {
			cli.TLI[i].Copy(*tlis[i])
		}
		log.Debugf("after copy: %+v", cli.TLI[0].Components[0])
	} else {
		isPlaying := len(cli.TLI) > 0 && cli.TLI[0].IsPlaying()
		if len(cli.TLI) > 0 && isPlaying {
			cli.Stop()
		}
		cli.TLI = make([]*tli.TLI, len(tlis))
		for i, tliParsed := range tlis {
			cli.TLI[i] = tliParsed
		}
		if isPlaying {
			cli.Play()
		}
	}

	log.Debug("Loaded TLI from file.")
	return nil
}

func (cli *CLI) Stop() (err error) {
	for _, tli := range cli.TLI {
		tli.Stop()
	}
	return
}

func (cli *CLI) Play() (err error) {
	for _, tli := range cli.TLI {
		if err = tli.Play(); err != nil {
			log.Error("Error playing TLI:", err)
			return
		}
	}
	return
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
		cli.load()
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
