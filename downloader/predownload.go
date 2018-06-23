package downloader

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const shaPrefix = "SHA256: "

var selectors = map[string]string{
	"root":    "#server .platform.mb-5.linux",
	"version": " .version",
	"dlUrl":   " .clipboard",
	"dlAttr":  "data-clipboard-text",
	"sha":     " .checksum",
}

func (d *Downloader) gatherInformation() (map[string]string, error) {
	res, err := http.Get(d.url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	mainNode := doc.Find(selectors["root"]).First()
	return parseNodesInformation(mainNode)
}

func parseNodesInformation(mainNode *goquery.Selection) (map[string]string, error) {
	updateInfo := make(map[string]string, 3)

	version := mainNode.Find(selectors["version"]).First().Text()
	updateInfo["version"] = formatVersion(version)

	downloadURL, exists := mainNode.Find(selectors["dlUrl"]).Attr(selectors["dlAttr"])
	if !exists {
		return updateInfo, fmt.Errorf("could not access download url")
	}
	updateInfo["url"] = downloadURL

	sha := mainNode.Find(selectors["sha"]).First().Text()
	updateInfo["sha"] = formatSha(sha)

	return updateInfo, nil
}

func formatVersion(version string) string {
	return strings.TrimSpace(version)
}

func formatSha(sha string) string {
	return strings.TrimPrefix(sha, shaPrefix)
}

// res, err := http.Get("https://www.teamspeak.com/en/downloads#server")
// if err != nil {
//   // handle error
// }
// defer res.Body.Close()
//
// doc, err := goquery.NewDocumentFromReader(res.Body)
// if err != nil {
//   log.Fatal(err)
// }
// fmt.Println(doc)
//
// s := doc.Find("#server .platform.mb-5.linux").First()
// v := strings.TrimSpace(s.Find(" .version").First().Text())
// sha := strings.TrimPrefix(s.Find(" .checksum").First().Text(), "SHA256: ")
// dl, _ := s.Find(" .clipboard").Attr("data-clipboard-text")
// fmt.Println(dl, " VERSION ", v, sha)
// // x, exists := doc.Find("#server .platform.mb-5.linux .clipboard").First().Attr("data-clipboard-text")
// // fmt.Println(x, exists)