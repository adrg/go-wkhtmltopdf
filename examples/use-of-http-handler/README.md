## The What

This example consists of a **HTTP server** who accepts request to generate PDF from given URLs.

## How To

At one terminal:

    $ go run main.go

And on a second one:

```
$ curl -i -XPOST -o sample1.pdf http://localhost:7070 -d'
{
	"converterOpts": {
	    "marginLeft": "10mm",
	    "marginRight": "10mm",
	    "marginTop": "10mm",
	    "marginBottom": "10mm",
        "outlineDepth": 0
	},
	"objectOpts": {
        "location": "sample1.html",
		"windowStatus": "ready"
	}
}
'
```