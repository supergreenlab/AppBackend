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

import "github.com/SuperGreenLab/AppBackend/internal/services/notifications"

func merge(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

type NotificationDataPlantComment struct {
	notifications.NotificationBaseData

	PlantID     string `json:"plantID"`
	FeedEntryID string `json:"feedEntryID"`
	CommentType string `json:"commentType"`
}

func (n NotificationDataPlantComment) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return merge(m, map[string]string{
		"plantID":     n.PlantID,
		"feedEntryID": n.FeedEntryID,
		"commentType": n.CommentType,
	})
}

type NotificationDataReminder struct {
	notifications.NotificationBaseData

	PlantID string `json:"plantID"`
}

func (n NotificationDataReminder) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return merge(m, map[string]string{
		"plantID": n.PlantID,
	})
}

type NotificationDataAlert struct {
	notifications.NotificationBaseData

	PlantID string `json:"plantID"`
}

func (n NotificationDataAlert) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return merge(m, map[string]string{
		"plantID": n.PlantID,
	})
}

type NotificationDataLikePlantComment struct {
	notifications.NotificationBaseData

	PlantID     string `json:"plantID"`
	FeedEntryID string `json:"feedEntryID"`
	CommentID   string `json:"commentID"`
	ReplyTo     string `json:"replyTo"`
}

func (n NotificationDataLikePlantComment) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return merge(m, map[string]string{
		"plantID":     n.PlantID,
		"feedEntryID": n.FeedEntryID,
		"commentID":   n.CommentID,
		"replyTo":     n.ReplyTo,
	})
}

type NotificationDataLikePlantFeedEntry struct {
	notifications.NotificationBaseData

	PlantID     string `json:"plantID"`
	FeedEntryID string `json:"feedEntryID"`
}

func (n NotificationDataLikePlantFeedEntry) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return merge(m, map[string]string{
		"plantID":     n.PlantID,
		"feedEntryID": n.FeedEntryID,
	})
}
