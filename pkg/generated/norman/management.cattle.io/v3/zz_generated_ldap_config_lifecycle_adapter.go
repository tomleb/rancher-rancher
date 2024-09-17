package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type ldapConfigLifecycleConverter struct {
	lifecycle LdapConfigLifecycle
}

func (w *ldapConfigLifecycleConverter) CreateContext(_ context.Context, obj *v3.LdapConfig) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *ldapConfigLifecycleConverter) RemoveContext(_ context.Context, obj *v3.LdapConfig) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *ldapConfigLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.LdapConfig) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type LdapConfigLifecycle interface {
	Create(obj *v3.LdapConfig) (runtime.Object, error)
	Remove(obj *v3.LdapConfig) (runtime.Object, error)
	Updated(obj *v3.LdapConfig) (runtime.Object, error)
}

type LdapConfigLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.LdapConfig) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.LdapConfig) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.LdapConfig) (runtime.Object, error)
}

type ldapConfigLifecycleAdapter struct {
	lifecycle LdapConfigLifecycleContext
}

func (w *ldapConfigLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *ldapConfigLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *ldapConfigLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *ldapConfigLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.LdapConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *ldapConfigLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *ldapConfigLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.LdapConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *ldapConfigLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *ldapConfigLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.LdapConfig))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewLdapConfigLifecycleAdapter(name string, clusterScoped bool, client LdapConfigInterface, l LdapConfigLifecycle) LdapConfigHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(LdapConfigGroupVersionResource)
	}
	adapter := &ldapConfigLifecycleAdapter{lifecycle: &ldapConfigLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.LdapConfig) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewLdapConfigLifecycleAdapterContext(name string, clusterScoped bool, client LdapConfigInterface, l LdapConfigLifecycleContext) LdapConfigHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(LdapConfigGroupVersionResource)
	}
	adapter := &ldapConfigLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.LdapConfig) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
