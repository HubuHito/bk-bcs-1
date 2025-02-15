/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	tkexv1alpha1 "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-gamedeployment-operator/pkg/apis/tkex/v1alpha1"
	versioned "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-gamedeployment-operator/pkg/client/clientset/versioned"
	internalinterfaces "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-gamedeployment-operator/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-gamedeployment-operator/pkg/client/listers/tkex/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GameDeploymentInformer provides access to a shared informer and lister for
// GameDeployments.
type GameDeploymentInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.GameDeploymentLister
}

type gameDeploymentInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGameDeploymentInformer constructs a new informer for GameDeployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGameDeploymentInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGameDeploymentInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGameDeploymentInformer constructs a new informer for GameDeployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGameDeploymentInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TkexV1alpha1().GameDeployments(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TkexV1alpha1().GameDeployments(namespace).Watch(options)
			},
		},
		&tkexv1alpha1.GameDeployment{},
		resyncPeriod,
		indexers,
	)
}

func (f *gameDeploymentInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGameDeploymentInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gameDeploymentInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&tkexv1alpha1.GameDeployment{}, f.defaultInformer)
}

func (f *gameDeploymentInformer) Lister() v1alpha1.GameDeploymentLister {
	return v1alpha1.NewGameDeploymentLister(f.Informer().GetIndexer())
}
