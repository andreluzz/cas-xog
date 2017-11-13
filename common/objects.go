package common

import "os"

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
	Type   string `xml:"type,attr"`
	XPath  string `xml:"xpath,attr"`
	Code   string `xml:"code,attr"`
	Action string `xml:"action,attr"`
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
	UpdateProgram     bool          `xml:"updateProgram,attr"`
	CopyPermissions   string        `xml:"copyPermissions,attr"`
	Template          string        `xml:"template,attr"`
	ExcelFile         string        `xml:"excel,attr"`
	ExcelStartRow     string        `xml:"startRow,attr"`
	InstanceTag       string        `xml:"instance,attr"`
	ExportToExcel     bool          `xml:"exportToExcel,attr"`
	Sections          []Section     `xml:"section"`
	Elements          []Element     `xml:"element"`
	Replace           []FileReplace `xml:"replace"`
	MatchExcel        []MatchExcel  `xml:"match"`
}

type Driver struct {
	Version       string
	Files         []DriverFile
	Info          os.FileInfo
	PackageDriver bool
	FilePath      string
}

func (d *Driver) Clear() {
	d.Version = ""
	d.Files = []DriverFile{}
	d.PackageDriver = false
	d.FilePath = ""
	d.Info = nil
}

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
	Migrations                []DriverFile `xml:"migrate"`
}

type XOGOutput struct {
	Code  string
	Debug string
}

type Definition struct {
	Action         string `xml:"action,attr"`
	Description    string `xml:"description,attr"`
	Default        string `xml:"default,attr"`
	TransformTypes string `xml:"transformTypes"`
	From           string `xml:"from"`
	To             string `xml:"to"`
	Value          string
}

type Version struct {
	Name           string       `xml:"name,attr"`
	Folder         string       `xml:"folder,attr"`
	DriverFileName string       `xml:"driver,attr"`
	Definitions    []Definition `xml:"definition"`
}

type Package struct {
	Name           string    `xml:"name,attr"`
	Folder         string    `xml:"folder,attr"`
	DriverFileName string    `xml:"driver,attr"`
	Versions       []Version `xml:"version"`
}

const (
	VERSION = 2.0

	LOOKUP    = "Lookups"
	PORTLET   = "Portlets"
	QUERY     = "Queries"
	PROCESS   = "Processes"
	PAGE      = "Pages"
	GROUP     = "Groups"
	MENU      = "Menus"
	OBS       = "Obs"
	OBJECT    = "Objects"
	VIEW      = "Views"
	MIGRATION = "Migrations"

	CUSTOM_OBJECT_INSTANCE     = "CustomObjectInstances"
	RESOURCE_CLASS_INSTANCE    = "ResourceClassInstances"
	WIP_CLASS_INSTANCE         = "WipClassInstances"
	INVESTMENT_CLASS_INSTANCE  = "InvestmentClassInstances"
	TRANSACTION_CLASS_INSTANCE = "TransactionClassInstances"
	RESOURCE_INSTANCE          = "ResourceInstances"
	USER_INSTANCE              = "UserInstances"
	PROJECT_INSTANCE           = "ProjectInstances"
	IDEA_INSTANCE              = "IdeaInstances"
	APPLICATION_INSTANCE       = "ApplicationInstances"
	ASSET_INSTANCE             = "AssetInstances"
	OTHER_INVESTMENT_INSTANCE  = "OtherInvestmentInstances"
	PRODUCT_INSTANCE           = "ProductInstances"
	SERVICE_INSTANCE           = "ServiceInstances"

	ACTION_REPLACE = "replace"
	ACTION_UPDATE  = "update"
	ACTION_REMOVE  = "remove"
	ACTION_INSERT  = "insert"

	FOLDER_READ      = "_read/"
	FOLDER_WRITE     = "_write/"
	FOLDER_MIGRATION = "_migration/"
	FOLDER_DEBUG     = "_debug/"
	FOLDER_PACKAGE   = "_packages/"
	FOLDER_MOCK      = "mock/"

	UNDEFINED      = ""
	OUTPUT_ERROR   = "error"
	OUTPUT_WARNING = "warning"
	OUTPUT_SUCCESS = "success"

	COLUMN_LEFT  = "left"
	COLUMN_RIGHT = "right"
)
