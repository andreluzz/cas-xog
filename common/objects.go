package common

import "os"

type SectionLink struct {
	Code	string `xml:"code,attr"`
}

type SectionField struct {
	Code         string `xml:"code,attr"`
	Column       string `xml:"column,attr"`
	Remove       bool   `xml:"remove,attr"`
	InsertBefore string `xml:"insertBefore,attr"`
}

type Section struct {
	Code 			string 			`xml:"code,attr"`
	SourcePosition 	string 			`xml:"sourcePosition,attr"`
	TargetPosition 	string 			`xml:"targetPosition,attr"`
	Action 			string 			`xml:"action,attr"`
	Links			[]SectionLink 	`xml:"link"`
	Fields 			[]SectionField 	`xml:"field"`
}

type Element struct {
	Type   string `xml:"type,attr"`
	XPath  string `xml:"xpath,attr"`
	Code   string `xml:"code,attr"`
	Action string `xml:"action,attr"`
}

type FileReplace struct {
	From string `xml:"from"`
	To 	 string `xml:"to"`
}

type MatchExcel	struct {
	Col 			int		`xml:"col,attr"`
	Tag 			string	`xml:"tag,attr"`
	XPath 			string	`xml:"xpath,attr"`
	AttributeName  	string	`xml:"attribute,attr"`
	AttributeValue 	string	`xml:"attributeValue,attr"`
	IsAttribute    	bool  	`xml:"isAttribute,attr"`
	MultiValued    	bool  	`xml:"multiValued,attr"`
	Separator      	string 	`xml:"separator,attr"`
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
	CopyPermissions   string		`xml:"copyPermissions,attr"`
	Template 		  string 		`xml:"template,attr"`
	ExcelFile 		  string 		`xml:"excel,attr"`
	ExcelStartRow 	  string 		`xml:"startRow,attr"`
	InstanceTag 	  string 		`xml:"instance,attr"`
	ExportToExcel 	  bool			`xml:"exportToExcel,attr"`
	Sections 		  []Section		`xml:"section"`
	Elements 		  []Element		`xml:"element"`
	Replace			  []FileReplace `xml:"replace"`
	MatchExcel		  []MatchExcel	`xml:"match"`
}

type Driver struct {
	Files 			[]DriverFile `xml:"file"`
	Info  			os.FileInfo
	PackageDriver 	bool
	FilePath 		string
}

type XOGOutput struct {
	Code	string
	Debug	string
}

type Definition struct {
	Action 			string	`xml:"action,attr"`
	Description 	string	`xml:"description,attr"`
	Default 		string	`xml:"default,attr"`
	TransformTypes	string	`xml:"transformTypes"`
	From 			string	`xml:"from"`
	To				string	`xml:"to"`
	Value	 		string
}

type Version struct {
	Name 			string `xml:"name,attr"`
	Folder 			string `xml:"folder,attr"`
	DriverFileName	string `xml:"driver,attr"`
	Definitions		[]Definition `xml:"definition"`
}

type Package struct {
	Name 			string `xml:"name,attr"`
	Folder 			string `xml:"folder,attr"`
	DriverFileName	string `xml:"driver,attr"`
	Versions 		[]Version `xml:"version"`
}

const LOOKUP 	= "lookups"
const PORTLET 	= "portlets"
const QUERY 	= "queries"
const PROCESS 	= "processes"
const PAGE 		= "pages"
const GROUP 	= "groups"
const MENU 		= "menus"
const OBS 		= "obs"
const OBJECT 	= "objects"
const VIEW 		= "views"
const MIGRATION	= "migrations"

const ACTION_REPLACE = "replace"
const ACTION_UPDATE  = "update"
const ACTION_REMOVE  = "remove"
const ACTION_INSERT  = "insert"

const CUSTOM_OBJECT_INSTANCE     = "customObjectInstances"
const RESOURCE_CLASS_INSTANCE    = "resourceClassInstances"
const WIP_CLASS_INSTANCE         = "wipClassInstances"
const INVESTMENT_CLASS_INSTANCE  = "investmentClassInstances"
const TRANSACTION_CLASS_INSTANCE = "transactionClassInstances"
const RESOURCE_INSTANCE          = "resourceInstances"
const USER_INSTANCE              = "userInstances"
const PROJECT_INSTANCE           = "projectInstances"
const IDEA_INSTANCE              = "ideaInstances"
const APPLICATION_INSTANCE       = "applicationInstances"
const ASSET_INSTANCE             = "assetInstances"
const OTHER_INVESTMENT_INSTANCE  = "otherInvestmentInstances"
const PRODUCT_INSTANCE           = "productInstances"
const SERVICE_INSTANCE           = "serviceInstances"

const FOLDER_READ 		= "_read/"
const FOLDER_WRITE 		= "_write/"
const FOLDER_MIGRATION 	= "_migration/"
const FOLDER_DEBUG 		= "_debug/"
const FOLDER_PACKAGE 	= "_packages/"
const FOLDER_MOCK 		= "mock/"

const UNDEFINED 		= ""
const OUTPUT_ERROR 		= "error"
const OUTPUT_WARNING 	= "warning"
const OUTPUT_SUCCESS 	= "success"

const (
	COLUMN_LEFT = "left"
	COLUMN_RIGHT = "right"
)