package model

import (
	"errors"
	"fmt"
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var docXogReadXML, soapEnvelope *etree.Document

//LoadXMLReadList loads the list of different types of xog read so it can be used during execution
func LoadXMLReadList(path string) {
	docXogReadXML = etree.NewDocument()
	docXogReadXML.ReadFromFile(path)
	soapEnvelope = etree.NewDocument()
	soapEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='envelope']/soapenv:Envelope")
	soapEnvelope.SetRoot(soapEnvelopeElement.Copy())
}

//SectionLink defines the fields for a link on a view section
type SectionLink struct {
	Code string `xml:"code,attr"`
}

//SectionField defines the fields for an attribute on a section attribute
type SectionField struct {
	Code         string `xml:"code,attr"`
	Column       string `xml:"column,attr"`
	Remove       bool   `xml:"remove,attr"`
	InsertBefore string `xml:"insertBefore,attr"`
}

//Section defines the fields to describe a section on a view
type Section struct {
	Code           string         `xml:"code,attr"`
	SourcePosition string         `xml:"sourcePosition,attr"`
	TargetPosition string         `xml:"targetPosition,attr"`
	Action         string         `xml:"action,attr"`
	Links          []SectionLink  `xml:"link"`
	Fields         []SectionField `xml:"field"`
}

//Element defines the fields to remove and insert elements from the xog xml file
type Element struct {
	Type         string `xml:"type,attr"`
	XPath        string `xml:"xpath,attr"`
	Code         string `xml:"code,attr"`
	Action       string `xml:"action,attr"`
	InsertBefore string `xml:"insertBefore,attr"`
}

//FileReplace defines the fields to replace strings on the xog xml file
type FileReplace struct {
	From string `xml:"from"`
	To   string `xml:"to"`
}

//MatchExcel defines the fields to map the cols on an excel file to attributes and data on the xog xml
type MatchExcel struct {
	Col           int    `xml:"col,attr"`
	XPath         string `xml:"xpath,attr"`
	AttributeName string `xml:"attribute,attr"`
	MultiValued   bool   `xml:"multiValued,attr"`
	Separator     string `xml:"separator,attr"`
}

//DriverFile defines the fields to manipulate the xog xml
type DriverFile struct {
	Code             string        `xml:"code,attr"`
	Path             string        `xml:"path,attr"`
	Type             string        `xml:"type,attr"`
	ObjCode          string        `xml:"objectCode,attr"`
	IgnoreReading    bool          `xml:"ignoreReading,attr"`
	SourcePartition  string        `xml:"sourcePartition,attr"`
	TargetPartition  string        `xml:"targetPartition,attr"`
	PartitionModel   string        `xml:"partitionModel,attr"`
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

//InitXML loads the properly xog xml to update the environment
func (d *DriverFile) InitXML(action, folder string) error {
	xml, err := parserReadXML(d)
	if action != constant.Read {
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

//SetXML fills the variable with the xog xml
func (d *DriverFile) SetXML(xml string) {
	d.xogXML = xml
}

//GetXML return the xog xml
func (d *DriverFile) GetXML() string {
	return d.xogXML
}

//GetAuxXML return the auxiliary xog xml
func (d *DriverFile) GetAuxXML() string {
	return d.auxXML
}

//RunXML executes a soap call to the properly xml (principal or auxiliary) depending on the action and the driver type
func (d *DriverFile) RunXML(action, sourceFolder string, environments *Environments, soapFunc util.Soap) error {
	d.Write(sourceFolder)
	if action == constant.Read {
		err := d.RunXogXML(environments.Source, soapFunc)
		if d.NeedAuxXML() {
			auxEnv := environments.Target
			if d.Type == constant.TypeProcess {
				auxEnv = environments.Source
			}
			err = d.RunAuxXML(auxEnv, soapFunc)
		}
		return err
	}
	return d.RunXogXML(environments.Target, soapFunc)
}

//RunAuxXML executes a soap call to the auxiliary xog xml
func (d *DriverFile) RunAuxXML(env *EnvType, soapFunc util.Soap) error {
	result, err := executeSoapCall(d.auxXML, env, soapFunc)
	d.auxXML = result
	return err
}

//RunXogXML executes a soap call to the principal xog xml
func (d *DriverFile) RunXogXML(env *EnvType, soapFunc util.Soap) error {
	result, err := executeSoapCall(d.xogXML, env, soapFunc)
	d.xogXML = result
	return err
}

//Write saves to the file system the content of the principal xog xml
func (d *DriverFile) Write(folder string) {
	tag := "NikuDataBus"
	if folder == constant.FolderDebug {
		tag = "XOGOutput"
	}
	r, _ := regexp.Compile("(?s)<" + tag + "(.*)</" + tag + ">")
	str := r.FindString(d.xogXML)
	if str == constant.Undefined {
		str = d.xogXML
	}
	ioutil.WriteFile(folder+d.Type+"/"+d.Path, []byte(str), os.ModePerm)
}

//NeedAuxXML validates if the driver needs to use an auxiliary xog xml
func (d *DriverFile) NeedAuxXML() bool {
	return (d.Type == constant.TypeObject && len(d.Elements) > 0) || (d.Type == constant.TypeView && d.Code != "*") || (d.Type == constant.TypeProcess && d.CopyPermissions != constant.Undefined) || (d.Type == constant.TypeMenu && len(d.Sections) > 0)
}

//NeedPackageTransform validates if a package driver needs to be transformed before install to an environment
func (d *DriverFile) NeedPackageTransform() bool {
	return (d.Type == constant.TypeObject && len(d.Elements) > 0) || (d.Type == constant.TypeView && d.Code != "*") || (d.Type == constant.TypeMenu && len(d.Sections) > 0)
}

//TagCDATA returns the correct tags where should be inserted the CDATA depending on the driver type
func (d *DriverFile) TagCDATA() (string, string) {
	switch d.Type {
	case constant.TypeProcess:
		return `<([^/].*):(query|update)(.*)"\s*>`, `</(.*):(query|update)>`
	case constant.TypeLookup:
		if !d.OnlyStructure {
			return `<nsql(.*)"\s*>`, `</nsql>`
		}
	}
	return constant.Undefined, constant.Undefined
}

//GetSplitWriteFilesPath returns a list with write files paths when an instance xog xml is split
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

//GetDummyLookup returns the xml that defines a simple lookup to avoid cross dependencies between objects and their attributes
func (d *DriverFile) GetDummyLookup() *etree.Element {
	return docXogReadXML.FindElement("//xogtype[@type='DummyLookup']/NikuDataBus").Copy()
}

//GetInstanceTag returns the instance tag according to the type of driver
func (d *DriverFile) GetInstanceTag() string {
	switch d.Type {
	case "CustomObjectInstances":
		return "instance"
	case "ResourceClassInstances":
		return "resourceclass"
	case "WipClassInstances":
		return "wipclass"
	case "InvestmentClassInstances":
		return "investmentClass"
	case "TransactionClassInstances":
		return "transactionclass"
	case "ResourceInstances":
		return "Resource"
	case "UserInstances":
		return "User"
	case "ProjectInstances":
		return "Project"
	case "IdeaInstances":
		return "Idea"
	case "ApplicationInstances":
		return "Application"
	case "AssetInstances":
		return "Asset"
	case "OtherInvestmentInstances":
		return "OtherInvestment"
	case "ProductInstances":
		return "Product"
	case "ServiceInstances":
		return "Service"
	}
	return constant.Undefined
}

//GetXMLType returns the constant value according to the type of driver
func (d *DriverFile) GetXMLType() string {
	switch d.Type {
	case "Files", "Objects", "Views", "Lookups", "Portlets", "Pages", "Menus", "Groups":
		return strings.ToLower(d.Type[:len(d.Type)-1])
	case "Processes":
		return "process"
	case "Queries":
		return "query"
	case "Obs":
		return "obs"
	case "CustomObjectInstances", "ResourceClassInstances", "WipClassInstances", "InvestmentClassInstances", "TransactionClassInstances", "ResourceInstances", "UserInstances", "ProjectInstances", "IdeaInstances", "ApplicationInstances", "AssetInstances", "OtherInvestmentInstances", "ProductInstances", "ServiceInstances", "Migrations":
		return strings.ToLower(d.Type[:1]) + d.Type[1:len(d.Type)-1]
	}
	return constant.Undefined
}

func executeSoapCall(body string, env *EnvType, soapFunc util.Soap) (string, error) {
	bodyWithSession := strings.Replace(body, "<xog:SessionID/>", "<xog:SessionID>"+env.Session+"</xog:SessionID>", -1)
	return soapFunc(bodyWithSession, env.URL)
}

func getAuxDriverFile(d *DriverFile) *DriverFile {
	switch d.Type {
	case constant.TypeProcess:
		return &DriverFile{Code: d.CopyPermissions, Path: "aux_" + d.CopyPermissions + ".xml", Type: d.Type}
	case constant.TypeObject:
		partition := d.SourcePartition
		if d.TargetPartition != constant.Undefined {
			partition = d.TargetPartition
		}
		return &DriverFile{Code: d.Code, Path: "aux_" + d.Path + ".xml", Type: d.Type, SourcePartition: partition}
	case constant.TypeView:
		partition := d.SourcePartition
		if d.TargetPartition != constant.Undefined {
			partition = d.TargetPartition
		}
		return &DriverFile{Code: d.Code, ObjCode: d.ObjCode, Path: "aux_" + d.Path + ".xml", SourcePartition: partition, Type: d.Type}
	case constant.TypeMenu:
		return &DriverFile{Code: d.Code, Path: "aux_" + d.Path + ".xml", Type: d.Type}
	}
	return nil
}

func parserReadXML(d *DriverFile) (string, error) {
	if d.Code == constant.Undefined {
		return constant.Undefined, errors.New("no attribute code defined")
	}

	if d.Path == constant.Undefined {
		return constant.Undefined, errors.New("no attribute path defined")
	}

	nikuDataBusElement := docXogReadXML.FindElement("//xogtype[@type='" + d.Type + "']/NikuDataBus")
	if nikuDataBusElement == nil {
		return constant.Undefined, errors.New("invalid object type")
	}
	envelope := soapEnvelope.Root().Copy()
	envelope.FindElement("//soapenv:Body").AddChild(nikuDataBusElement.Copy())

	req := etree.NewDocument()
	req.SetRoot(envelope)

	err := checkObjectCodeDefined(d)

	insertDefaultFiltersToReadXML(d, req)

	documentLocationElement := req.FindElement("//args[@name='documentLocation']")
	if documentLocationElement != nil {
		folder := "./" + constant.FolderWrite + "_" + d.Type + "/_document"
		documentLocationElement.CreateAttr("value", folder)
	}

	req.IndentTabs()
	str, err := req.WriteToString()
	return str, err
}

func checkObjectCodeDefined(d *DriverFile) error {
	if (d.Type == constant.TypeView || d.Type == constant.TypeCustomObjectInstance) && d.ObjCode == constant.Undefined {
		return fmt.Errorf("no attribute objectCode defined on tag <%s>", d.GetXMLType())
	}
	return nil
}

func insertDefaultFiltersToReadXML(d *DriverFile, req *etree.Document) {
	switch d.Type {
	case constant.TypeLookup:
		req.FindElement("//Filter[@name='code']").SetText(strings.ToUpper(d.Code))
	case constant.TypePortlet, constant.TypeQuery, constant.TypeProcess, constant.TypePage, constant.TypeGroup, constant.TypeMenu, constant.TypeObs:
		req.FindElement("//Filter[@name='code']").SetText(d.Code)
	case constant.TypeObject:
		req.FindElement("//Filter[@name='object_code']").SetText(d.Code)
	case constant.TypeView:
		req.FindElement("//Filter[@name='code']").SetText(d.Code)
		req.FindElement("//Filter[@name='object_code']").SetText(d.ObjCode)
		if d.SourcePartition == constant.Undefined {
			filter := req.FindElement("//Filter[@name='partition_code']")
			filter.Parent().RemoveChild(filter)
		} else {
			req.FindElement("//Filter[@name='partition_code']").SetText(d.SourcePartition)
		}
	case constant.TypeCustomObjectInstance:
		req.FindElement("//Filter[@name='instanceCode']").SetText(d.Code)
		req.FindElement("//Filter[@name='objectCode']").SetText(d.ObjCode)
	case constant.TypeResourceClassInstance:
		req.FindElement("//Filter[@name='resource_class']").SetText(d.Code)
	case constant.TypeWipClassInstance:
		req.FindElement("//Filter[@name='wipclass']").SetText(d.Code)
	case constant.TypeInvestmentClassInstance:
		req.FindElement("//Filter[@name='investmentclass']").SetText(d.Code)
	case constant.TypeTransactionClassInstance:
		req.FindElement("//Filter[@name='transclass']").SetText(d.Code)
	case constant.TypeResourceInstance:
		req.FindElement("//Filter[@name='resourceID']").SetText(d.Code)
	case constant.TypeUserInstance:
		req.FindElement("//Filter[@name='userName']").SetText(d.Code)
	case constant.TypeProjectInstance:
		req.FindElement("//Filter[@name='projectID']").SetText(d.Code)
	case constant.TypeIdeaInstance, constant.TypeApplicationInstance, constant.TypeAssetInstance, constant.TypeOtherInvestmentInstance, constant.TypeProductInstance, constant.TypeServiceInstance:
		req.FindElement("//Filter[@name='objectID']").SetText(d.Code)
	}
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

//Driver defines the file with a list of drivers to run
type Driver struct {
	Version       string
	Files         []DriverFile
	PackageDriver bool
	FilePath      string
	Info          os.FileInfo
}

//Clear reset the contents of the driver
func (d *Driver) Clear() {
	d.Version = constant.Undefined
	d.Files = []DriverFile{}
	d.PackageDriver = false
	d.FilePath = constant.Undefined
	d.Info = nil
}

//MaxTypeNameLen returns the largest size, number of characters in the type name, from the driver's list
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

//ByExecutionOrder used to order the drivers according to the users defined sequence
type ByExecutionOrder []DriverFile

func (d ByExecutionOrder) Len() int           { return len(d) }
func (d ByExecutionOrder) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ByExecutionOrder) Less(i, j int) bool { return d[i].ExecutionOrder < d[j].ExecutionOrder }

//DriverTypesPattern stores each driver type in an array to make it easier to read the xml file
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
