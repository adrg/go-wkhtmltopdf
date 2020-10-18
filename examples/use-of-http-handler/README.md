## The What

This example consists of a **HTTP server** who accepts request to generate PDF from given URLs.

## How To

At one terminal:

    $ go run main.go

And on a second one:

```
$ curl -i -XPOST -o sample1.pdf http://localhost:7070 -d'
{
    "url": "sample1.html",
	"converterOptions": {
	    "margin.left": "10mm",
	    "margin.right": "10mm",
	    "margin.top": "10mm",
	    "margin.bottom": "10mm"
	},
	"objectOptions": {
		"load.windowStatus": "ready"
	}
}
'
```