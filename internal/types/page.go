package types

var DefaultPageTitle = "" // Set in main.go

// StructuredPageData provides a structured way to pass data to a page
// Meta: SEO information for the page
// Data: Any data required by the page
type StructuredPageData struct {
	Meta *MetaData
	Data map[string]interface{}
}

// MetaData contains SEO information for a page
// PageTitle: The title of the page (defaults to website name)
// PageDescription: The description of the page (optional)
// PageImage: The image to be used when sharing the page (optional)
type MetaData struct {
	Title       string
	Description string
	Image       string
}

func NewStructurePageData(data map[string]interface{}, meta *MetaData) StructuredPageData {
	// Ensure required values are set
	if meta == nil {
		meta = &MetaData{
			Title: DefaultPageTitle,
		}
	} else if meta.Title == "" {
		meta.Title = DefaultPageTitle
	}

	// Return structured data
	return StructuredPageData{
		Meta: meta,
		Data: data,
	}
}
