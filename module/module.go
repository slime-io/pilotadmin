package module

import (
	"os"

	"slime.io/slime/framework/model/module"
	"slime.io/slime/modules/pilotadmin/model"

	"github.com/golang/protobuf/proto"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	istionetworkingapi "slime.io/slime/framework/apis/networking/v1alpha3"
	"slime.io/slime/framework/bootstrap"
	modapi "slime.io/slime/modules/pilotadmin/api/v1alpha1"
	"slime.io/slime/modules/pilotadmin/controllers"
)

var log = model.ModuleLog

type Module struct {
}

func (m *Module) Name() string {
	return model.ModuleName
}

func (m *Module) Config() proto.Message {
	return nil
}

func (m *Module) InitScheme(scheme *runtime.Scheme) error {
	for _, f := range []func(*runtime.Scheme) error{
		clientgoscheme.AddToScheme,
		modapi.AddToScheme,
		istionetworkingapi.AddToScheme,
	} {
		if err := f(scheme); err != nil {
			return err
		}
	}
	return nil
}

func (m *Module) InitManager(mgr manager.Manager, env bootstrap.Environment, cbs module.InitCallbacks) error {
	var err error
	if err = (controllers.NewReconciler(mgr, &env)).SetupWithManager(mgr); err != nil {
		log.Errorf("unable to create pilotadmin controller, %+v", err)
		os.Exit(1)
	}

	return nil
}
