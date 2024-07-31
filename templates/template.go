package templates

import "fmt"

func CreateTemplate(progress, problems, plans, insights string) string {
	return fmt.Sprintf(`
+---------------------+---------------------+
|        Progress     |        Problems     |
+---------------------+---------------------+
| %s                  | %s                  |
+---------------------+---------------------+
|          Plans      |       Insights      |
+---------------------+---------------------+
| %s                  | %s                  |
+---------------------+---------------------+
`, progress, problems, plans, insights)
}
