/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v2

import (
	v2 "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// HorizontalPodAutoscalerLister helps list HorizontalPodAutoscalers.
// All objects returned here must be treated as read-only.
type HorizontalPodAutoscalerLister interface {
	// List lists all HorizontalPodAutoscalers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.HorizontalPodAutoscaler, err error)
	// HorizontalPodAutoscalers returns an object that can list and get HorizontalPodAutoscalers.
	HorizontalPodAutoscalers(namespace string) HorizontalPodAutoscalerNamespaceLister
	HorizontalPodAutoscalerListerExpansion
}

// horizontalPodAutoscalerLister implements the HorizontalPodAutoscalerLister interface.
type horizontalPodAutoscalerLister struct {
	listers.ResourceIndexer[*v2.HorizontalPodAutoscaler]
}

// NewHorizontalPodAutoscalerLister returns a new HorizontalPodAutoscalerLister.
func NewHorizontalPodAutoscalerLister(indexer cache.Indexer) HorizontalPodAutoscalerLister {
	return &horizontalPodAutoscalerLister{listers.New[*v2.HorizontalPodAutoscaler](indexer, v2.Resource("horizontalpodautoscaler"))}
}

// HorizontalPodAutoscalers returns an object that can list and get HorizontalPodAutoscalers.
func (s *horizontalPodAutoscalerLister) HorizontalPodAutoscalers(namespace string) HorizontalPodAutoscalerNamespaceLister {
	return horizontalPodAutoscalerNamespaceLister{listers.NewNamespaced[*v2.HorizontalPodAutoscaler](s.ResourceIndexer, namespace)}
}

// HorizontalPodAutoscalerNamespaceLister helps list and get HorizontalPodAutoscalers.
// All objects returned here must be treated as read-only.
type HorizontalPodAutoscalerNamespaceLister interface {
	// List lists all HorizontalPodAutoscalers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2.HorizontalPodAutoscaler, err error)
	// Get retrieves the HorizontalPodAutoscaler from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v2.HorizontalPodAutoscaler, error)
	HorizontalPodAutoscalerNamespaceListerExpansion
}

// horizontalPodAutoscalerNamespaceLister implements the HorizontalPodAutoscalerNamespaceLister
// interface.
type horizontalPodAutoscalerNamespaceLister struct {
	listers.ResourceIndexer[*v2.HorizontalPodAutoscaler]
}