package model

import (
	"github.com/sirupsen/logrus"
	frameworkmodel "slime.io/slime/framework/model"
)

const ModuleName = "pilotadmin"

var ModuleLog = logrus.WithField(frameworkmodel.LogFieldKeyModule, ModuleName)
