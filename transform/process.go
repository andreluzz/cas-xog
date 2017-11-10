package transform

import (
	"errors"
	"regexp"
	"strings"
	"github.com/beevik/etree"
	"github.com/andreluzz/cas-xog/common"
)

func specificProcessTransformations(xog, aux *etree.Document, file common.DriverFile) error {
	removeElementFromParent(xog, "//lookups")

	if file.CopyPermissions != "" {
		securityElement, err := copyProcessPermissions(aux)
		if err != nil {
			return err
		}
		removeElementFromParent(xog, "//Security")
		process := xog.FindElement("//Process")
		if process == nil {
			return errors.New("process element not found")
		}
		process.AddChild(securityElement)
	}

	return nil
}

func copyProcessPermissions(xog *etree.Document) (*etree.Element, error) {
	element := xog.FindElement("//Security")

	if element == nil {
		return nil, errors.New("auxiliary xog to copy security from has no security element")
	}

	return element.Copy(), nil
}

func IncludeCDATA(xog *etree.Document) ([]byte, error) {
	xogQueryTagString, _ := xog.WriteToString()
	sqlQueryTagRegexp, _ := regexp.Compile(`(<[^/].*):(query|update)`)
	sqlTags := sqlQueryTagRegexp.FindAllString(xogQueryTagString, -1)

	if len(sqlTags) <= 0 {
		return []byte(xogQueryTagString), nil
	}

	for _, tag := range sqlTags {
		for _, e := range xog.FindElements("//" + tag[1:]) {
			e.CreateAttr("escapeText", "false")
		}
	}

	xogString, _ := xog.WriteToString()

	iniTagRegexp, _ := regexp.Compile(`<([^/].*):(query|update)(.*)>`)
	endTagRegexp, _ := regexp.Compile(`</(.*):(query|update)>`)

	iniIndex := iniTagRegexp.FindAllStringIndex(xogString, -1)
	endIndex := endTagRegexp.FindAllStringIndex(xogString, -1)

	shiftIndex := 0

	for i := 0; i < len(iniIndex); i++ {
		index := iniIndex[i][1] + shiftIndex
		xogString = xogString[:index] + "<![CDATA[" + xogString[index:]

		sqlString := xogString[index:endIndex[i][1]]

		paramRegexp, _ := regexp.Compile(`<(.*):param(.*)/>`)
		paramIndex := paramRegexp.FindStringIndex(sqlString)

		shiftIndex += 9

		eIndex := endIndex[i][0] + shiftIndex
		if len(paramIndex) > 0 {
			eIndex = endIndex[i][0] + 12 - (len(sqlString) - paramIndex[0])
		}

		xogString = xogString[:eIndex] + "]]>" + xogString[eIndex:]

		shiftIndex += 3
	}

	replacer := strings.NewReplacer("&gt;", ">", "&lt;", "<", "&apos;", "'", "&quot;", "\"",)
	xogString = replacer.Replace(xogString)

	return []byte(xogString), nil
}