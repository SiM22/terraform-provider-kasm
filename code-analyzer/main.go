package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type CodeStructure struct {
	Package    string
	Imports    []string
	Functions  []Function
	Types      []Type
	Variables  []Variable
	ApiMethods []ApiMethod
	Relations  []Relation
}

type Function struct {
	Name      string
	Signature string
	Doc       string
	Calls     []string // Functions that this function calls
}

type Type struct {
	Name    string
	Fields  []string
	Methods []Function
}

type Variable struct {
	Name  string
	Type  string
	Value string
}

type ApiMethod struct {
	Name         string
	Endpoint     string
	HttpMethod   string
	RequestType  string
	ResponseType string
	Doc          string
}

type Relation struct {
	Source      string
	Target      string
	Type        string // "calls", "implements", "uses", etc.
	Description string
}

// Helper function to extract a function's signature as a string
func getFunctionSignature(fn *ast.FuncDecl, fset *token.FileSet) string {
	if fn.Type == nil {
		return ""
	}

	var params []string
	if fn.Type.Params != nil {
		for _, p := range fn.Type.Params.List {
			var paramType string
			// Get the source code for the parameter type
			// Just use fmt.Sprintf to get a string representation of the type
			paramType = fmt.Sprintf("%s", p.Type)

			// Add parameter names
			var paramNames []string
			for _, name := range p.Names {
				paramNames = append(paramNames, name.Name)
			}

			if len(paramNames) > 0 {
				params = append(params, fmt.Sprintf("%s %s", strings.Join(paramNames, ", "), paramType))
			} else {
				params = append(params, paramType)
			}
		}
	}

	var results []string
	if fn.Type.Results != nil {
		for _, r := range fn.Type.Results.List {
			var resultType string
			// Get the source code for the result type
			// Just use fmt.Sprintf to get a string representation of the type
			resultType = fmt.Sprintf("%s", r.Type)

			// Add result names if any
			var resultNames []string
			for _, name := range r.Names {
				resultNames = append(resultNames, name.Name)
			}

			if len(resultNames) > 0 {
				results = append(results, fmt.Sprintf("%s %s", strings.Join(resultNames, ", "), resultType))
			} else {
				results = append(results, resultType)
			}
		}
	}

	signature := fmt.Sprintf("func(%s)", strings.Join(params, ", "))
	if len(results) > 0 {
		if len(results) == 1 {
			signature += " " + results[0]
		} else {
			signature += fmt.Sprintf(" (%s)", strings.Join(results, ", "))
		}
	}

	return signature
}

// Helper function to extract struct fields
func extractStructFields(structType *ast.StructType) []string {
	var fields []string
	if structType.Fields != nil {
		for _, field := range structType.Fields.List {
			var fieldType string
			fieldType = fmt.Sprintf("%s", field.Type)

			for _, name := range field.Names {
				fields = append(fields, fmt.Sprintf("%s %s", name.Name, fieldType))
			}
		}
	}
	return fields
}

// Helper function to extract variables
func extractVariables(decl *ast.GenDecl) []Variable {
	var variables []Variable
	if decl.Tok != token.VAR && decl.Tok != token.CONST {
		return variables
	}

	for _, spec := range decl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			var typeStr string
			if valueSpec.Type != nil {
				typeStr = fmt.Sprintf("%s", valueSpec.Type)
			}

			for i, name := range valueSpec.Names {
				var valueStr string
				if i < len(valueSpec.Values) {
					valueStr = fmt.Sprintf("%s", valueSpec.Values[i])
				}

				variables = append(variables, Variable{
					Name:  name.Name,
					Type:  typeStr,
					Value: valueStr,
				})
			}
		}
	}
	return variables
}

func main() {
	fset := token.NewFileSet()

	// Find all Go files in the current directory and subdirectories
	var goFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			goFiles = append(goFiles, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	if len(goFiles) == 0 {
		log.Println("No Go files found in the directory")
		json.NewEncoder(os.Stdout).Encode([]CodeStructure{})
		return
	}

	// Map to store package information
	packageMap := make(map[string]*CodeStructure)

	// Process each Go file
	for _, file := range goFiles {
		// Parse the Go file
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Error parsing file %s: %v\n", file, err)
			continue
		}

		// Get or create the package structure
		pkgName := f.Name.Name
		cs, exists := packageMap[pkgName]
		if !exists {
			cs = &CodeStructure{
				Package: pkgName,
				// Initialize all slices to empty slices instead of nil
				Imports:    []string{},
				Functions:  []Function{},
				Types:      []Type{},
				Variables:  []Variable{},
				ApiMethods: []ApiMethod{},
				Relations:  []Relation{},
			}
			packageMap[pkgName] = cs
		}

		// Extract imports
		for _, imp := range f.Imports {
			cs.Imports = append(cs.Imports, imp.Path.Value)
		}

		// Visit all nodes in the AST
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				// Extract function information
				f := Function{
					Name:      x.Name.Name,
					Signature: getFunctionSignature(x, fset),
					Calls:     []string{},
				}

				if x.Doc != nil {
					f.Doc = x.Doc.Text()
				}

				// Extract API method information if it looks like a client method
				if strings.Contains(f.Doc, "API") || strings.Contains(f.Doc, "endpoint") ||
					strings.Contains(f.Doc, "HTTP") || strings.Contains(f.Doc, "request") {
					// Try to extract endpoint and method from doc
					endpoint := ""
					httpMethod := ""
					requestType := ""
					responseType := ""

					// Look for endpoint pattern in doc
					if strings.Contains(f.Doc, "endpoint") || strings.Contains(f.Doc, "Endpoint") {
						endpointStart := strings.Index(f.Doc, "/")
						if endpointStart >= 0 {
							endpointEnd := strings.IndexAny(f.Doc[endpointStart:], " \n\t")
							if endpointEnd > 0 {
								endpoint = f.Doc[endpointStart : endpointStart+endpointEnd]
							} else {
								endpoint = f.Doc[endpointStart:]
							}
						}
					}

					// Look for HTTP method in doc
					for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH"} {
						if strings.Contains(f.Doc, method) {
							httpMethod = method
							break
						}
					}

					// Look for request/response types in function signature
					if strings.Contains(f.Signature, "Request") {
						parts := strings.Split(f.Signature, "Request")
						if len(parts) > 1 {
							// Extract the request type from before "Request"
							requestType = "Request"
							for i := len(parts[0]) - 1; i >= 0; i-- {
								if parts[0][i] == ' ' || parts[0][i] == '*' || parts[0][i] == '(' {
									requestType = parts[0][i+1:] + requestType
									break
								}
							}
						}
					}

					if strings.Contains(f.Signature, "Response") {
						parts := strings.Split(f.Signature, "Response")
						if len(parts) > 1 {
							// Extract the response type from before "Response"
							responseType = "Response"
							for i := len(parts[0]) - 1; i >= 0; i-- {
								if parts[0][i] == ' ' || parts[0][i] == '*' || parts[0][i] == '(' {
									responseType = parts[0][i+1:] + responseType
									break
								}
							}
						}
					}

					// If we found API-related information, add it as an API method
					if endpoint != "" || httpMethod != "" || requestType != "" || responseType != "" {
						apiMethod := ApiMethod{
							Name:         x.Name.Name,
							Endpoint:     endpoint,
							HttpMethod:   httpMethod,
							RequestType:  requestType,
							ResponseType: responseType,
							Doc:          f.Doc,
						}
						cs.ApiMethods = append(cs.ApiMethods, apiMethod)
					}
				}

				// Extract function calls by examining the function body
				if x.Body != nil {
					ast.Inspect(x.Body, func(n ast.Node) bool {
						if call, ok := n.(*ast.CallExpr); ok {
							if ident, ok := call.Fun.(*ast.Ident); ok {
								// Direct function call
								f.Calls = append(f.Calls, ident.Name)

								// Add a relation
								relation := Relation{
									Source:      x.Name.Name,
									Target:      ident.Name,
									Type:        "calls",
									Description: fmt.Sprintf("%s calls %s", x.Name.Name, ident.Name),
								}
								cs.Relations = append(cs.Relations, relation)
							} else if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
								// Method call or package function call
								if ident, ok := sel.X.(*ast.Ident); ok {
									callName := fmt.Sprintf("%s.%s", ident.Name, sel.Sel.Name)
									f.Calls = append(f.Calls, callName)

									// Add a relation
									relation := Relation{
										Source:      x.Name.Name,
										Target:      callName,
										Type:        "calls",
										Description: fmt.Sprintf("%s calls %s", x.Name.Name, callName),
									}
									cs.Relations = append(cs.Relations, relation)
								}
							}
						}
						return true
					})
				}

				// Check if this is a method of a type
				if x.Recv != nil && len(x.Recv.List) > 0 {
					// This is a method, not a standalone function
					recv := x.Recv.List[0]
					var typeName string

					// Handle pointer receivers
					switch rt := recv.Type.(type) {
					case *ast.StarExpr:
						if ident, ok := rt.X.(*ast.Ident); ok {
							typeName = ident.Name
						}
					case *ast.Ident:
						typeName = rt.Name
					}

					if typeName != "" {
						// Find the type and add this method to it
						for i, t := range cs.Types {
							if t.Name == typeName {
								cs.Types[i].Methods = append(cs.Types[i].Methods, f)
								return true
							}
						}
					}
				} else {
					// This is a standalone function
					cs.Functions = append(cs.Functions, f)
				}

			case *ast.GenDecl:
				if x.Tok == token.TYPE {
					// Extract type information
					for _, spec := range x.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							t := Type{
								Name:    typeSpec.Name.Name,
								Fields:  []string{},
								Methods: []Function{},
							}

							// Extract struct fields if this is a struct type
							if structType, ok := typeSpec.Type.(*ast.StructType); ok {
								t.Fields = extractStructFields(structType)
							}

							cs.Types = append(cs.Types, t)
						}
					}
				} else if x.Tok == token.VAR || x.Tok == token.CONST {
					// Extract variable information
					vars := extractVariables(x)
					cs.Variables = append(cs.Variables, vars...)
				}
			}
			return true
		})
	}

	// Convert the map to a slice for JSON output
	var structures []CodeStructure
	for _, cs := range packageMap {
		structures = append(structures, *cs)
	}

	// If no structures were found, return an empty array instead of null
	if len(structures) == 0 {
		structures = []CodeStructure{}
	}

	json.NewEncoder(os.Stdout).Encode(structures)
}
