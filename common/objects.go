package common

type Menu struct {
	Code           string `xml:"code,attr"`
	Action         string `xml:"action,attr"`
	TargetPosition int    `xml:"targetPosition,attr"`
	Links          []struct {
		Code           string `xml:"code,attr"`
		TargetPosition int    `xml:"targetPosition,attr"`
	} `xml:"link"`
}

type ViewSection struct {
	SourcePosition string `xml:"sourcePosition,attr"`
	TargetPosition string `xml:"targetPosition,attr"`
	Action string `xml:"action,attr"`
	Fields []struct {
		Code         string `xml:"code,attr"`
		Column       string `xml:"column,attr"`
		Remove       bool   `xml:"remove,attr"`
		InsertBefore string `xml:"insertBefore,attr"`
	}`xml:"field"`
}

type Element struct {
	Type   string `xml:"type,attr"`
	XPath  string `xml:"xpath,attr"`
	Code   string `xml:"code,attr"`
	Action string `xml:"action,attr"`
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
	RemoveObjAssoc    bool			`xml:"removeObjectsAssociation,attr"`
	RemoveSecurity    bool			`xml:"removeSecurity,attr"`
	Template 		  string 		`xml:"template,attr"`
	ExcelFile 		  string 		`xml:"excel,attr"`
	ExcelStartRow 	  string 		`xml:"startRow,attr"`
	InstanceTag 	  string 		`xml:"instance,attr"`
	ExportToExcel 	  bool			`xml:"exportToExcel,attr"`
	Menus             []Menu        `xml:"menu"`
	Sections 		  []ViewSection `xml:"section"`
	Elements 		  []Element		`xml:"element"`
	Replace			  []struct {
		From string `xml:"from"`
		To 	 string `xml:"to"`
	} `xml:"replace"`
	MatchExcel		  []struct {
		Col 			int		`xml:"col,attr"`
		Tag 			string	`xml:"tag,attr"`
		XPath 			string	`xml:"xpath,attr"`
		AttributeName  	string	`xml:"attribute,attr"`
		AttributeValue 	string	`xml:"attributeValue,attr"`
		IsAttribute    	bool  	`xml:"isAttribute,attr"`
		MultiValued    	bool  	`xml:"multiValued,attr"`
		Separator      	string 	`xml:"separator,attr"`
	} `xml:"match"`
}

type Driver struct {
	Files []DriverFile `xml:"file"`
}

type XOGOutput struct {
	Code	string
	Debug	string
}

const LOOKUP 	string 	= "lookups"
const PORTLET 	string 	= "portlets"
const QUERY 	string 	= "queries"
const PROCESS 	string 	= "processes"
const PAGE 		string 	= "pages"
const MENU 		string 	= "menus"
const OBS 		string 	= "obs"
const OBJECT 	string 	= "objects"
const VIEW 		string	= "views"
const MIGRATION	string	= "migrations"

const ACTION_REPLACE string = "replace"
const ACTION_UPDATE  string = "update"
const ACTION_REMOVE  string = "remove"
const ACTION_INSERT  string = "insert"

const CUSTOM_OBJECT_INSTANCE     string = "customObjectInstances"
const RESOURCE_CLASS_INSTANCE    string = "resourceClassInstances"
const WIP_CLASS_INSTANCE         string = "wipClassInstances"
const INVESTMENT_CLASS_INSTANCE  string = "investmentClassInstances"
const TRANSACTION_CLASS_INSTANCE string = "transactionClassInstances"
const RESOURCE_INSTANCE          string = "resourceInstances"
const USER_INSTANCE              string = "userInstances"
const PROJECT_INSTANCE           string = "projectInstances"
const IDEA_INSTANCE              string = "ideaInstances"
const APPLICATION_INSTANCE       string = "applicationInstances"
const ASSET_INSTANCE             string = "assetInstances"
const OTHER_INVESTMENT_INSTANCE  string = "otherInvestmentInstances"
const PRODUCT_INSTANCE           string = "productInstances"
const SERVICE_INSTANCE           string = "serviceInstances"

const FOLDER_READ 		string = "_read/"
const FOLDER_WRITE 		string = "_write/"
const FOLDER_MIGRATION 	string = "_migration/"
const FOLDER_DEBUG 		string = "_debug/"