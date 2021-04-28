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

package feeds

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sirupsen/logrus"

	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"

	"github.com/julienschmidt/httprouter"
)

type timelapseUploadURLParams struct {
	FileName string `json:"fileName"`
}

type timelapseUploadURLResult struct {
	UploadPath string `json:"uploadPath"`
}

func timelapseUploadURLHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmup := timelapseUploadURLParams{}
	if err := tools.DecodeJSONBody(w, r, &fmup); err != nil {
		logrus.Errorf("tools.DecodeJSONBody in timelapseUploadURLHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := timelapseUploadURLResult{}
	expiry := time.Second * 60

	path := fmt.Sprintf("frame-%s.jpg", uuid.Must(uuid.NewV4()).String())
	url2, err := storage.Client.PresignedPutObject("timelapses", path, expiry)
	if err != nil {
		logrus.Errorf("minioClient.PresignedPutObject in timelapseUploadURLHandler %q - %s", err, path)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res.UploadPath = url2.RequestURI()

	if err := json.NewEncoder(w).Encode(res); err != nil {
		logrus.Errorf("json.NewEncoder in timelapseUploadURLHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
