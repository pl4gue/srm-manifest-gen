package srm

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

// Maps a directory and all of it's child executables
type DirExecs map[string][]string

// Searches for empty arrays and empty strings on the arrays that aren't empty and removes them from
// the map
func (self DirExecs) Clean() {
	for dir, execs := range self {
		var _new []string

		// If the array is empty deletes it's entry
		if len(execs) == 0 {
			delete(self, dir)
			continue
		}

		// If removes empty strings from the array
		for _, s := range execs {
			if s != "" {
				_new = append(_new, s)
			}
		}

		self[dir] = _new
	}
}

type App struct {
	Root_dir          string
	found_executables DirExecs
	chosen_execs      Games
}

// Prompts the user for the root directory to search for
func (a *App) PromptRoot() error {
	rd, err := Input(
		"Enter the directory to search for: ",
		func(s string) error {
			info, err := os.Stat(s)
			if err != nil {
				log.With("err", err).Debugf("Encountered an error while checking if %s is a directory", s)
				return err
			}

			if !info.IsDir() {
				err := fmt.Errorf("Given path %s is not a directory", s)
				log.Debug(err.Error())
				return err
			}

			return nil
		})
	if err != nil {
		return err
	}

	a.Root_dir = rd
	log.With("root_dir", a.Root_dir).Infof("App root dir was set.")
	return nil
}

// Gets the directories right under the root directory as the keys of a map
func (a *App) get_games_as_keys() (DirExecs, error) {
	parents, err := os.ReadDir(a.Root_dir)
	if err != nil {
		log.With("err", err).Debugf("Couldn't read %s", a.Root_dir)
		return nil, err
	}

	dir_execs := make(DirExecs)

	for _, parent := range parents {
		if !parent.IsDir() {
			continue
		}

		log.Debugf("Found game %s in root directory.", parent.Name())
		dir_execs[parent.Name()] = []string{}
	}

	return dir_execs, nil
}

// Populates the app with directories and executables
func (a *App) Populate() error {
	game_execs, err := a.get_games_as_keys()
	if err != nil {
		return err
	}

	var tasks []func() error
	var mu sync.Mutex

	// Creates the tasks that get the executables for each game and add them in the
	for game := range game_execs {
		tasks = append(tasks, func() error {
			log.Debugf("Getting executables for %s", game)
			found_execs, err := GetExecsIn(filepath.Join(a.Root_dir, game))

			// If no executables are returned
			// Should return the possible cause right away
			if len(found_execs) == 0 {
				log.With("err", err).Debugf("Found no executables for %s", game)
				return err
			}

			mu.Lock()
			game_execs[game] = found_execs
			mu.Unlock()

			return err
		})
	}

	// Use concurrency to get executables
	err = Spin(
		"Getting executables...",
		func() {
			// Runs tasks concurrently and prints any error found
			errs := RunConcurrently(tasks...)

			for i, err := range errs {
				log.Debug("Failed task:", "task_id", i, "err", err)
			}
		},
	)

	game_execs.Clean()
	a.found_executables = game_execs

	return err
}

func (a *App) ChooseExecutables() error {
	log.Debug("Selecting executables.")

	chosen_execs := make(Games)

	for game, execs := range a.found_executables {
		log.Debugf("Selecting executables for %s", game)
		if len(execs) == 0 {
			continue
		}

		// If only one executable is found prompts the user to confirm if they wants to use that
		// executable
		if len(execs) == 1 {
			confirm, err := Confirm("Is " + execs[0] + " the correct executable for " + game + "?")
			if err != nil {
				log.Debugf("Found an error while confirming %s's executable.", game)
				return err
			}

			if confirm {
				chosen_execs[game] = execs[0]
			}
			continue
		}

		// Creates the options to select and marks the first as selected,
		// also adds a "None" option to skip the game
		options := append(huh.NewOptions(execs...), huh.NewOption("None", ""))
		for i, op := range options {
			options[i] = op.Selected(i == 0)
		}

		selected, err := Select("Found multiple executables in "+game, options)
		if err != nil {
			log.With("selected", selected, "err", err).Debugf("Found an error while selecting executables for %s.", game)
			return err
		}

		// If the option "None" is selected then don't add the game in chosen_execs
		if selected == "" {
			continue
		}

		chosen_execs[game] = selected
	}

	a.chosen_execs = chosen_execs

	return nil
}

func (a *App) JSONify() ([]byte, error) {
	return a.chosen_execs.JSONify()
}
