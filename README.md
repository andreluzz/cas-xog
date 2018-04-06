[![codecov](https://codecov.io/gh/andreluzz/cas-xog/branch/master/graph/badge.svg)](https://codecov.io/gh/andreluzz/cas-xog)
[![Build status](https://ci.appveyor.com/api/projects/status/4lixe3lc9c9pxoa5?svg=true)](https://ci.appveyor.com/project/andreluzz/cas-xog)
[![Go Report Card](https://goreportcard.com/badge/github.com/andreluzz/cas-xog)](https://goreportcard.com/report/github.com/andreluzz/cas-xog)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Github All Releases](https://img.shields.io/github/downloads/andreluzz/cas-xog/total.svg)](https://github.com/andreluzz/cas-xog/releases/latest)
[![GitHub release](https://img.shields.io/github/release/andreluzz/cas-xog.svg)](https://github.com/andreluzz/cas-xog/releases/latest)

# CAS-XOG
Execute XOG files reading and writing in a more easy way

This is a new method of creating XOG files. Using a Driver XML file, you can define with objects you would like to read, write and migrate.


### How to use

1. Download the [lastest stable release](https://github.com/andreluzz/cas-xog/releases/latest) (cas-xog.exe and xogRead.xml);
2. Create a file called [xogEnv.xml](#xog-environment-example), in the same folder of the cas-xog.exe, to define the environments connections configuration;
3. Create a folder called "drivers" with all driver files (.driver) you need to defining the objects you want to read and write;
4. Execute the cas-xog.exe and follow the instructions in the screen.

### Other contents
* [Global attributes](#global-attributes) 
* [Global Sub Tags](#global-sub-tags) 
* [Package creation and deploy](#package-creation-and-deploy)
* [Data migration](#data-migration)

### Description of structure Driver tags

| Tag | Description |
| ------ | ------ |
| [`object`](#tag-object) | Used to read and write objects attributes, actions and links. |
| [`view`](#tag-view) | Used to read and write views. |
| [`process`](#tag-process) | Used to read and write processes. |
| [`lookup`](#tag-lookup) | Used to read and write lookups. |
| [`portlet`](#tag-portlet) | Used to read and write portlets. |
| [`query`](#tag-query) | Used to read and write queries. |
| [`page`](#tag-page) | Used to read and write pages. |
| [`menu`](#tag-menu) | Used to read and write menus. |

### Description of instance Driver tags

| Tag | Description |
| ------ | ------ |
| [`customObjectInstance`](#tag-customobjectinstance) | Used to read and write customObject instances. |
| [`resourceClassInstance`](#tag-resourceclassinstance) | Used to read and write resourceClass instances. |
| [`wipClassInstance`](#tag-wipclassinstance) | Used to read and write wipClass instances. |
| [`investmentClassInstance`](#tag-investmentclassinstance) | Used to read and write investmentClass instances. |
| [`transactionClassInstance`](#tag-transactionclassinstance) | Used to read and write transactionClass instances. |
| [`resourceInstance`](#tag-resourceinstance) | Used to read and write resource instances. |
| [`userInstance`](#tag-userinstance) | Used to read and write user instances. |
| [`projectInstance`](#tag-projectinstance) | Used to read and write project instances. |
| [`ideaInstance`](#tag-ideainstance) | Used to read and write idea instances. |
| [`applicationInstance`](#tag-applicationinstance) | Used to read and write application instances. |
| [`assetInstance`](#tag-assetinstance) | Used to read and write asset instances. |
| [`otherInvestmentInstance`](#tag-otherinvestmentinstance) | Used to read and write otherInvestment instances. |
| [`productInstance`](#tag-productinstance) | Used to read and write product instances. |
| [`serviceInstance`](#tag-serviceinstance) | Used to read and write service instances. |
| [`obsInstance`](#tag-obsinstance) | Used to read and write OBS instances. |
| [`themeInstance`](#tag-themeinstance) | Used to read and write UI Theme instances. |
| [`groupInstance`](#tag-groupinstance) | Used to read and write groups. |

## Tag `object`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Object code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 
| `partitionModel` | Used when you need to set a new partitionModel or change the current one. | no |
| `targetPartition` | Used to change elements partition code to the defined value. When uses alone without sourcePartition replaces the tag partitionCode on all elements. | no |
| `sourcePartition` | Used to read only elements from this partition code. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <object code="idea" path="idea.xml" />
    <object code="application" path="application.xml" partitionModel="new-corp" />
    <object code="systems" path="systems.xml" sourcePartition="IT" targetPartition="NIKU.ROOT" />
    <object code="inv" path="inv.xml" targetPartition="NIKU.ROOT" />
</xogdriver>
```

### Sub tag `element`
Used to read only the selected elements from the object.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `type` | Defines what element to read. Availables types: `attribute`, `action` and `link`. | yes | 
| `code` | Code of the element that you want to include. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <object code="test_subobj" path="test_subobj.xml">
        <element type="attribute" code="attr_auto_number" />
        <element type="attribute" code="novo_atr" />
        <element type="action" code="tst_run_proc_lnk" />
        <element type="link" code="test_subobj.link_test" />
    </object>
</xogdriver>
```

## Tag `view`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Defines the views that should be read. Use `*` for reading all views from object; Use `*pattern_string` to choose all the views where the code contains the pattern string; Use the view code without `*` if you want a single view. | yes |
| `objectCode` | Object code. | yes |
| `path` | Path where the file will be saved on the file system. | yes | 
| `sourcePartition` | When defined reads only views from this partition code. | no |
| `targetPartition` | Used to replaces the source value tag partitionCode of elements with the defined value. If you want to use this feature, the `sourcePartition` tag is required. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <view code="*" objectCode="obj_system" path="view_0.xml" />
    <view code="*" objectCode="obj_system" path="view_1.xml" sourcePartition="IT" />
    <view code="*" objectCode="obj_system" path="view_2.xml" sourcePartition="IT" targetPartition="HR" />
    <view code="*project" objectCode="project" path="project.xml" sourcePartition="IT" />
    <view code="obj_system.audit" objectCode="obj_system" path="view_3.xml" sourcePartition="HR" />
    <view code="obj_system.audit" objectCode="obj_system" path="view_4.xml" sourcePartition="HR" targetPartition="IT" />
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
    <view code="obj_system.audit" objectCode="obj_system" path="view_3.xml" sourcePartition="HR">
        <section action="insert" sourcePosition="1" targetPosition="1" />
        <section action="replace" sourcePosition="1" targetPosition="3" />
        <section action="remove" targetPosition="3" />
    </view>
</xogdriver>
```

### Sub tag `field`
Used to read and transform only the selected fields from the section. Only sections with action `update` can use sub tag `field`.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the field to transform. | yes | 
| `remove` | Use `true` if the field should be removed from the target view. Default value is considered `false`. | no |
| `column` | The section's column where the field will be inserted in the target view. Required if `remove` tag is not defined as `true`. | no |
| `insertBefore` | If an attribute code is defined in this tag, the new field will be positioned before this attributte in the target view. If not, will insert the field as in the last position of the column. | no | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <view code="obj_system.audit" objectCode="obj_system" path="view_3.xml" sourcePartition="HR">
        <section action="update" sourcePosition="1" targetPosition="1" >
            <field code="analist" column="left" insertBefore="created_by" />
            <field code="status" column="left" insertBefore="created_by" />
            <field code="new_status" column="right" />
            <field code="created_date" remove="true" />
        </section>
    </view>
</xogdriver>
```

### Sub tag `element`
Used to read view actions and actions group.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the element that you want to include or remove. | yes | 
| `type` | Defines what element to read. Availables types: `actionGroup` and `action`. | yes | 
| `action` | Defines what to do with the element in the target environment. Availables actions: `insert` and `remove`. | yes |
| `insertBefore` | When a code is defined in this tag, the action or group will be positioned before this element in the target view. If not, it will be inserted as the last element. | no | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <view code="cas_environmentProperties" objectCode="cas_environment" path="cas_environment_view.xml" sourcePartition="NIKU.ROOT">
        <element type="actionGroup" code="group_code" action="remove" />
        <element type="actionGroup" code="group_code" action="insert" insertBefore="target_group_code"/>
        <element type="action" code="action_code" action="insert" />
        <element type="action" code="action_code" action="insert" insertBefore="target_action_code"/>
        <element type="action" code="action_code" action="remove" />
    </object>
</xogdriver>
```

## Tag `process`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Process code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes |
| `copyPermissions` | The code of the process you want to copy the permissions from. | yes |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <process code="PRC_0001" path="PRC_0001.xml" />
    <process code="PRC_0002" path="PRC_0002.xml" copyPermissions="PRC_0001" />
</xogdriver>
```

## Tag `lookup`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Lookup code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 
| `onlyStructure` | Used to create a lookup with a fake query to prevent error of attributes that have not yet been imported. Only available for dynamic lookups. | no | 
| `sourcePartition` | When defined changes only elements from this partition code. Should be used together with targetPartition tag. Only available for static lookups. | no |
| `targetPartition` | Used to change the partition code. Used alone without sourcePartition replaces the tag partitionCode of all lookup values with the defined value. Only available for static lookups. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <lookup code="INV_APPLICATION_CATEGORY_TYPE" path="INV_APPLICATION_CATEGORY_TYPE.xml" />
    <lookup code="LOOKUP_FIN_CHARGECODES" path="LOOKUP_FIN_CHARGECODES.xml" onlyStructure="true" />
    <lookup code="LOOKUP_CAS_XOG_1" path="LOOKUP_CAS_XOG_1.xml" targetPartition="NIKU.ROOT" />
    <lookup code="LOOKUP_CAS_XOG_2" path="LOOKUP_CAS_XOG_2.xml" sourcePartition="IT" targetPartition="NIKU.ROOT" />
</xogdriver>
```

### Sub tag `nsql`
Used to replace the nsql query inside an dynamic lookup.

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
	<lookup code="LOOKUP_CAS_XOG_1" path="LOOKUP_CAS_XOG_1.xml" targetPartition="NIKU.ROOT">
		<nsql>
            SELECT @SELECT:RESOURCES.ID:ID@,
            @SELECT:RESOURCES.LAST_NAME:LAST_NAME@,
            @SELECT:RESOURCES.FIRST_NAME:FIRST_NAME@,
            @SELECT:RESOURCES.FULL_NAME:FULL_NAME@,
            @SELECT:RESOURCES.UNIQUE_NAME:UNIQUE_NAME@,
            @SELECT:RESOURCES.UNIQUE_NAME:UNIQUE_CODE@
            FROM SRM_RESOURCES RESOURCES
            WHERE @FILTER@
            AND @WHERE:SECURITY:RESOURCE:RESOURCES.ID@
            @BROWSE-ONLY:
            AND RESOURCES.IS_ACTIVE = 1
            :BROWSE-ONLY@
		</nsql>
	</lookup>
</xogdriver>
```

## Tag `portlet`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Portlet code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <portlet code="cop.teamCapacityLinkable" path="cop.teamCapacityLinkable.xml" />
</xogdriver>
```

## Tag `query`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Query code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <query code="cop.projectCostsPhaseLinkable" path="cop.projectCostsPhaseLinkable.xml" />
</xogdriver>
```

## Tag `page`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Page code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <page code="pma.ideaFrame" path="pma.ideaFrame.xml" />
</xogdriver>
```

## Tag `menu`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Menu code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <menu code="application" path="menu_result.xml" />
</xogdriver>
```

### Sub tag `section`
Used to read only the selected section from the menu.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the attribute that you want to include. | yes | 
| `action` | Defines what to do in the target menu. Available actions: `insert` and `update`. To use the update action, you need to include the sub tag [`link`](#sub-tag-link). | yes | 
| `targetPosition` | Position where you want to insert the section in the target menu. | no | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <menu code="application" path="menu_result_section_link.xml">
        <section action="insert" code="menu_sec_cas_xog" targetPosition="2" />
    </menu>
</xogdriver>
```

### Sub tag `link`
Used to read only the selected links inside a section tag from the menu.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Code of the link that you want to include. | yes |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <menu code="application" path="menu_result_section_link.xml">
        <section action="update" code="npt.personal">
            <link code="odf.obj_testeList" />
        </section>
        <section action="insert" code="menu_sec_cas_xog" targetPosition="2">
            <link code="cas_proc_running_tab" />
        </section>
    </menu>
</xogdriver>
```

## Tag `customObjectInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `objectCode` | Defines the code of the custom object you want to read the instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <customObjectInstance code="*" objectCode="obj_system" path="instances.xml" />
</xogdriver>
```

## Tag `resourceClassInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <resourceClassInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `wipClassInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <wipClassInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `investmentClassInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <investmentClassInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `transactionClassInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <transactionClassInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `resourceInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <resourceInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `userInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <userInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `projectInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <projectInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `ideaInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <ideaInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `applicationInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <applicationInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `assetInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <assetInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `otherInvestmentInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <otherInvestmentInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `productInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <productInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `serviceInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Instance codes. Use code equals * to get all instances. | yes |
| `path` | Path where the file will be saved on the file system. | yes |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <serviceInstance code="*" path="instances.xml" />
</xogdriver>
```

## Tag `obsInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | OBS code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <obsInstance code="department" path="obs_department.xml" />
</xogdriver>
```

## Tag `themeInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | UI Theme code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <themeInstance code="tealgrey" path="tealgrey.xml" />
</xogdriver>
```

## Tag `groupInstance`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Group code. | yes | 
| `path` | Path where the file will be saved on the file system. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <groupInstance code="cop.systemAdministrator" path="systemAdministrator.xml" />
</xogdriver>
```

# Global Attributes
Attributes that can be used in any [structure](#description-of-structure-driver-tags) and [instance](#description-of-instance-driver-tags) tags.

### Attribute `ignoreReading`
Used to ignore the reading from source environment. Use it to avoid reading more than once the same structure. Intended to be used when you need to write the same structure more than once to resolve cross dependencies issues.

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <object code="idea" path="idea.xml" ignoreReading="true" />
</xogdriver>
```

# Global Sub Tags
Sub tags that can be used in any [structure](#description-of-structure-driver-tags) and [instance](#description-of-instance-driver-tags) tags.

### Sub Tag `replace`
Used to do a replace one string with another one in the xog result.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `from` | Defines which string should be replaced. | yes | 
| `to` | String that will replace the one defined in the `from` tag. | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <process code="PRC_0001" path="PRC_0001.xml">
        <replace>
            <from>endpoint="http://development.server.com"</from>
            <to>endpoint="http://production.server.com"</to>
        </replace>
        <replace>
            <from>set var="xogUser" value="adminXogUser"</from>
            <to>set var="xogUser" value="anotherAdminXogUser"</to>
        </replace>
    </process>
</xogdriver>
```

### Sub Tag `element`
Used to do a remove elements or element attribute from the xog result using xpath.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `action` | Define what action should be done. `insert` and `remove` are available. `insert` may be used to create or replace. | yes | 
| `xpath` | String that defines the path in the XML to the element you want to transform. | yes | 
| `attribute` | String that defines the attribute from the element define in the xpath. | no |
| `value` | String that defines the value to insert or replace in the attribute from the element define in the xpath. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <page code="projmgr.projectPageFrame" path="page_project.xml">
        <element action="remove" xpath="//OBSAssocs" />
        <element action="remove" xpath="//Security" />
    </page>
    <obs code="department" path="obs_department.xml">
        <element action="remove" xpath="//associatedObject" />
        <element action="remove" xpath="//Security" />
        <element action="remove" xpath="//rights" />
    </obs>
    <groupInstance code="ObjectAdmin" path="ObjectAdmin.xml">
        <element action="remove" xpath="/NikuDataBus/groups/group/members" />
    </groupInstance>
    <projectInstance code="PR1126" path="projects.xml">
        <element action="remove" xpath="//Resource" attribute="projectRoleID" />
    </projectInstance>
    <ideaInstance code="ID1062" path="ID1062.xml">
        <element action="remove" xpath="//Idea" attribute="entityCode" />
        <element action="insert" xpath="//Idea" attribute="financialLocation" value="br" />
        <element action="insert" xpath="//OBSAssocs">
            <xml>
                <![CDATA[
                    <OBSAssoc id="teste_1" name="OBS 1" unitPath="/br/rj"/>
                    <OBSAssoc id="teste_2" name="OBS 2" unitPath="/pres/dir/ger/coord"/>
                ]]>
            </xml>
        </element>
    </ideaInstance>
</xogdriver>
```

### Sub Tag `filter`
Used to read instances using custom filter values. When defined all standard filters will be removed and only the defined ones will be used.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Name of the object attribute used to filter. | yes | 
| `criteria` | How the filter should be used. Can be used: OR, EQUALS, BETWEEN, BEFORE, AFTER.  | yes | 
| `customAttribute` | Defines if this is a custom attribute filter or not. | no | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <projectInstance path="prj_filtered.xml">
        <filter name="start" criteria="BETWEEN">2015-01-07,2017-01-15</filter>
        <filter name="custom_status" customAttribute="true" criteria="EQUALS">1</filter>
    </projectInstance>
</xogdriver>
```

### Sub Tag `args`
Used to read instances using custom header args values. When defined all standard header args will be removed and only the defined ones will be used.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Name of the header arg. | yes | 
| `value` | Value to the header arg.  | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <projectInstance path="prj_filtered.xml">
        <args name="include_tasks" value="false" />
    </projectInstance>
</xogdriver>
```

# Package creation and deploy
This feature should be used to deploy structures and instances in a more consolidated and organized way. You need to create a zip containing: a package file (.package), one or more driver files (.driver) and folders for versions and the XOG xml files.

### Tag `package`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Defines what will be displayed to the user as the package name. | yes | 
| `folder` | The first level folder inside the zip file that represents the package. | yes |
| `driver` | The default driver for all package versions. If the version has no driver this one will be used.  | yes |

### Sub tag `version`
This tag is required, as every package should have at least one version. If there is only one version, it will be chosen automatically.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Defines what will be displayed to the user as the version name. | yes |
| `folder` | The folder that represents the files for this version. | yes |
| `driver` | Defines the driver for this version. Can be used to define a version with demo data and other with only structure for example. | no |

### Sub tag `definition`
This tag is not required and should be used to define questions to the user to answer. The answers will be used to change specific definitions withih XOG xml files of the package.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `action` | Defines the desired action. Available actions: `changePartitionModel`, `changePartition` and `replaceString`.  | yes | 
| `description` | The question text that is asked to the user when installing the package. | yes |
| `default` | The default value for this definition.  | no |
| `transformTypes` | Define in what types of files this action should be performed, separated by single commas. If not defined the action will be performed in all XOG xml files. Use the same types defined in [`driver types`](#description-of-driver-types)  | no |
| `from` | Defines which string should be replaced in the XOG xml files. Required when action is `replaceString`.  | no |
| `to` | String that will replace the one defined in the `from` tag. Use the special string `"##DEFINITION_VALUE##` to set the position for the value defined by the user. Required when action is `replaceString`.  | no |

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

To install the package the user should save the zip file inside a folder named `packages` in the same directory of the `cas-xog.exe` file.

# Data migration 
This feature is used to export instances to an excel file and read data from excel file to a XOG template creating an xml to import data to the environment.

## Export data to excel
Should be used with a [driver instance type](#description-of-instance-driver-types) to read data from the environment and save the match attributes to an excel file.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Defines the name that will be displayed to the user. | yes | 
| `path` | Path where the file will be saved on the file system. | yes |
| `exportToExcel` | If set to true creates an excel file with the matched data. | yes |
| `excel` | The name of the file to export the data.  | yes |
| `instance` | The name of the main tag that represents the instance object that is being read.  | yes |

### Sub tag `match`
This tag is required for export to excel data.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `attribute` | Defines the attribute in the element where you want to get the data from. If no xpath is defined then we get the value from the main instance element defined. | no | 
| `xpath` | A string representing the path to the element you want to get the data from. If no attribute value is defined then we get the value from the tag text. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <resourceInstance code="*" path="res.xml" exportToExcel="true" excel="res.xlsx" instance="Resource">
        <match attribute="resourceId" />
        <match xpath="//PersonalInformation" attribute="displayName" />
        <match xpath="//PersonalInformation" attribute="emailAddress" />
        <match xpath="//PersonalInformation" attribute="firstName" />
        <match xpath="//PersonalInformation" attribute="lastName" />
        <match xpath="//OBSAssoc[@id='corpLocationOBS']" attribute="unitPath" />
        <match xpath="//OBSAssoc[@id='resourcePool']" attribute="unitPath" />
        <match xpath="//ColumnValue[@name='partition_code']" />
    </resourceInstance>
</xogdriver>
```

## Read data from excel to create XOG instances xml 
Should be used with to create an XOG xml file with an instance for each line in the excel file.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `path` | Path where the file will be saved on the file system. | yes |
| `template` | Path to the template that should be used to create the XOG xml file. | yes |
| `instance` | The name of the main tag that represents the instance object that should be created. | yes |
| `excel` | Path to the excel file with the data. | yes |
| `startRow` | The line number in the excel file that we will start reading to create the instances. Default value is 1. | no |
| `instancesPerFile` | Defines the amout of instances in each write xog file. If not defined only one file should be created with all instances. | no |

### Sub tag `match`
This tag is required for export to excel data.

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `col` | Defines from which column of excel we'll get the data to include in the XOG xml file. | yes |
| `attribute` | Defines which attribute in the element will receive the data. If no xpath is defined then we set this attribute in the main element instance. | no |
| `xpath` | A string representing the path to the element you want to set the data. If no attribute value is defined then we set the value as a tag text. | no |
| `removeIfNull` | If set to true and the value in excel is null, the element associated with xpath is removed. | no |
| `multiValued` | If set to true this element will be treated as multi-valued. | no |
| `separator` | Defines what character is being used to separate the options in the multi-valued data. Default value is ';'. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="2.0">
    <migration path="subs.xml" template="template.xml" instance="instance" excel="dados.xlsx" startRow="2" >
        <match col="1" attribute="instanceCode" />
        <match col="1" xpath="//ColumnValue[@name='code']" />
        <match col="2" xpath="//ColumnValue[@name='name']" />
        <match col="3" xpath="//ColumnValue[@name='status']" />
        <match col="4" xpath="//ColumnValue[@name='multivalue_status']" multiValued="true" separator=";" />
        <match col="5" xpath="//ColumnValue[@name='analista']" />
    </migration>
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
This is an example of configuring the environments file.

If the user does not have XOG access to an environment, simply remove the username and password tags. In this case, the system will prompt you to login at run time. Allowing someone else to enter with the necessary credentials.

This information is stored in memory and will be requested again only if the system is restarted.

If the URL has a non-default port (80/443) it should be informed as follows: `http://development.server.com:8888`

| Attribute | Description | Required |
| ------ | ------ | ------ |
| `name` | Defines an unique identifier that will be displayed in the application for choice of actions. | yes |
| `username` | Username with permission to execute XOG in the environment. | no |
| `password` | Password associated with username. | no |
| `endpoint` | Defines the environment's URL. | yes |

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
