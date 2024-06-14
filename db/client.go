package db

type Client interface{}

type Mongo struct{}

func NewMongo(URI string) *Mongo {
	return &Mongo{}
}
