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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sirupsen/logrus"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"

	"github.com/julienschmidt/httprouter"
)

type timelapseUploadRequest struct {
}

type timelapseUploadURLResult struct {
	UploadPath string `json:"uploadPath"`
}

func timelapseUploadURLHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmup := timelapseUploadRequest{}
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

func timelapseLatestPic(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
	timelapseIDStr := p.ByName("id")
	timelapseID, err := uuid.FromString(timelapseIDStr)
	if err != nil {
		logrus.Errorf("uuid.FromString in timelapseLatestPic %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	frame, err := db.GetTimelapseFrame(timelapseID)
	if err != nil {
		logrus.Errorf("db.GetTimelapseFrame in timelapseLatestPic %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if frame.UserID != uid {
		errorMsg := "Access denied"
		logrus.Errorf("frame.UserID.UUID in timelapseLatestPic uid: %s", errorMsg, err, uid)
		http.Error(w, errorMsg, http.StatusUnauthorized)
		return
	}

	err = tools.LoadFeedMediaPublicURLs(&frame)
	if err != nil {
		logrus.Errorf("tools.LoadFeedMediaPublicURLs in timelapseLatestPic %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(frame); err != nil {
		logrus.Errorf("json.NewEncoder in timelapseLatestPic %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type SGLOverlayParams struct {
	Box   appbackend.Box         `json:"box"`
	Plant appbackend.Plant       `json:"plant"`
	Meta  appbackend.MetricsMeta `json:"meta"`
	URL   string                 `json:"url"`
	Host  string                 `json:"host"`
}

func sglOverlayHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sop := SGLOverlayParams{}
	if err := tools.DecodeJSONBody(w, r, &sop); err != nil {
		logrus.Errorf("tools.DecodeJSONBody in sglOverlayHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", sop.URL, nil)
	if err != nil {
		logrus.Errorf("http.NewRequest in sglOverlayHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.Host = sop.Host

	resp, err := client.Do(request)
	if err != nil {
		logrus.Errorf("client.Do in sglOverlayHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	picBuffer := &bytes.Buffer{}
	if _, err := picBuffer.ReadFrom(resp.Body); err != nil {
		logrus.Errorf("picBuffer.ReadFrom in sglOverlayHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	picBuffer, err = appbackend.AddSGLOverlays(sop.Box, sop.Plant, sop.Meta, picBuffer)

	w.Write(picBuffer.Bytes())
}
