/*
 * Copyright (C) 2020  SuperGreenLab <towelie@supergreenlab.com>
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

package feeds

import (
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
)

func loadFeedMediaPublicURLs(fm db.FeedMedia) (db.FeedMedia, error) {
	expiry := time.Second * 60 * 60
	url1, err := storage.Client.PresignedGetObject("feedmedias", fm.FilePath, expiry, nil)
	if err != nil {
		return fm, err
	}
	fm.FilePath = url1.RequestURI()

	url2, err := storage.Client.PresignedGetObject("feedmedias", fm.ThumbnailPath, expiry, nil)
	if err != nil {
		return fm, err
	}
	fm.ThumbnailPath = url2.RequestURI()
	return fm, nil
}
