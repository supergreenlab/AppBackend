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
	"fmt"

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
		feedEntry, err := db.GetFeedEntry(com.FeedEntryID)
		if err != nil {
			logrus.Errorf("listenCommentsAdded db.GetFeedEntry: %q\n", err)
			continue
		}
		plant, err := db.GetPlantForFeedEntryID(com.FeedEntryID)
		if err != nil {
			logrus.Errorf("listenCommentsAdded db.GetPlantForFeedEntryID: %q\n", err)
			continue
		}
		user, err := db.GetUser(com.UserID)
		if err != nil {
			logrus.Errorf("listenCommentsAdded db.GetUser: %q\n", err)
			continue
		}

		if com.ReplyTo.Valid {
			comReplied, err := db.GetComment(com.ReplyTo.UUID)
			if err != nil {
				logrus.Errorf("listenCommentsAdded db.GetComment: %q\n", err)
			}
			if com.UserID != comReplied.UserID {
				title := fmt.Sprintf("%s replied to your comment on the diary %s!", user.Nickname, plant.Name)
				data, notif := NewNotificationDataPlantCommentReply(title, com.Text, "", plant.ID.UUID, feedEntry.ID.UUID, comReplied.ID.UUID)
				notifications.SendNotificationToUser(comReplied.UserID, data, &notif)
			}
		} else if com.UserID != feedEntry.UserID {
			title := fmt.Sprintf("%s posted a message on your diary %s!", user.Nickname, plant.Name)
			data, notif := NewNotificationDataPlantComment(title, com.Text, "", plant.ID.UUID, feedEntry.ID.UUID, com.Type)
			notifications.SendNotificationToUser(feedEntry.UserID, data, &notif)
		}
	}
}

func initComments() {
	go listenCommentsAdded()
}
