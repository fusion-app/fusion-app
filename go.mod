module github.com/fusion-app/fusion-app

require (
	github.com/NYTimes/gziphandler v1.0.1 // indirect
	github.com/Shopify/sarama v1.19.0
	github.com/bsm/sarama-cluster v2.1.15+incompatible
	github.com/go-openapi/spec v0.19.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.6.2
	github.com/iancoleman/strcase v0.0.0-20180726023541-3605ed457bf7
	github.com/imdario/mergo v0.3.7
	github.com/kubeflow/pytorch-operator v0.4.0
	github.com/kubeflow/tf-operator v0.5.1
	github.com/kubeless/kubeless v1.0.5
	github.com/lxs137/fusion-app v0.0.0-20191108124951-c3fc93c5960a // indirect
	github.com/operator-framework/operator-sdk v0.10.2-0.20191010224636-fd8747add695
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	k8s.io/api v0.0.0-20190612125737-db0771252981
	k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery v0.0.0-20190612125636-6a5db36e93ad
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20190603182131-db7b694dc208
	sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools v0.1.10
)

// Pinned to kubernetes-1.13.4
replace (
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190228180357-d002e88f6236
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190228174230-b40b2a5939e4
)

replace (
	github.com/coreos/prometheus-operator => github.com/coreos/prometheus-operator v0.29.0
	// Pinned to v2.9.2 (kubernetes-1.13.1) so https://proxy.golang.org can
	// resolve it correctly.
	github.com/prometheus/prometheus => github.com/prometheus/prometheus v0.0.0-20190424153033-d3245f150225
	k8s.io/kube-state-metrics => k8s.io/kube-state-metrics v1.6.0
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.1.12
	sigs.k8s.io/controller-tools => sigs.k8s.io/controller-tools v0.1.11-0.20190411181648-9d55346c2bde
)

replace github.com/operator-framework/operator-sdk => github.com/operator-framework/operator-sdk v0.10.1
