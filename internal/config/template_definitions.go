package config

type DefaultTemplateDefinition struct {
	Definition string
	Content    string
}

// DefaultTemplateDefinitions Go templating requires all called templates to be defined
// This is a list of default template items and it's content that will be used if not defined
var DefaultTemplateDefinitions = []DefaultTemplateDefinition{
	{"head", ""},                            // in <head>
	{"content", "This page has no content"}, // Page content
	{"scripts", ""},                         // Bottom of body
}
