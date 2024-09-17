package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type multiClusterAppRevisionLifecycleConverter struct {
	lifecycle MultiClusterAppRevisionLifecycle
}

func (w *multiClusterAppRevisionLifecycleConverter) CreateContext(_ context.Context, obj *v3.MultiClusterAppRevision) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *multiClusterAppRevisionLifecycleConverter) RemoveContext(_ context.Context, obj *v3.MultiClusterAppRevision) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *multiClusterAppRevisionLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.MultiClusterAppRevision) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type MultiClusterAppRevisionLifecycle interface {
	Create(obj *v3.MultiClusterAppRevision) (runtime.Object, error)
	Remove(obj *v3.MultiClusterAppRevision) (runtime.Object, error)
	Updated(obj *v3.MultiClusterAppRevision) (runtime.Object, error)
}

type MultiClusterAppRevisionLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.MultiClusterAppRevision) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.MultiClusterAppRevision) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.MultiClusterAppRevision) (runtime.Object, error)
}

type multiClusterAppRevisionLifecycleAdapter struct {
	lifecycle MultiClusterAppRevisionLifecycleContext
}

func (w *multiClusterAppRevisionLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *multiClusterAppRevisionLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *multiClusterAppRevisionLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *multiClusterAppRevisionLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.MultiClusterAppRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *multiClusterAppRevisionLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *multiClusterAppRevisionLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.MultiClusterAppRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *multiClusterAppRevisionLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *multiClusterAppRevisionLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.MultiClusterAppRevision))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewMultiClusterAppRevisionLifecycleAdapter(name string, clusterScoped bool, client MultiClusterAppRevisionInterface, l MultiClusterAppRevisionLifecycle) MultiClusterAppRevisionHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(MultiClusterAppRevisionGroupVersionResource)
	}
	adapter := &multiClusterAppRevisionLifecycleAdapter{lifecycle: &multiClusterAppRevisionLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.MultiClusterAppRevision) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewMultiClusterAppRevisionLifecycleAdapterContext(name string, clusterScoped bool, client MultiClusterAppRevisionInterface, l MultiClusterAppRevisionLifecycleContext) MultiClusterAppRevisionHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(MultiClusterAppRevisionGroupVersionResource)
	}
	adapter := &multiClusterAppRevisionLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.MultiClusterAppRevision) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
