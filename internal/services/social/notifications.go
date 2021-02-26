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
	"github.com/SuperGreenLab/AppBackend/internal/services/notifications"
	"github.com/gofrs/uuid"
)

var (
	NotificationTypePlantComment       = "PLANT_COMMENT"
	NotificationTypePlantCommentReply  = "PLANT_COMMENT_REPLY"
	NotificationTypeReminder           = "REMINDER"
	NotificationTypeAlert              = "ALERT"
	NotificationTypeLikePlantComment   = "LIKE_PLANT_COMMENT"
	NotificationTypeLikePlantFeedEntry = "LIKE_PLANT_FEEDENTRY"
)

type NotificationDataPlantComment struct {
	notifications.NotificationBaseData

	PlantID     uuid.UUID `json:"plantID"`
	FeedEntryID uuid.UUID `json:"feedEntryID"`
	CommentType string    `json:"commentType"`
}

func (n NotificationDataPlantComment) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return n.Merge(m, map[string]string{
		"plantID":     n.PlantID.String(),
		"feedEntryID": n.FeedEntryID.String(),
		"commentType": n.CommentType,
	})
}

func NewNotificationDataPlantComment(title, body, imageUrl string, plantID, feedEntryID uuid.UUID, commentType string) (NotificationDataPlantComment, messaging.Notification) {
	return NotificationDataPlantComment{
			NotificationBaseData: notifications.NotificationBaseData{
				Type:  NotificationTypePlantComment,
				Title: title,
				Body:  body,
			},
			PlantID:     plantID,
			FeedEntryID: feedEntryID,
			CommentType: commentType,
		},
		messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageUrl,
		}
}

type NotificationDataPlantCommentReply struct {
	notifications.NotificationBaseData

	PlantID     uuid.UUID `json:"plantID"`
	FeedEntryID uuid.UUID `json:"feedEntryID"`
	CommentID   uuid.UUID `json:"commentID"`
}

func (n NotificationDataPlantCommentReply) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return n.Merge(m, map[string]string{
		"plantID":     n.PlantID.String(),
		"feedEntryID": n.FeedEntryID.String(),
		"commentID":   n.CommentID.String(),
	})
}

func NewNotificationDataPlantCommentReply(title, body, imageUrl string, plantID, feedEntryID uuid.UUID, commentID uuid.UUID) (NotificationDataPlantCommentReply, messaging.Notification) {
	return NotificationDataPlantCommentReply{
			NotificationBaseData: notifications.NotificationBaseData{
				Type:  NotificationTypePlantCommentReply,
				Title: title,
				Body:  body,
			},
			PlantID:     plantID,
			FeedEntryID: feedEntryID,
			CommentID:   commentID,
		},
		messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageUrl,
		}
}

type NotificationDataReminder struct {
	notifications.NotificationBaseData

	PlantID uuid.UUID `json:"plantID"`
}

func (n NotificationDataReminder) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return n.Merge(m, map[string]string{
		"plantID": n.PlantID.String(),
	})
}

type NotificationDataLikePlantComment struct {
	notifications.NotificationBaseData

	PlantID     uuid.UUID     `json:"plantID"`
	FeedEntryID uuid.UUID     `json:"feedEntryID"`
	CommentID   uuid.UUID     `json:"commentID"`
	ReplyTo     uuid.NullUUID `json:"replyTo"`
}

func (n NotificationDataLikePlantComment) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	m2 := map[string]string{
		"plantID":     n.PlantID.String(),
		"feedEntryID": n.FeedEntryID.String(),
		"commentID":   n.CommentID.String(),
	}
	if n.ReplyTo.Valid {
		m2["replyTo"] = n.ReplyTo.UUID.String()
	}
	return n.Merge(m, m2)
}

func NewNotificationDataLikePlantComment(title, body, imageUrl string, plantID, feedEntryID uuid.UUID, commentID uuid.UUID, replyTo uuid.NullUUID) (NotificationDataLikePlantComment, messaging.Notification) {
	return NotificationDataLikePlantComment{
			NotificationBaseData: notifications.NotificationBaseData{
				Type:  NotificationTypeLikePlantComment,
				Title: title,
				Body:  body,
			},
			PlantID:     plantID,
			FeedEntryID: feedEntryID,
			CommentID:   commentID,
			ReplyTo:     replyTo,
		},
		messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageUrl,
		}
}

type NotificationDataLikePlantFeedEntry struct {
	notifications.NotificationBaseData

	PlantID     uuid.UUID `json:"plantID"`
	FeedEntryID uuid.UUID `json:"feedEntryID"`
}

func (n NotificationDataLikePlantFeedEntry) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return n.Merge(m, map[string]string{
		"plantID":     n.PlantID.String(),
		"feedEntryID": n.FeedEntryID.String(),
	})
}

func NewNotificationDataLikePlantFeedEntry(title, body, imageUrl string, plantID, feedEntryID uuid.UUID) (NotificationDataLikePlantFeedEntry, messaging.Notification) {
	return NotificationDataLikePlantFeedEntry{
			NotificationBaseData: notifications.NotificationBaseData{
				Type:  NotificationTypeLikePlantFeedEntry,
				Title: title,
				Body:  body,
			},
			PlantID:     plantID,
			FeedEntryID: feedEntryID,
		},
		messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageUrl,
		}
}
