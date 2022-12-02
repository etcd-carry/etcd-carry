package schema

import (
	arv1 "k8s.io/api/admissionregistration/v1"
	arv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	asv1 "k8s.io/api/autoscaling/v1"
	asv2beta1 "k8s.io/api/autoscaling/v2beta1"
	asv2beta2 "k8s.io/api/autoscaling/v2beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	certv1beta1 "k8s.io/api/certificates/v1beta1"
	cov1 "k8s.io/api/coordination/v1"
	cov1beta1 "k8s.io/api/coordination/v1beta1"
	corev1 "k8s.io/api/core/v1"
	dcv1alpha1 "k8s.io/api/discovery/v1alpha1"
	dcv1beta1 "k8s.io/api/discovery/v1beta1"
	ev1beta1 "k8s.io/api/events/v1beta1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	fv1alpha1 "k8s.io/api/flowcontrol/v1alpha1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	nv1alpha1 "k8s.io/api/node/v1alpha1"
	nv1beta1 "k8s.io/api/node/v1beta1"
	pv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	schv1 "k8s.io/api/scheduling/v1"
	schv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	schv1beta1 "k8s.io/api/scheduling/v1beta1"
	sv1 "k8s.io/api/storage/v1"
	sv1alpha1 "k8s.io/api/storage/v1alpha1"
	sv1beta1 "k8s.io/api/storage/v1beta1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	aarv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	aarv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	//external
	aarv1.AddToScheme(Scheme)
	aarv1beta1.AddToScheme(Scheme)
	apiextv1.AddToScheme(Scheme)
	apiextv1beta1.AddToScheme(Scheme)
	arv1.AddToScheme(Scheme)
	arv1beta1.AddToScheme(Scheme)
	appsv1.AddToScheme(Scheme)
	appsv1beta1.AddToScheme(Scheme)
	appsv1beta2.AddToScheme(Scheme)
	asv1.AddToScheme(Scheme)
	asv2beta1.AddToScheme(Scheme)
	asv2beta2.AddToScheme(Scheme)
	batchv1.AddToScheme(Scheme)
	batchv1beta1.AddToScheme(Scheme)
	batchv2alpha1.AddToScheme(Scheme)
	certv1beta1.AddToScheme(Scheme)
	cov1.AddToScheme(Scheme)
	cov1beta1.AddToScheme(Scheme)
	corev1.AddToScheme(Scheme)
	dcv1alpha1.AddToScheme(Scheme)
	dcv1beta1.AddToScheme(Scheme)
	ev1beta1.AddToScheme(Scheme)
	extv1beta1.AddToScheme(Scheme)
	fv1alpha1.AddToScheme(Scheme)
	netv1.AddToScheme(Scheme)
	netv1beta1.AddToScheme(Scheme)
	nv1alpha1.AddToScheme(Scheme)
	nv1beta1.AddToScheme(Scheme)
	pv1beta1.AddToScheme(Scheme)
	rbacv1.AddToScheme(Scheme)
	rbacv1alpha1.AddToScheme(Scheme)
	rbacv1beta1.AddToScheme(Scheme)
	schv1.AddToScheme(Scheme)
	schv1alpha1.AddToScheme(Scheme)
	schv1beta1.AddToScheme(Scheme)
	sv1.AddToScheme(Scheme)
	sv1alpha1.AddToScheme(Scheme)
	sv1beta1.AddToScheme(Scheme)
}
