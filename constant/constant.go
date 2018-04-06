package constant

//Defines the constants for the system
const (
	Version = 2.0

	TypeLookup    = "Lookups"
	TypePortlet   = "Portlets"
	TypeQuery     = "Queries"
	TypeProcess   = "Processes"
	TypePage      = "Pages"
	TypeMenu      = "Menus"
	TypeObject    = "Objects"
	TypeView      = "Views"
	TypeMigration = "Migrations"

	TypeCustomObjectInstance     = "CustomObjectInstances"
	TypeResourceClassInstance    = "ResourceClassInstances"
	TypeWipClassInstance         = "WipClassInstances"
	TypeInvestmentClassInstance  = "InvestmentClassInstances"
	TypeTransactionClassInstance = "TransactionClassInstances"
	TypeResourceInstance         = "ResourceInstances"
	TypeUserInstance             = "UserInstances"
	TypeProjectInstance          = "ProjectInstances"
	TypeProgramInstance          = "ProgramInstances"
	TypeGroupInstance            = "GroupInstances"
	TypeOBSInstance              = "OBSInstances"
	TypeThemeInstance            = "ThemeInstances"
	TypeDocumentInstance         = "DocumentInstances"

	ActionReplace 		  = "replace"
	ActionUpdate  		  = "update"
	ActionRemove  		  = "remove"
	ActionInsert  		  = "insert"

	Read    = "r"
	Write   = "w"
	Migrate = "m"

	FolderRead      = "_read/"
	FolderWrite     = "_write/"
	FolderMigration = "_migration/"
	FolderDebug     = "_debug/"
	FolderPackage   = "_packages/"
	FolderMock      = "mock/"

	Undefined     = ""
	OutputError   = "error"
	OutputWarning = "warning"
	OutputSuccess = "success"
	OutputIgnored = "ignored"

	PackageActionChangePartitionModel = "changePartitionModel"
	PackageActionChangePartition      = "changePartition"
	PackageActionReplaceString        = "replaceString"

	ColumnLeft  = "left"
	ColumnRight = "right"

	ElementTypeLink        = "link"
	ElementTypeAction      = "action"
	ElementTypeAttribute   = "attribute"
	ElementTypeActionGroup = "actionGroup"

	Source = "source"
	Target = "target"

	DefaultInstanceTag = "instance"
)
