/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 *  Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *  Licensed under the MIT License (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *  http://opensource.org/licenses/MIT
 *  Unless required by applicable law or agreed to in writing, software distributed under
 *  the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 *  either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package worker

import (
	"context"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-data-manager/pkg/prom"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-data-manager/pkg/types"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-data-manager/pkg/utils"
	"sync"
	"time"

	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	"github.com/Tencent/bk-bcs/bcs-common/common/codec"
	"github.com/Tencent/bk-bcs/bcs-common/pkg/bcsapi"
	"github.com/Tencent/bk-bcs/bcs-common/pkg/msgqueue"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-data-manager/pkg/cmanager"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-data-manager/pkg/common"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-data-manager/pkg/datajob"
	"github.com/micro/go-micro/v2/broker"
	"github.com/robfig/cron/v3"
)

// Producer produce data job
type Producer struct {
	msgQueue        msgqueue.MessageQueue
	cron            *cron.Cron
	CMClient        cmanager.ClusterManagerClient
	k8sStorageCli   bcsapi.Storage
	mesosStorageCli bcsapi.Storage
	ctx             context.Context
	cancel          context.CancelFunc
	resourceGetter  common.GetterInterface
	concurrency     int
}

// NewProducer new producer
func NewProducer(rootCtx context.Context, msgQueue msgqueue.MessageQueue, cron *cron.Cron,
	cmClient cmanager.ClusterManagerClient, k8sStorageCli, mesosStorageCli bcsapi.Storage,
	getter common.GetterInterface, concurrency int) *Producer {
	ctx, cancel := context.WithCancel(rootCtx)
	return &Producer{
		msgQueue:        msgQueue,
		cron:            cron,
		CMClient:        cmClient,
		k8sStorageCli:   k8sStorageCli,
		mesosStorageCli: mesosStorageCli,
		ctx:             ctx,
		cancel:          cancel,
		resourceGetter:  getter,
		concurrency:     concurrency,
	}
}

// Stop stop producer
func (p *Producer) Stop() {
	p.cron.Stop()
}

// Run run producer
func (p *Producer) Run() {
	defer func() {
		if r := recover(); r != nil {
			blog.Errorf("internal error: %v", p)
		}
	}()
	p.cron.Start()
}

// InitCronList get all cron func
func (p *Producer) InitCronList() error {
	minSpec := "0-59/1 * * * * "
	if _, err := p.cron.AddFunc(minSpec, func() {
		p.WorkloadProducer(types.DimensionMinute)
	}); err != nil {
		return err
	}

	tenMinSpec := "0-59/10 * * * * "
	if _, err := p.cron.AddFunc(tenMinSpec, func() {
		p.NamespaceProducer(types.DimensionMinute)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(tenMinSpec, func() {
		p.ClusterProducer(types.DimensionMinute)
	}); err != nil {
		return err
	}

	hourSpec := "10 * * * * "
	if _, err := p.cron.AddFunc(hourSpec, func() {
		p.WorkloadProducer(types.DimensionHour)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(hourSpec, func() {
		p.NamespaceProducer(types.DimensionHour)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(hourSpec, func() {
		p.ClusterProducer(types.DimensionHour)
	}); err != nil {
		return err
	}

	daySpec := "30 0 * * *"
	if _, err := p.cron.AddFunc(daySpec, func() {
		p.WorkloadProducer(types.DimensionDay)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(daySpec, func() {
		p.NamespaceProducer(types.DimensionDay)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(daySpec, func() {
		p.ClusterProducer(types.DimensionDay)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(daySpec, func() {
		p.ProjectProducer(types.DimensionDay)
	}); err != nil {
		return err
	}
	if _, err := p.cron.AddFunc(daySpec, func() {
		p.PublicProducer(types.DimensionDay)
	}); err != nil {
		return err
	}
	blog.Infof("init cron list")
	return nil
}

// PublicProducer is the function to produce public data job and send to message queue
func (p *Producer) PublicProducer(dimension string) {
	opts := types.JobCommonOpts{
		Dimension:   dimension,
		ObjectType:  types.PublicType,
		CurrentTime: utils.FormatTime(time.Now(), dimension),
	}
	err := p.SendJob(opts)
	if err != nil {
		blog.Errorf("send public job to msg queue error, opts: %v, err: %v", opts, err)
		return
	}
	blog.Infof("[producer] send public job success")
}

// ProjectProducer is the function to produce project data job and send to message queue
func (p *Producer) ProjectProducer(dimension string) {
	startTime := time.Now()
	var err error
	defer func() {
		prom.ReportProduceJobLatencyMetric(types.ProjectType, dimension, err, startTime)
	}()
	jobTime := utils.FormatTime(time.Now(), dimension)
	cmConn, err := p.CMClient.GetClusterManagerConn()
	if err != nil {
		blog.Errorf("get cm conn error:%v", err)
		return
	}
	defer cmConn.Close()
	cliWithHeader := p.CMClient.NewGrpcClientWithHeader(p.ctx, cmConn)
	projectList, err := p.resourceGetter.GetProjectIDList(cliWithHeader.Ctx, cliWithHeader.Cli)
	if err != nil || projectList == nil {
		blog.Errorf("get projectIDList error: %v", err)
		return
	}
	for _, project := range projectList {
		opts := types.JobCommonOpts{
			ProjectID:   project.ProjectID,
			BusinessID:  project.BusinessID,
			CurrentTime: jobTime,
			Dimension:   dimension,
			ObjectType:  types.ProjectType,
		}
		err := p.SendJob(opts)
		if err != nil {
			blog.Errorf("send project job to msg queue error, opts: %v, err: %v", opts, err)
			return
		}
	}
	blog.Infof("[producer] send project job success, count: %d, jobTime:%v, startTime:%v, currentTime:%v, cost:%v",
		len(projectList), jobTime, startTime, time.Now(), time.Now().Sub(startTime))
}

// ClusterProducer is the function to produce cluster data job and send to message queue
func (p *Producer) ClusterProducer(dimension string) {
	startTime := time.Now()
	jobTime := utils.FormatTime(time.Now(), dimension)
	var err error
	defer func() {
		prom.ReportProduceJobLatencyMetric(types.ClusterType, dimension, err, startTime)
	}()
	cmConn, err := p.CMClient.GetClusterManagerConn()
	if err != nil {
		blog.Errorf("get cm conn error:%v", err)
		return
	}
	defer cmConn.Close()
	cliWithHeader := p.CMClient.NewGrpcClientWithHeader(p.ctx, cmConn)
	clusterList, err := p.resourceGetter.GetClusterIDList(cliWithHeader.Ctx, cliWithHeader.Cli)
	if err != nil || clusterList == nil {
		blog.Errorf("get clusterList error: %v", err)
		return
	}
	for _, cluster := range clusterList {
		opts := types.JobCommonOpts{
			ProjectID:   cluster.ProjectID,
			BusinessID:  cluster.BusinessID,
			ClusterID:   cluster.ClusterID,
			ClusterType: cluster.ClusterType,
			CurrentTime: jobTime,
			Dimension:   dimension,
			ObjectType:  types.ClusterType,
		}
		err := p.SendJob(opts)
		if err != nil {
			blog.Errorf("send cluster job to msg queue error, opts: %v, err: %v", opts, err)
			return
		}
	}
	blog.Infof("[producer] send cluster job success, count: %d, jobTime:%v, startTime:%v, currentTime:%v, cost:%v",
		len(clusterList), jobTime, startTime, time.Now(), time.Now().Sub(startTime))
}

// NamespaceProducer is the function to produce namespace data job and send to message queue
func (p *Producer) NamespaceProducer(dimension string) {
	startTime := time.Now()
	jobTime := utils.FormatTime(time.Now(), dimension)
	var err error
	defer func() {
		prom.ReportProduceJobLatencyMetric(types.NamespaceType, dimension, err, startTime)
	}()
	cmConn, err := p.CMClient.GetClusterManagerConn()
	if err != nil {
		blog.Errorf("get cm conn error:%v", err)
		return
	}
	defer cmConn.Close()
	cliWithHeader := p.CMClient.NewGrpcClientWithHeader(p.ctx, cmConn)
	namespaceList, err := p.resourceGetter.GetNamespaceList(cliWithHeader.Ctx, cliWithHeader.Cli,
		p.k8sStorageCli, p.mesosStorageCli)
	if err != nil || namespaceList == nil {
		blog.Errorf("get namespace list error: %v", err)
		return
	}
	for _, namespace := range namespaceList {
		opts := types.JobCommonOpts{
			ClusterID:   namespace.ClusterID,
			ProjectID:   namespace.ProjectID,
			BusinessID:  namespace.BusinessID,
			ClusterType: namespace.ClusterType,
			Namespace:   namespace.Name,
			CurrentTime: jobTime,
			Dimension:   dimension,
			ObjectType:  types.NamespaceType,
		}
		err := p.SendJob(opts)
		if err != nil {
			blog.Errorf("send namespace job to msg queue error, opts: %v, err: %v", opts, err)
			return
		}
	}
	blog.Infof("[producer] send all namespace job, count:%d, jobTime:%v, startTime:%v, "+
		"currentTime:%v, cost:%v", len(namespaceList), jobTime, startTime, time.Now(), time.Now().Sub(startTime))
}

// WorkloadProducer is the function to produce workload data job and send to message queue
func (p *Producer) WorkloadProducer(dimension string) {
	startTime := time.Now()
	jobTime := utils.FormatTime(time.Now(), dimension)
	var err error
	defer func() {
		prom.ReportProduceJobLatencyMetric(types.WorkloadType, dimension, err, startTime)
	}()
	cmConn, err := p.CMClient.GetClusterManagerConn()
	if err != nil {
		blog.Errorf("get cm conn error:%v", err)
		return
	}
	defer cmConn.Close()
	cliWithHeader := p.CMClient.NewGrpcClientWithHeader(p.ctx, cmConn)
	clusterList, err := p.resourceGetter.GetClusterIDList(cliWithHeader.Ctx, cliWithHeader.Cli)
	if err != nil || clusterList == nil {
		blog.Errorf("get clusterList error: %v", err)
		return
	}
	var totalWorkload int
	countCh := make(chan int, 200)
	go func() {
		for count := range countCh {
			totalWorkload = totalWorkload + count
		}
	}()
	chPool := make(chan struct{}, p.concurrency)
	blog.Infof("[producer] concurrency:%d", p.concurrency)
	wg := sync.WaitGroup{}
	for key := range clusterList {
		chPool <- struct{}{}
		wg.Add(1)
		go func(clusterMeta *types.ClusterMeta) {
			workloadList := make([]*types.WorkloadMeta, 0)
			defer func() {
				wg.Done()
				<-chPool
				countCh <- len(workloadList)
			}()
			switch clusterMeta.ClusterType {
			case types.Kubernetes:
				namespaceList, err := p.resourceGetter.GetNamespaceListByCluster(clusterMeta, p.k8sStorageCli, p.mesosStorageCli)
				if err != nil {
					blog.Errorf("get workload list error: %v", err)
					return
				}
				if workloadList, err = p.resourceGetter.GetK8sWorkloadList(namespaceList, p.k8sStorageCli); err != nil {
					blog.Errorf("get workload list error: %v", err)
					return
				}
			case types.Mesos:
				if workloadList, err = p.resourceGetter.GetMesosWorkloadList(clusterMeta, p.mesosStorageCli); err != nil {
					blog.Errorf("get workload list error: %v", err)
					return
				}
			}
			for _, workload := range workloadList {
				opts := types.JobCommonOpts{
					ProjectID:    workload.ProjectID,
					BusinessID:   workload.BusinessID,
					ClusterID:    workload.ClusterID,
					ClusterType:  workload.ClusterType,
					Namespace:    workload.Namespace,
					WorkloadType: workload.ResourceType,
					Name:         workload.Name,
					CurrentTime:  jobTime,
					Dimension:    dimension,
					ObjectType:   types.WorkloadType,
				}
				if err = p.SendJob(opts); err != nil {
					blog.Errorf("send workload job to msg queue error, opts: %v, err: %v", opts, err)
					return
				}
			}
			blog.Infof("[producer] send workload job success, count: %d", len(workloadList))
		}(clusterList[key])
	}
	wg.Wait()
	close(chPool)
	time.Sleep(100 * time.Microsecond)
	close(countCh)
	blog.Infof("[producer] send all workload job, count:%d, jobTime:%v, startTime:%v, "+
		"currentTime:%v, cost:%v", totalWorkload, jobTime, startTime, time.Now(), time.Now().Sub(startTime))
}

// SendJob is the function to send data job to msg queue
func (p *Producer) SendJob(opts types.JobCommonOpts) error {
	var err error
	defer func() {
		prom.ReportProduceJobTotalMetric(opts.ObjectType, opts.Dimension, err)
	}()
	dataJob := datajob.DataJob{Opts: opts}
	msg := &broker.Message{Header: map[string]string{
		"resourceType": types.DataJobQueue,
		"clusterId":    "dataManager",
	}}
	err = codec.EncJson(dataJob, &msg.Body)
	if err != nil {
		blog.Errorf("transfer dataJob to msg body error, dataJob: %v, error: %v", dataJob, err)
		return err
	}
	err = p.msgQueue.Publish(msg)
	if err != nil {
		blog.Errorf("send message error: %v", err)
		return err
	}
	return nil
}
