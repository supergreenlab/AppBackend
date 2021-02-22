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
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/api/option"
)

var (
	client        *firebase.App
	ch            chan NotificationObject
	fcmConfigPath = pflag.String("fcmconfigpath", "/etc/appbackend/fcmconfig.json", "Url to the redis instance")
)

type NotificationData interface {
	ToMap() map[string]string
}

type NotificationBaseData struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (n NotificationBaseData) ToMap() map[string]string {
	return map[string]string{
		"type":  n.Type,
		"title": n.Title,
		"body":  n.Body,
	}
}

type NotificationObject struct {
	user         db.User
	data         NotificationData
	notification *messaging.Notification
}

func SendNotificationToUser(userID uuid.UUID, data NotificationData, notification *messaging.Notification) {
	userends, err := db.GetUserEndsForUserID(userID)
	if err != nil {
		logrus.Errorf("SendNotificationToUser: %q\n", err)
		return
	}
	cli, err := client.Messaging(context.Background())
	if err != nil {
		logrus.Errorf("SendNotificationToUser: %q\n", err)
		return
	}
	for _, userend := range userends {
		if userend.NotificationToken.Valid && userend.NotificationToken.String != "" {
			logrus.Infof("Sending notification to %s\n", userend.NotificationToken.String)
			msg := &messaging.Message{Data: data.ToMap(), Notification: notification, Token: userend.NotificationToken.String}
			if str, err := cli.Send(context.Background(), msg); err != nil {
				logrus.Errorf("cli.Send: %q\n", err)
			} else {
				logrus.Info(str)
			}
		}
	}
}

func Init() {
	var err error
	ctx := context.Background()
	config := &firebase.Config{ProjectID: "supergreenlab-6cd05"}
	opt := option.WithCredentialsFile(*fcmConfigPath)
	client, err = firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatal(err)
	}
}
