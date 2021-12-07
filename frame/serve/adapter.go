package serve

// WebFrameWork is an interface which is used as an adapter of
// framework and goAdmin. It must implement two methods. Use registers
// the routes and the corresponding handlers. Content writes the
// response to the corresponding context of framework.
type WebFrameWork interface {
	Name() string
	Use(app interface{}) error
	SetApp(app interface{}) error
}

// BaseAdapter is a base adapter contains some helper functions.
type BaseAdapter struct {
}

// GetUse is a helper function adds the plugins to the framework.
func (base *BaseAdapter) GetUse(app interface{}, wf WebFrameWork) error {
	if err := wf.SetApp(app); err != nil {
		return err
	}

	return nil
}
