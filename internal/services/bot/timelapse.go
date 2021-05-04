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
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	"github.com/SuperGreenLab/AppBackend/internal/services/cron"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/dgrijalva/jwt-go"
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
	ID     uuid.UUID                   `json:"id"`
	Token  string                      `json:"token"`
	Frames []appbackend.TimelapseFrame `json:"timelapseFrames"`

	Plant  appbackend.Plant   `json:"plant"`
	Box    appbackend.Box     `json:"box"`
	Device *appbackend.Device `json:"device,omitempty"`
}

func timelapseJob(from, to time.Time) func() {
	return func() {
		timelapses, err := db.GetTimelapses()
		if err != nil {
			logrus.Errorf("db.GetTimelapses in timelapseJob %q", err)
			return
		}

		for _, timelapse := range timelapses {
			plant, err := db.GetPlant(timelapse.PlantID)
			if err != nil {
				logrus.Errorf("db.GetPlant in timelapseJob %q", err)
				time.Sleep(1 * time.Second)
				continue
			}

			box, err := db.GetBox(plant.BoxID)
			if err != nil {
				logrus.Errorf("db.GetBox in timelapseJob %q", err)
				time.Sleep(1 * time.Second)
				continue
			}

			var device *appbackend.Device
			if box.DeviceID.Valid {
				d, err := db.GetDevice(plant.BoxID)
				if err != nil {
					logrus.Errorf("db.GetDevice in timelapseJob %q", err)
					time.Sleep(1 * time.Second)
					continue
				}
				device = &d
			}

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

			requestID := uuid.Must(uuid.NewV4())

			// TODO DRY with internal/server/routes/users/login.go
			hmacSampleSecret := []byte(viper.GetString("JWTSecret"))
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"type":   "timelapse_worker",
				"userID": timelapse.UserID.String(),
			})
			tokenString, err := token.SignedString(hmacSampleSecret)
			if err != nil {
				logrus.Errorf("token.SignedString in loginHandler %q", err)
				continue
			}

			req := TimelapseRequest{
				ID:     requestID,
				Token:  tokenString,
				Frames: frames,
				Box:    box,
				Plant:  plant,
				Device: device,
			}
			if err := sendTimelapseRequests(req); err != nil {
				logrus.Errorf("sendTimelapseRequests in timelapseUploadURLHandler %q - %+v", err, req)
				continue
			}
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

func scheduleDailyTimelapse() {
	t := time.Now()
	from := t.Add(-24 * time.Hour)
	to := t

	cron.SetJob("timelapse", "0 0 * * *", timelapseJob(from, to))
}

func scheduleWeeklyTimelapse() {
	t := time.Now()
	from := t.Add(-7 * 24 * time.Hour)
	to := t

	cron.SetJob("timelapse", "0 0 * * sun", timelapseJob(from, to))
}

func initTimelapse() {
	scheduleDailyTimelapse()
	scheduleWeeklyTimelapse()
}
