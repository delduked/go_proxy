package templates

import "html/template"

var Tmpl = template.Must(template.New("errorPage").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Domain Not Found</title>
</head>
<body>
    <h1>Error: Domain Not Found</h1>
    <p>The domain '{{.Domain}}' is not configured for redirection.</p>
</body>
</html>
`)) 
