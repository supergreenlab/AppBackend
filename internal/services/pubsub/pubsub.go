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

type ControllerIntMetric struct {
	ControllerID string
	Key          string
	Value        float64
}

func SubscribeControllerIntMetric(topic string) chan ControllerIntMetric {
	ch := make(chan ControllerIntMetric, 100)
	rps := r.PSubscribe(topic)
	go func() {
		for msg := range rps.Channel() {
			v, err := strconv.ParseFloat(msg.Payload, 64)
			if err != nil {
				logrus.Errorf("strconv.ParseFloat in SubscribeControllerIntMetric %q - %+v", err.Error(), msg)
				continue
			}

			keyParts := strings.Split(msg.Channel, ".")

			ch <- ControllerIntMetric{ControllerID: keyParts[1], Key: keyParts[3], Value: v}
		}
		close(ch)
	}()
	return ch
}

func Init() {
	initRedis()
	initPubsub()
}
