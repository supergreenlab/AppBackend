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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/storage"

	"github.com/gofrs/uuid"

	"github.com/sirupsen/logrus"

	"github.com/SuperGreenLab/AppBackend/internal/server/tools"

	"github.com/julienschmidt/httprouter"
)

type feedMediaUploadURLParams struct {
	FileName string `json:"fileName"`
}

type feedMediaUploadURLResult struct {
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`
}

func feedMediaUploadURLHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmup := feedMediaUploadURLParams{}
	if err := tools.DecodeJSONBody(w, r, &fmup); err != nil {
		logrus.Errorf("tools.DecodeJSONBody in feedMediaUploadURLHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path := ""
	if strings.HasSuffix(fmup.FileName, ".mp4") {
		path = fmt.Sprintf("videos-%s.mp4", uuid.Must(uuid.NewV4()).String())
	} else if strings.HasSuffix(fmup.FileName, ".jpg") || strings.HasSuffix(fmup.FileName, ".jpeg") {
		path = fmt.Sprintf("pictures-%s.jpg", uuid.Must(uuid.NewV4()).String())
	} else {
		logrus.Errorf("Unknown file type %s", fmup.FileName)
		http.Error(w, "Unknown file type", http.StatusBadRequest)
		return
	}

	res := feedMediaUploadURLResult{}
	minioClient := storage.CreateMinioClient()
	expiry := time.Second * 60

	url1, err := minioClient.PresignedPutObject("feedmedias", path, expiry)
	if err != nil {
		logrus.Errorf("minioClient.PresignedPutObject in feedMediaUploadURLHandler %q - %s", err, path)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.FilePath = url1.RequestURI()

	path = fmt.Sprintf("thumbnail-%s.jpg", uuid.Must(uuid.NewV4()).String())
	url2, err := minioClient.PresignedPutObject("feedmedias", path, expiry)
	if err != nil {
		logrus.Errorf("minioClient.PresignedPutObject in feedMediaUploadURLHandler %q - %s", err, path)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.ThumbnailPath = url2.RequestURI()

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logrus.Errorf("json.NewEncoder in feedMediaUploadURLHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
