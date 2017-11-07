# CAS-XOG
Execute XOG files reading and writing in a more easy way

This is a new method of creating XOG files. Using a Driver XML file, you can define with objects you would like to read, write and migrate.


### How to use

1. Download the [lastest stable release](https://github.com/andreluzz/cas-xog/releases/latest) (cas-xog.exe and xogRead.xml);
2. Create a file called [xogEnv.xml](#xog-environment-example), in the same folder of the cas-xog.exe, to define the environments connections configuration;
3. Create a folder called "drivers" with all driver files (.driver) you need to defining the objects you want to read and write;
4. Execute the cas-xog.exe and follow the instructions in the screen.

### Description of structure Driver types

| Type | Description |
| ------ | ------ |
| [`objects`](#type-objects) | Used to read and write objects attributes, actions and links |
| [`views`](#type-views) | Used to read and write views |
| [`processes`](#type-processes) | Used to read and write processes |
| [`lookups`](#type-lookups) | Used to read and write lookups |
| [`portlets`](#type-portlets) | Used to read and write portlets |
| [`queries`](#type-queries) | Used to read and write queries |
| [`pages`](#type-pages) | Used to read and write pages |
| [`menus`](#type-menus) | Used to read and write menus |
| [`obs`](#type-obs) | Used to read and write obs |
| [`groups`](#type-groups) | Used to read and write groups |

### Description of instance Driver types

| Type | Description |
| ------ | ------ |
| [`customObjectInstances`](#type-customObjectInstances) | Used to read and write customObject instances |
| [`resourceClassInstances`](#type-resourceClassInstances) | Used to read and write resourceClass instances |
| [`wipClassInstances`](#type-wipClassInstances) | Used to read and write wipClass instances |
| [`investmentClassInstances`](#type-investmentClassInstances) | Used to read and write investmentClass instances |
| [`transactionClassInstances`](#type-transactionClassInstances) | Used to read and write transactionClass instances |
| [`resourceInstances`](#type-resourceInstances) | Used to read and write resource instances |
| [`userInstances`](#type-userInstances) | Used to read and write user instances |
| [`projectInstances`](#type-projectInstances) | Used to read and write project instances |
| [`ideaInstances`](#type-ideaInstances) | Used to read and write idea instances |
| [`applicationInstances`](#type-applicationInstances) | Used to read and write application instances |
| [`assetInstances`](#type-assetInstances) | Used to read and write asset instances |
| [`otherInvestmentInstances`](#type-otherInvestmentInstances) | Used to read and write otherInvestment instances |
| [`productInstances`](#type-productInstances) | Used to read and write product instances |
| [`serviceInstances`](#type-serviceInstances) | Used to read and write service instances |

# Type `objects`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Object code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 
| `partitionModel` | Defines a new partitionModel if you want to set or change | no |
| `sourcePartition` | When defined reads only elements from this partition code | no |
| `targetPartition` | Used to change the current partition code. Used alone without sourcePartition replaces the tag partitionCode of all xog elements with the defined value. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="idea" path="idea.xml" type="objects" />
    <file code="application" path="application.xml" type="objects" partitionModel="new-corp" />
    <file code="systems" path="systems.xml" type="objects" sourcePartition="IT" targetPartition="NIKU.ROOT" />
    <file code="inv" path="inv.xml" type="objects" targetPartition="NIKU.ROOT" />
</xogdriver>
```

### Sub tag `element`
Used to read only the selected elements from the object

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `type` | Defines what element to read. Available: `attribute`, `action` and `link` | yes | 
| `code` | Code of the attribute that you want to include | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="test_subobj" path="test_subobj.xml" type="objects">
        <element type="attribute" code="attr_auto_number" />
        <element type="attribute" code="novo_atr" />
        <element type="action" code="tst_run_proc_lnk" />
        <element type="link" code="test_subobj.link_test" />
    </file>
</xogdriver>
```

# Type `views`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | View code, use `*` if you want to get all views from one object or the view code if you want a single view. | yes |
| `objectCode` | Object code | yes |
| `path` | Path to the file to be saved on the file system | yes | 
| `sourcePartition` | When defined reads only views from this partition code | no |
| `targetPartition` | Used to replaces the source value tag partitionCode of elements with the defined value. Required sourcePartition tag to use these feature. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="views" code="*" objectCode="obj_system" path="view_0.xml" />
    <file type="views" code="*" objectCode="obj_system" path="view_1.xml" sourcePartition="IT" />
    <file type="views" code="*" objectCode="obj_system" path="view_2.xml" sourcePartition="IT" targetPartition="HR" />
    <file type="views" code="obj_system.audit" objectCode="obj_system" path="view_3.xml" sourcePartition="HR" />
    <file type="views" code="obj_system.audit" objectCode="obj_system" path="view_4.xml" sourcePartition="HR" targetPartition="IT" />
</xogdriver>
```

### Sub tag `section`
Used to read and transform only the selected section from the view. Only single views can use the sub tag `section`, cant be used with `code='*'`.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `action` | Defines what to do in the section. Available: `remove`, `replace`, `insert` and `update`. To use action `update` is required to include the sub tag `field`. | yes | 
| `sourcePosition` | Position of the view in the source. Required for actions: `replace`, `insert` and `update`. | no | 
| `targetPosition` | Position where you want to insert the section in the target view. Required for action `remove`. | no | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="views" code="obj_system.audit" objectCode="obj_system" path="view_3.xml" sourcePartition="HR">
        <section action="insert" sourcePosition="1" targetPosition="1" />
        <section action="replace" sourcePosition="1" targetPosition="3" />
        <section action="remove" targetPosition="3" />
    </file>
</xogdriver>
```

### Sub tag `field`
Used to read and transform only the selected fields from the section. Only sections with action `update` can use sub tag `field`.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the field to transform. | yes | 
| `column` | Section column where to insert the field in the target view. | yes |
| `insertBefore` | Code of the field in the target view where you want to position the new field. | no | 
| `remove` | Defines if the field should be removed from the target view. Use `true` or `false`. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="views" code="obj_system.audit" objectCode="obj_system" path="view_3.xml" sourcePartition="HR">
        <section action="update" sourcePosition="1" targetPosition="1" >
            <field code="analist" insertBefore="created_by" column="left" />
            <field code="status" insertBefore="created_by" column="left" />
            <field code="new_status" column="right" />
            <field code="created_date" remove="true" />
        </section>
    </file>
</xogdriver>
```

# Type `processes`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Process code | yes | 
| `path` | Path to the file to be saved on the file system | yes |
| `copyPermissions` | Defnies the code of the process you want to copy the permissions from. | yes |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="processes" code="PRC_0001" path="PRC_0001.xml" />
    <file type="processes" code="PRC_0002" path="PRC_0002.xml" copyPermissions="PRC_0001" />
</xogdriver>
```

# Type `lookups`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Lookup code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 
| `onlyStructure` | Used to create a lookup with a fake query to prevent error of attributes that have not yet been imported. Only available for dynamic lookups. | no | 
| `sourcePartition` | When defined changes only elements from this partition code. should be used together with targetPartition tag. Only available for static lookups. | no |
| `targetPartition` | Used to change the partition code. Used alone without sourcePartition replaces the tag partitionCode of all lookup values with the defined value. Only available for static lookups. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="INV_APPLICATION_CATEGORY_TYPE" path="INV_APPLICATION_CATEGORY_TYPE.xml" type="lookups" />
    <file code="LOOKUP_FIN_CHARGECODES" path="LOOKUP_FIN_CHARGECODES.xml" onlyStructure="true" type="lookups" />
    <file code="LOOKUP_CAS_XOG_1" path="LOOKUP_CAS_XOG_1.xml" type="lookups" targetPartition="NIKU.ROOT" />
    <file code="LOOKUP_CAS_XOG_2" path="LOOKUP_CAS_XOG_2.xml" type="lookups" sourcePartition="IT" targetPartition="NIKU.ROOT" />
</xogdriver>
```

# Type `portlets`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Portlet code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="cop.teamCapacityLinkable" path="cop.teamCapacityLinkable.xml" type="portlets"/>
</xogdriver>
```

# Type `queries`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Query code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="cop.projectCostsPhaseLinkable" path="cop.projectCostsPhaseLinkable.xml" type="queries"/>
</xogdriver>
```

# Type `pages`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Page code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="pma.ideaFrame" path="pma.ideaFrame.xml" type="pages"/>
</xogdriver>
```

# Type `menus`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Page code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="application" path="menu_result.xml" type="menus" />
</xogdriver>
```

### Sub tag `section`
Used to read only the selected section from the menu

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the attribute that you want to include. | yes | 
| `action` | Defines what to do in the target menu. Available: `insert` and `update`. To use action update is required to include the sub tag `link`. | yes | 
| `targetPosition` | Position where you want to insert the section in the target menu. | no | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="application" path="menu_result_section_link.xml" type="menus">
        <section action="insert" code="menu_sec_cas_xog" targetPosition="2" />
    </file>
</xogdriver>
```

### Sub tag `link`
Used to read only the selected links inside a section tag from the menu

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the link that you want to include. | yes |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="application" path="menu_result_section_link.xml" type="menus">
        <section action="update" code="npt.personal">
            <link code="odf.obj_testeList" />
        </section>
        <section action="insert" code="menu_sec_cas_xog" targetPosition="2">
            <link code="cas_proc_running_tab" />
        </section>
    </file>
</xogdriver>
```

# Type `obs`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | OBS code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="department" path="obs_department.xml" type="obs" />
</xogdriver>
```

# Type `groups`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Group code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="cop.systemAdministrator" path="systemAdministrator.xml" type="groups" />
</xogdriver>
```

# Global Sub Tags
Sub tags that can be used in any type of `file` tag.

### Sub Tag `replace`
Used to do a replace from one string to another in the xog result.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `from` | Defines which string should be searched for to be changed | yes | 
| `to` | String that will replace what was defined in the `from` tag | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file code="PRC_0001" path="PRC_0001.xml" type="processes">
        <replace>
            <from>endpoint="http://development.server.com"</from>
            <to>endpoint="http://production.server.com"</to>
        </replace>
        <replace>
            <from>set var="xogUser" value="adminXogUser"</from>
            <to>set var="xogUser" value="anotherAdminXogUser"</to>
        </replace>
    </file>
</xogdriver>
```

### Sub Tag `element`
Used to do a remove elements in the xog result using xpath.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `action` | Action remove in the element tag. Only action `remove` is available. | yes | 
| `xpath` | String that defines the path in the XML to the element you want to remove. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="pages" code="projmgr.projectPageFrame" path="page_project.xml">
        <element action="remove" xpath="//OBSAssocs" />
        <element action="remove" xpath="//Security" />
    </file>
    <file type="obs" code="department" path="obs_department.xml">
        <element action="remove" xpath="//associatedObject" />
        <element action="remove" xpath="//Security" />
        <element action="remove" xpath="//rights" />
    </file>
    <file type="groups" code="ObjectAdmin" path="ObjectAdmin.xml">
        <element action="remove" xpath="/NikuDataBus/groups/group/members" />
    </file>
</xogdriver>
```

# Package creation and deploy
This feature should be used to deploy structures and instances in a more consolidated way. Should be created a zip containing: a package file (.package), one or more driver files (.driver) and folders for versions and the XOG xml files.

### Tag `package`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Defines the name that will be displayed to the user. | yes | 
| `folder` | The first level folder inside the zip file that represents the package. | yes |
| `driver` | The default driver for all package versions. If the version has no driver this will be used.  | yes |

### Sub tag `version`
This tag is required every package should have at least one version

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Defines the name that will be displayed to the user to select the version. If there is only one will be chosen automatically. | yes | 
| `folder` | The folder that represents the files for this version. | yes |
| `driver` | Defines the driver for this version. Can be used to define a version with demo data and other with only structure for example.  | no |

### Sub tag `definition`
This tag is not required, should be used to define questions to the user to answer and use the result to replace values in the XOG xml files.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `action` | Defines the action to the replace. Available: `changePartitionModel`, `changePartition` and `replaceString`.  | yes | 
| `description` | The texto to display the question to the user whrn installing the package. | yes |
| `default` | The default value for this definition.  | no |
| `transformTypes` | Define in what types of files this transformation should be performed. If not defined the replace will be used in all XOG xml files. Use the same types as in [`driver types`](#description-of-driver-types)  | no |
| `from` | The string that should be found in the XOG xml files. Required when action is `replaceString`.  | no |
| `to` | The string that should replace the current value. Use the special string `"##DEFINITION_VALUE##` to set the position for the value defined by the user. Required when action is `replaceString`.  | no |

### Package file example
```xml
<?xml version="1.0" encoding="utf-8"?>
<package name="CAS-FIN" folder="cas-fin/" driver="cas_fin.driver">
    <version name="Version Oracle" folder="oracle/">
        <definition action="changePartitionModel" description="Target partition model" default="corporate" />
        <definition action="changePartition" description="Target partition" default="NIKU.ROOT" />
        <definition action="replaceString" description="Set processes scripts XOG user" default="xogadmin">
            <transformTypes>processes</transformTypes>
            <from>set value="xogadmin" var="user_param"</from>
            <to>set value="##DEFINITION_VALUE##" var="user_param"</to>
        </definition>
    </version>
    <version name="Version MS SQL Server" folder="sqlserver/" driver="cas_fin-sql.driver">
        <definition action="changePartitionModel" description="Target partition model" default="corporate" />
        <definition action="changePartition" description="Target partition" default="NIKU.ROOT" />
        <definition action="replaceString" description="Set processes scripts XOG user" default="xogadmin">
            <transformTypes>processes</transformTypes>
            <from>set value="xog_admin" var="user_param"</from>
            <to>set value="##DEFINITION_VALUE##" var="user_param"</to>
        </definition>
    </version>
</package>
```

### ZIP folders and files structure
```
└── cas-fin
    ├── oracle
    │   ├── lookups
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   └── cas_fin_lkp_lista_filas_exec.xml
    │   ├── objects
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   └── cas_fin_lkp_lista_filas_exec.xml
    │   └── processes
    │       ├── cas_fin_prc_config_trans.xml
    │       ├── cas_fin_prc_cria_trans_cronog.xml
    │       ├── cas_fin_prc_fila_trans_aguarda.xml
    │       └── cas_fin_prc_processa_fila.xml
    ├── sqlserver
    │   ├── lookups
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   └── cas_fin_lkp_lista_filas_exec.xml
    │   ├── objects
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   ├── cas_fin_lkp_list_taks_proj.xml
    │   │   └── cas_fin_lkp_lista_filas_exec.xml
    │   ├── processes
    │   │   ├── cas_fin_prc_config_trans.xml
    │   │   ├── cas_fin_prc_cria_trans_cronog.xml
    │   │   ├── cas_fin_prc_fila_trans_aguarda.xml
    │   │   └── cas_fin_prc_processa_fila.xml
    │   └── cas_fin-sql.driver
    ├── cas_fin.driver
    └── cas-fin.package
```

To install the package the user should save the zip file inside a folder named `package` in the same directory of the `cas-xog.exe` file.

# Data migration 
This feature should be used to export instances to an excel file and read data from excel file to a XOG template creating an xml to import data to the environment.

### Export data to excel
Should be used with a [driver instance type](#description-of-instance-driver-types) to read data from the environment and save the match attributes to an excel file.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Defines the name that will be displayed to the user. | yes | 
| `path` | Path to the file to be saved on the file system. | yes |
| `exportToExcel` | If set to true creates an excel file with the matched data. | yes |
| `excel` | The name of the file to export the data.  | yes |
| `instance` | The name of the main tag that represents the instance object that is being read.  | yes |

### Sub tag `match`
This tag is required for export to excel data.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `attribute` | Defines the attribute from the element you want to get the data. If no xpath is defined then we get the value from the main instance element defined.  | no | 
| `xpath` | A string representing the path to the element you want to get the data. If no attribute value is defined then we get the value from the tag text. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="resourceInstances" code="*" path="res.xml" exportToExcel="true" excel="res.xlsx" instance="Resource">
        <match attribute="resourceId" />
        <match xpath="//PersonalInformation" attribute="displayName" />
        <match xpath="//PersonalInformation" attribute="emailAddress" />
        <match xpath="//PersonalInformation" attribute="firstName" />
        <match xpath="//PersonalInformation" attribute="lastName" />
        <match xpath="//OBSAssoc[@id='corpLocationOBS']" attribute="unitPath" />
        <match xpath="//OBSAssoc[@id='resourcePool']" attribute="unitPath" />
        <match xpath="//ColumnValue[@name='partition_code']" />
    </file>
</xogdriver>
```

### Read data from excel to create XOG instances xml 
Should be used with to create an XOG xml file with many instances as lines in the excel file.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `path` | Path to the file to be saved on the file system. | yes |
| `template` | Path to the template that should be used to create the XOG xml file. | yes |
| `instance` | The name of the main tag that represents the instance object that should be created. | yes |
| `excel` | The path to the excel file with the data. | yes |
| `startRow` | The line number in the excel file that we will start reading to create the instances. Default value is 1. | no |

### Sub tag `match`
This tag is required for export to excel data.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `col` | Define which column of excel we'll get the data to include in the XOG xml file. | yes |
| `attribute` | Defines the attribute in the element you want to set the data. If no xpath is defined then we set this attribute in the main element instance.  | no |
| `xpath` | A string representing the path to the element you want to set the data. If no attribute value is defined then we set the value as a tag text. | no |
| `multiValued` | If set to true defines that this element should be treated as multi value. | no |
| `separator` | Defines what character were used to separate the options in the multi value on the excel data. Default value is ';'. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <file type="migrations" path="subs.xml" template="template.xml" instance="instance" excel="dados.xlsx" startRow="2" >
        <match col="1" attribute="instanceCode" />
        <match col="1" xpath="//ColumnValue[@name='code']" />
        <match col="2" xpath="//ColumnValue[@name='name']" />
        <match col="3" xpath="//ColumnValue[@name='status']" />
        <match col="4" xpath="//ColumnValue[@name='multivalue_status']" multiValued="true" separator=";" />
        <match col="5" xpath="//ColumnValue[@name='analista']" />
    </file>
</xogdriver>
```

### Template file example
```xml
<NikuDataBus xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="../xsd/nikuxog_customObjectInstance.xsd">
	<Header action="write" externalSource="NIKU" objectType="customObjectInstance" version="8.0" />
	<customObjectInstances objectCode="obj_sistema">
		<instance instanceCode="" objectCode="obj_sistema">
			<CustomInformation>
				<ColumnValue name="code"></ColumnValue>
				<ColumnValue name="name"></ColumnValue>
				<ColumnValue name="status_novo"></ColumnValue>
				<ColumnValue name="analista"></ColumnValue>
				<ColumnValue name="multivalue_status"></ColumnValue>
			</CustomInformation>
		</instance>
	</customObjectInstances>
</NikuDataBus>
```

# XOG Environment example:

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogenvs version="2.0">
    <env name="Development">
        <username>username</username>
        <password>12345</password>
        <endpoint>http://development.server.com</endpoint>
    </env>
    <env name="Quality">
        <username>username</username>
        <password>12345</password>
        <endpoint>http://quality.server.com</endpoint>
    </env>
    <env name="Production">
        <username>username</username>
        <password>12345</password>
        <endpoint>https://production.server.com</endpoint>
    </env>
</xogenvs>
```
