# Go Webserver Template
## Features
- Go webserver with Gin
- HTMX with Go templates
- API handlers
- WebSocket system
- PostgreSQL with sqlx and async helper functions
- Tailwind CSS
- Database migration system

## Template Setup

To finalize the setup, follow these quick steps:

1. **Clone your repo:**

    ```bash
    git clone https://github.com/[YourUsername]/[YourRepoName].git
    ```

2. **Navigate to the repo:**

    ```bash
    cd [YourRepoName]
    ```

3. **Run the setup script:**

    ```bash
    python repo-template-setup.py
    ```

4. **Commit and push:**

    ```bash
    git add -A
    git commit -m "Finalize setup"
    git push
    ```

Replace `[YourUsername]` and `[YourRepoName]` with your GitHub username and repository name.

That's it! Your repo is now ready to use.

## Requirements
- [Go](https://golang.org/)
- [PostgreSQL](https://www.postgresql.org/) or [Docker](https://www.docker.com/) (for database)
- [Node.js](https://nodejs.org/en/) (for frontend dependencies)
- [Python](https://www.python.org/) (for template setup script)
- [VSCode](https://code.visualstudio.com/) (recommended)

### Recommended VSCode Extensions
- [Go](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [gotemplate-syntax](https://marketplace.visualstudio.com/items?itemName=casualjim.gotemplate)

## Usage

### Change the layout
The style of the application is defined in the `templates/layouts/default_base.html` file.  
You can change this as needed, but mind the imports and templating to maintain all functionality.  

### Creating a page

#####  Create a new html file in the `templates` directory
##### Add the following to any page template:

```html
{{ template "default_base.html" . }}
{{ define "content" }}
Your page content should go here
{{ end }}
```

##### Create a handler for the page in `./handlers` with this structure:
```go
type HomePageHandler struct{}

// Implements PageRouteRegistrar interface
func (h *HomePageHandler) Handler(engine *gin.Engine) {
    engine.GET("/", h.get)
}

func (h *HomePageHandler) get(c *gin.Context) {
    // Set SEO meta data
    meta := &page.PageMetaData{
        Title:       "Demo Home Page",
        Description: "This is a demo home page showing off the boilerplate.",
    }

    // Any additional data required by the page
    data := &map[string]interface{}{
        "Counter": CounterData{
            Value: 0,
            Color: "#fff",
        },
    }

    // Turn it into structured data
    structuredData := page.StructurePageData(&data, meta)

    // Render page
    c.HTML(http.StatusOK, "home_page.html", structuredData)
}
```

You can register multiple related handlers in one file. For more details see the  [Gin documentation](https://gin-gonic.com/docs/)

##### Register the handler in `./config/handlers.go`:
Add the handler to the `RouteHandlers` slice:
```go
&page_handlers.SomePageHandler{},
```

### Creating an HTMX component
The instructions are the same as for creating a page, but you want to add the route in the relevant page or system's handler file.  
Additionally, htmx component templates should not include ```{{ template "default_base.html" . }}```.

You can still use `page.StructurePageData(data, nil)` without defining meta for between page and component, if the page loads the component as well.

### API Handlers
API handlers are defined in `./handlers/api_handlers`.  
They are registered in `./config/handlers.go` alongside page handlers.  
API handlers should return JSON data. Included in the package is a structured API response outputting responses like so:
```json
{
    "success": true,
    "data": {
        "someData": "Some data"
    }
}
```
or for errors:

```json
{
    "success": false,
    "error": "Error message"
}
```

You can generate this data using `common.ApiResponseFactory.Ok(data)` and `common.ApiResponseFactory.Error(err)` respectively.

An API handler will generally look like this:
```go
type HomeApiHandler struct{}

type HomePageData struct {
    Time string `json:"time"`
}

// Implements PageRouteRegistrar interface
func (h *HomeApiHandler) Handler(engine *gin.Engine) {
    engine.GET("/api/home/get-server-time", h.get)
}

func (h *HomeApiHandler) get(c *gin.Context) {
    timeStr := time.Now().Format("2006-01-02 15:04:05")
    resp := common.ApiResponseFactory.Ok(
        &HomePageData{Time: timeStr})

    // Render page
    c.JSON(http.StatusOK, resp)
}
```

## Default Template Definitions
These ensure that all pages have the required blocks and that the layout is loaded correctly.  
You can add additional blocks to the layout as needed in `config/template_definitions.go`.  
The default definitions are:

- `content`: Each page expects a `content` block to be defined. HTMX components do not use this block.
- `scripts`: Optional, at the bottom of the page. Can be used to add additional scripts to the page. `<script>` tags are required.
- `head`: Optional, will be loaded in `<head>` of the page.
- `page_title`: Will be used as the page title. Defaults to website name.
