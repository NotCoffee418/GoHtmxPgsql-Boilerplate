package form

import (
	"bytes"
	"net/http"

	"github.com/NotCoffee418/GoWebsite-Boilerplate/internal/access"
	"github.com/gin-gonic/gin"
)

type formResponse struct {
	IsSuccessful bool
	Message      string
}

type validationData struct {
	Messages         []formResponse
	DivId            string
	Visible          bool
	ContainerClasses string
	MessageClasses   string
}

type FormResponseBuilder struct {
	data validationData
}

// FormResponse is a struct that represents a validation message
// divId should be an empty div already on the page when using OOB swap
func NewFormResponseBuilder(divId string) *FormResponseBuilder {
	// Clean Id
	if len(divId) > 0 && divId[0] == '#' {
		divId = divId[1:]
	}

	// Create builder
	return &FormResponseBuilder{
		data: validationData{
			Messages: []formResponse{},
			DivId:    divId,
			Visible:  false, // Assgn on build
		},
	}
}

// AddMessage adds a message to the validation response
func (b *FormResponseBuilder) AddMessage(isSuccessful bool, message string) *FormResponseBuilder {
	b.data.Messages = append(b.data.Messages, formResponse{
		IsSuccessful: isSuccessful,
		Message:      message,
	})
	return b
}

// SetContainerClasses sets any classes for the outer div
func (b *FormResponseBuilder) SetContainerClasses(classes string) *FormResponseBuilder {
	b.data.ContainerClasses = classes
	return b
}

// SetMessageClasses sets the classes for div around the validation message
func (b *FormResponseBuilder) SetMessageClasses(classes string) *FormResponseBuilder {
	b.data.MessageClasses = classes
	return b
}

// BuildHtmlRenderer renders the validation response using the provided gin context
func (b *FormResponseBuilder) BuildHtmlRenderer(c *gin.Context) {
	b.finalizeBuild()
	c.HTML(http.StatusOK, "form_response.html", gin.H{
		"ValidationMessages": b.data,
	})
}

// BuildHtmlString returns the HTML string for the validation response
func (b *FormResponseBuilder) BuildHtmlString() (string, error) {
	b.finalizeBuild()
	var renderedHTML bytes.Buffer
	err := access.Tmpl.ExecuteTemplate(&renderedHTML, "form_response.html", b.data)
	if err != nil {
		return "", err
	}
	return renderedHTML.String(), nil
}

func (b *FormResponseBuilder) finalizeBuild() {
	b.data.Visible = len(b.data.Messages) > 0
}
