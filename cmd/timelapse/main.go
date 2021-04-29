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
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/services/bot"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_            = pflag.String("accesskey", "", "Access key")
	_            = pflag.String("storageurl", "http://192.168.1.87:9000", "SGL Backend storage url")
	_            = pflag.String("storagehost", "minio:9000", "SGL Backend storage host name")
	downloadChan = make(chan bot.TimelapseRequest)
	generateChan = make(chan bot.TimelapseRequest)
)

func init() {
	viper.SetDefault("StorageUrl", "http://192.168.1.87:9000")
	viper.SetDefault("StorageHost", "minio:9000")
}

func downloadFrame(url, dst string) error {
	url = fmt.Sprintf("%s%s", viper.GetString("StorageUrl"), url)
	logrus.Infof("Downloading %s to %s", url, dst)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Host = viper.GetString("StorageHost")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	io.Copy(dstFile, resp.Body)
	return nil
}

func downloadTimelapses() {
	for tr := range downloadChan {
		dir := fmt.Sprintf("/var/timelapse/render-%s", tr.ID.String())
		if err := os.Mkdir(dir, 0700); err != nil {
			logrus.Errorf("os.Mkdir in downloadTimelapses %q", err)
			continue
		}
		var i = 0
		for _, frame := range tr.Frames {
			dst := fmt.Sprintf("%s/%010d.jpg", dir, i)
			if err := downloadFrame(frame.FilePath, dst); err != nil {
				continue
			}
			i += 1
		}
		doneDownloadTimelapses(tr)
	}
}

func generateTimelapses() {
	for tr := range generateChan {
		doneGenerateTimelapses(tr)
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
	jsonStr, err := json.Marshal(tr)
	if err != nil {
		return err
	}
	err = SetString(fmt.Sprintf("download-%s", tr.ID), string(jsonStr))
	if err != nil {
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
	if err := UnsetString(fmt.Sprintf("generate-%s", tr.ID)); err != nil {
		return err
	}
	return nil
}

func main() {
	dir := "/var/timelapse"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0700); err != nil {
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

	logrus.Info("Timelapse worker started")
	select {}
}
