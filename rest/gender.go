package rest

import snsgo "github.com/rifqoi/sns-go"

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

func (g Gender) Convert() snsgo.Gender {
	switch g {
	case Male:
		return snsgo.Male
	case Female:
		return snsgo.Female
	}

	return snsgo.GenderUndefined
}
