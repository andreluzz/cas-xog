package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"strings"
	"testing"
)

var packageMockFolder string

func init() {
	packageMockFolder = "../" + constant.FolderMock + "transform/"
}

func TestExecuteToReturnErrorNoHeaderElement(t *testing.T) {
	xog := etree.NewDocument()
	err := Execute(xog, nil, &model.DriverFile{})
	if err == nil {
		t.Fatalf("Error executing transformation. Not testing if xog has element head.")
	}
}

func TestExecuteToReturnPage(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypePage,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "page_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming page XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "page_result.xml") == false {
		t.Errorf("Error transforming page XOG file. Invalid result XML.")
	}
}

func TestExecuteToInsertElement(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypePage,
		Elements: []model.Element{
			{
				Action:    constant.ActionInsert,
				XPath:     "//tabbedPage",
				Attribute: "testAttribute",
				Value:     "test12345",
			},
			{
				Action:    constant.ActionInsert,
				XPath:     "//tab/OBSAssocs",
				XMLString: `<OBSAssoc id="test_obs" name="Test OBS" unitPath="/All/New/Node"/>`,
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "page_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming page XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "page_insert_element_result.xml") == false {
		t.Errorf("Error transforming page XOG file. Invalid result XML.")
	}
}

func TestExecuteToRemoveAllAttributesBut(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypePage,
		Elements: []model.Element{
			{
				Action:    constant.ActionRemoveAllButNot,
				XPath:     "//obs",
				Attribute: "name,code",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "obs_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming page XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "obs_remove_all_attributes_result.xml") == false {
		t.Errorf("Error transforming page XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnPageWithoutElementOBSandSecurity(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypePage,
		Elements: []model.Element{
			{
				Action: "remove",
				XPath:  "//OBSAssocs",
			},
			{
				Action: "remove",
				XPath:  "//Security",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "page_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming page XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "page_no_element_result.xml") == false {
		t.Errorf("Error transforming page XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnGroup(t *testing.T) {
	file := model.DriverFile{
		Code: "ObjectAdmin",
		Type: constant.TypeGroupInstance,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "group_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming group XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "group_result.xml") == false {
		t.Errorf("Error transforming group XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnGroupWithoutMembers(t *testing.T) {
	file := model.DriverFile{
		Code: "ObjectAdmin",
		Type: constant.TypeGroupInstance,
		Elements: []model.Element{
			{
				Action: "remove",
				XPath:  "//members",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "group_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming group XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "group_no_members_result.xml") == false {
		t.Errorf("Error transforming group XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnPortletFromQuery(t *testing.T) {
	file := model.DriverFile{
		Code: "apm.appByQuadrant",
		Type: constant.TypePortlet,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "portlet_query_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming query portlet XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "portlet_query_result.xml") == false {
		t.Errorf("Error transforming query portlet XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnPortletFromObject(t *testing.T) {
	file := model.DriverFile{
		Code: "test_cas_xog",
		Type: constant.TypePortlet,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "portlet_object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object portlet XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "portlet_object_result.xml") == false {
		t.Errorf("Error transforming object portlet XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnQuery(t *testing.T) {
	file := model.DriverFile{
		Code: "cop.processBottlenecks",
		Type: constant.TypeQuery,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "query_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming query XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "query_result.xml") == false {
		t.Errorf("Error transforming query XOG file. Invalid result XML.")
	}
}

func readMockResultAndCompare(xog *etree.Document, compareXML string) bool {
	xog.Indent(2)
	xogString, _ := xog.WriteToString()
	xogString = strings.Replace(xogString, " ", "", -1)

	xogProcessedToCompare := etree.NewDocument()
	xogProcessedToCompare.ReadFromFile(packageMockFolder + compareXML)
	xogProcessedToCompare.Indent(2)

	xogProcessedToCompareString, _ := xogProcessedToCompare.WriteToString()
	xogProcessedToCompareString = strings.Replace(xogProcessedToCompareString, " ", "", -1)
	if xogString != xogProcessedToCompareString {
		xog.WriteToFile("../" + constant.FolderDebug + "go_test_debug.xml")
		return false
	}
	return true
}

func TestExecuteToReturnOBS(t *testing.T) {
	file := model.DriverFile{
		Code: "strategic_plan",
		Type: constant.TypeOBSInstance,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "obs_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming OBS XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "obs_result.xml") == false {
		t.Errorf("Error transforming OBS XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnOBSWithoutSecurityAndObject(t *testing.T) {
	file := model.DriverFile{
		Code: "strategic_plan",
		Type: constant.TypeOBSInstance,
		Elements: []model.Element{
			{
				Action: "remove",
				XPath:  "//associatedObject",
			},
			{
				Action: "remove",
				XPath:  "//Security",
			},
			{
				Action: "remove",
				XPath:  "//rights",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "obs_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming OBS XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "obs_no_object_and_security_result.xml") == false {
		t.Errorf("Error transforming OBS XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnInstanceCorrectHeader(t *testing.T) {
	file := model.DriverFile{
		Type: constant.TypeResourceClassInstance,
	}
	xog := etree.NewDocument()
	xog.ReadFromString("<NikuDataBus><Header action=\"write\" externalSource=\"NIKU\" objectType=\"contentPack\" version=\"8.0\"/></NikuDataBus>")
	err := Execute(xog, nil, &file)
	if err != nil {
		t.Fatalf("Error transforming instance(RESOURCE_CLASS_INSTANCE) XOG file. Debug: %s", err.Error())
	}

	headerElement := xog.FindElement("//Header[@version='12.0']")
	if headerElement == nil {
		t.Errorf("Error transforming instance(RESOURCE_CLASS_INSTANCE) XOG file. Header wrong version number")
	}

	file = model.DriverFile{
		Type: constant.TypeWipClassInstance,
	}
	xog = etree.NewDocument()
	xog.ReadFromString("<NikuDataBus><Header action=\"write\" externalSource=\"NIKU\" objectType=\"contentPack\" version=\"8.0\"/></NikuDataBus>")
	err = Execute(xog, nil, &file)
	if err != nil {
		t.Fatalf("Error transforming instance(WIP_CLASS_INSTANCE) XOG file. Debug: %s", err.Error())
	}

	headerElement = xog.FindElement("//Header[@version='12.0']")
	if headerElement == nil {
		t.Errorf("Error transforming instance(WIP_CLASS_INSTANCE) XOG file. Header wrong version number")
	}

	file = model.DriverFile{
		Type: constant.TypeTransactionClassInstance,
	}
	xog = etree.NewDocument()
	xog.ReadFromString("<NikuDataBus><Header action=\"write\" externalSource=\"NIKU\" objectType=\"contentPack\" version=\"8.0\"/></NikuDataBus>")
	err = Execute(xog, nil, &file)
	if err != nil {
		t.Fatalf("Error transforming instance(TRANSACTION_CLASS_INSTANCE) XOG file. Debug: %s", err.Error())
	}

	headerElement = xog.FindElement("//Header[@version='12.0']")
	if headerElement == nil {
		t.Errorf("Error transforming instance(TRANSACTION_CLASS_INSTANCE) XOG file. Header wrong version number")
	}

	file = model.DriverFile{
		Type: constant.TypeInvestmentClassInstance,
	}
	xog = etree.NewDocument()
	xog.ReadFromString("<NikuDataBus><Header action=\"write\" externalSource=\"NIKU\" objectType=\"contentPack\" version=\"8.0\"/></NikuDataBus>")
	err = Execute(xog, nil, &file)
	if err != nil {
		t.Fatalf("Error transforming instance(INVESTMENT_CLASS_INSTANCE) XOG file. Debug: %s", err.Error())
	}

	headerElement = xog.FindElement("//Header[@version='14.1']")
	if headerElement == nil {
		t.Errorf("Error transforming instance(INVESTMENT_CLASS_INSTANCE) XOG file. Header wrong version number")
	}

	file = model.DriverFile{
		Type: constant.TypeThemeInstance,
	}
	xog = etree.NewDocument()
	xog.ReadFromString("<NikuDataBus><Header action=\"write\" externalSource=\"NIKU\" objectType=\"contentPack\" version=\"8.0\"/></NikuDataBus>")
	err = Execute(xog, nil, &file)
	if err != nil {
		t.Fatalf("Error transforming (THEME_UI) XOG file. Debug: %s", err.Error())
	}

	headerElement = xog.FindElement("//Header[@version='13.0']")
	if headerElement == nil {
		t.Errorf("Error transforming (THEME_UI) XOG file. Header wrong version number")
	}
}

func TestIncludeCDATAToReturnString(t *testing.T) {
	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog_cdata.xml")

	xogString, _ := xog.WriteToString()
	iniTagRegexp := `<([^/].*):(query|update)(.*)"\s*>`
	endTagRegexp := `</(.*):(query|update)>`

	XOGString := IncludeCDATA(xogString, iniTagRegexp, endTagRegexp)

	result := etree.NewDocument()
	result.ReadFromString(XOGString)

	if readMockResultAndCompare(result, "process_result_cdata.xml") == false {
		t.Errorf("Error including CDATA tag to process XOG file. Invalid result XML.")
	}
}

func TestIncludeCDATAWithoutQueryToReturnXML(t *testing.T) {
	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "process_full_xog.xml")

	xogString, _ := xog.WriteToString()
	iniTagRegexp := `<([^/].*):(query|update)(.*)"\s*>`
	endTagRegexp := `</(.*):(query|update)>`

	XOGString := IncludeCDATA(xogString, iniTagRegexp, endTagRegexp)

	result := etree.NewDocument()
	result.ReadFromString(XOGString)

	if readMockResultAndCompare(result, "process_full_xog.xml") == false {
		t.Errorf("Error including escapeText attribute to process XOG file. Invalid result XML.")
	}
}

func TestExecuteToRemoveAttributeFromElement(t *testing.T) {
	file := model.DriverFile{
		Code: "obj_sistema",
		Type: constant.TypeObject,
		Elements: []model.Element{
			{
				Action:    "remove",
				XPath:     "//customAttribute",
				Attribute: "column",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "object_full_xog.xml")
	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming object XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "object_element_remove_attribute_result.xml") == false {
		t.Errorf("Error transforming object XOG file. Invalid result XML.")
	}

}
