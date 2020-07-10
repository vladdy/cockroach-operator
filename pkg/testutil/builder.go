package testutil

import (
	api "github.com/cockroachdb/cockroach-operator/api/v1alpha1"
	"github.com/cockroachdb/cockroach-operator/pkg/resource"
	corev1 "k8s.io/api/core/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	amtypes "k8s.io/apimachinery/pkg/types"
)

type ClusterBuilder struct {
	cluster api.CrdbCluster
}

func NewBuilder(name string) ClusterBuilder {
	b := ClusterBuilder{
		cluster: api.CrdbCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Labels:      make(map[string]string),
				Annotations: make(map[string]string),
			},
			Spec: api.CrdbClusterSpec{
				Image: "cockroachdb/cockroach:v19.2.6",
			},
		},
	}

	return b
}

func (b ClusterBuilder) Namespaced(namespace string) ClusterBuilder {
	b.cluster.Namespace = namespace
	return b
}

func (b ClusterBuilder) WithUID(uid string) ClusterBuilder {
	b.cluster.ObjectMeta.UID = amtypes.UID(uid)
	return b
}

func (b ClusterBuilder) WithNodeCount(c int32) ClusterBuilder {
	b.cluster.Spec.Nodes = c
	return b
}

func (b ClusterBuilder) WithEmptyDirDataStore() ClusterBuilder {
	b.cluster.Spec.DataStore = api.Volume{EmptyDir: &corev1.EmptyDirVolumeSource{}}
	return b
}

func (b ClusterBuilder) WithPVDataStore(size, storageClass string) ClusterBuilder {
	quantity, _ := apiresource.ParseQuantity(size)

	volumeMode := corev1.PersistentVolumeFilesystem
	b.cluster.Spec.DataStore = api.Volume{
		VolumeClaim: &api.VolumeClaim{
			PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: quantity,
					},
				},
				StorageClassName: &storageClass,
				VolumeMode:       &volumeMode,
			},
		},
	}

	return b
}

func (b ClusterBuilder) WithHTTPPort(port int32) ClusterBuilder {
	b.cluster.Spec.HTTPPort = &port
	return b
}

func (b ClusterBuilder) WithTLS() ClusterBuilder {
	b.cluster.Spec.TLSEnabled = true
	return b
}

func (b ClusterBuilder) WithNodeTLS(secret string) ClusterBuilder {
	b.cluster.Spec.NodeTLSSecret = secret
	return b
}

func (b ClusterBuilder) WithImage(image string) ClusterBuilder {
	b.cluster.Spec.Image = image
	return b
}

func (b ClusterBuilder) Cr() *api.CrdbCluster {
	cluster := b.cluster.DeepCopy()

	return cluster
}

func (b ClusterBuilder) Cluster() *resource.Cluster {
	cluster := resource.NewCluster(b.Cr())
	return &cluster
}
