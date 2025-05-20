package enums

type Stage string

const (
	StageI   Stage = "I"
	StageII  Stage = "II"
	StageIII Stage = "III"
	StageF   Stage = "F"
)

func (s Stage) IsValid() bool {
	switch s {
	case StageI, StageII, StageIII, StageF:
		return true
	}
	return false
}
