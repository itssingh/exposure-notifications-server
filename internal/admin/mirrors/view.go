// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, softwar
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mirrors

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/exposure-notifications-server/internal/admin"
	mirrordatabase "github.com/google/exposure-notifications-server/internal/mirror/database"
	mirrormodel "github.com/google/exposure-notifications-server/internal/mirror/model"
	"github.com/google/exposure-notifications-server/internal/serverenv"
)

type viewController struct {
	config *admin.Config
	env    *serverenv.ServerEnv
}

func NewView(c *admin.Config, env *serverenv.ServerEnv) admin.Controller {
	return &viewController{config: c, env: env}
}

func (v *viewController) Execute(c *gin.Context) {
	ctx := c.Request.Context()
	m := admin.TemplateMap{}

	db := mirrordatabase.New(v.env.Database())
	mirror := &mirrormodel.Mirror{}
	if idParam := c.Param("id"); idParam != "0" {
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			admin.ErrorPage(c, "unable to parse `id` param.")
			return
		}
		mirror, err = db.GetMirror(ctx, id)
		if err != nil {
			admin.ErrorPage(c, fmt.Sprintf("Error loading mirror: %v", err))
			return
		}
	}

	var mirrorFiles []*mirrormodel.MirrorFile
	if mirror.ID != 0 {
		var err error
		mirrorFiles, err = db.ListFiles(ctx, mirror.ID)
		if err != nil {
			admin.ErrorPage(c, fmt.Sprintf("Error loading mirror files: %v", err))
			return
		}
	}

	m["mirror"] = mirror
	m["mirrorFiles"] = mirrorFiles
	c.HTML(http.StatusOK, "mirror", m)
}
