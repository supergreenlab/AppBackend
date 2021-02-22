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
	"upper.io/db.v3/postgresql"
)

func GetFeedEntry(feedEntryID uuid.UUID) (FeedEntry, error) {
	feedEntry := FeedEntry{}

	sess, err := postgresql.Open(Settings)
	if err != nil {
		return feedEntry, err
	}
	defer sess.Close()

	selector := sess.Select("*").From("feedentries").Where("id = ?", feedEntryID)
	if err := selector.One(&feedEntry); err != nil {
		return feedEntry, err
	}

	return feedEntry, nil
}

func GetComment(commentID uuid.UUID) (Comment, error) {
	comment := Comment{}

	sess, err := postgresql.Open(Settings)
	if err != nil {
		return comment, err
	}
	defer sess.Close()

	selector := sess.Select("*").From("comments").Where("id = ?", commentID)
	if err := selector.One(&comment); err != nil {
		return comment, err
	}

	return comment, nil
}

func GetUserEndsForUserID(userID uuid.UUID) ([]UserEnd, error) {
	userends := []UserEnd{}

	sess, err := postgresql.Open(Settings)
	if err != nil {
		return userends, err
	}
	defer sess.Close()

	selector := sess.Select("*").From("userends").Where("userid = ?", userID)
	if err := selector.All(&userends); err != nil {
		return userends, err
	}

	return userends, nil
}
