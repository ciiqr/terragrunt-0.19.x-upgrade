package main

import (
	"os"
	"io"
	"log"
	"path"
	"io/ioutil"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/parser"
	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/hcl/token"
)

func main() {
	// convert all supplied files
	terraformTfvars := os.Args[1:]
	for _, terraformTfvar := range terraformTfvars {
		// paths
		directory := path.Dir(terraformTfvar)
		terragruntHcl := path.Join(directory, "terragrunt.hcl")

		// terragrunt.hcl
		f, err := os.Create(terragruntHcl)
		if err != nil {
		    log.Fatal(err)
		    os.Exit(1)
		}

		// write to terragrunt.hcl
		ConvertToTf12(f, terraformTfvar)

		// remove terraform.tfvars file
		os.Remove(terraformTfvar)
	}
}

func ConvertToTf12(output io.Writer, path string) {
	// read file
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// read hcl
	fileNode := &ast.File{}
	fileNode, err = parser.Parse(bytes)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// get terragrunt node
	terragrunt := fileNode.Node.(*ast.ObjectList).Items[0]

	// inputs index
	inputsIndex := 0

	// TODO: we should use this to filter the list and pull the terragrunt block out no matter whether it's the first block, but I haven't bothered because all my terragrunt files start with the terragrunt block
	terragruntKey := terragrunt.Keys[0].Token.Text
	if terragruntKey == "terragrunt" {
		inputsIndex = 1

		// print terragrunt comments
		WriteComments(output, terragrunt)

		// print terragrunt values
		printer.Fprint(output, terragrunt.Val.(*ast.ObjectType).List)
		io.WriteString(output, "\n\n")
	}

	// print rest of the values
	values := &ast.ObjectList{Items: fileNode.Node.(*ast.ObjectList).Items[inputsIndex:]}
	if len(values.Items) > 0 {
		inputs := &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{Token: token.Token{Type: token.IDENT, Text: "inputs"}},
			},
			Val: &ast.ObjectType{
				List: values,
			},
			Assign: token.Pos{Line: 1, Column: 1},
		}

		// print values
		printer.Fprint(output, inputs)
		io.WriteString(output, "\n\n")
	}
}

func WriteComments(output io.Writer, obj *ast.ObjectItem) {
	if obj.LeadComment != nil {
		for _, comment := range obj.LeadComment.List {
			io.WriteString(output, comment.Text + "\n")
		}
	}
	// If key and val are on different lines, treat line comments like lead comments.
	if obj.LineComment != nil && obj.Val.Pos().Line != obj.Keys[0].Pos().Line {
		for _, comment := range obj.LineComment.List {
			io.WriteString(output, comment.Text + "\n")
		}
	}
}
