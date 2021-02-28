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
	"github.com/gofrs/uuid"
)

func GetFeedEntry(feedEntryID uuid.UUID) (FeedEntry, error) {
	feedEntry := FeedEntry{}
	err := GetObjectWithID(feedEntryID, "feedentries", &feedEntry)
	return feedEntry, err
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
