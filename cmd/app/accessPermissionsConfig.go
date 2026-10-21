package app

import (
	"fmt"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/database/migrations"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var rbacModel string = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`

func AccessPermissionsConfig(db *gorm.DB) *casbin.Enforcer {
	if db == nil {
		fmt.Println("error.inicialize.casbin database connection is nil")
		return nil
	}

	gormadapter.TurnOffAutoMigrate(db)

	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(db, &models.PermissoesAcesso{}, "permissoes_acesso")
	if err != nil {
		fmt.Println("error.inicialize.casbin " + err.Error())
		return nil
	}

	model, err := model.NewModelFromString(rbacModel)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	enforcer, err := casbin.NewEnforcer(model, adapter)
	if err != nil {
		fmt.Println("error.inicialize.casbin" + err.Error())
		return nil
	}

	if err := enforcer.LoadPolicy(); err != nil {
		fmt.Println("error.inicialize.casbin " + err.Error())
		return nil
	}

	migrations.InicializaPermissoesAcesso(enforcer)

	return enforcer
}
