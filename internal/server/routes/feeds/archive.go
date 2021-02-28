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
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/data/db"
	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	fmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/routes/feeds/middlewares"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"upper.io/db.v3/lib/sqlbuilder"
)

func archivePlantHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid := r.Context().Value(middlewares.UserIDContextKey{}).(uuid.UUID)
	sess := r.Context().Value(middlewares.SessContextKey{}).(sqlbuilder.Database)
	ueid := r.Context().Value(fmiddlewares.UserEndIDContextKey{}).(uuid.UUID)

	id := p.ByName("id")

	o := &db.Plant{}
	err := sess.Collection("plants").Find("id", id).One(o)
	if err != nil {
		logrus.Errorf("sess.Collection('plants') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if uid != o.GetUserID() {
		errorMsg := "Plant is owned by another user"
		logrus.Errorf("uid != o.GetUserID() in archivePlantHandler %q - uid: %s o: %+v", errorMsg, uid, o)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	if _, err := sess.Update("plants").Set("archived", true).Where("id = ?", o.GetID()).Exec(); err != nil {
		logrus.Errorf("sess.Update('plants') in archivePlantHandler %q - uid: %s o: %+v", err, uid, o)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := sess.Update("userend_plants").Set("dirty", true).Where("plantid", id).And("userendid != ?", ueid).And("userendid in (select id from userends where userid = ?)", uid).Exec(); err != nil {
		logrus.Warningf("sess.Update('userend_plants') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := sess.DeleteFrom("userend_plants").Where("plantid = ?", id).And("userendid = ?", ueid).Exec(); err != nil {
		logrus.Errorf("sess.DeleteFrom('userend_plants') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := sess.DeleteFrom("userend_timelapses").Where("timelapseid in (select id from timelapses where timelapses.plantid = ?)", id).Exec(); err != nil {
		logrus.Errorf("sess.DeleteFrom('userend_timelapses') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := sess.DeleteFrom("userend_feeds").Where("feedid = (select feedid from plants where plants.id = ?)", id).Exec(); err != nil {
		logrus.Errorf("sess.DeleteFrom('userend_feeds') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := sess.DeleteFrom("userend_feedentries").Where("feedentryid in (select id from feedentries where feedentries.feedid = (select feedid from plants where plants.id = ?))", id).Exec(); err != nil {
		logrus.Errorf("sess.DeleteFrom('userend_feedentries') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := sess.DeleteFrom("userend_feedmedias").Where("feedmediaid in (select id from feedmedias where feedentryid in (select id from feedentries where feedentries.feedid = (select feedid from plants where plants.id = ?)))", id).Exec(); err != nil {
		logrus.Errorf("sess.DeleteFrom('userend_feedmedias') in archivePlantHandler %q - id: %s uid: %s ueid: %s", err, id, uid, ueid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
