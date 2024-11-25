package config

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	innerLog "github.com/hyzx-go/common-b2c/log"
	"github.com/spf13/viper"
	"log"
	"path"
	"strings"
)

type DefaultParserLoader struct {
	viper  *viper.Viper
	parser *parser
}

func NewDefaultParserLoader(parser *parser) ParserLoader {
	return &DefaultParserLoader{parser: parser, viper: viper.New()}
}

func (d *DefaultParserLoader) Load() {
	if err := d.loadConfig(); err != nil {
		log.Fatal("load config failed", err.Error())
	}
	innerLog.GetLogger().Info("init parser success")
}

func (d *DefaultParserLoader) Destroy() {
	for _, fb := range d.parser.factoryBeans {
		if fb == nil {
			continue
		}

		if err := fb.Destroy(); err != nil && !errors.Is(err, ErrNotFind) {
			innerLog.GetLogger().Error("destroy failed", fb, err.Error())
		}
	}

	d.parser.beanKeys = nil
	d.parser.factoryBeans = nil
}

func (d *DefaultParserLoader) loadConfig() error {

	file := fmt.Sprintf("config-%s.yaml", d.parser.env)
	if d.parser.options.confFilepath.file != "" {
		str := strings.Split(d.parser.options.confFilepath.file, ".")
		if len(str) != 2 {
			return errors.New(fmt.Sprintf("file name error, %s", d.parser.options.confFilepath.file))
		}
		file = str[0] + "-" + d.parser.env + "." + str[1]
	}

	// Tells viper to look at the Environment Variables.
	d.viper.AutomaticEnv()
	d.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	d.viper.SetConfigFile(path.Join(d.parser.options.confFilepath.dir, file))

	// read config
	if err := d.viper.ReadInConfig(); err != nil {
		return err
	}

	// override the settings
	for k, v := range d.viper.AllSettings() {
		d.viper.Set(k, v)
	}

	// watch config
	if d.parser.options.watchConfigSwitch {
		d.viper.WatchConfig()
	}

	// change config
	d.viper.OnConfigChange(func(e fsnotify.Event) {
		// Don't do it yet
		innerLog.GetLogger().Info("On Config Changed", e)
	})

	// apply config
	if err := d.initConfig(); err != nil {
		return err
	}

	return nil
}

func (d *DefaultParserLoader) initConfig() (err error) {
	d.parser.initBeanKeys()
	for _, key := range d.parser.beanKeys {

		configBean := getBeanFactory(key)
		inConfig := d.viper.InConfig(key)

		if inConfig {
			if err := d.viper.UnmarshalKey(key, configBean); err != nil {
				return err
			}
		}

		if err := configBean.Initialize(inConfig, d.parser); err != nil {
			return err
		}

		d.parser.factoryBeans = append(d.parser.factoryBeans, configBean)
	}

	d.parser.options.mu.Lock()
	for k, v := range d.parser.options.rawVal {
		if err := d.viper.UnmarshalKey(k, v); err != nil {
			return fmt.Errorf("unmarshal key, key: %s, err: %w", k, err)
		}
	}
	d.parser.options.mu.Unlock()
	return nil
}
