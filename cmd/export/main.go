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
	"bytes"
	_ "embed"
	"encoding/json"
	"html/template"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/jeremywohl/flatten"
	"github.com/sirupsen/logrus"
)

type Plant appbackend.Plant

func (p Plant) SettingsFlat() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal([]byte(p.Settings), &m); err != nil {
		logrus.Fatalf("json.Unmarshal in SettingsFlat %q", err)
	}
	return flatten.Flatten(m, "", flatten.DotStyle)
}

type IndexTemplateData struct {
	Plants []Plant
}

type PlantTemplateData struct {
	Plant Plant
}

//go:embed templates/index.tmpl.html
var indexTemplateString string

//go:embed templates/plant.tmpl.html
var plantTemplateString string

func main() {
	InitConfig()
	db.Init()

	plants := make([]Plant, 0)
	err := db.Sess.Select("*").From("plants").Where("deleted=false").And("is_public=true").All(&plants)
	if err != nil {
		logrus.Fatalf("%q", err)
	}

	indexData := IndexTemplateData{
		Plants: plants,
	}
	indexTmpl, err := template.New("index").Parse(indexTemplateString)
	if err != nil {
		logrus.Fatalf("%q", err)
	}
	var indexHtml bytes.Buffer
	if err := indexTmpl.Execute(&indexHtml, indexData); err != nil {
		logrus.Fatalf("%q", err)
	}
	logrus.Infof("%s", indexHtml.String())

	plantTmpl, err := template.New("plant").Parse(plantTemplateString)
	if err != nil {
		logrus.Fatalf("%q", err)
	}

	for _, plant := range plants {
		plantData := PlantTemplateData{
			Plant: plant,
		}
		var plantHtml bytes.Buffer
		if err := plantTmpl.Execute(&plantHtml, plantData); err != nil {
			logrus.Fatalf("%q", err)
		}
		//logrus.Infof("%s", plantHtml.String())
	}
}
