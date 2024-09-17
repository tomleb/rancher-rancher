package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type authConfigLifecycleConverter struct {
	lifecycle AuthConfigLifecycle
}

func (w *authConfigLifecycleConverter) CreateContext(_ context.Context, obj *v3.AuthConfig) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *authConfigLifecycleConverter) RemoveContext(_ context.Context, obj *v3.AuthConfig) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *authConfigLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.AuthConfig) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type AuthConfigLifecycle interface {
	Create(obj *v3.AuthConfig) (runtime.Object, error)
	Remove(obj *v3.AuthConfig) (runtime.Object, error)
	Updated(obj *v3.AuthConfig) (runtime.Object, error)
}

type AuthConfigLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.AuthConfig) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.AuthConfig) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.AuthConfig) (runtime.Object, error)
}

type authConfigLifecycleAdapter struct {
	lifecycle AuthConfigLifecycleContext
}

func (w *authConfigLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *authConfigLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *authConfigLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *authConfigLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.AuthConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *authConfigLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *authConfigLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.AuthConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *authConfigLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *authConfigLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.AuthConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewAuthConfigLifecycleAdapter(name string, clusterScoped bool, client AuthConfigInterface, l AuthConfigLifecycle) AuthConfigHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(AuthConfigGroupVersionResource)
	}
	adapter := &authConfigLifecycleAdapter{lifecycle: &authConfigLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.AuthConfig) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewAuthConfigLifecycleAdapterContext(name string, clusterScoped bool, client AuthConfigInterface, l AuthConfigLifecycleContext) AuthConfigHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(AuthConfigGroupVersionResource)
	}
	adapter := &authConfigLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.AuthConfig) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
