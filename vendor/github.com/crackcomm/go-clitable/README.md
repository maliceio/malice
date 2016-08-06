# go-clitable

### Deprecated in favor of [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter)

[![GoDoc](https://godoc.org/github.com/crackcomm/go-clitable?status.svg)](https://godoc.org/github.com/crackcomm/go-clitable)

ASCII and Markdown tables in console for golang.

## Usage

### Print table

```Go
table := New([]string{"Name", "Host", "..."})
table.AddRow(map[string]interface{}{"Name": "..."})
table.Print()
```

```
|----------------------------------------------------------------------------------------|
| Name              | Host                 | Type             | _id                      |
|----------------------------------------------------------------------------------------|
| MongoLab          | mongolab.com         | MongoDB Provider | 52518c5d56357d17ec000002 |
|----------------------------------------------------------------------------------------|
| Google App Engine | appengine.google.com | App Engine       | 52518ff356357d17ec000004 |
|----------------------------------------------------------------------------------------|
| Heroku            | heroku.com           | App Engine       | 5251918e56357d17ec000005 |
|----------------------------------------------------------------------------------------|
```

### Horizontal table

```Go
table.PrintHorizontal(map[string]interface{}{
	"Name": "MongoLab",
	"Host": "mongolab.com",
})
```

```
|---------------------------------|
| Name | MongoLab                 |
|---------------------------------|
| Host | mongolab.com             |
|---------------------------------|
| Type | MongoDB Provider         |
|---------------------------------|
| _id  | 52518c5d56357d17ec000002 |
|---------------------------------|
```

### Markdown table

```Go
table := New([]string{"Name", "Host", "..."})
table.AddRow(map[string]interface{}{"Name": "..."})
table.Markdown = true
table.Print()
```

```
| Name              | Host                 | Type             | _id                      |
| ----------------- | -------------------- | ---------------- | ------------------------ |
| MongoLab          | mongolab.com         | MongoDB Provider | 52518c5d56357d17ec000002 |
| Google App Engine | appengine.google.com | App Engine       | 52518ff356357d17ec000004 |
| Heroku            | heroku.com           | App Engine       | 5251918e56357d17ec000005 |
```

| Name              | Host                 | Type             | _id                      |
| ----------------- | -------------------- | ---------------- | ------------------------ |
| MongoLab          | mongolab.com         | MongoDB Provider | 52518c5d56357d17ec000002 |
| Google App Engine | appengine.google.com | App Engine       | 52518ff356357d17ec000004 |
| Heroku            | heroku.com           | App Engine       | 5251918e56357d17ec000005 |


## License

Unlicensed. For more information, please refer to http://unlicense.org.
