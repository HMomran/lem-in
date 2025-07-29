package Parse

type path struct {
	roomNames []string
	ISPath    bool
}

type Room struct {
	Name      string
	x         int
	y         int
	IsStart   bool
	IsEnd     bool
	Neighbors []*Room
}
