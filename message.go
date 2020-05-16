package goslack

import (
	"fmt"
)

// AddConfig appends a new configuration for sending slack messages
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
	item := config.getItem(configItem.Level, configItem.URL)
	if item == nil {
		config.config = append(config.config, configItem)
		return nil
	}
	item.Push = configItem.Push
	return nil
}

func send(content string, level severityLevel) error {
	ciList := config.getItemsByLevel(level)
	if ciList == nil {
		return fmt.Errorf("No slack config for %s found", level)
	}
	var result error
	for i := range ciList {
		if err := disp.send(ciList[i], content); err != nil {
			if result == nil {
				result = err
				continue
			}
			result = fmt.Errorf("%v\n%v", result, err)
		}
	}
	return result
}

// Infof sends a text to all slack channels added with AddConfig()
// and severity level "INFO". This messages have the green colored bar
// on the left side
func Infof(text string, params ...interface{}) error {
	content := fmt.Sprintf(text, params...)
	return send(content, info)
}

// Warningf sends a text to all slack channels added with AddConfig()
// and severity level "WARNING". This messages have the yellow colored bar
// on the left side
func Warningf(text string, params ...interface{}) error {
	content := fmt.Sprintf(text, params...)
	return send(content, warning)
}

// Errorf sends a text to all slack channels added with AddConfig()
// and severity level "ERROR". This messages have the red colored bar
// on the left side
func Errorf(text string, params ...interface{}) error {
	content := fmt.Sprintf(text, params...)
	return send(content, fault)
}
