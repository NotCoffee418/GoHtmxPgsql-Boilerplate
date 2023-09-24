package page

var DefaultPageTitle = "" // Set in main.go

// StructuredPageData provides a structured way to pass data to a page
// Meta: SEO information for the page
// Data: Any data required by the page
type StructuredPageData[T any] struct {
	Meta *PageMetaData
	Data *T
}

// PageMetaData contains SEO information for a page
// PageTitle: The title of the page (defaults to website name)
// PageDescription: The description of the page (optional)
// PageImage: The image to be used when sharing the page (optional)
type PageMetaData struct {
	Title       string
	Description string
	Image       string
}

func StructurePageData[T any](data *T, meta *PageMetaData) StructuredPageData[T] {
	// Ensure required values are set
	if meta == nil {
		meta = &PageMetaData{
			Title: DefaultPageTitle,
		}
	} else if meta.Title == "" {
		meta.Title = DefaultPageTitle
	}

	// Return structured data
	return StructuredPageData[T]{
		Meta: meta,
		Data: data,
	}
}
