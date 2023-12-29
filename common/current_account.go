package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type Requester interface {
	OAuth
	Account
}

type Account interface {
	AccountID() primitive.ObjectID
}

type OAuth interface {
	OAuthID() string
}

type currentAccount struct {
	OAuth
	Account
}

func CurrentAccount(t OAuth, u Account) *currentAccount {
	return &currentAccount{
		t, u,
	}
}
