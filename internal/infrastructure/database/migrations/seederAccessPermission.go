package migrations

import (
	"log"
	"middleware/internal/domain/constants"

	"github.com/casbin/casbin/v3"
)

var accessPolicies = [][]string{
	{constants.UserProfileAdministrador, "/*", "(GET)|(POST)|(PATCH)|(DELETE)|(PUT)|(OPTIONS)"},
	{constants.UserProfileOperador, "/clps", "(GET)|(PATCH)"},
	{constants.UserProfileOperador, "/clps/*", "(GET)|(PATCH)"},
	{constants.UserProfileOperador, "/tags", "(GET)|(PATCH)"},
	{constants.UserProfileOperador, "/tags/*", "(GET)|(PATCH)"},
	{constants.UserProfileOperador, "/type-clps", "GET"},
	{constants.UserProfileOperador, "/type-clps/", "GET"},
	{constants.UserProfileOperador, "/type-tags", "GET"},
	{constants.UserProfileOperador, "/type-tags/", "GET"},
	{constants.UserProfileOperador, "/type-operations", "GET"},
	{constants.UserProfileOperador, "/type-operations/", "GET"},
	{constants.UserProfileOperador, "/swaps", "GET"},
	{constants.UserProfileOperador, "/swaps/", "GET"},
	{constants.UserProfileOperador, "/users/password", "PATCH"},
	{constants.UserProfileOperador, "/users/password/", "PATCH"},
}

var legacyAccessPolicies = [][]string{
	{constants.UserProfileOperador, "/clps", "(GET)|(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileOperador, "/clps/*", "(GET)|(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileOperador, "/tags", "(GET)|(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileOperador, "/tags/*", "(GET)|(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileAdmLegacy, "/*", "(GET)|(POST)|(PATCH)|(DELETE)|(PUT)|(OPTIONS)"},
	{constants.UserProfileOperador, "/clps", "(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileOperador, "/clps/*", "(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileOperador, "/tags", "(POST)|(PATCH)|(DELETE)"},
	{constants.UserProfileOperador, "/tags/*", "(POST)|(PATCH)|(DELETE)"},
}

func InicializaPermissoesAcesso(enforcer *casbin.Enforcer) {
	removerPoliticas(legacyAccessPolicies, enforcer)
	adicionarPoliticas(accessPolicies, enforcer)
}

func adicionarPoliticas(politicas [][]string, enforcer *casbin.Enforcer) {
	for _, politica := range politicas {
		if existePolitica, _ := enforcer.HasPolicy(politica); !existePolitica {
			ok, err := enforcer.AddPolicy(politica)
			if err != nil {
				// utils.ErroLog.Printf("Erro ao adicionar politica: %v", err)
			} else if ok {
				log.Printf("Politica adicionada: %v", politica)
			}
		}
	}
}

func removerPoliticas(politicas [][]string, enforcer *casbin.Enforcer) {
	for _, politica := range politicas {
		if existePolitica, _ := enforcer.HasPolicy(politica); existePolitica {
			ok, err := enforcer.RemovePolicy(politica)
			if err != nil {
				// utils.ErroLog.Printf("Erro ao remover politica: %v", err)
			} else if ok {
				log.Printf("Politica removida: %v", politica)
			}
		}
	}
}
