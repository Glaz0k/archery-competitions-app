package enums

type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

func (s Gender) IsValid() bool {
	switch s {
	case Male, Female:
		return true
	}
	return false
}
