package goslack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type severityLevel string

var config = slackConfig{}

const (
	info    = "INFO"
	warning = "WARNING"
	fault   = "ERROR"
	solved  = "SOLVED"
)

type slackAttachment struct {
	MrkDown []string `json:"mrkdwn_in"`
	Text    string   `json:"text"`
	Color   string   `json:"color"`
	Title   string   `json:"pretext"`
}

type slackMessage struct {
	Attachments []slackAttachment `json:"attachments"`
}

func (sm *slackMessage) create(sl severityLevel, content string) {
	sm.Attachments = []slackAttachment{
		{
			MrkDown: []string{"text"},
			Color:   sm.getColorBySeveretyLevel(sl),
			Text:    content,
			Title:   sm.getTitleBySeveretyLevel(sl),
		},
	}
}

func (sm *slackMessage) getText() string {
	if len(sm.Attachments) == 0 {
		return ""
	}
	return sm.Attachments[0].Text
}

func (sm *slackMessage) getColorBySeveretyLevel(sl severityLevel) string {
	switch sl {
	case info, solved:
		return "good"
	case warning:
		return "warning"
	case fault:
		return "danger"
	}
	return ""
}

func (sm *slackMessage) getTitleBySeveretyLevel(sl severityLevel) string {
	switch sl {
	case info:
		return "Info:"
	case warning:
		return "Warning:"
	case fault:
		return "Error:"
	case solved:
		return "Solved:"
	}
	return ""
}

func (sl severityLevel) isValid() bool {
	switch sl {
	case info, warning, fault:
		return true
	}
	return false
}

type message struct {
	content                   slackMessage
	timestamp                 time.Time
	avgSecondsBetweenMessages int
	count                     int
	ConfigItem
}

func (m *message) isTimeout() bool {
	secondsBeetweenMessages := int(time.Now().Unix() - m.timestamp.Unix())
	if m.avgSecondsBetweenMessages == 0 && secondsBeetweenMessages > minSecondsBetweenMessages {
		return true
	}
	if m.avgSecondsBetweenMessages > 0 && secondsBeetweenMessages+15 > m.avgSecondsBetweenMessages {
		return true
	}
	return false
}

func (m *message) update() {
	if m.avgSecondsBetweenMessages == 0 {
		m.avgSecondsBetweenMessages = int(time.Now().Sub(m.timestamp))
	} else {
		m.avgSecondsBetweenMessages = (m.avgSecondsBetweenMessages*m.count + int(time.Now().Sub(m.timestamp))) / (m.count + 1)
	}
	m.count++
	m.timestamp = time.Now()
}

func (m *message) hasConfig(content slackMessage, cfg ConfigItem) bool {
	return (m.content.getText() == content.getText() &&
		m.Level == cfg.Level &&
		m.URL == cfg.URL)
}

func (m *message) isSameAgain(msg message) bool {
	if !m.hasConfig(msg.content, msg.ConfigItem) {
		return false
	}

	secondsBeetweenMessages := int(m.timestamp.Unix() - msg.timestamp.Unix())
	isSecondMessageAndSendAgain := (m.avgSecondsBetweenMessages == 0 && secondsBeetweenMessages > minSecondsBetweenMessages)
	isRepeatedMessageAndSendAgain := (m.avgSecondsBetweenMessages > 0 && secondsBeetweenMessages > (m.avgSecondsBetweenMessages+10))
	return isSecondMessageAndSendAgain || isRepeatedMessageAndSendAgain
}

func (m *message) send() error {
	data, err := json.Marshal(m.content)
	if err != nil {
		return err
	}

	payload := bytes.NewReader(data)

	req, err := http.NewRequest("POST", m.URL, payload)
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
