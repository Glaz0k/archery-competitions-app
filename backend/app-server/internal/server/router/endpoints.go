package router

// IndividualGroups

const (
	CreateIndividualGroup = "api/individual_groups"
)

// Competitions
const (
	EndCompetition    = "api/competitions/{competition_id}/end"
	CreateCompetition = "api/competition"
	EditCompetition   = "api/competitions/{competition_id}"
)

// Cups
const (
	CreateCup = "api/cups"
	Cup       = "api/cups/{cup_id}"
)

const (
	CreateRangeGroup           = "api/range_group"
	CreateQualification        = "api/qualification"
	CreateQualificationRound   = "api/qualification_round"
	CreateQualificationSection = "api/qualification_section"
	CreateRange                = "api/range"
	CreateShot                 = "api/shot"
)
