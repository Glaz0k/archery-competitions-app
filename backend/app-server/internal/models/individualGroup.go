package models

type IndividualGroup struct {
	ID            int     `json:"id"`
	CompetitionID int     `json:"competition_id"`
	Bow           *string `json:"bow"`
	Identity      *string `json:"identity"`
	State         *string `json:"state"`
}

type FinalGrid struct {
	GroupID      int          `json:"group_id"`
	Quarterfinal Quarterfinal `json:"quarterfinal"`
	Semifinal    Semifinal    `json:"semifinal"`
	Final        Final        `json:"final"`
}

type Quarterfinal struct {
	Sparring1 Sparring `json:"sparring_1"`
	Sparring2 Sparring `json:"sparring_2"`
	Sparring3 Sparring `json:"sparring_3"`
	Sparring4 Sparring `json:"sparring_4"`
}

type Semifinal struct {
	Sparring5 Sparring `json:"sparring_5"`
	Sparring6 Sparring `json:"sparring_6"`
}

type Final struct {
	SparringGold   Sparring `json:"sparring_gold"`
	SparringBronze Sparring `json:"sparring_bronze"`
}
