package model

import (
	"errors"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var docXogReadXML, soapEnvelope *etree.Document

func LoadXMLReadList(path string) {
	docXogReadXML = etree.NewDocument()
	docXogReadXML.ReadFromFile(path)
	soapEnvelope = etree.NewDocument()
	soapEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='envelope']/soapenv:Envelope")
	soapEnvelope.SetRoot(soapEnvelopeElement.Copy())
}

type SectionLink struct {
	Code string `xml:"code,attr"`
}

type SectionField struct {
	Code         string `xml:"code,attr"`
	Column       string `xml:"column,attr"`
	Remove       bool   `xml:"remove,attr"`
	InsertBefore string `xml:"insertBefore,attr"`
}

type Section struct {
	Code           string         `xml:"code,attr"`
	SourcePosition string         `xml:"sourcePosition,attr"`
	TargetPosition string         `xml:"targetPosition,attr"`
	Action         string         `xml:"action,attr"`
	Links          []SectionLink  `xml:"link"`
	Fields         []SectionField `xml:"field"`
}

type Element struct {
	Type         string `xml:"type,attr"`
	XPath        string `xml:"xpath,attr"`
	Code         string `xml:"code,attr"`
	Action       string `xml:"action,attr"`
	InsertBefore string `xml:"insertBefore,attr"`
}

type FileReplace struct {
	From string `xml:"from"`
	To   string `xml:"to"`
}

type MatchExcel struct {
	Col           int    `xml:"col,attr"`
	XPath         string `xml:"xpath,attr"`
	AttributeName string `xml:"attribute,attr"`
	MultiValued   bool   `xml:"multiValued,attr"`
	Separator     string `xml:"separator,attr"`
}

type DriverFile struct {
	Code              string        `xml:"code,attr"`
	Path              string        `xml:"path,attr"`
	Type              string        `xml:"type,attr"`
	ObjCode           string        `xml:"objectCode,attr"`
	SingleView        bool          `xml:"singleView,attr"`
	CopyToView        string        `xml:"copyToView,attr"`
	IgnoreReading     bool          `xml:"ignoreReading,attr"`
	SourcePartition   string        `xml:"sourcePartition,attr"`
	TargetPartition   string        `xml:"targetPartition,attr"`
	PartitionModel    string        `xml:"partitionModel,attr"`
	InsertBefore      string        `xml:"insertBefore,attr"`
	InsertBeforeIndex string        `xml:"insertBeforeIndex,attr"`
	UpdateProgram    bool          `xml:"updateProgram,attr"`
	CopyPermissions  string        `xml:"copyPermissions,attr"`
	Template         string        `xml:"template,attr"`
	ExcelFile        string        `xml:"excel,attr"`
	ExcelStartRow    string        `xml:"startRow,attr"`
	InstanceTag      string        `xml:"instance,attr"`
	ExportToExcel    bool          `xml:"exportToExcel,attr"`
	OnlyStructure    bool          `xml:"onlyStructure,attr"`
	PackageTransform bool          `xml:"packageTransform,attr"`
	InstancesPerFile int           `xml:"instancesPerFile,attr"`
	NSQL             string        `xml:"nsql"`
	Sections         []Section     `xml:"section"`
	Elements         []Element     `xml:"element"`
	Replace          []FileReplace `xml:"replace"`
	MatchExcel       []MatchExcel  `xml:"match"`
	ExecutionOrder   int
	xogXML           string
	auxXML           string
}

func (d *DriverFile) InitXML(action, folder string) error {
	xml, err := parserReadXML(d)
	if action != constant.READ {
		xml, err = parserWriteXML(d, folder)
	}

	d.xogXML = xml
	if d.NeedAuxXML() {
		d.auxXML, err = parserReadXML(getAuxDriverFile(d))
		if err != nil {
			return err
		}
	}
	return err
}

func (d *DriverFile) SetXML(xml string) {
	d.xogXML = xml
}

func (d *DriverFile) GetXML() string {
	return d.xogXML
}

func (d *DriverFile) GetAuxXML() string {
	return d.auxXML
}

func (d *DriverFile) RunXML(action, sourceFolder string, environments *Environments, soapFunc util.Soap) error {
	d.Write(sourceFolder)
	if action == constant.READ {
		err := d.RunXogXML(environments.Source, soapFunc)
		if d.NeedAuxXML() {
			auxEnv := environments.Target
			if d.Type == constant.PROCESS {
				auxEnv = environments.Source
			}
			err = d.RunAuxXML(auxEnv, soapFunc)
		}
		return err
	} else {
		return d.RunXogXML(environments.Target, soapFunc)
	}
	return nil
}

func (d *DriverFile) RunAuxXML(env *EnvType, soapFunc util.Soap) error {
	result, err := executeSoapCall(d.auxXML, env, soapFunc)
	d.auxXML = result
	return err
}

func (d *DriverFile) RunXogXML(env *EnvType, soapFunc util.Soap) error {
	result, err := executeSoapCall(d.xogXML, env, soapFunc)
	d.xogXML = result
	return err
}

func (d *DriverFile) Write(folder string) {
	tag := "NikuDataBus"
	if folder == constant.FOLDER_DEBUG {
		tag = "XOGOutput"
	}
	r, _ := regexp.Compile("(?s)<" + tag + "(.*)</" + tag + ">")
	str := r.FindString(d.xogXML)
	if str == constant.UNDEFINED {
		str = d.xogXML
	}
	ioutil.WriteFile(folder+d.Type+"/"+d.Path, []byte(str), os.ModePerm)
}

func (d *DriverFile) NeedAuxXML() bool {
	return (d.Type == constant.OBJECT && len(d.Elements) > 0) || (d.Type == constant.VIEW && d.Code != "*") || (d.Type == constant.PROCESS && d.CopyPermissions != constant.UNDEFINED) || (d.Type == constant.MENU && len(d.Sections) > 0)
}

func (d *DriverFile) NeedPackageTransform() bool {
	return (d.Type == constant.OBJECT && len(d.Elements) > 0) || (d.Type == constant.VIEW && d.Code != "*") || (d.Type == constant.MENU && len(d.Sections) > 0)
}

func (d *DriverFile) TagCDATA() (string, string) {
	switch d.Type {
	case constant.PROCESS:
		return `<([^/].*):(query|update)(.*)"\s*>`, `</(.*):(query|update)>`
	case constant.LOOKUP:
		if !d.OnlyStructure {
			return `<nsql(.*)"\s*>`, `</nsql>`
		}
	}
	return constant.UNDEFINED, constant.UNDEFINED
}

func (d *DriverFile) GetSplitWriteFilesPath(folder string) ([]string, error) {
	files, err := ioutil.ReadDir(folder + d.Type)
	if err != nil {
		return nil, err
	}
	var splitPath []string
	matchFilename := util.GetPathWithoutExtension(d.Path)
	for _, filename := range files {
		if matchFilename == filename.Name()[:len(matchFilename)] {
			splitPath = append(splitPath, filename.Name())
		}
	}
	return splitPath, nil
}

func (d *DriverFile) GetDummyLookup() *etree.Element {
	return docXogReadXML.FindElement("//xogtype[@type='DummyLookup']/NikuDataBus").Copy()
}

func (d *DriverFile) GetInstanceTag() string {
	switch d.Type {
	case "CustomObjectInstances":
		return "instance"
	case "ResourceClassInstances":
		return "resourceClass"
	case "WipClassInstances":
		return "wipClass"
	case "InvestmentClassInstances":
		return "investmentClass"
	case "TransactionClassInstances":
		return "transactionClass"
	case "ResourceInstances":
		return "resource"
	case "UserInstances":
		return "user"
	case "ProjectInstances":
		return "project"
	case "IdeaInstances":
		return "idea"
	case "ApplicationInstances":
		return "application"
	case "AssetInstances":
		return "asset"
	case "OtherInvestmentInstances":
		return "otherInvestment"
	case "ProductInstances":
		return "product"
	case "ServiceInstances":
		return "service"
	}
	return constant.UNDEFINED
}

func (d *DriverFile) GetXMLType() string {
	switch d.Type {
	case "Files":
		return "file"
	case "Objects":
		return "object"
	case "Views":
		return "view"
	case "Processes":
		return "process"
	case "Lookups":
		return "lookup"
	case "Portlets":
		return "portlet"
	case "Queries":
		return "query"
	case "Pages":
		return "page"
	case "Menus":
		return "menu"
	case "Obs":
		return "obs"
	case "Groups":
		return "group"
	case "CustomObjectInstances":
		return "customObjectInstance"
	case "ResourceClassInstances":
		return "resourceClassInstance"
	case "WipClassInstances":
		return "wipClassInstance"
	case "InvestmentClassInstances":
		return "investmentClassInstance"
	case "TransactionClassInstances":
		return "transactionClassInstance"
	case "ResourceInstances":
		return "resourceInstance"
	case "UserInstances":
		return "userInstance"
	case "ProjectInstances":
		return "projectInstance"
	case "IdeaInstances":
		return "ideaInstance"
	case "ApplicationInstances":
		return "applicationInstance"
	case "AssetInstances":
		return "assetInstance"
	case "OtherInvestmentInstances":
		return "otherInvestmentInstance"
	case "ProductInstances":
		return "productInstance"
	case "ServiceInstances":
		return "serviceInstance"
	case "Migrations":
		return "migration"
	}
	return constant.UNDEFINED
}

func executeSoapCall(body string, env *EnvType, soapFunc util.Soap) (string, error) {
	bodyWithSession := strings.Replace(body, "<xog:SessionID/>", "<xog:SessionID>"+env.Session+"</xog:SessionID>", -1)
	return soapFunc(bodyWithSession, env.URL)
}

func getAuxDriverFile(d *DriverFile) *DriverFile {
	switch d.Type {
	case constant.PROCESS:
		return &DriverFile{Code: d.CopyPermissions, Path: "aux_" + d.CopyPermissions + ".xml", Type: d.Type}
	case constant.OBJECT:
		partition := d.SourcePartition
		if d.TargetPartition != constant.UNDEFINED {
			partition = d.TargetPartition
		}
		return &DriverFile{Code: d.Code, Path: "aux_" + d.Path + ".xml", Type: d.Type, SourcePartition: partition}
	case constant.VIEW:
		partition := d.SourcePartition
		if d.TargetPartition != constant.UNDEFINED {
			partition = d.TargetPartition
		}
		return &DriverFile{Code: d.Code, ObjCode: d.ObjCode, Path: "aux_" + d.Path + ".xml", SourcePartition: partition, Type: d.Type}
	case constant.MENU:
		return &DriverFile{Code: d.Code, Path: "aux_" + d.Path + ".xml", Type: d.Type}
	}
	return nil
}

func parserReadXML(d *DriverFile) (string, error) {
	if d.Code == constant.UNDEFINED {
		return constant.UNDEFINED, errors.New("no attribute code defined")
	}

	if d.Path == constant.UNDEFINED {
		return constant.UNDEFINED, errors.New("no attribute path defined")
	}

	nikuDataBusElement := docXogReadXML.FindElement("//xogtype[@type='" + d.Type + "']/NikuDataBus")
	if nikuDataBusElement == nil {
		return constant.UNDEFINED, errors.New("invalid object type")
	}
	envelope := soapEnvelope.Root().Copy()
	envelope.FindElement("//soapenv:Body").AddChild(nikuDataBusElement.Copy())

	req := etree.NewDocument()
	req.SetRoot(envelope)

	switch d.Type {
	case constant.LOOKUP:
		req.FindElement("//Filter[@name='code']").SetText(strings.ToUpper(d.Code))
	case constant.PORTLET, constant.QUERY, constant.PROCESS, constant.PAGE, constant.GROUP, constant.MENU, constant.OBS:
		req.FindElement("//Filter[@name='code']").SetText(d.Code)
	case constant.OBJECT:
		req.FindElement("//Filter[@name='object_code']").SetText(d.Code)
	case constant.VIEW:
		if d.ObjCode == constant.UNDEFINED {
			return constant.UNDEFINED, errors.New("no attribute objectCode defined on tag <view>")
		}
		req.FindElement("//Filter[@name='code']").SetText(d.Code)
		req.FindElement("//Filter[@name='object_code']").SetText(d.ObjCode)
		if d.SourcePartition == constant.UNDEFINED {
			filter := req.FindElement("//Filter[@name='partition_code']")
			filter.Parent().RemoveChild(filter)
		} else {
			req.FindElement("//Filter[@name='partition_code']").SetText(d.SourcePartition)
		}
	case constant.CUSTOM_OBJECT_INSTANCE:
		if d.ObjCode == constant.UNDEFINED {
			return constant.UNDEFINED, errors.New("no attribute objectCode defined on tag <customObjectInstance>")
		}
		req.FindElement("//Filter[@name='instanceCode']").SetText(d.Code)
		req.FindElement("//Filter[@name='objectCode']").SetText(d.ObjCode)
	case constant.RESOURCE_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='resource_class']").SetText(d.Code)
	case constant.WIP_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='wipclass']").SetText(d.Code)
	case constant.INVESTMENT_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='investmentclass']").SetText(d.Code)
	case constant.TRANSACTION_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='transclass']").SetText(d.Code)
	case constant.RESOURCE_INSTANCE:
		req.FindElement("//Filter[@name='resourceID']").SetText(d.Code)
	case constant.USER_INSTANCE:
		req.FindElement("//Filter[@name='userName']").SetText(d.Code)
	case constant.PROJECT_INSTANCE:
		req.FindElement("//Filter[@name='projectID']").SetText(d.Code)
	case constant.IDEA_INSTANCE, constant.APPLICATION_INSTANCE, constant.ASSET_INSTANCE, constant.OTHER_INVESTMENT_INSTANCE, constant.PRODUCT_INSTANCE, constant.SERVICE_INSTANCE:
		req.FindElement("//Filter[@name='objectID']").SetText(d.Code)
	}

	documentLocationElement := req.FindElement("//args[@name='documentLocation']")
	if documentLocationElement != nil {
		folder := "./" + constant.FOLDER_WRITE + "_" + d.Type + "/_document"
		documentLocationElement.CreateAttr("value", folder)
	}

	req.IndentTabs()
	str, err := req.WriteToString()
	return str, err
}

func parserWriteXML(d *DriverFile, folder string) (string, error) {
	nikuDataBusXML := etree.NewDocument()
	nikuDataBusXML.ReadFromFile(folder + d.Type + "/" + d.Path)

	req := etree.NewDocument()
	req.SetRoot(soapEnvelope.Root().Copy())

	req.FindElement("//soapenv:Body").AddChild(nikuDataBusXML.Root())
	req.IndentTabs()
	return req.WriteToString()
}

type Driver struct {
	Version       string
	Files         []DriverFile
	PackageDriver bool
	FilePath      string
	Info          os.FileInfo
}

func (d *Driver) Clear() {
	d.Version = constant.UNDEFINED
	d.Files = []DriverFile{}
	d.PackageDriver = false
	d.FilePath = constant.UNDEFINED
	d.Info = nil
}

func (d *Driver) MaxTypeNameLen() int {
	max := 0
	for _, f := range d.Files {
		strLen := len(f.GetXMLType())
		if strLen > max {
			max = strLen
		}
	}
	return max
}

type ByExecutionOrder []DriverFile

func (d ByExecutionOrder) Len() int           { return len(d) }
func (d ByExecutionOrder) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByExecutionOrder) Less(i, j int) bool { return d[i].ExecutionOrder < d[j].ExecutionOrder }

type DriverTypesPattern struct {
	Version                   string       `xml:"version,attr"`
	Files                     []DriverFile `xml:"file"`
	Objects                   []DriverFile `xml:"object"`
	Views                     []DriverFile `xml:"view"`
	Processes                 []DriverFile `xml:"process"`
	Lookups                   []DriverFile `xml:"lookup"`
	Portlets                  []DriverFile `xml:"portlet"`
	Queries                   []DriverFile `xml:"query"`
	Pages                     []DriverFile `xml:"page"`
	Menus                     []DriverFile `xml:"menu"`
	Obs                       []DriverFile `xml:"obs"`
	Groups                    []DriverFile `xml:"group"`
	CustomObjectInstances     []DriverFile `xml:"customObjectInstance"`
	ResourceClassInstances    []DriverFile `xml:"resourceClassInstance"`
	WipClassInstances         []DriverFile `xml:"wipClassInstance"`
	InvestmentClassInstances  []DriverFile `xml:"investmentClassInstance"`
	TransactionClassInstances []DriverFile `xml:"transactionClassInstance"`
	ResourceInstances         []DriverFile `xml:"resourceInstance"`
	UserInstances             []DriverFile `xml:"userInstance"`
	ProjectInstances          []DriverFile `xml:"projectInstance"`
	IdeaInstances             []DriverFile `xml:"ideaInstance"`
	ApplicationInstances      []DriverFile `xml:"applicationInstance"`
	AssetInstances            []DriverFile `xml:"assetInstance"`
	OtherInvestmentInstances  []DriverFile `xml:"otherInvestmentInstance"`
	ProductInstances          []DriverFile `xml:"productInstance"`
	ServiceInstances          []DriverFile `xml:"serviceInstance"`
	Migrations                []DriverFile `xml:"migration"`
}
