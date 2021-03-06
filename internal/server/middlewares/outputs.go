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

package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SuperGreenLab/AppBackend/internal/data/kv"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type CacheKeyContextKey struct{}

// OutputObjectID - returns the inserted object ID
func OutputObjectID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := r.Context().Value(InsertedIDContextKey{}).(uuid.UUID)
	response := struct {
		ID string `json:"id"`
	}{id.String()}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("json.NewEncoder in OutputObjectID %q - %+v", err, response)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func OutputResult(name string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		cacheKey, okCache := r.Context().Value(CacheKeyContextKey{}).(string)
		results := r.Context().Value(SelectResultContextKey{}).(interface{})

		response := map[string]interface{}{}
		response[name] = results

		var cacheData bytes.Buffer
		var mw io.Writer = w
		if okCache {
			mw = io.MultiWriter(w, &cacheData)
		}

		if err := json.NewEncoder(mw).Encode(response); err != nil {
			logrus.Errorf("json.NewEncoder in OutputResult %q - %+v", err, response)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if okCache {
			kv.SetStringWithExpiration(cacheKey, cacheData.String(), 1*time.Minute)
			kv.SetStringWithExpiration(fmt.Sprintf("%s.last", cacheKey), cacheData.String(), 20*time.Minute)
		}
	}
}

// OutputSelectResult - returns the list of data
func OutputSelectResult(collection string) httprouter.Handle {
	return OutputResult(collection)
}

// OutputSelectOneResult - returns the data
func OutputSelectOneResult() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		result := r.Context().Value(SelectResultContextKey{}).(interface{})
		if err := json.NewEncoder(w).Encode(result); err != nil {
			logrus.Errorf("json.NewEncoder in OutputSelectOneResult %q - %+v", err, result)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// OutputOK - returns the OK response
func OutputOK(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	response := struct {
		Status string `json:"status"`
	}{"OK"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Errorf("json.NewEncoder in OutputOK %q - %+v", err, response)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
