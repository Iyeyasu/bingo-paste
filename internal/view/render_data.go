package view

// RenderData contains the data rendered to the view template.
type RenderData struct {
	Title      string
	ReadOnly   bool
	IconURL    string
	JavaScript string
	CSS        string
}
