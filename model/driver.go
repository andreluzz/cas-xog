package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/util"
	"github.com/beevik/etree"
)

var docXogReadXML, soapEnvelope *etree.Document

//LoadXMLReadList loads the list of different types of xog read so it can be used during execution
func LoadXMLReadList(path string) {
	docXogReadXML = etree.NewDocument()
	docXogReadXML.ReadFromFile(path)
	soapEnvelope = etree.NewDocument()
	soapEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='envelope']/soapenv:Envelope")
	soapEnvelope.SetRoot(soapEnvelopeElement.Copy())
	initInstancesTagByType()
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
	Attribute    string `xml:"attribute,attr"`
	Value        string `xml:"value,attr"`
	XMLString    string `xml:"xml"`
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
	RemoveIfNull  bool   `xml:"removeIfNull,attr"`
	Separator     string `xml:"separator,attr"`
}

//Filter defines the fields to filter the XOG read
type Filter struct {
	Criteria string `xml:"criteria,attr"`
	Name     string `xml:"name,attr"`
	Custom   bool   `xml:"customAttribute,attr"`
	Value    string `xml:",chardata"`
}

//HeaderArg defines the fields to include in the XOG read header
type HeaderArg struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

//DriverFile defines the fields to manipulate the xog xml
type DriverFile struct {
	Code             string        `xml:"code,attr"`
	Path             string        `xml:"path,attr"`
	Type             string        `xml:"type,attr"`
	ID               string        `xml:"id,attr"`
	TargetID         string        `xml:"targetId,attr"`
	ObjCode          string        `xml:"objectCode,attr"`
	ObjType          string        `xml:"objectType,attr"`
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
	Action           string        `xml:"action,attr"`
	NSQL             string        `xml:"nsql"`
	Sections         []Section     `xml:"section"`
	Elements         []Element     `xml:"element"`
	Replace          []FileReplace `xml:"replace"`
	MatchExcel       []MatchExcel  `xml:"match"`
	Filters          []Filter      `xml:"filter"`
	HeaderArgs       []HeaderArg   `xml:"args"`
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

//RestAPI validates if the driver uses the rest api
func (d *DriverFile) RestAPI() bool {
	return strings.HasPrefix(d.Type, "API")
}

//APIType returns the type without api prefix
func (d *DriverFile) APIType() string {
	return d.Type[3:]
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
	var splitPath []string
	if d.InstancesPerFile <= 0 {
		return splitPath, nil
	}

	files, err := ioutil.ReadDir(folder + d.Type)
	if err != nil {
		return nil, err
	}

	matchFilename := util.GetPathWithoutExtension(d.Path)
	for _, filename := range files {
		lengthMatchFilename := len(matchFilename)
		lengthFilename := len(filename.Name())
		if lengthFilename >= lengthMatchFilename && matchFilename == filename.Name()[:len(matchFilename)] {
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
	if d.Type == constant.TypeMigration {
		if d.InstanceTag == constant.Undefined {
			return constant.DefaultInstanceTag
		}
		return d.InstanceTag
	}
	if value, ok := instancesTag[d.Type]; ok {
		return value
	}
	return constant.Undefined
}

//GetXMLType returns the constant value according to the type of driver
func (d *DriverFile) GetXMLType() string {
	switch d.Type {
	case "APIBlueprints", "APITeams":
		return "api." + strings.ToLower(d.Type[3:len(d.Type)-1])
	case "Files", "Objects", "Views", "Lookups", "Portlets", "Pages", "Menus":
		return strings.ToLower(d.Type[:len(d.Type)-1])
	case "Processes":
		return "process"
	case "Queries":
		return "query"
	case "OBSInstances":
		return "obsInstance"
	case "CustomObjectInstances", "ResourceClassInstances", "WipClassInstances", "InvestmentClassInstances", "TransactionClassInstances",
		"ResourceInstances", "UserInstances", "ProjectInstances", "IdeaInstances", "ApplicationInstances", "AssetInstances", "OtherInvestmentInstances",
		"ProductInstances", "ServiceInstances", "BenefitPlanInstances", "BudgetPlanInstances", "CategoryInstances", "ChangeInstances",
		"ChargeCodeInstances", "CompanyClassInstances", "CostPlanInstances", "CostPlusCodeInstances", "DepartmentInstances", "EntityInstances",
		"GroupInstances", "IncidentInstances", "IssueInstances", "PortfolioInstances", "ProgramInstances", "ReleaseInstances",
		"ReleasePlanInstances", "RequirementInstances", "RequisitionInstances", "RiskInstances", "RoleInstances", "ThemeInstances", "VendorInstances", "DocumentInstances", "Migrations":
		return strings.ToLower(d.Type[:1]) + d.Type[1:len(d.Type)-1]
	}
	return constant.Undefined
}

func executeSoapCall(body string, env *EnvType, soapFunc util.Soap) (string, error) {
	bodyWithSession := strings.Replace(body, "<xog:SessionID/>", "<xog:SessionID>"+env.Session+"</xog:SessionID>", -1)
	return soapFunc(bodyWithSession, env.URL, env.Proxy)
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
	if d.Code == constant.Undefined && len(d.Filters) <= 0 {
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

	if len(d.HeaderArgs) > 0 {
		headerElement := req.FindElement("//args").Parent()
		for _, a := range req.FindElements("//args") {
			a.Parent().RemoveChild(a)
		}
		for _, a := range d.HeaderArgs {
			args := etree.NewElement("args")
			args.CreateAttr("name", a.Name)
			args.CreateAttr("value", a.Value)
			headerElement.AddChild(args)
		}
	}

	if len(d.Filters) > 0 {
		insertCustomFiltersToReadXML(d, req)
	} else {
		insertDefaultFiltersToReadXML(d, req)
	}

	documentLocationElement := req.FindElement("//args[@name='documentLocation']")
	if documentLocationElement != nil && len(d.HeaderArgs) <= 0 {
		ex, _ := os.Executable()
		exPath := filepath.Dir(ex)
		folder := exPath + "/" + constant.FolderWrite + d.Type + "/_document"
		util.ValidateFolder(folder)
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

func insertCustomFiltersToReadXML(d *DriverFile, req *etree.Document) {
	filterParentElement := req.FindElement("//Filter").Parent()
	for _, f := range req.FindElements("//Filter") {
		f.Parent().RemoveChild(f)
	}
	for _, f := range d.Filters {
		tag := "Filter"
		if f.Custom {
			tag = "FilterByCustomInfo"
		}

		filter := etree.NewElement(tag)
		filter.CreateAttr("criteria", f.Criteria)
		filter.CreateAttr("name", f.Name)
		filter.SetText(f.Value)
		filterParentElement.AddChild(filter)
	}
}

func insertDefaultFiltersToReadXML(d *DriverFile, req *etree.Document) {
	switch d.Type {
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
	case constant.TypeUserInstance:
		if d.Code != "*" {
			req.FindElement("//Filter[@name='userName']").SetText(d.Code)
		} else {
			for _, f := range req.FindElements("//Filter") {
				f.Parent().RemoveChild(f)
			}
		}
	case constant.TypeProjectInstance, constant.TypeProgramInstance:
		req.FindElement("//Filter[@name='projectID']").SetText(d.Code)
	case constant.TypeLookup:
		req.FindElement("//Filter").SetText(strings.ToUpper(d.Code))
	case constant.TypeDocumentInstance:
		req.FindElement("//Filter[@name='parentObjectID']").SetText(d.Code)
		req.FindElement("//Filter[@name='parentObjectType']").SetText(d.ObjType)
	default:
		req.FindElement("//Filter").SetText(d.Code)
	}
}

func parserWriteXML(d *DriverFile, folder string) (string, error) {
	nikuDataBusXML := etree.NewDocument()
	err := nikuDataBusXML.ReadFromFile(folder + d.Type + "/" + d.Path)
	if err != nil {
		return constant.Undefined, fmt.Errorf("File not found")
	}

	req := etree.NewDocument()
	req.SetRoot(soapEnvelope.Root().Copy())

	req.FindElement("//soapenv:Body").AddChild(nikuDataBusXML.Root())
	req.IndentTabs()
	return req.WriteToString()
}

var instancesTag map[string]string

func initInstancesTagByType() {
	instancesTag = make(map[string]string)
	instancesTag["CustomObjectInstances"] = "instance"
	instancesTag["ResourceClassInstances"] = "resourceclass"
	instancesTag["WipClassInstances"] = "wipclass"
	instancesTag["InvestmentClassInstances"] = "investmentClass"
	instancesTag["TransactionClassInstances"] = "transactionclass"
	instancesTag["ResourceInstances"] = "Resource"
	instancesTag["UserInstances"] = "User"
	instancesTag["ProjectInstances"] = "Project"
	instancesTag["IdeaInstances"] = "Idea"
	instancesTag["ApplicationInstances"] = "Application"
	instancesTag["AssetInstances"] = "Asset"
	instancesTag["OtherInvestmentInstances"] = "OtherInvestment"
	instancesTag["ProductInstances"] = "Product"
	instancesTag["ServiceInstances"] = "Service"
	instancesTag["BenefitPlanInstances"] = "BenefitPlan"
	instancesTag["BudgetPlanInstances"] = "BudgetPlan"
	instancesTag["CategoryInstances"] = "category"
	instancesTag["ChangeInstances"] = "changeRequest"
	instancesTag["ChargeCodeInstances"] = "chargeCode"
	instancesTag["CompanyClassInstances"] = "companyclass"
	instancesTag["CostPlanInstances"] = "CostPlan"
	instancesTag["CostPlusCodeInstances"] = "costPlusCode"
	instancesTag["DepartmentInstances"] = "Department"
	instancesTag["EntityInstances"] = "Entity"
	instancesTag["GroupInstances"] = "group"
	instancesTag["IncidentInstances"] = "incident"
	instancesTag["IssueInstances"] = "issue"
	instancesTag["PortfolioInstances"] = "pfmPortfolio"
	instancesTag["ProgramInstances"] = "Project"
	instancesTag["ReleaseInstances"] = "release"
	instancesTag["ReleasePlanInstances"] = "releaseplan"
	instancesTag["RequirementInstances"] = "requirement"
	instancesTag["RequisitionInstances"] = "requisition"
	instancesTag["RiskInstances"] = "risk"
	instancesTag["RoleInstances"] = "Role"
	instancesTag["ThemeInstances"] = "UITheme"
	instancesTag["VendorInstances"] = "vendor"
	instancesTag["DocumentInstances"] = "document"
}

//Driver defines the file with a list of drivers to run
type Driver struct {
	Version        string
	Files          []DriverFile
	PackageDriver  bool
	FilePath       string
	Info           os.FileInfo
	Folder         string
	AutomaticWrite bool
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
	AutomaticWrite            bool         `xml:"autoWrite,attr"`
	Files                     []DriverFile `xml:"file"`
	Objects                   []DriverFile `xml:"object"`
	Views                     []DriverFile `xml:"view"`
	Processes                 []DriverFile `xml:"process"`
	Lookups                   []DriverFile `xml:"lookup"`
	Portlets                  []DriverFile `xml:"portlet"`
	Queries                   []DriverFile `xml:"query"`
	Pages                     []DriverFile `xml:"page"`
	Menus                     []DriverFile `xml:"menu"`
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
	BenefitPlanInstances      []DriverFile `xml:"benefitPlanInstance"`
	BudgetPlanInstances       []DriverFile `xml:"budgetPlanInstance"`
	CategoryInstances         []DriverFile `xml:"categoryInstance"`
	ChangeInstances           []DriverFile `xml:"changeInstance"`
	ChargeCodeInstances       []DriverFile `xml:"chargeCodeInstance"`
	CompanyClassInstances     []DriverFile `xml:"companyClassInstance"`
	CostPlanInstances         []DriverFile `xml:"costPlanInstance"`
	CostPlusCodeInstances     []DriverFile `xml:"costPlusCodeInstance"`
	DepartmentInstances       []DriverFile `xml:"departmentInstance"`
	EntityInstances           []DriverFile `xml:"entityInstance"`
	GroupInstances            []DriverFile `xml:"groupInstance"`
	IncidentInstances         []DriverFile `xml:"incidentInstance"`
	IssueInstances            []DriverFile `xml:"issueInstance"`
	OBSInstances              []DriverFile `xml:"obsInstance"`
	PortfolioInstances        []DriverFile `xml:"portfolioInstance"`
	ProgramInstances          []DriverFile `xml:"programInstance"`
	ReleaseInstances          []DriverFile `xml:"releaseInstance"`
	ReleasePlanInstances      []DriverFile `xml:"releasePlanInstance"`
	RequirementInstances      []DriverFile `xml:"requirementInstance"`
	RequisitionInstances      []DriverFile `xml:"requisitionInstance"`
	RiskInstances             []DriverFile `xml:"riskInstance"`
	RoleInstances             []DriverFile `xml:"roleInstance"`
	ThemeInstances            []DriverFile `xml:"themeInstance"`
	VendorInstances           []DriverFile `xml:"vendorInstance"`
	DocumentInstances         []DriverFile `xml:"documentInstance"`
	APIBlueprints             []DriverFile `xml:"api.blueprint"`
	APITeams                  []DriverFile `xml:"api.team"`
}
