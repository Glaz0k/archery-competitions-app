package router

const (
	IndividualGroupEndpoint            = "/api/individual-groups/{individual_group_id}"
	IndividualGroupCompetitorsEndpoint = "/api/individual-groups/{group_id}/competitors"
	StartQualification                 = "/api/individual-groups/{individual_group_id}/qualification/start"

	EndCompetition        = "/api/competitions/{competition_id}/end"
	CompetitionEndpoint   = "/api/competitions/{competition_id}"
	CreateIndividualGroup = "/api/competitions/{competition_id}/individual-group"

	CupsEndpoint         = "/api/cups"
	CupEndpoint          = "/api/cups/{cup_id}"
	CompetitionsEndpoint = "/api/cups/{cup_id}/competitions"

	CreateRangeGroup           = "/api/range-group"
	CreateQualificationRound   = "/api/qualification-round"
	CreateQualificationSection = "/api/qualification-section"
	CreateRange                = "/api/range"
	CreateShot                 = "/api/range/{range_id}/shots"

	Competitor                    = "/api/competitors/{competitor_id}"
	GetCompetitorsFromCompetition = "/api/competitors/{competition_id}/competitors"
	GetQualificationSections      = "/api/qualification_sections/{id}"

	RegisterCompetitor = "/api/competitors/registration"
	GetCompetitor      = "/api/competitors/{competitor_id}"
)
