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

package explorer

import (
	"context"
	"net/http"

	"github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/SuperGreenLab/AppBackend/internal/server/tools"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func loadFeedMedias(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		fmus := r.Context().Value(middlewares.SelectResultContextKey{}).(tools.FeedMediasURLs)
		feedMedias := fmus.AsFeedMediasArray()
		for i, fm := range feedMedias {
			err := tools.LoadFeedMediaPublicURLs(fm)
			if err != nil {
				logrus.Errorf("tools.LoadFeedMediaPublicURLs in fetchPublicFeedEntries %q - p: %+v", err, p)
				continue
			}
			feedMedias[i] = fm
		}
		ctx := context.WithValue(r.Context(), middlewares.SelectResultContextKey{}, &feedMedias)
		fn(w, r.WithContext(ctx), p)
	}
}

func loadFeedMedia(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		feedMedia := r.Context().Value(middlewares.SelectResultContextKey{}).(tools.FeedMediasURL)
		err := tools.LoadFeedMediaPublicURLs(feedMedia)
		if err != nil {
			logrus.Errorf("tools.LoadFeedMediaPublicURLs in fetchPublicFeedEntries %q - p: %+v", err, p)
		}
		ctx := context.WithValue(r.Context(), middlewares.SelectResultContextKey{}, feedMedia)
		fn(w, r.WithContext(ctx), p)
	}
}
