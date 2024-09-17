package v3

import (
	"context"

	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type dockerCredentialLifecycleConverter struct {
	lifecycle DockerCredentialLifecycle
}

func (w *dockerCredentialLifecycleConverter) CreateContext(_ context.Context, obj *v3.DockerCredential) (runtime.Object, error) {
	return w.lifecycle.Create(obj)
}

func (w *dockerCredentialLifecycleConverter) RemoveContext(_ context.Context, obj *v3.DockerCredential) (runtime.Object, error) {
	return w.lifecycle.Remove(obj)
}

func (w *dockerCredentialLifecycleConverter) UpdatedContext(_ context.Context, obj *v3.DockerCredential) (runtime.Object, error) {
	return w.lifecycle.Updated(obj)
}

type DockerCredentialLifecycle interface {
	Create(obj *v3.DockerCredential) (runtime.Object, error)
	Remove(obj *v3.DockerCredential) (runtime.Object, error)
	Updated(obj *v3.DockerCredential) (runtime.Object, error)
}

type DockerCredentialLifecycleContext interface {
	CreateContext(ctx context.Context, obj *v3.DockerCredential) (runtime.Object, error)
	RemoveContext(ctx context.Context, obj *v3.DockerCredential) (runtime.Object, error)
	UpdatedContext(ctx context.Context, obj *v3.DockerCredential) (runtime.Object, error)
}

type dockerCredentialLifecycleAdapter struct {
	lifecycle DockerCredentialLifecycleContext
}

func (w *dockerCredentialLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *dockerCredentialLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *dockerCredentialLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	return w.CreateContext(context.Background(), obj)
}

func (w *dockerCredentialLifecycleAdapter) CreateContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.CreateContext(ctx, obj.(*v3.DockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dockerCredentialLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	return w.FinalizeContext(context.Background(), obj)
}

func (w *dockerCredentialLifecycleAdapter) FinalizeContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.RemoveContext(ctx, obj.(*v3.DockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *dockerCredentialLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	return w.UpdatedContext(context.Background(), obj)
}

func (w *dockerCredentialLifecycleAdapter) UpdatedContext(ctx context.Context, obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.UpdatedContext(ctx, obj.(*v3.DockerCredential))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewDockerCredentialLifecycleAdapter(name string, clusterScoped bool, client DockerCredentialInterface, l DockerCredentialLifecycle) DockerCredentialHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(DockerCredentialGroupVersionResource)
	}
	adapter := &dockerCredentialLifecycleAdapter{lifecycle: &dockerCredentialLifecycleConverter{lifecycle: l}}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.DockerCredential) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}

func NewDockerCredentialLifecycleAdapterContext(name string, clusterScoped bool, client DockerCredentialInterface, l DockerCredentialLifecycleContext) DockerCredentialHandlerContextFunc {
	if clusterScoped {
		resource.PutClusterScoped(DockerCredentialGroupVersionResource)
	}
	adapter := &dockerCredentialLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapterContext(name, clusterScoped, adapter, client.ObjectClient())
	return func(ctx context.Context, key string, obj *v3.DockerCredential) (runtime.Object, error) {
		newObj, err := syncFn(ctx, key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
