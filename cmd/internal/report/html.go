package report

import (
  "fmt"
)

func index(aName, bName, generalTable, secTable, symTable string) string {
  return fmt.Sprintf(`
 <!DOCTYPE html>
<html>
<head>
  <style>
    table, th, td {
      border: 1px solid;
    }
  </style>
</head>
<body>

  <p>A file: %s</p>
  <p>B file: %s</p>

  <h1>General</h1>

 %s

  <h1>Sections</h1>

 %s

  <h1>Symbols</h1>

 %s

</body>
</html> 
`,aName, bName, generalTable, secTable, symTable)
}
