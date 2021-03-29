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

	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
)

type FeedMediasURL interface {
	SetURLs(filePath, thumbnailPath string)
	GetURLs() (filePath, thumbnailPath string)
}

func loadFeedMediaPublicURLs(fm FeedMediasURL) error {
	filePath, thumbnailPath := fm.GetURLs()
	expiry := time.Second * 60 * 60
	if filePath != "" {
		url1, err := storage.Client.PresignedGetObject("feedmedias", filePath, expiry, nil)
		if err != nil {
			return err
		}
		filePath = url1.RequestURI()
	}

	if thumbnailPath != "" {
		url2, err := storage.Client.PresignedGetObject("feedmedias", thumbnailPath, expiry, nil)
		if err != nil {
			return err
		}
		thumbnailPath = url2.RequestURI()
	}
	fm.SetURLs(filePath, thumbnailPath)
	return nil
}
