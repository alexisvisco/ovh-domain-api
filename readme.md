### Simple ovh domain api

Get the api:
```
go get github.com/alexisvisco/ovh-domain-api
```

Thenn just use it:

```go
import (
	"fmt"
	"github.com/alexisvisco/ovh-domain-api/domain"
	"github.com/alexisvisco/ovh-domain-api/domain/subsidiary"
)

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