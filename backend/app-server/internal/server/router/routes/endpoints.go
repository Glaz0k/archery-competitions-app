package routes

const (
	IndividualGroupEndpoint            = "/api/individual_groups/{group_id}"
	IndividualGroupCompetitorsEndpoint = "/api/individual_groups/{group_id}/competitors"
	StartQualification                 = "/api/individual_groups/{group_id}/qualification/start"
	EndQualification                   = "/api/individual_groups/{group_id}/qualification/end"
	UpdateIndividualGroup              = "/api/individual_groups/{group_id}/competitors/sync"
	FinalGrid                          = "/api/individual_groups/{group_id}/final_grid"
	StartQuarterfinal                  = "/api/individual_groups/{group_id}/quarterfinal/start"

	EndCompetition              = "/api/competitions/{competition_id}/end"
	CompetitionEndpoint         = "/api/competitions/{competition_id}"
	IndividualGroupsCompetition = "/api/competitions/{competition_id}/individual_groups"

	CupsEndpoint         = "/api/cups"
	CupEndpoint          = "/api/cups/{cup_id}"
	CompetitionsEndpoint = "/api/cups/{cup_id}/competitions"

	Competitor         = "/api/competitors/{competitor_id}"
	RegisterCompetitor = "/api/competitors/registration"
	GetCompetitor      = "/api/competitors/{competitor_id}"

	CompetitorsCompetitionEndpoint = "/api/competitions/{competition_id}/competitors"
	CompetitorCompetitionEndpoint  = "/api/competitions/{competition_id}/competitors/{competitor_id}"

	GetQualificationSections           = "/api/qualification_sections/{id}"
	QualificationSectionRangesEndpoint = "/api/qualification_sections/{id}/rounds/{round_ordinal}/ranges"
	EndRange                           = "/api/qualification_sections/{id}/rounds/{round_ordinal}/ranges/{range_ordinal}/end"

	QualificationTable = "/api/individual_groups/{group_id}/qualification"
	SparringPlace      = "/api/sparring_places/{id}"
	Ranges             = "/api/sparring_places/{id}/ranges"
	RangeEnd           = "api/sparring_places/{id}/ranges/{range_original}/end"
)
