package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	store "kmodules.xyz/objectstore-api/api/v1"
)

const (
	ResourceKindRecovery     = "Recovery"
	ResourceSingularRecovery = "recovery"
	ResourcePluralRecovery   = "recoveries"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Recovery struct {
	metav1.TypeMeta   `json:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RecoverySpec   `json:"spec,omitempty"`
	Status            RecoveryStatus `json:"status,omitempty"`
}

type RecoverySpec struct {
	Repository core.ObjectReference `json:"repository"`
	// Snapshot to recover. Default is latest snapshot.
	// +optional
	Snapshot            string         `json:"snapshot,omitempty"`
	Paths               []string       `json:"paths,omitempty"`
	RecoverTo           RecoveryTarget `json:"recoverTo,omitempty"`
	RecoveryPolicy      `json:"recoveryPolicy,omitempty"`
	ContainerAttributes *core.Container `json:"containerAttributes,omitempty"`

	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	NodeSelector     map[string]string           `json:"nodeSelector,omitempty"`
	ImagePullSecrets []core.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

type RecoveryTarget struct {
	Volume    *store.LocalSpec      `json:"volume,omitempty"`
	ObjectRef *core.ObjectReference `json:"objectRef,omitempty"`
}

type RecoveryPolicy string

const (
	RecoveryPolicyIfNotRecovered RecoveryPolicy = "IfNotRecovered"
	RecoveryPolicyAlways         RecoveryPolicy = "Always"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RecoveryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Recovery `json:"items,omitempty"`
}

type RecoveryPhase string

const (
	RecoveryPending   RecoveryPhase = "Pending"
	RecoveryRunning   RecoveryPhase = "Running"
	RecoverySucceeded RecoveryPhase = "Succeeded"
	RecoveryFailed    RecoveryPhase = "Failed"
	RecoveryUnknown   RecoveryPhase = "Unknown"
)

type RecoveryStatus struct {
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
	Phase              RecoveryPhase  `json:"phase,omitempty"`
	Stats              []RestoreStats `json:"stats,omitempty"`
}

type RestoreStats struct {
	Path     string        `json:"path,omitempty"`
	Phase    RecoveryPhase `json:"phase,omitempty"`
	Duration string        `json:"duration,omitempty"`
}
