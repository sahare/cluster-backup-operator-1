module github.com/open-cluster-management/cluster-backup-operator

go 1.17

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/open-cluster-management/api v0.0.0-20210527013639-a6845f2ebcb1
	github.com/open-cluster-management/library-go v0.0.0-20210805092456-f9c5ab6c4fc1 // indirect
	github.com/open-cluster-management/multicloud-operators-channel v1.2.2-2-20201130-37b47
	github.com/openshift/api v0.0.0-20210521075222-e273a339932a //Openshift 4.6
	github.com/openshift/hive/apis v0.0.0-20210915004009-18827f64c00e // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.0
	github.com/vmware-tanzu/velero v1.6.1
	k8s.io/api v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v12.0.0+incompatible
	open-cluster-management.io/api v0.0.0-20210908005819-815ac23c7308
	sigs.k8s.io/controller-runtime v0.9.1

)

replace (
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.10
	k8s.io/api => k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.3
	k8s.io/client-go => k8s.io/client-go v0.21.3
)
