package transform

import (
	"errors"
	"regexp"

	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
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

func IncludeEscapeText(xog *etree.Document) (string, error) {
	xogQueryTagString, _ := xog.WriteToString()
	sqlQueryTagRegexp, _ := regexp.Compile(`(<[^/].*):(query|update)`)
	sqlTags := sqlQueryTagRegexp.FindAllString(xogQueryTagString, -1)

	if len(sqlTags) <= 0 {
		return xogQueryTagString, nil
	}

	for _, tag := range sqlTags {
		for _, e := range xog.FindElements("//" + tag[1:]) {
			e.CreateAttr("escapeText", "false")
		}
	}

	xogString, _ := xog.WriteToString()

	return xogString, nil
}
