package format

import (
	"regexp"
	"strings"
)

func reaplaceAllRegex(s string, old string, new string) string {
	regex := regexp.MustCompile(old)
	return regex.ReplaceAllString(s, new)
}

func SqlFormat(sql string) (str string, linhas int) {
	tab := "  "
	parenthesisLevel := 0
	between := false
	deep := 0
	shift := createShiftArr("    ")

	sqlArray := split_sql(sql, tab)

	for i, line := range sqlArray {
		parenthesisLevel = isSubquery(line, parenthesisLevel)

		if ok, _ := regexp.MatchString(`\s{0,}\s{0,}SELECT\s{0,}`, line); ok {
			line = reaplaceAllRegex(line, `\,`, ",\n")
		}

		if ok, _ := regexp.MatchString(`\s{0,}\s{0,}SET\s{0,}`, line); ok {
			line = reaplaceAllRegex(line, `\,`, ",\n")
		}

		//Ajuste para que a identacao do between fique correta
		if ok, _ := regexp.MatchString(`AND`, line); ok {
			if ok, _ := regexp.MatchString(`BETWEEN`, sqlArray[i-1]); ok {
				line = reaplaceAllRegex(line, `\s{0,}AND`, "AND")
				between = true
			}
		}

		if ok, _ := regexp.MatchString(`\s{0,}\(\s{0,}SELECT\s{0,}`, line); ok {
			deep++
			str += shift[deep] + line
		} else if ok, _ := regexp.MatchString(`\'`, line); ok || between {
			between = false
			if parenthesisLevel < 1&deep {
				deep--
			}
			str += line
		} else {
			str += shift[deep] + line
			if parenthesisLevel < 1&deep {
				deep--
			}
		}
	}

	str = reaplaceAllRegex(str, `^\n{1,}`, "")
	str = reaplaceAllRegex(str, `\n{1,}`, "\n")

	linhas = strings.Count(str, "\n") + 1
	// fmt.Println(str)
	// fmt.Println(linhas)
	// fmt.Println("")
	return
}

func split_sql(sql, tab string) []string {
	sql = reaplaceAllRegex(sql, `\s{1,}`, " ")
	sql = reaplaceAllRegex(sql, `^SELECT`, "SELECT~::~ "+tab)
	sql = strings.ReplaceAll(sql, "AND", "~::~"+tab+"AND")
	sql = strings.ReplaceAll(sql, "CASE", "~::~"+tab+"CASE")
	sql = strings.ReplaceAll(sql, "ELSE", "~::~"+tab+"ELSE")
	sql = strings.ReplaceAll(sql, "END", "~::~"+tab+"END")
	sql = strings.ReplaceAll(sql, "FROM", "~::~FROM")
	sql = reaplaceAllRegex(sql, `GROUP\s{1,}BY`, "~::~GROUP BY")
	sql = strings.ReplaceAll(sql, "HAVING", "~::~HAVING")
	sql = strings.ReplaceAll(sql, " IN ", "IN")
	sql = strings.ReplaceAll(sql, "JOIN", "~::~JOIN")
	sql = strings.ReplaceAll(sql, "CROSS ~::~JOIN", "~::~CROSS JOIN")
	sql = strings.ReplaceAll(sql, "INNER ~::~JOIN", "~::~INNER JOIN")
	sql = strings.ReplaceAll(sql, "LEFT ~::~JOIN", "~::~LEFT JOIN")
	sql = strings.ReplaceAll(sql, "RIGHT ~::~JOIN", "~::~RIGHT JOIN")
	sql = strings.ReplaceAll(sql, " OR ", "~::~"+tab+tab+"OR ")
	sql = reaplaceAllRegex(sql, `ORDER\s{1,}BY`, "~::~ORDER BY ")
	sql = strings.ReplaceAll(sql, "OVER", "~::~"+tab+"OVER ")
	sql = reaplaceAllRegex(sql, `\(\s{0,}SELECT`, "~::~(SELECT")
	sql = reaplaceAllRegex(sql, `\)\s{0,}SELECT`, ")~::~SELECT")
	sql = strings.ReplaceAll(sql, "OVER", "~::~"+tab+"OVER ")
	sql = strings.ReplaceAll(sql, "THEN", " THEN~::~"+tab+"")
	sql = strings.ReplaceAll(sql, "UNION", "~::~UNION~::~")
	sql = strings.ReplaceAll(sql, "USING", "~::~USING")
	sql = strings.ReplaceAll(sql, "WHEN", "~::~"+tab+"WHEN")
	sql = strings.ReplaceAll(sql, "WHERE", "~::~WHERE")
	sql = strings.ReplaceAll(sql, "WITH", "~::~WITH")
	sql = reaplaceAllRegex(sql, `\,\s{0,}\(`, ",~::~( ")
	sql = reaplaceAllRegex(sql, `\,`, ",~::~ "+tab)
	sql = strings.ReplaceAll(sql, "ALL", "ALL")
	sql = strings.ReplaceAll(sql, " AS ", "AS")
	sql = strings.ReplaceAll(sql, "ASC", "ASC")
	sql = strings.ReplaceAll(sql, "DESC", "DESC")
	sql = strings.ReplaceAll(sql, "DISTINCT", "DISTINCT")
	sql = strings.ReplaceAll(sql, "EXISTS", "EXISTS")
	sql = strings.ReplaceAll(sql, "NOT", "NOT")
	sql = strings.ReplaceAll(sql, "NULL", "NULL")
	sql = strings.ReplaceAll(sql, "LIKE", "LIKE")
	sql = reaplaceAllRegex(sql, `ORDER\s{1,}BY`, "~::~ORDER BY")
	sql = reaplaceAllRegex(sql, `\s{0,}SELECT`, "SELECT")
	sql = reaplaceAllRegex(sql, `\s{0,}UPDATE`, "UPDATE")
	sql = reaplaceAllRegex(sql, `~::~{1,}`, "~::~")
	sql = reaplaceAllRegex(sql, `~::~\s{0,}~::~`, "~::~")
	sql = strings.ReplaceAll(sql, "=", " = ")
	sql = reaplaceAllRegex(sql, `\s{1,}=\s{1,}`, " = ")
	sql = strings.ReplaceAll(sql, "SET", "SET")

	return strings.Split(sql, "~::~")
}

func createShiftArr(space string) []string {
	shift := []string{"\n"} // array of shifts

	for i := 0; i < 100; i++ {
		shift = append(shift, shift[i]+space)
	}

	return shift
}

func isSubquery(str string, parenthesisLevel int) int {
	return parenthesisLevel - (len(reaplaceAllRegex(str, `\(`, "")) - len(reaplaceAllRegex(str, `\)`, "")))
}
