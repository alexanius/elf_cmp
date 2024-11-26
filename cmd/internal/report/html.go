package report

import (
  "fmt"
)

func index(aName, bName, generalTable, secTable string) string {
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

</body>
</html> 
`,aName, bName, generalTable, secTable)
}

func singleSection(sectionName, fileName, syms string) string {
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
  <p>Section %s report</p>

  <h1>General</h1>

%s

</body>
</html> 
`, fileName, sectionName, syms)
}

func compareSections(sectionName, aFileName, bFileName string) string {
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
  <p>Section %s report</p>

</body>
</html> 
`, aFileName, bFileName, sectionName)
}
