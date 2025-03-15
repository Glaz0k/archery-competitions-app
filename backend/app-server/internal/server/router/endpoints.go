package router

const (
	IndividualGroupEndpoint = "/api/individual-groups/{individual_group_id}"

	EndCompetition        = "/api/competitions/{competition_id}/end"
	CompetitionEndpoint   = "/api/competitions/{competition_id}"
	CreateIndividualGroup = "/api/competitions/{competition_id}/individual-group"

	CreateCup         = "/api/cups"
	CupEndpoint       = "/api/cups/{cup_id}"
	CreateCompetition = "/api/cups/{cup_id}/competitions"

	CreateRangeGroup           = "/api/range-group"
	CreateQualification        = "/api/qualification"
	CreateQualificationRound   = "/api/qualification-round"
	CreateQualificationSection = "/api/qualification-section"
	CreateRange                = "/api/range"
	CreateShot                 = "/api/shot"
)
