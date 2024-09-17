package v1

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

type apiServiceLifecycleConverter struct {
	lifecycle APIServiceLifecycle
}

func (w *apiServiceLifecycleConverter) CreateContext(_ context.Context, obj *v1.APIService) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *apiServiceLifecycleConverter) RemoveContext(_ context.Context, obj *v1.APIService) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *apiServiceLifecycleConverter) UpdatedContext(_ context.Context, obj *v1.APIService) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type APIServiceLifecycle interface {
	Create(obj *v1.APIService) (runtime.Object, error)
	Remove(obj *v1.APIService) (runtime.Object, error)
	Updated(obj *v1.APIService) (runtime.Object, error)
}

type APIServiceLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v1.APIService) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v1.APIService) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v1.APIService) (runtime.Object, error)
}

type apiServiceLifecycleAdapter struct {
	lifecycle APIServiceLifecycleContext
}

func (w *apiServiceLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *apiServiceLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *apiServiceLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *apiServiceLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v1.APIService))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *apiServiceLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *apiServiceLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v1.APIService))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *apiServiceLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *apiServiceLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v1.APIService))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewAPIServiceLifecycleAdapter(name string, clusterScoped bool, client APIServiceInterface, l APIServiceLifecycle) APIServiceHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(APIServiceGroupVersionResource)
	}
	adapter := &apiServiceLifecycleAdapter{lifecycle: &apiServiceLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v1.APIService) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewAPIServiceLifecycleAdapterContext(name string, clusterScoped bool, client APIServiceInterface, l APIServiceLifecycleContext) APIServiceHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(APIServiceGroupVersionResource)
	}
	adapter := &apiServiceLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v1.APIService) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
