/*
Copyright 2018 Weaveworks Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1beta1

import (
	flux_weave_works_v1beta1 "github.com/weaveworks/flux/integrations/apis/flux.weave.works/v1beta1"
	versioned "github.com/weaveworks/flux/integrations/client/clientset/versioned"
	internalinterfaces "github.com/weaveworks/flux/integrations/client/informers/externalversions/internalinterfaces"
	v1beta1 "github.com/weaveworks/flux/integrations/client/listers/flux.weave.works/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// HelmReleaseInformer provides access to a shared informer and lister for
// HelmReleases.
type HelmReleaseInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.HelmReleaseLister
}

type helmReleaseInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewHelmReleaseInformer constructs a new informer for HelmRelease type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewHelmReleaseInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredHelmReleaseInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredHelmReleaseInformer constructs a new informer for HelmRelease type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredHelmReleaseInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.FluxV1beta1().HelmReleases(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.FluxV1beta1().HelmReleases(namespace).Watch(options)
			},
		},
		&flux_weave_works_v1beta1.HelmRelease{},
		resyncPeriod,
		indexers,
	)
}

func (f *helmReleaseInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredHelmReleaseInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *helmReleaseInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&flux_weave_works_v1beta1.HelmRelease{}, f.defaultInformer)
}

func (f *helmReleaseInformer) Lister() v1beta1.HelmReleaseLister {
	return v1beta1.NewHelmReleaseLister(f.Informer().GetIndexer())
}
