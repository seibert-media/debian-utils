package package_creator

type Creator interface {
	CreatePackage() error
}

type creator struct {

}

func New() *creator {
	return new (creator)
}

func (c *creator) CreatePackage() error {
	return nil
}