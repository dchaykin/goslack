package goslack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var config = slackConfig{}

const (
	info    = "INFO"
	warning = "WARNING"
	fault   = "ERROR"
)

type slackMessage struct {
	Text string `json:"text"`
}

type severityLevel string

func (sl severityLevel) isValid() bool {
	switch sl {
	case info, warning, fault:
		return true
	}
	return false
}

type sender struct {
	content slackMessage
	ConfigItem
}

func (s sender) send() error {
	data, err := json.Marshal(s.content)
	if err != nil {
		return err
	}

	payload := bytes.NewReader(data)

	req, err := http.NewRequest("POST", s.URL, payload)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error while sending a slack message: %s. Code: %d", resp.Status, resp.StatusCode)
	}

	return nil
}

type ConfigItem struct {
	Level severityLevel `json:"level"`
	URL   string        `json:"slack_url"`
	Push  *bool         `json:"push,omitempty"`
}

type slackConfig struct {
	config []ConfigItem
}

func (sc *slackConfig) getItemByLevel(level severityLevel) *ConfigItem {
	for i := range config.config {
		if config.config[i].Level == level {
			return &(config.config[i])
		}
	}
	return nil
}
