package p2mfmt_test

import (
	"fmt"
	"github.com/shiyou0130011/p2mfmt"
)

func ExampleTable() {
	table := p2mfmt.Table{}
	table.ParseRow(`|~Name|~Age|~Sex|`)
	table.ParseRow(`| John | 27 | M |`)
	table.ParseRow(`| Anny | 16 | F |`)
	table.ParseRow(`| Mike | 19 | M |`)
	
	fmt.Print(table)
}

func ExampleConvert() {
	puki := `* Title
- list 1
- list 2
- list 3
-- Sub List(( This is a reference ))

** Table
|~Name|~Age|~Sex|
| Anny | 16 | F |
| John | 27 | M |
| Mike | 19 |~|`
	
	mediawiki, category := p2mfmt.Convert(puki)
	fmt.Printf("Mediawiki: \n%s\n", mediawiki)
	
	fmt.Printf("Categories: %v", category)
}