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
	"context"
	"net/http"

	"github.com/gofrs/uuid"

	cmiddlewares "github.com/SuperGreenLab/AppBackend/internal/server/middlewares"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"
)

// AuthStackWithUserEnd - Adds userEndID from JWT claim, sets it as required
func AuthStackWithUserEnd() middleware.Stack {
	auth := cmiddlewares.AuthStack()
	auth.Use(JwtTokenUserEndID)
	auth.Use(UserEndIDRequired)
	return auth
}

// UserEndIDContextKey - context key which stores the request's userEndID
type UserEndIDContextKey struct{}

// JwtTokenUserEndID - Sets the userEndID from the claim, does nothing if missing
func JwtTokenUserEndID(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		claims := r.Context().Value(cmiddlewares.JwtClaimsContextKey{}).(jwt.MapClaims)
		if userEndID, ok := claims["userEndID"]; ok == true {
			ctx := context.WithValue(r.Context(), UserEndIDContextKey{}, uuid.FromStringOrNil(userEndID.(string)))
			fn(w, r.WithContext(ctx), p)
		} else {
			fn(w, r, p)
		}
	}
}
