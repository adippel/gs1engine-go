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
		if line == "" || strings.HasPrefix(line, "#") {
			continue // skip comments or blank lines
		}

		// Find the title by splitting on " #"
		var title string
		titleSplit := strings.SplitN(line, " #", 2)
		if len(titleSplit) == 2 {
			title = strings.TrimSpace(titleSplit[1])
			line = strings.TrimSpace(titleSplit[0])
		}

		fields := strings.Fields(line)
		entry := ApplicationIdentifierSpec{
			Title: title,
		}

		for i, field := range fields {
			if slices.Contains([]rune("0123456789"), rune(field[0])) {
				entry.AI = field
			} else if slices.Contains([]rune("*!?\"$%&'()+,-./:;<=>@[\\]^_`{|}~"), rune(field[0])) {
				entry.Flags = field
			} else if slices.Contains([]rune("NXY"), rune(field[0])) {
				entry.Specification = strings.Split(fields[i], ",")
			} else {
				entry.Attributes = append(entry.Attributes, field)
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
