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

package alerts

import (
	"firebase.google.com/go/v4/messaging"
	"github.com/SuperGreenLab/AppBackend/internal/services/notifications"
	"github.com/gofrs/uuid"
)

var (
	NotificationTypeAlert    = "ALERT"
	NotificationTypeReminder = "REMINDER"
)

type NotificationDataAlert struct {
	notifications.NotificationBaseData

	PlantID uuid.UUID `json:"plantID"`
}

func (n NotificationDataAlert) ToMap() map[string]string {
	m := n.NotificationBaseData.ToMap()
	return n.Merge(m, map[string]string{
		"plantID": n.PlantID.String(),
	})
}

func NewNotificationDataAlert(title, body, imageUrl string, plantID uuid.UUID) (NotificationDataAlert, messaging.Notification) {
	return NotificationDataAlert{
			NotificationBaseData: notifications.NotificationBaseData{
				Type:  NotificationTypeAlert,
				Title: title,
				Body:  body,
			},
			PlantID: plantID,
		}, messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageUrl,
		}
}
