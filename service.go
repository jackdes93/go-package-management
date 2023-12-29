package go_package_manager

type service struct {
	name        string
	version     string
	env         string
	opts        []Option
	subServices []Runnable
}
