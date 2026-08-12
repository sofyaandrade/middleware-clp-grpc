package migrations

import (
	"log"

	"github.com/casbin/casbin/v3"
)

var admPolicy = [][]string{
	{"ADMINISTRADOR", "/*", "(GET)|(POST)|(PATCH)|(DELETE)|(PUT)|(OPTIONS)"},
}

func InicializaPermissoesAcesso(enforcer *casbin.Enforcer) {
	adicionarPoliticas(admPolicy, enforcer)
}

func adicionarPoliticas(politicas [][]string, enforcer *casbin.Enforcer) {
	for _, politica := range politicas {
		if existePolitica, _ := enforcer.HasPolicy(politica); !existePolitica {
			ok, err := enforcer.AddPolicy(politica)
			if err != nil {
				// utils.ErroLog.Printf("Erro ao adicionar política: %v", err)
			} else if ok {
				log.Printf("Política adicionada: %v", politica)
			}
		}
	}
}
