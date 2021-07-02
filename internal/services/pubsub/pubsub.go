/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package pubsub

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func PublishObject(topic string, obj interface{}) error {
	/*b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	r.Publish(fmt.Sprintf("pub.%s", topic), string(b))*/
	ps.Pub(obj, topic)
	return nil
}

func SubscribeOject(topic string) chan interface{} {
	return ps.Sub(topic)
}

type ControllerStringMetric struct {
	Type         string `json:"type"`
	ControllerID string `json:"controllerID"`
	Key          string `json:"key"`
	Value        string `json:"value"`
}

type ControllerIntMetric struct {
	Type         string  `json:"type"`
	ControllerID string  `json:"controllerID"`
	Key          string  `json:"key"`
	Value        float64 `json:"value"`
}

type ControllerLog struct {
	Type         string `json:"type"`
	ControllerID string `json:"controllerID"`
	Module       string `json:"module"`
	Msg          string `json:"msg"`
}

func SubscribeControllerIntMetric(topic string) (chan ControllerIntMetric, chan bool) {
	stop := make(chan bool)
	ch := make(chan ControllerIntMetric, 100)
	rps := r.PSubscribe(topic)
	subCh := rps.Channel()
	go func() {
		defer close(stop)
		defer close(ch)
		defer rps.Close()
		for {
			select {
			case msg := <-subCh:
				v, err := strconv.ParseFloat(msg.Payload, 64)
				if err != nil {
					logrus.Errorf("strconv.ParseFloat in SubscribeControllerIntMetric %q - %+v", err.Error(), msg)
					continue
				}

				keyParts := strings.Split(msg.Channel, ".")

				ch <- ControllerIntMetric{ControllerID: keyParts[1], Key: keyParts[3], Value: v}
			case <-stop:
				logrus.Infof("Closing channel for %s", topic)
				return
			}
		}
	}()
	return ch, stop
}

func SubscribeControllerLogs(topic string) (chan interface{}, chan bool) {
	stop := make(chan bool)
	ch := make(chan interface{}, 100)
	rps := r.PSubscribe(topic)
	subCh := rps.Channel()
	go func() {
		defer close(stop)
		defer close(ch)
		defer rps.Close()
		for {
			select {
			case msg := <-subCh:
				keyParts := strings.Split(msg.Channel, ".")
				if len(keyParts) < 3 {
					logrus.Errorf("Unknown channel identifier: %q", msg.Channel)
					continue
				}

				if keyParts[len(keyParts)-1] == "cmd" {
					continue
				}
				if keyParts[len(keyParts)-1] == "log" {
					ch <- ControllerLog{Type: "log", ControllerID: keyParts[1], Module: keyParts[2], Msg: msg.Payload}
					continue
				}
				v, err := strconv.ParseFloat(msg.Payload, 64)
				if err != nil {
					ch <- ControllerStringMetric{Type: "string", ControllerID: keyParts[1], Key: keyParts[3], Value: msg.Payload}
				} else {
					ch <- ControllerIntMetric{Type: "int", ControllerID: keyParts[1], Key: keyParts[3], Value: v}
				}
			case <-stop:
				logrus.Infof("Closing channel for %s", topic)
				return
			}
		}
	}()
	return ch, stop
}

func PublicRemoteCmd(identifier, cmd string) {
	r.Publish(fmt.Sprintf("pub.%s.cmd", identifier), cmd)
}

func Init() {
	initRedis()
	initPubsub()
}
