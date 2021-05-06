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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/services/bot"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// TODO remove dependencies on external commands

type feedMediaUploadURLParams struct {
	FileName string `json:"fileName"`
}

type feedMediaUploadURLResult struct {
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath"`
}

type insertResponse struct {
	ID uuid.UUID `json:"id"`
}

const STORAGE_DIR = "/var/timelapse"

var (
	downloadChan = make(chan bot.TimelapseRequest)
	generateChan = make(chan bot.TimelapseRequest)
)

func downloadTimelapses() {
	for tr := range downloadChan {
		dir := fmt.Sprintf("%s/render-%s", STORAGE_DIR, tr.ID.String())
		downloadedFrames := make([]appbackend.TimelapseFrame, len(tr.Frames)*2)
		i := 0
		for _, frame := range tr.Frames {
			dst := fmt.Sprintf("%s/frame-%d.jpg", dir, i)
			if _, err := os.Stat(dst); os.IsNotExist(err) {
				if err := appbackend.DownloadFile(frame.FilePath, dst); err != nil {
					logrus.Errorf("appbackend.DownloadFile in downloadTimelapses %q", err)
					continue
				}
			} else if err != nil {
				logrus.Errorf("os.Stat in downloadTimelapses %q", err)
				continue
			}
			frame.FilePath = dst
			downloadedFrames[i] = frame

			var meta appbackend.MetricsMeta
			if err := json.Unmarshal([]byte(frame.Meta), &meta); err != nil {
				logrus.Errorf("json.Unmarshal in downloadTimelapses %q", err)
				continue
			}

			if err := appbackend.AddSGLOverlaysForFile(tr.Box, tr.Plant, meta, dst); err != nil {
				logrus.Errorf("appbackend.AddSGLOverlaysForFile in generateTimelapses %q", err)
				continue
			}

			if i != 0 && i%2 == 0 {
				image1 := fmt.Sprintf("%s/frame-%d.jpg", dir, i-2)
				image2 := fmt.Sprintf("%s/frame-%d.jpg", dir, i)
				dstBlur := fmt.Sprintf("%s/frame-%d.jpg", dir, i-1)
				if _, err := os.Stat(dstBlur); os.IsNotExist(err) {
					cmd := exec.Command("composite", "-blend", "50", image1, image2, "-matte", dstBlur)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Start(); err != nil {
						logrus.Errorf("cmd.Start in generateTimelapses %q", err)
					}
					cmd.Wait()
				}
				blurFrame := frame
				blurFrame.FilePath = dstBlur
				downloadedFrames[i-1] = blurFrame
			}
			i += 2
		}
		tr.Frames = downloadedFrames
		if err := doneDownloadTimelapses(tr); err != nil {
			logrus.Errorf("doneDownloadTimelapses in downloadTimelapses %q", err)
			continue
		}
	}
}

func generateTimelapses() {
	for tr := range generateChan {
		dir := fmt.Sprintf("%s/render-%s", STORAGE_DIR, tr.ID.String())
		dst := fmt.Sprintf("%s/render-%s.mp4", dir, tr.ID.String())

		cmd := exec.Command("/usr/bin/ffmpeg", "-r", "40", "-i", fmt.Sprintf("%s/frame-%%d.jpg", dir), "-vf", "scale=1440:-2", dst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			logrus.Errorf("cmd.Start in generateTimelapses %q", err)
		}
		cmd.Wait()

		dstThumbnail := fmt.Sprintf("%s/render-%s.jpg", dir, tr.ID.String())
		cmd = exec.Command("/usr/bin/ffmpeg", "-i", dst, "-ss", "00:00:00.000", "-frames:v", "1", dstThumbnail)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			logrus.Errorf("cmd.Start in generateTimelapses %q", err)
		}
		cmd.Wait()

		if err := doneGenerateTimelapses(tr); err != nil {
			logrus.Errorf("doneGenerateTimelapses in generateTimelapses %q", err)
			continue
		}
	}
}

func handleTimelapse(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if r.Header.Get("Authentication") != fmt.Sprintf("Bearer %s", viper.GetString("AccessKey")) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	tr := bot.TimelapseRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tr); err != nil {
		logrus.Errorf("json.Unmarshal in handleTimelapse %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := startDownloadTimelapses(tr); err != nil {
		logrus.Errorf("startDownloadTimelapses in handleTimelapse %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "OK")
}

func startDownloadTimelapses(tr bot.TimelapseRequest) error {
	logrus.Infof("Starting download render %s", tr.ID.String())
	jsonStr, err := json.Marshal(tr)
	if err != nil {
		return err
	}
	err = SetString(fmt.Sprintf("download-%s", tr.ID), string(jsonStr))
	if err != nil {
		return err
	}

	dir := fmt.Sprintf("%s/render-%s", STORAGE_DIR, tr.ID.String())
	if err := os.Mkdir(dir, 0700); err != nil {
		logrus.Errorf("os.Mkdir in downloadTimelapses %q", err)
		return err
	}

	downloadChan <- tr
	return nil
}

func doneDownloadTimelapses(tr bot.TimelapseRequest) error {
	if err := UnsetString(fmt.Sprintf("download-%s", tr.ID)); err != nil {
		return err
	}
	return startGenerateTimelapses(tr)
}

func startGenerateTimelapses(tr bot.TimelapseRequest) error {
	logrus.Infof("Starting generate render %s", tr.ID.String())
	jsonStr, err := json.Marshal(tr)
	if err != nil {
		return err
	}
	err = SetString(fmt.Sprintf("generate-%s", tr.ID), string(jsonStr))
	if err != nil {
		return err
	}

	generateChan <- tr
	return nil
}

func doneGenerateTimelapses(tr bot.TimelapseRequest) error {
	dir := fmt.Sprintf("%s/render-%s", STORAGE_DIR, tr.ID.String())
	fileName := fmt.Sprintf("render-%s.mp4", tr.ID.String())

	fmup := feedMediaUploadURLParams{
		FileName: fileName,
	}
	fmur := feedMediaUploadURLResult{}
	if err := appbackend.POSTSGLObject(tr.Token, "/feedMediaUploadURL", fmup, &fmur); err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", dir, fileName)
	if err := appbackend.UploadSGLObjectFile(fmur.FilePath, filePath); err != nil {
		return err
	}
	thumbnailPath := fmt.Sprintf("%s/render-%s.jpg", dir, tr.ID.String())
	if err := appbackend.UploadSGLObjectFile(fmur.ThumbnailPath, thumbnailPath); err != nil {
		return err
	}

	fe := appbackend.FeedEntry{
		FeedID: tr.Plant.FeedID,
		Date:   time.Now(),
		Type:   "FE_TIMELAPSE",

		Params: "{}",
	}

	ir := insertResponse{}
	if err := appbackend.POSTSGLObject(tr.Token, "/feedEntry", fe, &ir); err != nil {
		return err
	}

	filePathS := strings.Split(fmur.FilePath, "/")
	filePathS = strings.Split(filePathS[2], "?")

	thumbnailPathS := strings.Split(fmur.ThumbnailPath, "/")
	thumbnailPathS = strings.Split(thumbnailPathS[2], "?")
	fm := appbackend.FeedMedia{
		FeedEntryID:   ir.ID,
		FilePath:      filePathS[0],
		ThumbnailPath: thumbnailPathS[0],

		Params: "{}",
	}
	if err := appbackend.POSTSGLObject(tr.Token, "/feedMedia", fm, nil); err != nil {
		return err
	}

	if err := UnsetString(fmt.Sprintf("generate-%s", tr.ID)); err != nil {
		return err
	}

	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return nil
}

func loadInitialDownloads() {
	iter := db.NewIterator(util.BytesPrefix([]byte("download-")), nil)
	for iter.Next() {
		tr := bot.TimelapseRequest{}
		if err := json.Unmarshal(iter.Value(), &tr); err != nil {
			logrus.Errorf("json.Unmarshal in loadInitialDownloads %q - iter.Value: %s", err, iter.Value)
			continue
		}
		logrus.Infof("Starting download render %s", tr.ID.String())
		downloadChan <- tr
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		logrus.Errorf("iter.Error in loadInitialDownloads %q", err)
		return
	}
}

func loadInitialGenerate() {
	iter := db.NewIterator(util.BytesPrefix([]byte("generate-")), nil)
	for iter.Next() {
		tr := bot.TimelapseRequest{}
		if err := json.Unmarshal(iter.Value(), &tr); err != nil {
			logrus.Errorf("json.Unmarshal in loadInitialGenerate %q - iter.Value: %s", err, iter.Value)
			continue
		}
		generateChan <- tr
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		logrus.Errorf("iter.Error in loadInitialGenerate %q", err)
		return
	}
}

func main() {
	if _, err := os.Stat(STORAGE_DIR); os.IsNotExist(err) {
		if err := os.Mkdir(STORAGE_DIR, 0700); err != nil {
			logrus.Fatalf("os.Mkdir in main %q", err)
		}
	}
	InitConfig()
	InitKV()

	router := httprouter.New()

	router.POST("/timelapse", handleTimelapse)

	go func() {
		corsOpts := cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: false,
		}

		log.Fatal(http.ListenAndServe(":8083", cors.New(corsOpts).Handler(router)))
	}()

	go downloadTimelapses()
	go generateTimelapses()

	loadInitialDownloads()
	loadInitialGenerate()

	logrus.Info("Timelapse worker started")
	select {}
}
