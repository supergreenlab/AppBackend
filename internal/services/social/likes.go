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
	"github.com/SuperGreenLab/AppBackend/internal/services/slack"
	"github.com/sirupsen/logrus"
)

func listenLikesAdded() {
	ch := pubsub.SubscribeOject("insert.likes")
	for c := range ch {
		like := c.(middlewares.InsertMessage).Object.(*db.Like)
		if like.CommentID.Valid {
			com, err := db.GetComment(like.CommentID.UUID)
			if err != nil {
				logrus.Errorf("db.GetComment in listenLikesAdded %q - %+v", err, like)
				continue
			}

			if com.UserID == like.UserID {
				continue
			}

			plant, err := db.GetPlantForFeedEntryID(com.FeedEntryID)
			if err != nil {
				logrus.Errorf("db.GetPlantForFeedEntryID in listenLikesAdded %q - %+v", err, com)
				continue
			}

			user, err := db.GetUser(like.UserID)
			if err != nil {
				logrus.Errorf("db.GetUser listenLikesAdded %q - %+v", err, like)
				continue
			}

			title := fmt.Sprintf("%s liked your comment on the diary %s!", user.Nickname, plant.Name)
			data, notif := NewNotificationDataLikePlantComment(title, "Tap to view comment", "", plant.ID.UUID, com.FeedEntryID, like.CommentID.UUID, com.ReplyTo)
			notifications.SendNotificationToUser(com.UserID, data, &notif)
			slack.CommentLikeAdded(*like, com, plant, user)
		} else if like.FeedEntryID.Valid {
			/*feedEntry, err := db.GetFeedEntry(like.FeedEntryID.UUID)
			if err != nil {
				logrus.Errorf("listenLikesAdded db.GetFeedEntry: %q\n", err)
				continue
			}*/
			plant, err := db.GetPlantForFeedEntryID(like.FeedEntryID.UUID)
			if err != nil {
				logrus.Errorf("db.GetPlantForFeedEntryID in listenLikesAdded %q - %+v", err, like)
				continue
			}

			if plant.UserID == like.UserID {
				continue
			}

			user, err := db.GetUser(like.UserID)
			if err != nil {
				logrus.Errorf("db.GetUser in listenLikesAdded %q - %+v", err, like)
				continue
			}
			title := fmt.Sprintf("%s liked your growlog on the diary %s!", user.Nickname, plant.Name)
			data, notif := NewNotificationDataLikePlantFeedEntry(title, "Tap to view growlog", "", plant.ID.UUID, like.FeedEntryID.UUID)
			notifications.SendNotificationToUser(plant.UserID, data, &notif)
			slack.PostLikeAdded(*like, plant, user)
		}
	}
}

func initLikes() {
	go listenLikesAdded()
}
