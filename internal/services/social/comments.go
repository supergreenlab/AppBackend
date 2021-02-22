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

package social

import (
	"firebase.google.com/go/v4/messaging"
	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/notifications"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/sirupsen/logrus"
)

func listenCommentsAdded() {
	ch := pubsub.SubscribeOject("insert.comments")
	for c := range ch {
		com := c.(middlewares.InsertMessage).Object.(*db.Comment)
		if feedEntry, err := db.GetFeedEntry(com.FeedEntryID); err == nil {
			if com.UserID != feedEntry.UserID {
				notifications.SendNotificationToUser(feedEntry.UserID, NotificationDataPlantComment{}, &messaging.Notification{})
			}
		} else {
			logrus.Errorf("listenCommentsAdded: %q\n", err)
			continue
		}

		if com.ReplyTo.Valid {
			if comReplied, err := db.GetComment(com.ReplyTo.UUID); err == nil {
				if com.UserID != comReplied.UserID {
					notifications.SendNotificationToUser(comReplied.UserID, NotificationDataPlantComment{}, &messaging.Notification{})
				}
			} else {
				logrus.Errorf("listenCommentsAdded: %q\n", err)
			}
		}
	}
}

func initComments() {
	go listenCommentsAdded()
}
