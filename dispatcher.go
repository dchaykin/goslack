package goslack

import "time"

var disp = dispatcher{}

func init() {
	go runDispatcher()
}

func runDispatcher() {
	for {
		disp.removeOldMessages()
		time.Sleep(15 * time.Second)
	}
}

const minSecondsBetweenMessages = 20 * 30

type dispatcher struct {
	messages []message
}

func (d *dispatcher) removeOldMessages() {
	for {
		i := d.removeOneTimedoutMessage()
		if i < 0 {
			return
		}
		d.messages[i] = d.messages[len(d.messages)-1]
		d.messages = d.messages[:len(d.messages)-1]
	}
}

func (d *dispatcher) removeOneTimedoutMessage() int {
	if d.messages == nil {
		return -1
	}

	for i, msg := range d.messages {
		if msg.isTimeout() {
			msg.setLevel(solved)
			msg.send()
			return i
		}
	}
	return -1
}

func (d *dispatcher) findSameMessage(msg message) *message {
	if d.messages == nil {
		return nil
	}

	for i := range d.messages {
		if d.messages[i].isSameAgain(msg) {
			return &d.messages[i]
		}
	}
	return nil
}

func (d *dispatcher) findMessageWithConfig(msg message) *message {
	if d.messages == nil {
		return nil
	}

	for i := range d.messages {
		if d.messages[i].hasConfig(msg.content, msg.ConfigItem) {
			return &d.messages[i]
		}
	}
	return nil
}

func (d *dispatcher) send(ci ConfigItem, content string) error {
	msg := message{
		content:                   slackMessage{},
		timestamp:                 time.Now(),
		avgSecondsBetweenMessages: 0,
		count:                     0,
	}
	msg.content.create(ci.Level, content)
	msg.Level = ci.Level
	msg.Push = ci.Push
	msg.URL = ci.URL

	previousMessage := d.findSameMessage(msg)
	if previousMessage != nil {
		previousMessage.update()
		return nil
	}

	msg2send := d.findMessageWithConfig(msg)
	if msg2send != nil {
		msg2send.update()
	} else {
		msg2send = &msg
		d.messages = append(d.messages, msg)
	}

	return msg2send.send()
}
