package config

import "path"

func SetConfigFilePath(p string) Option {
	return func(o *Options) {
		o.confFilepath.dir, o.confFilepath.file = path.Split(p)
	}
}

func SetWatchConfigSwitch(on bool) Option {
	return func(o *Options) {
		o.watchConfigSwitch = on
	}
}

// SetRawVal Register the configuration structure.
// out must be pointer.
func SetRawVal(key string, out interface{}) Option {
	return func(o *Options) {
		o.mu.Lock()
		if o.rawVal == nil {
			o.rawVal = make(map[string]interface{})
		}
		o.rawVal[key] = out
		o.mu.Unlock()
	}
}
