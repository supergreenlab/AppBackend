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
	"regexp"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/services/notifications"
	"github.com/SuperGreenLab/AppBackend/internal/services/pubsub"
	"github.com/SuperGreenLab/AppBackend/internal/services/slack"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

var (
	mentionRegexp = regexp.MustCompile(`@([a-zA-Z0-9_-]*)`)
)

func listenCommentsAdded() {
	ch := pubsub.SubscribeOject("insert.comments")
	for c := range ch {
		com := c.(middlewares.InsertMessage).Object.(*db.Comment)
		id := c.(middlewares.InsertMessage).ID
		feedEntry, err := db.GetFeedEntry(com.FeedEntryID)
		if err != nil {
			logrus.Errorf("db.GetFeedEntry in listenCommentsAdded %q - %+v", err, com)
			continue
		}
		plant, err := db.GetPlantForFeedEntryID(com.FeedEntryID)
		if err != nil {
			logrus.Errorf("db.GetPlantForFeedEntryID in listenCommentsAdded %q - %+v", err, com)
			continue
		}
		user, err := db.GetUser(com.UserID)
		if err != nil {
			logrus.Errorf("db.GetUser in listenCommentsAdded %q - %+v", err, com)
			continue
		}

		var userIDNotif uuid.UUID
		if com.ReplyTo.Valid {
			comReplied, err := db.GetComment(com.ReplyTo.UUID)
			if err != nil {
				logrus.Errorf("db.GetComment in listenCommentsAdded %q - %+v", err, com)
			}
			if com.UserID != comReplied.UserID {
				title := fmt.Sprintf("%s replied to your comment on the diary %s!", user.Nickname, plant.Name)
				data, notif := NewNotificationDataPlantCommentReply(title, com.Text, "", plant.ID.UUID, feedEntry.ID.UUID, comReplied.ID.UUID)
				notifications.SendNotificationToUser(comReplied.UserID, data, &notif)
				userIDNotif = comReplied.UserID
			}
		} else if com.UserID != feedEntry.UserID {
			title := fmt.Sprintf("%s posted a message on your diary %s!", user.Nickname, plant.Name)
			data, notif := NewNotificationDataPlantComment(title, com.Text, "", plant.ID.UUID, feedEntry.ID.UUID, com.Type)
			notifications.SendNotificationToUser(feedEntry.UserID, data, &notif)
			userIDNotif = feedEntry.UserID
		}

		mentions := mentionRegexp.FindAllStringSubmatch(com.Text, -1)
		for _, m := range mentions {
			userMentionned, err := db.GetUserForNickname(m[1])
			if err != nil {
				logrus.Errorf("db.GetUserForNickname in listenCommentsAdded %q - %+v", err, m)
				continue
			}
			if userMentionned.ID.UUID == userIDNotif {
				continue
			}
			title := fmt.Sprintf("%s mentionned you in a comment on the diary %s!", user.Nickname, plant.Name)
			comID := id
			if com.ReplyTo.Valid {
				comID = com.ReplyTo.UUID
			}
			data, notif := NewNotificationDataPlantCommentReply(title, com.Text, "", plant.ID.UUID, feedEntry.ID.UUID, comID)
			notifications.SendNotificationToUser(userMentionned.ID.UUID, data, &notif)
		}
		slack.CommentPosted(id, *com, plant, user)
	}
}

func initComments() {
	go listenCommentsAdded()
}
