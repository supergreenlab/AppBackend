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

package db

import (
	"time"

	"github.com/gofrs/uuid"
)

func GetFeedEntry(feedEntryID uuid.UUID) (FeedEntry, error) {
	feedEntry := FeedEntry{}
	err := GetObjectWithID(feedEntryID, "feedentries", &feedEntry)
	return feedEntry, err
}

func GetFeedEntriesBetweenDates(from, to time.Time) ([]FeedEntry, error) {
	feedEntries := []FeedEntry{}
	selector := Sess.Select("*").From("feedentries").Where("cat >= ?", from).And("cat <= ?", to)
	if err := selector.All(&feedEntries); err != nil {
		return feedEntries, err
	}
	return feedEntries, nil
}

func GetComment(commentID uuid.UUID) (Comment, error) {
	comment := Comment{}
	err := GetObjectWithID(commentID, "comments", &comment)
	return comment, err
}

func GetBox(boxID uuid.UUID) (Box, error) {
	box := Box{}
	err := GetObjectWithID(boxID, "boxes", &box)
	return box, err
}

func GetBoxFromPlantFeed(feedID uuid.UUID) (Box, error) {
	box := Box{}
	selector := Sess.Select("boxes.*").From("boxes").Join("plants").On("plants.boxid = boxes.id").Where("plants.feedid = ?", feedID)
	if err := selector.One(&box); err != nil {
		return box, err
	}
	return box, nil
}

func GetDevice(deviceID uuid.UUID) (Device, error) {
	device := Device{}
	err := GetObjectWithID(deviceID, "devices", &device)
	return device, err
}

func GetDeviceFromPlantFeed(feedID uuid.UUID) (Device, error) {
	device := Device{}
	selector := Sess.Select("device.*").From("devices").Join("boxes").On("boxes.deviceid = devices.id").Join("plants").On("plants.boxid = boxes.id").Where("plants.feedid = ?", feedID)
	if err := selector.One(&selector); err != nil {
		return device, err
	}
	return device, nil
}

func GetUserEndsForUserID(userID uuid.UUID) ([]UserEnd, error) {
	userends := []UserEnd{}

	err := GetObjectsWithField("userid", userID, "userends", &userends)
	return userends, err
}

func GetPlantForFeedID(feedID uuid.UUID) (Plant, error) {
	plant := Plant{}
	err := GetObjectWithField("feedid", feedID, "plants", &plant)
	return plant, err
}

func GetPlantForFeedEntryID(feedEntryID uuid.UUID) (Plant, error) {
	plant := Plant{}
	selector := Sess.Select("plants.*").From("plants").Join("feedentries").On("plants.feedid = feedentries.feedid").Where("feedentries.id = ?", feedEntryID)
	if err := selector.One(&plant); err != nil {
		return plant, err
	}

	return plant, nil
}

func GetActivePlantsForControllerIdentifier(controllerID string, boxSlotID int) ([]Plant, error) {
	plants := []Plant{}
	selector := Sess.Select("plants.*").From("plants").Join("boxes").On("boxes.id = plants.boxid").Join("devices").On("devices.id = boxes.deviceid").Where("devices.identifier = ?", controllerID).And("boxes.devicebox = ?", boxSlotID).And("plants.deleted = false").And("plants.archived = false").And("devices.deleted = false")
	if err := selector.All(&plants); err != nil {
		return plants, err
	}

	return plants, nil

}
