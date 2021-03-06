package goslack

import (
	"fmt"
)

func AddConfig(configItem ...ConfigItem) error {
	for _, ci := range configItem {
		if err := addConfig(ci); err != nil {
			return err
		}
	}
	return nil
}

func addConfig(configItem ConfigItem) error {
	if !configItem.Level.isValid() {
		return fmt.Errorf("Unknown severity level: %s. Allowed: INFO, WARNING, ERROR", configItem.Level)
	}
	item := config.getItemByLevel(configItem.Level)
	if item == nil {
		config.config = append(config.config, configItem)
		return nil
	}
	item.Push = configItem.Push
	item.URL = configItem.URL
	return nil
}

func send(content string, level severityLevel) error {
	configItem := config.getItemByLevel(level)
	if configItem == nil {
		return fmt.Errorf("No slack config for %s found", level)
	}
	return disp.send(*configItem, content)
}

func Infof(text string, params ...interface{}) error {
	content := fmt.Sprintf(text, params...)
	return send(content, info)
}

func Warningf(text string, params ...interface{}) error {
	content := fmt.Sprintf(text, params...)
	return send(content, warning)
}

func Errorf(text string, params ...interface{}) error {
	content := fmt.Sprintf(text, params...)
	return send(content, fault)
}
