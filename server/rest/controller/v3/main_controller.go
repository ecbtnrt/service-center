//Copyright 2017 Huawei Technologies Co., Ltd
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package v3

import (
	"encoding/json"
	"github.com/ServiceComb/service-center/pkg/rest"
	"github.com/ServiceComb/service-center/server/rest/controller/v4"
	"github.com/ServiceComb/service-center/version"
	"net/http"
)

const API_VERSION = "3.0.0"

type MainService struct {
	v4.MainService
}

func (this *MainService) URLPatterns() []rest.Route {
	return []rest.Route{
		{rest.HTTP_METHOD_GET, "/version", this.GetVersion},
		{rest.HTTP_METHOD_GET, "/health", this.ClusterHealth},
	}
}

func (this *MainService) GetVersion(w http.ResponseWriter, r *http.Request) {
	result := v4.Result{
		VersionSet: version.Ver(),
		ApiVersion: API_VERSION,
	}
	resultJSON, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(resultJSON)
}
