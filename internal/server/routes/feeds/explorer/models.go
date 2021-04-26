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

	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
)

type publicPlant struct {
	ID            string    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	ThumbnailPath string    `db:"thumbnailpath" json:"thumbnailPath"`
	LastUpdate    time.Time `db:"lastupdate" json:"lastUpdate"`

	Followed bool `db:"followed" json:"followed"`
	NFollows int  `db:"nfollows" json:"nFollows"`

	Settings    string `db:"settings" json:"settings"`
	BoxSettings string `db:"boxsettings" json:"boxSettings"`
}

func (r *publicPlant) SetURLs(paths []string) {
	r.ThumbnailPath = paths[0]
}

func (r publicPlant) GetURLs() []tools.S3Path {
	return []tools.S3Path{
		tools.S3Path{
			Path:   &r.ThumbnailPath,
			Bucket: "feedmedias",
		},
	}
}

type publicPlants []*publicPlant

func (pfe *publicPlants) AsFeedMediasArray() []tools.S3FileHolder {
	res := make([]tools.S3FileHolder, len(*pfe))
	for i, fe := range *pfe {
		res[i] = fe
	}
	return res
}

type publicFeedEntry struct {
	appbackend.FeedEntry

	Liked      bool `db:"liked" json:"liked"`
	Bookmarked bool `db:"bookmarked" json:"bookmarked"`
	NComments  int  `db:"ncomments" json:"nComments"`
	NLikes     int  `db:"nlikes" json:"nLikes"`

	// Split model?
	PlantID            *uuid.NullUUID `db:"plantid,omitempty" json:"plantID,omitempty"`
	PlantName          *string        `db:"plantname,omitempty" json:"plantName,omitempty"`
	PlantThumbnailPath *string        `db:"plantthumbnailpath,omitempty" json:"plantThumbnailPath,omitempty"`
	PlantSettings      *string        `db:"plantsettings,omitempty" json:"plantSettings,omitempty"`
	BoxSettings        *string        `db:"boxsettings,omitempty" json:"boxSettings,omitempty"`
	Followed           *bool          `db:"followed,omitempty" json:"followed,omitempty"`
	NFollows           *int           `db:"nfollows,omitempty" json:"nFollows,omitempty"`

	Nickname    *string        `db:"nickname" json:"nickname"`
	Pic         *string        `db:"pic" json:"pic"`
	CommentID   *uuid.NullUUID `db:"commentid,omitempty" json:"commentID,omitempty"`
	Comment     *string        `db:"comment,omitempty" json:"comment,omitempty"`
	ReplyTo     *string        `db:"commentreplyto,omitempty" json:"commentReplyTo,omitempty"`
	CommentType *string        `db:"commenttype,omitempty" json:"commentType,omitempty"`
	CommentDate *time.Time     `db:"commentdate,omitempty" json:"commentDate,omitempty"`

	LikeDate      *time.Time `db:"likecat,omitempty" json:"likeDate,omitempty"`
	ThumbnailPath *string    `db:"thumbnailpath,omitempty" json:"thumbnailPath,omitempty"`
}

func (r *publicFeedEntry) SetURLs(paths []string) {
	if paths[0] != "" {
		*r.ThumbnailPath = paths[0]
	}
	if paths[1] != "" {
		*r.Pic = paths[1]
	}
	if paths[2] != "" {
		*r.PlantThumbnailPath = paths[2]
	}
}

func (r publicFeedEntry) GetURLs() (paths []tools.S3Path) {
	return []tools.S3Path{
		tools.S3Path{
			Path:   r.ThumbnailPath,
			Bucket: "feedmedias",
		},
		tools.S3Path{
			Path:   r.Pic,
			Bucket: "users",
		},
		tools.S3Path{
			Path:   r.PlantThumbnailPath,
			Bucket: "feedmedias",
		},
	}
}

type publicFeedEntries []*publicFeedEntry

func (pfe *publicFeedEntries) AsFeedMediasArray() []tools.S3FileHolder {
	res := make([]tools.S3FileHolder, len(*pfe))
	for i, fe := range *pfe {
		res[i] = fe
	}
	return res
}

type publicFeedMedia struct {
	appbackend.FeedMedia
}

func (r *publicFeedMedia) SetURLs(paths []string) {
	r.FilePath = paths[0]
	r.ThumbnailPath = paths[1]
}

func (r publicFeedMedia) GetURLs() []tools.S3Path {
	return []tools.S3Path{
		tools.S3Path{
			Path:   &r.FilePath,
			Bucket: "feedmedias",
		},
		tools.S3Path{
			Path:   &r.ThumbnailPath,
			Bucket: "feedmedias",
		},
	}
}

type publicFeedMedias []*publicFeedMedia

func (pfe *publicFeedMedias) AsFeedMediasArray() []tools.S3FileHolder {
	res := make([]tools.S3FileHolder, len(*pfe))
	for i, fe := range *pfe {
		res[i] = fe
	}
	return res
}
