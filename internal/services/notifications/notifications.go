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

package notifications

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/services/prometheus"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	client *firebase.App
	ch     chan UserNotification
	_      = pflag.String("fcmconfigpath", "/etc/appbackend/fcmconfig.json", "Url to the redis instance")
)

type NotificationData interface {
	ToMap() map[string]string
	GetType() string
}

type NotificationBaseData struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (n NotificationBaseData) GetType() string {
	return n.Type
}

func (n NotificationBaseData) ToMap() map[string]string {
	return map[string]string{
		"type":  n.Type,
		"title": n.Title,
		"body":  n.Body,
	}
}

type UserNotification struct {
	userID       uuid.UUID
	data         NotificationData
	notification *messaging.Notification
}

func handleUserNotifications() {
	for un := range ch {
		userends, err := db.GetUserEndsForUserID(un.userID)
		if err != nil {
			logrus.Errorf("db.GetUserEndsForUserID in handleUserNotifications %q - %+v", err, un)
			return
		}
		cli, err := client.Messaging(context.Background())
		if err != nil {
			logrus.Errorf("client.Messaging in handleUserNotifications %q", err)
			return
		}
		tokensMap := map[string]bool{}
		for _, userend := range userends {
			if userend.NotificationToken.Valid && userend.NotificationToken.String != "" {
				tokensMap[userend.NotificationToken.String] = true
			}
		}
		tokens := []string{}
		for k := range tokensMap {
			tokens = append(tokens, k)
		}
		if len(tokens) > 0 {
			logrus.Infof("Sending notification to %q\n", tokens)
			msg := &messaging.MulticastMessage{Data: un.data.ToMap(), Notification: un.notification, Tokens: tokens}
			prometheus.NotificationSent(un.data.GetType())
			if _, err := cli.SendMulticast(context.Background(), msg); err != nil {
				prometheus.NotificationError(un.data.GetType())
				logrus.Errorf("cli.SendMulticast in handleUserNotifications %q - %+v", err, un)
			}
		}
	}
}

func SendNotificationToUser(userID uuid.UUID, data NotificationData, notification *messaging.Notification) {
	ch <- UserNotification{userID, data, notification}
}

func Init() {
	var err error
	ctx := context.Background()
	config := &firebase.Config{ProjectID: "supergreenlab-6cd05"}
	opt := option.WithCredentialsFile(viper.GetString("FCMConfigPath"))
	client, err = firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalf("firebase.NewApp in Init %q", err)
	}

	ch = make(chan UserNotification, 100)
	go handleUserNotifications()
}
