//go:build generator

package data

import "git.sr.ht/~danihall/foss-eco-simcity-hex-back/internal/domain/core"

// SeedOrdinances is the authoritative data for Ordinances. Update with caution.
var SeedOrdinances = []core.Ordinance{
	{
		Name:        core.DrugPrevention,
		Label:       "campagnes de prévention sur les drogues",
		Description: "Informe les habitants sur les risques liés à la consommation de drogues. Réduit la criminalité liée et améliore le niveau de santé",
	},
	{
		Name:        core.PassCulture,
		Label:       "Pass culture",
		Description: "Alloue un budget pour des achats culturels à chaque jeune habitant. Améliore leur niveau d'éducation",
	},
	{
		Name:        core.VSSCampaign,
		Label:       "campagnes d'information sur les agressions sexistes et sexuelles",
		Description: "Sensibilise les habitants aux problématiques de comportement et actes relevant d'agressions sexistes et sexuelles. Améliore le niveau d'éducation et réduit la criminalité liée aux VSS",
	},
	{
		Name:        core.FreeBicycles,
		Label:       "vélos en libre-service",
		Description: "Fournit des vélos en libre-service. Peut potentiellement réduire la pollution liée à la circulation",
	},
}
