package typesafe

import "encoding/json"

// Question is one of Choice, Score, or Noul.
type Question interface {
	questionType() string
}

// Choice picks one option from named criteria.
type Choice struct {
	Instructions string
	Criteria     map[string]string
}

func (Choice) questionType() string { return "choice" }

func (q Choice) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         string            `json:"type"`
		Instructions string            `json:"instructions"`
		Criteria     map[string]string `json:"criteria"`
	}{"choice", q.Instructions, q.Criteria})
}

// Score rates on a scale described by ordered criteria.
type Score struct {
	Instructions string
	Criteria     []string
}

func (Score) questionType() string { return "score" }

func (q Score) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         string   `json:"type"`
		Instructions string   `json:"instructions"`
		Criteria     []string `json:"criteria"`
	}{"score", q.Instructions, q.Criteria})
}

// Noul answers a true/false question with a probability.
type Noul struct {
	Instructions string
	Criteria     map[string]string
}

func (Noul) questionType() string { return "noul" }

func (q Noul) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type         string            `json:"type"`
		Instructions string            `json:"instructions"`
		Criteria     map[string]string `json:"criteria"`
	}{"noul", q.Instructions, q.Criteria})
}
