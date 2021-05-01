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

package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/data/storage"
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	"github.com/SuperGreenLab/AppBackend/internal/services/cron"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("timelapseworkeraccesskey", "", "")
	_ = pflag.String("timelapseworkes", "", "List of base urls for timelapse workers, delimited by comas")
)

type TimelapseRequest struct {
	ID         uuid.UUID                   `json:"id"`
	UploadPath string                      `json:"uploadPath`
	Frames     []appbackend.TimelapseFrame `json:"timelapseFrames"`
}

func timelapseJob() {
	timelapses, err := db.GetTimelapses()
	if err != nil {
		logrus.Errorf("db.GetTimelapses in timelapseJob %q", err)
		return
	}

	t := time.Now()
	from := time.Date(t.Year(), t.Month(), t.Day()-7, 0, 0, 0, 0, time.UTC)
	to := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC)

	for _, timelapse := range timelapses {
		frames, err := db.GetTimelapseFrames(timelapse.ID.UUID, from, to)
		if err != nil {
			logrus.Errorf("db.GetTimelapses in timelapseJob %q", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(frames) == 0 {
			continue
		}

		for i, frame := range frames {
			err = tools.LoadFeedMediaPublicURLs(&frame)
			if err != nil {
				logrus.Errorf("tools.LoadFeedMediaPublicURLs in timelapseJob %q - frame: %+v", err, frame)
				continue
			}

			frames[i] = frame
		}

		expiry := time.Hour * 3
		requestID := uuid.Must(uuid.NewV4())
		path := fmt.Sprintf("render-%s.mp4", requestID.String())
		url2, err := storage.Client.PresignedPutObject("timelapses", path, expiry)
		if err != nil {
			logrus.Errorf("minioClient.PresignedPutObject in timelapseUploadURLHandler %q - %s", err, path)
			continue
		}

		req := TimelapseRequest{
			ID:         requestID,
			UploadPath: url2.RequestURI(),
			Frames:     frames,
		}
		if err := sendTimelapseRequests(req); err != nil {
			logrus.Errorf("sendTimelapseRequests in timelapseUploadURLHandler %q - %s", err, path)
			continue
		}
	}
}

func sendTimelapseRequests(req TimelapseRequest) error {
	servers := strings.Split(viper.GetString("TimelapseWorkers"), ",")
	s := rand.Int() % len(servers)
	server := servers[s]
	url := fmt.Sprintf("%s/timelapse", server)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	jsonStr, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Authentication", fmt.Sprintf("Bearer %s", viper.GetString("TimelapseWorkerAccessKey")))
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func initTimelapse() {
	cron.SetJob("timelapse", "45 * * * *", timelapseJob)
}
