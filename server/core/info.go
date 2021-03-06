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
package core

import (
	"encoding/json"
	"github.com/ServiceComb/service-center/pkg/util"
	"github.com/ServiceComb/service-center/server/core/backend"
	pb "github.com/ServiceComb/service-center/server/core/proto"
	"github.com/ServiceComb/service-center/server/infra/registry"
	"github.com/ServiceComb/service-center/version"
	"github.com/astaxie/beego"
	"golang.org/x/net/context"
)

var ServerInfo *pb.ServerInformation

func init() {
	ServerInfo = newInfo()
}

func newInfo() *pb.ServerInformation {
	return &pb.ServerInformation{
		Version: "0",
		Config: &pb.ServerConfig{
			MaxHeaderBytes: int64(beego.AppConfig.DefaultInt("max_header_bytes", 16384)),
			MaxBodyBytes:   beego.AppConfig.DefaultInt64("max_body_bytes", 2097152),

			ReadHeaderTimeout: beego.AppConfig.DefaultString("read_header_timeout", "60s"),
			ReadTimeout:       beego.AppConfig.DefaultString("read_timeout", "60s"),
			IdleTimeout:       beego.AppConfig.DefaultString("idle_timeout", "60s"),
			WriteTimeout:      beego.AppConfig.DefaultString("write_timeout", "60s"),

			LimitTTLUnit:     beego.AppConfig.DefaultString("limit_ttl", "s"),
			LimitConnections: int64(beego.AppConfig.DefaultInt("limit_conns", 0)),
			LimitIPLookup: beego.AppConfig.DefaultString("limit_iplookups",
				"RemoteAddr,X-Forwarded-For,X-Real-IP"),

			SslEnabled:    beego.AppConfig.DefaultInt("ssl_mode", 1) != 0,
			SslMinVersion: beego.AppConfig.DefaultString("ssl_min_version", "TLSv1.2"),
			SslVerifyPeer: beego.AppConfig.DefaultInt("ssl_verify_client", 1) != 0,
			SslCiphers:    beego.AppConfig.String("ssl_ciphers"),

			AutoSyncInterval:  beego.AppConfig.DefaultString("auto_sync_interval", "30s"),
			CompactIndexDelta: beego.AppConfig.DefaultInt64("compact_index_delta", 100),
		},
	}
}

func LoadServerInformation() error {
	resp, err := backend.Registry().Do(context.Background(),
		registry.GET, registry.WithStrKey(GetSystemKey()))
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return nil
	}

	err = json.Unmarshal(resp.Kvs[0].Value, ServerInfo)
	if err != nil {
		util.Logger().Errorf(err, "load system config failed, maybe incompatible")
		return nil
	}
	return nil
}

func UpgradeServerVersion() error {
	ServerInfo.Version = version.Ver().Version

	bytes, err := json.Marshal(ServerInfo)
	if err != nil {
		return err
	}
	_, err = backend.Registry().Do(context.Background(),
		registry.PUT, registry.WithStrKey(GetSystemKey()), registry.WithValue(bytes))
	if err != nil {
		return err
	}
	return nil
}
