package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type tokenLifecycleConverter struct {
	lifecycle TokenLifecycle
}

func (w *tokenLifecycleConverter) CreateContext(_ context.Context, obj *v3.Token) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *tokenLifecycleConverter) RemoveContext(_ context.Context, obj *v3.Token) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *tokenLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.Token) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type TokenLifecycle interface {
	Create(obj *v3.Token) (runtime.Object, error)
	Remove(obj *v3.Token) (runtime.Object, error)
	Updated(obj *v3.Token) (runtime.Object, error)
}

type TokenLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.Token) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.Token) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.Token) (runtime.Object, error)
}

type tokenLifecycleAdapter struct {
	lifecycle TokenLifecycleContext
}

func (w *tokenLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *tokenLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *tokenLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *tokenLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.Token))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *tokenLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *tokenLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.Token))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *tokenLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *tokenLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.Token))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewTokenLifecycleAdapter(name string, clusterScoped bool, client TokenInterface, l TokenLifecycle) TokenHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(TokenGroupVersionResource)
	}
	adapter := &tokenLifecycleAdapter{lifecycle: &tokenLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Token) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewTokenLifecycleAdapterContext(name string, clusterScoped bool, client TokenInterface, l TokenLifecycleContext) TokenHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(TokenGroupVersionResource)
	}
	adapter := &tokenLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.Token) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
