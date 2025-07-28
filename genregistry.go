package gs1

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// ApplicationIdentifierSpec is an intermediate representation of the parsed AI description.
type ApplicationIdentifierSpec struct {
	AI            string
	Flags         string
	Specification []string
	Attributes    []string
	Title         string
}

// DownloadSyntaxDictionary downloads the most recent GS1 Syntax Dictionary from
// https://github.com/gs1/gs1-syntax-dictionary/blob/main/gs1-syntax-dictionary.txt. Pass a tag to release to download
// the correct dictionary.
func DownloadSyntaxDictionary(release string) (bytes.Buffer, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/gs1/gs1-syntax-dictionary/refs/tags/%s/gs1-syntax-dictionary.txt", release)
	resp, err := http.Get(url)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error downloading syntax dictionary release %s: %w", release, err)
	}
	defer resp.Body.Close()

	data := bytes.Buffer{}

	_, err = io.Copy(&data, resp.Body)
	if err != nil {
		return bytes.Buffer{}, fmt.Errorf("error reading syntax dictionary release %s: %w", release, err)
	}

	return data, nil
}

// ParseSyntaxDictionary reads in the format of the GS1 Syntax Dictionary published here
// https://github.com/gs1/gs1-syntax-dictionary.
func ParseSyntaxDictionary(r io.Reader) ([]ApplicationIdentifierSpec, error) {
	var entries []ApplicationIdentifierSpec
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || len(line) == 0 {
			continue // skip comments or blank lines
		}

		columns := getColumns(line)
		entry := ApplicationIdentifierSpec{}

		for _, column := range columns {
			if slices.Contains([]rune("0123456789"), rune(column[0])) {
				entry.AI = column
			} else if slices.Contains([]rune("*!?\"$%&'()+,-./:;<=>@[\\]^_`{|}~"), rune(column[0])) {
				entry.Flags = strings.ReplaceAll(column, " ", "")
			} else if slices.Contains([]rune("NXYZ"), rune(column[0])) {
				entry.Specification = strings.Split(column, ",")
			} else if strings.HasPrefix(column, "# ") {
				entry.Title = column[2:]
			} else {
				entry.Attributes = strings.Split(column, " ")
			}
		}

		// That's a specification spanning several AIs
		if strings.Contains(entry.AI, "-") {
			aiRange := strings.Split(entry.AI, "-")
			startAI, err := strconv.Atoi(aiRange[0])
			if err != nil {
				return nil, fmt.Errorf("error parsing start AI range %s: %w", entry.AI, err)
			}
			stopAI, err := strconv.Atoi(aiRange[1])
			if err != nil {
				return nil, fmt.Errorf("error parsing stop AI range %s: %w", entry.AI, err)
			}

			for i := startAI; i <= stopAI; i++ {
				strconvAI := strconv.Itoa(i)
				entries = append(entries, ApplicationIdentifierSpec{
					AI:            strconvAI,
					Flags:         entry.Flags,
					Specification: entry.Specification,
					Attributes:    entry.Attributes,
					Title:         entry.Title,
				})
			}
		} else {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// getColumns is a custom parser that detects columns which are separated by two spaces.
func getColumns(line string) []string {
	inColumn := true
	currentColumnBuilder := strings.Builder{}
	columns := []string{}
	i := 0

	for {
		char := line[i]
		if char == ' ' && i+1 < len(line) && line[i+1] == ' ' && inColumn {
			// two white spaces found -> the current column is completed;
			columns = append(columns, currentColumnBuilder.String())
			currentColumnBuilder.Reset()
			inColumn = false
		}
		if !inColumn && char != ' ' {
			// after we were out of column, we found a first non-white space character -> new column begins
			inColumn = true
		}
		if inColumn {
			currentColumnBuilder.WriteByte(char)
		}

		i++
		if i == len(line) {
			columns = append(columns, currentColumnBuilder.String())
			break
		}
	}

	return columns
}
