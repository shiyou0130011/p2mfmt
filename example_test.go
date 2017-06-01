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

	// Output:
	// {| class="wikitable"
	// |-
	// ! Name !! Age !! Sex
	// |-
	// |  John  ||  27  ||  M
	// |-
	// |  Anny  ||  16  ||  F
	// |-
	// |  Mike  ||  19  ||  M
	// |}

}

func ExampleConvert() {
	puki := `* Title
#navi(Example)

** Main List

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
	fmt.Printf("Mediawiki: \n%s\n\n", mediawiki)

	fmt.Printf("Categories: %v", category)

	// Output:
	// Mediawiki:
	// = Title=
	//
	//
	// * list 1
	// * list 2
	// * list 3
	// :* Sub List<ref> This is a reference </ref>
	//
	// == Table==
	//
	// {| class="wikitable" style="margin: 0 auto;"
	// |-
	// ! Name !! Age  !! Sex
	// |-
	// |  Anny  ||  16  ||  F
	// |-
	// |  John  ||  27  || rowspan = "2" |  M
	// |-
	// |  Mike  ||  19
	// |}
	//
	//
	// == 備註 ==
	// <references/>
	//
	// Categories: [Example]
}
