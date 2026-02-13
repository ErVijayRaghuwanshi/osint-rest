package header

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

// WatchAndReload watches the header config file for changes and hot-reloads
// the configuration into the Manager without restarting the service.
// It runs in a blocking loop — call it in a goroutine.
// Send on stop to terminate the watcher.
func WatchAndReload(path string, m *Manager, log zerolog.Logger, stop <-chan struct{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file watcher")
		return
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to watch header config")
		return
	}

	log.Info().Str("path", path).Msg("Watching header config for changes")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				log.Info().Str("path", path).Msg("Header config changed, reloading...")

				newCfg, err := LoadHeaderConfig(path)
				if err != nil {
					log.Error().Err(err).Msg("Failed to parse updated header config, keeping old config")
					continue
				}

				m.Reload(newCfg)
				log.Info().Msg("Header config reloaded successfully")
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error().Err(err).Msg("File watcher error")

		case <-stop:
			log.Info().Msg("Stopping header config watcher")
			return
		}
	}
}
