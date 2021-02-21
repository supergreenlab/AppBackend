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
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/api/option"
)

var (
	client        *firebase.App
	ch            chan NotificationObject
	fcmConfigPath = pflag.String("fcmconfigpath", "/etc/appbackend/fcmconfig.json", "Url to the redis instance")
)

type NotificationObject struct {
	user         db.User
	data         map[string]string
	notification *messaging.Notification
}

func SendNotificationToUser(user db.User, data map[string]string, notification *messaging.Notification) {
	logrus.Infof("Sending notification %s", user.Nickname)
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
