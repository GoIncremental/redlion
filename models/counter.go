package models

import "github.com/goincremental/dal"

//Counter is a struct use to increment counts of things keyed by name
type Counter struct {
	ID       dal.ObjectID `bson:"_id" json:"id"`
	Name     string       `bson:"name" json:"name"`
	Number   int          `bson:"number" json:"number"`
	SystemID dal.ObjectID `bson:"systemId" json:"-,omitempty"`
}

// CounterModel defines methods available on the Counter struct
type CounterModel interface {
	GetNext(string) (int, error)
}

// NewCounter returns a CounterModel struct pre configured with a
// dal database connection and systemID ready to run the queries
func NewCounter(systemID *dal.ObjectID, db dal.Database) CounterModel {
	col := db.C("counters")
	return &counter{
		systemID: systemID,
		col:      col,
	}
}

type counter struct {
	systemID *dal.ObjectID
	col      dal.Collection
}

func (m *counter) GetNext(name string) (next int, err error) {
	next = -999
	var result Counter

	c := dal.Change{
		Update:    dal.BSON{"$inc": dal.BSON{"number": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	q := dal.BSON{
		"systemId": m.systemID,
		"name":     name,
	}
	_, err = m.col.Find(q).Apply(c, &result)
	if err == dal.ErrNotFound {
		next = 1
		err = nil
		return
	} else if err != nil {
		return
	}

	next = result.Number
	return
}
