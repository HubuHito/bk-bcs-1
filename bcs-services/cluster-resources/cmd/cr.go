/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * 	http://opensource.org/licenses/MIT
 *
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * cr.go ClusterResources 模块服务启动相关
 */

package cmd

import (
	"flag"
	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	"github.com/Tencent/bk-bcs/bcs-common/common/conf"
	"github.com/Tencent/bk-bcs/bcs-services/cluster-resources/pkg/common"
	"github.com/Tencent/bk-bcs/bcs-services/cluster-resources/pkg/options"
)

var confFilePath = flag.String("conf", common.DefaultConfPath, "配置文件路径")

// Start 初始化并启动 ClusterResources 服务
func Start() {
	flag.Parse()
	blog.Infof("Conf File Path: %s", *confFilePath)
	opts, err := options.LoadConf(*confFilePath)

	// 初始化日志相关配置
	// TODO 排查 LogDir 不生效原因，目前都是在 ./logs
	blog.InitLogs(conf.LogConfig{
		LogDir:          opts.Log.LogDir,
		LogMaxSize:      opts.Log.LogMaxSize,
		LogMaxNum:       opts.Log.LogMaxNum,
		ToStdErr:        opts.Log.ToStdErr,
		AlsoToStdErr:    opts.Log.AlsoToStdErr,
		Verbosity:       opts.Log.Verbosity,
		StdErrThreshold: opts.Log.StdErrThreshold,
		VModule:         opts.Log.VModule,
		TraceLocation:   opts.Log.TraceLocation,
	})
	defer blog.CloseLogs()

	if err != nil {
		blog.Fatalf("Load Cluster Resources Options Failed: %s", err.Error())
	}
	crSvc := newClusterResourcesService(opts)
	if err := crSvc.Init(); err != nil {
		blog.Fatalf("Init Cluster Resources Failed: %s", err.Error())
	}
	if err := crSvc.Run(); err != nil {
		blog.Fatalf("Run Cluster Resources Failed: %s", err.Error())
	}
}
