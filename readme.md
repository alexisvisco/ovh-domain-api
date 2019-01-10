### Simple ovh domain api

```go
func main() {
	client, e := domain.NewClient(subsidiary.FR)
	if e != nil {
		panic(e)
	}
	results, e := client.DomainInfo("google.fr")
	if e != nil {
		if e == domain.UnknownExtension {
			fmt.Println("UnknownExtension")
		} else {
			panic(e)
		}
	}
	fmt.Println(results.IsTaken())
}
```