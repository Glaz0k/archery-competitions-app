package routes

const (
	IndividualGroupEndpoint            = "/api/individual-groups/{individual_group_id}"
	IndividualGroupCompetitorsEndpoint = "/api/individual-groups/{group_id}/competitors"
	StartQualification                 = "/api/individual-groups/{individual_group_id}/qualification/start"
	QualificationTable                 = "/api/individual_groups/{group_id}/qualification"
	UpdateIndividualGroup              = "/api/individual-groups/{group_id}/competitors/sync"

	EndCompetition              = "/api/competitions/{competition_id}/end"
	CompetitionEndpoint         = "/api/competitions/{competition_id}"
	IndividualGroupsCompetition = "/api/competitions/{competition_id}/individual_groups"

	CupsEndpoint         = "/api/cups"
	CupEndpoint          = "/api/cups/{cup_id}"
	CompetitionsEndpoint = "/api/cups/{cup_id}/competitions"

	Competitor               = "/api/competitors/{competitor_id}"
	GetQualificationSections = "/api/qualification_sections/{id}"

	RegisterCompetitor = "/api/competitors/registration"
	GetCompetitor      = "/api/competitors/{competitor_id}"

	CompetitorsCompetitionEndpoint = "/api/competitions/{competition_id}/competitors"
	CompetitorCompetitionEndpoint  = "/api/competitions/{competition_id}/competitors/{competitor_id}"

	SparringPlace = "/api/sparring_places/{id}"
	Ranges        = "/api/sparring_places/{id}/ranges"
	RangeEnd      = "api/sparring_places/{id}/ranges/{range_original}/end"
)
