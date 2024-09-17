package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespacedDockerCredentialLifecycleConverter struct {
	lifecycle NamespacedDockerCredentialLifecycle
}

func (w *namespacedDockerCredentialLifecycleConverter) CreateContext(_ context.Context, obj *v3.NamespacedDockerCredential) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *namespacedDockerCredentialLifecycleConverter) RemoveContext(_ context.Context, obj *v3.NamespacedDockerCredential) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *namespacedDockerCredentialLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.NamespacedDockerCredential) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type NamespacedDockerCredentialLifecycle interface {
	Create(obj *v3.NamespacedDockerCredential) (runtime.Object, error)
	Remove(obj *v3.NamespacedDockerCredential) (runtime.Object, error)
	Updated(obj *v3.NamespacedDockerCredential) (runtime.Object, error)
}

type NamespacedDockerCredentialLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.NamespacedDockerCredential) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.NamespacedDockerCredential) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.NamespacedDockerCredential) (runtime.Object, error)
}

type namespacedDockerCredentialLifecycleAdapter struct {
	lifecycle NamespacedDockerCredentialLifecycleContext
}

func (w *namespacedDockerCredentialLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *namespacedDockerCredentialLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *namespacedDockerCredentialLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *namespacedDockerCredentialLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.NamespacedDockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespacedDockerCredentialLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *namespacedDockerCredentialLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.NamespacedDockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *namespacedDockerCredentialLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *namespacedDockerCredentialLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.NamespacedDockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewNamespacedDockerCredentialLifecycleAdapter(name string, clusterScoped bool, client NamespacedDockerCredentialInterface, l NamespacedDockerCredentialLifecycle) NamespacedDockerCredentialHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespacedDockerCredentialGroupVersionResource)
	}
	adapter := &namespacedDockerCredentialLifecycleAdapter{lifecycle: &namespacedDockerCredentialLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.NamespacedDockerCredential) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewNamespacedDockerCredentialLifecycleAdapterContext(name string, clusterScoped bool, client NamespacedDockerCredentialInterface, l NamespacedDockerCredentialLifecycleContext) NamespacedDockerCredentialHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(NamespacedDockerCredentialGroupVersionResource)
	}
	adapter := &namespacedDockerCredentialLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.NamespacedDockerCredential) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
