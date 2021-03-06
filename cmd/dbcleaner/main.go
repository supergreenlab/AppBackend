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

package main

import (
	"log"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/config"
	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type FeedEntriesDuplicateCount struct {
	N      int       `db:"n"`
	Date   time.Time `db:"createdat"`
	Type   string    `db:"etype"`
	FeedID uuid.UUID `db:"feedid" json:"feedID"`
}

func getDuplicatedFeedEntries() []FeedEntriesDuplicateCount {
	rows, err := db.Sess.Query("select count(*), createdat, etype, feedid from feedentries where deleted=false group by createdat, etype, feedid having count(*)>1 order by createdat desc")
	if err != nil {
		logrus.Fatal(err)
	}
	dups := []FeedEntriesDuplicateCount{}
	for rows.Next() {
		dup := FeedEntriesDuplicateCount{}
		if err := rows.Scan(&dup.N, &dup.Date, &dup.Type, &dup.FeedID); err != nil {
			log.Fatal(err)
		}
		dups = append(dups, dup)
	}
	return dups
}

func getFeedEntriesForDup(dup FeedEntriesDuplicateCount) []db.FeedEntry {
	feedEntries := []db.FeedEntry{}
	selector := db.Sess.Select("*").From("feedentries").Where("createdat = ?", dup.Date).And("etype = ?", dup.Type).And("feedid = ?", dup.FeedID).And("deleted=false").OrderBy("cat desc")
	if err := selector.All(&feedEntries); err != nil {
		logrus.Fatal(err)
	}

	return feedEntries
}

func getFeedMediasForFeedEntry(fe db.FeedEntry) []db.FeedMedia {
	feedMedias := []db.FeedMedia{}
	selector := db.Sess.Select("*").From("feedmedias").Where("feedentryid = ?", fe.ID).And("deleted=false")
	if err := selector.All(&feedMedias); err != nil {
		logrus.Fatal(err)
	}

	return feedMedias

}

func deleteFeedEntry(fe db.FeedEntry) {
	if _, err := db.Sess.Update("feedentries").Set("deleted", true).Where("id = ?", fe.ID.UUID).Exec(); err != nil {
		logrus.Fatal(err)
	}
	if _, err := db.Sess.Update("userend_feedentries").Set("dirty", true).Where("feedentryid = ?", fe.ID.UUID).Exec(); err != nil {
		logrus.Fatal(err)
	}
}

func deleteFeedMedia(fm db.FeedMedia) {
	if _, err := db.Sess.Update("feedmedias").Set("deleted", true).Where("id = ?", fm.ID.UUID).Exec(); err != nil {
		logrus.Fatal(err)
	}
	if _, err := db.Sess.Update("userend_feedmedias").Set("dirty", true).Where("feedmediaid = ?", fm.ID.UUID).Exec(); err != nil {
		logrus.Fatal(err)
	}
}

func updateFeedEntryDate(fe db.FeedEntry, d time.Time) {
	if _, err := db.Sess.Update("feedentries").Set("createdat", d).Where("id = ?", fe.ID.UUID).Exec(); err != nil {
		logrus.Fatal(err)
	}
}

type mediaTypes []string

func (mt mediaTypes) contains(t string) bool {
	for _, m := range mt {
		if m == t {
			return true
		}
	}
	return false
}

func main() {
	config.Init()
	db.Init()

	mt := mediaTypes{
		"FE_MEDIA",
		"FE_BENDING",
		"FE_DEFOLATION",
		"FE_TRANSPLANT",
		"FE_FIMMING",
		"FE_TOPPING",
		"FE_MEASURE",
	}

	dups := getDuplicatedFeedEntries()
	for _, dup := range dups {
		logrus.Infof("%d, %s, %s, %s", dup.N, dup.Date, dup.Type, dup.FeedID)
		fes := getFeedEntriesForDup(dup)
		if mt.contains(dup.Type) {
			kept := 0
			for _, fe := range fes {
				fms := getFeedMediasForFeedEntry(fe)
				if len(fms) == 0 {
					logrus.Infof("deleting %s %d %s", fe.ID.UUID, len(fms), fe.Type)
					for _, fm := range fms {
						deleteFeedMedia(fm)
					}
					deleteFeedEntry(fe)
				} else {
					logrus.Infof("keeping %s %d %s", fe.ID.UUID, len(fms), fe.Type)
					if kept > 0 {
						d := fe.Date.Add(time.Minute * time.Duration(kept))
						logrus.Infof("adding %dmin to created at, was %s, will be %s", kept, fe.Date, d)
						updateFeedEntryDate(fe, d)
					}
					kept += 1
				}
			}
		} else {
			logrus.Infof("keeping %s %s", fes[0].ID.UUID, fes[0].Type)
			for _, fe := range fes[1:] {
				logrus.Infof("deleting %s %s", fe.ID.UUID, fe.Type)
				deleteFeedEntry(fe)
			}
		}
	}
}
