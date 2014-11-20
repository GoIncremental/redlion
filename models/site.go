package models

import (
	"log"

	"github.com/goincremental/dal"
)

type Site struct {
	ID  dal.ObjectID `bson:"_id"`
	Env string
}

func GetSiteID(db dal.Database) (*dal.ObjectID, error) {
	var aSite Site
	sites := db.C("Sites")

	err := sites.Find(dal.BSON{"name": "redlion"}).One(&aSite)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &aSite.ID, nil
}
