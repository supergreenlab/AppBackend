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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func PublishObject(topic string, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	r.Publish(fmt.Sprintf("pub.%s", topic), string(b))
	ps.Pub(obj, topic)
	return nil
}

func SubscribeOject(topic string) chan interface{} {
	return ps.Sub(topic)
}

func SubscribeMetric(topic string) chan int {
	ch := make(chan int)
	rps := r.Subscribe(topic)
	go func() {
		for msg := range rps.Channel() {
			v, err := strconv.Atoi(msg.Payload)
			if err != nil {
				logrus.Error(err.Error())
				continue
			}

			ch <- v
		}
		close(ch)
	}()
	return ch
}
