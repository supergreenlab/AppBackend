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

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
)

func GetFeedEntry(feedEntryID uuid.UUID) (appbackend.FeedEntry, error) {
	feedEntry := appbackend.FeedEntry{}
	err := GetObjectWithID(feedEntryID, "feedentries", &feedEntry)
	return feedEntry, err
}

func GetFeedEntriesBetweenDates(from, to time.Time) ([]appbackend.FeedEntry, error) {
	feedEntries := []appbackend.FeedEntry{}
	selector := Sess.Select("*").From("feedentries").Where("cat >= ?", from).And("cat <= ?", to)
	if err := selector.All(&feedEntries); err != nil {
		return feedEntries, err
	}
	return feedEntries, nil
}

func SetFeedEntryMeta(feedEntryID uuid.UUID, meta string) error {
	if _, err := Sess.Update("feedentries").Set("meta", meta).Where("id = ?", feedEntryID).Exec(); err != nil {
		return err
	}
	return nil
}

func GetComment(commentID uuid.UUID) (Comment, error) {
	comment := Comment{}
	err := GetObjectWithID(commentID, "comments", &comment)
	return comment, err
}

func GetBox(boxID uuid.UUID) (appbackend.Box, error) {
	box := appbackend.Box{}
	err := GetObjectWithID(boxID, "boxes", &box)
	return box, err
}

func GetBoxFromPlantFeed(feedID uuid.UUID) (appbackend.Box, error) {
	box := appbackend.Box{}
	selector := Sess.Select("boxes.*").From("boxes").Join("plants").On("plants.boxid = boxes.id").Where("plants.feedid = ?", feedID)
	if err := selector.One(&box); err != nil {
		return box, err
	}
	return box, nil
}

func GetDevice(deviceID uuid.UUID) (appbackend.Device, error) {
	device := appbackend.Device{}
	err := GetObjectWithID(deviceID, "devices", &device)
	return device, err
}

func GetDeviceFromPlantFeed(feedID uuid.UUID) (appbackend.Device, error) {
	device := appbackend.Device{}
	selector := Sess.Select("devices.*").From("devices").Join("boxes").On("boxes.deviceid = devices.id").Join("plants").On("plants.boxid = boxes.id").Where("plants.feedid = ?", feedID)
	if err := selector.One(&device); err != nil {
		return device, err
	}
	return device, nil
}

func GetUserEndsForUserID(userID uuid.UUID) ([]UserEnd, error) {
	userends := []UserEnd{}

	err := GetObjectsWithField("userid", userID, "userends", &userends)
	return userends, err
}

func GetPlantForFeedID(feedID uuid.UUID) (appbackend.Plant, error) {
	plant := appbackend.Plant{}
	err := GetObjectWithField("feedid", feedID, "plants", &plant)
	return plant, err
}

func GetPlantForFeedEntryID(feedEntryID uuid.UUID) (appbackend.Plant, error) {
	plant := appbackend.Plant{}
	selector := Sess.Select("plants.*").From("plants").Join("feedentries").On("plants.feedid = feedentries.feedid").Where("feedentries.id = ?", feedEntryID)
	if err := selector.One(&plant); err != nil {
		return plant, err
	}

	return plant, nil
}

func GetActivePlantsForControllerIdentifier(controllerID string, boxSlotID int) ([]appbackend.Plant, error) {
	plants := []appbackend.Plant{}
	selector := Sess.Select("plants.*").From("plants").Join("boxes").On("boxes.id = plants.boxid").Join("devices").On("devices.id = boxes.deviceid").Where("devices.identifier = ?", controllerID).And("boxes.devicebox = ?", boxSlotID).And("plants.deleted = false").And("plants.archived = false").And("devices.deleted = false")
	if err := selector.All(&plants); err != nil {
		return plants, err
	}

	return plants, nil
}

func GetTimelapses() ([]appbackend.Timelapse, error) {
	timelapses := []appbackend.Timelapse{}

	selector := Sess.Select("timelapses.*").From("timelapses").Join("plants").On("plants.id = timelapses.plantid").Where("plants.deleted = false").And("timelapses.deleted = false")

	if err := selector.All(&timelapses); err != nil {
		return timelapses, err
	}

	return timelapses, nil
}

func GetTimelapseFrames(timelapseID uuid.UUID, from, to time.Time) ([]appbackend.TimelapseFrame, error) {
	timelapseFrames := []appbackend.TimelapseFrame{}

	selector := Sess.Select("timelapseframes.*").From("timelapseframes").Where("timelapseframes.timelapseid = ?", timelapseID).And("cat >= ?", from).And("cat <= ?", to)

	if err := selector.All(&timelapseFrames); err != nil {
		return timelapseFrames, err
	}

	return timelapseFrames, nil
}
