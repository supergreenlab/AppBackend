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

package explorer

import (
	"time"

	sgldb "github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/gofrs/uuid"
)

type publicPlantResult struct {
	ID            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	FilePath      string `db:"filepath" json:"filePath"`
	ThumbnailPath string `db:"thumbnailpath" json:"thumbnailPath"`

	Followed bool `db:"followed" json:"followed"`

	Settings    string `db:"settings" json:"settings"`
	BoxSettings string `db:"boxsettings" json:"boxSettings"`
}

func (r *publicPlantResult) SetURLs(filePath string, thumbnailPath string) {
	r.FilePath = filePath
	r.ThumbnailPath = thumbnailPath
}

func (r publicPlantResult) GetURLs() (filePath string, thumbnailPath string) {
	filePath, thumbnailPath = r.FilePath, r.ThumbnailPath
	return
}

type publicFeedEntry struct {
	sgldb.FeedEntry

	Liked      bool `db:"liked" json:"liked"`
	Bookmarked bool `db:"bookmarked" json:"bookmarked"`
	NComments  int  `db:"ncomments" json:"nComments"`
	NLikes     int  `db:"nlikes" json:"nLikes"`

	// TODO make an interface based middleware to unify the select* middlewares
	PlantID       uuid.NullUUID `db:"plantid,omitempty" json:"plantID,omitempty"`
	PlantName     string        `db:"name,omitempty" json:"plantName,omitempty"`
	CommentID     uuid.NullUUID `db:"commentid,omitempty" json:"commentID,omitempty"`
	Comment       *string       `db:"comment,omitempty" json:"comment,omitempty"`
	LikeDate      *time.Time    `db:"likecat,omitempty" json:"likeDate,omitempty"`
	ThumbnailPath *string       `db:"thumbnailpath,omitempty" json:"thumbnailpath,omitempty"`
}

func (r *publicFeedEntry) SetURLs(_, thumbnailPath string) {
	if thumbnailPath != "" {
		*r.ThumbnailPath = thumbnailPath
	}
}

func (r publicFeedEntry) GetURLs() (filePath, thumbnailPath string) {
	filePath, thumbnailPath = "", ""
	if r.ThumbnailPath != nil {
		thumbnailPath = *r.ThumbnailPath
	}
	return
}

type publicFeedEntriesResult struct {
	Entries []publicFeedEntry `json:"entries"`
}

type publicFeedEntryResult struct {
	Entry publicFeedEntry `json:"entry"`
}

type publicFeedMediasResult struct {
	Medias []sgldb.FeedMedia `json:"medias"`
}
