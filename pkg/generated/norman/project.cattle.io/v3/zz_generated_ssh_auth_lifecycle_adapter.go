package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type sshAuthLifecycleConverter struct {
	lifecycle SSHAuthLifecycle
}

func (w *sshAuthLifecycleConverter) CreateContext(_ context.Context, obj *v3.SSHAuth) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *sshAuthLifecycleConverter) RemoveContext(_ context.Context, obj *v3.SSHAuth) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *sshAuthLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.SSHAuth) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type SSHAuthLifecycle interface {
	Create(obj *v3.SSHAuth) (runtime.Object, error)
	Remove(obj *v3.SSHAuth) (runtime.Object, error)
	Updated(obj *v3.SSHAuth) (runtime.Object, error)
}

type SSHAuthLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.SSHAuth) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.SSHAuth) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.SSHAuth) (runtime.Object, error)
}

type sshAuthLifecycleAdapter struct {
	lifecycle SSHAuthLifecycleContext
}

func (w *sshAuthLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *sshAuthLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *sshAuthLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *sshAuthLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.SSHAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *sshAuthLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *sshAuthLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.SSHAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *sshAuthLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *sshAuthLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.SSHAuth))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewSSHAuthLifecycleAdapter(name string, clusterScoped bool, client SSHAuthInterface, l SSHAuthLifecycle) SSHAuthHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(SSHAuthGroupVersionResource)
	}
	adapter := &sshAuthLifecycleAdapter{lifecycle: &sshAuthLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.SSHAuth) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewSSHAuthLifecycleAdapterContext(name string, clusterScoped bool, client SSHAuthInterface, l SSHAuthLifecycleContext) SSHAuthHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(SSHAuthGroupVersionResource)
	}
	adapter := &sshAuthLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.SSHAuth) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
