package scene2d

type Group struct {
	Actor
}

func (self *Group) AddActor() {
	// self.Draw(batch, parentAlpha)
}

var hello []*Actor
var gp *Group = &Group{}

func init() {
	hello = []*Actor{
		&gp.Actor,
		&Actor{},
	}
}

type IAA interface {
	GetX() float32
}

type GG struct {
	x float32
}

func (self *GG) GetX() float32 {
	return self.x
}

type HH struct {
	GG
}

var zzz []IAA

// We achieve polymorhphism through interfaces
// As now I can have an array of IActors which can contains Actors, Groups, Widgets, UI etc..

func InitZZZ() {
	zzz = []IAA{
		&GG{},
		&HH{},
	}
}
