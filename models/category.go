package models

type Category struct {
	Value string
	Label string
}

func DefaultCategories() []Category {
	return []Category{
		{Value: "general", Label: "Général"},
		{Value: "tech", Label: "Tech"},
		{Value: "jeux", Label: "Jeux"},
		{Value: "business", Label: "Business"},
		{Value: "écologie", Label: "Écologie"},
		{Value: "santé", Label: "Santé"},
		{Value: "sport", Label: "Sport"},
		{Value: "culture", Label: "Culture"},
		{Value: "éducation", Label: "Éducation"},
	}
}
