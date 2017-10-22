# CAS-XOG
Execute XOG files reading and writing in a more easy way

This is a new method of creating XOG files. Using a Driver XML file, you can define with objects you would like to read, write and migrate.


### How to use

1. Download the [lastest stable release](https://github.com/andreluzz/cas-xog/releases/latest) (cas-xog.exe and xogRead.xml);
2. Create a file called [xogEnv.xml](#xog-environment-example), in the same folder of the cas-xog.exe, to define the environments connections configuration;
3. Create a folder called "drivers" with all [xog Driver xml](#xog-driver-example) files defining the objects you want to read and write;
4. Execute the cas-xog.exe and follow the instructions in the screen.

### Description of Driver types

| Type | Description |
| ------ | ------ |
| [`objects`](#type-objects) | Used to read and write objects attributes, actions and links  |
| [`processes`](#type-processes) | Used to read and write processes  |
| [`lookups`](#type-lookups) | Used to read and write lookups  |
| [`portlets`](#type-portlets) | Used to read and write portlets  |
| [`queries`](#type-queries) | Used to read and write queries  |
| [`pages`](#type-pages) | Used to read and write pages  |


# Type `objects`

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Object code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 
| `partitionModel` | Defines a new partitionModel if you want to set or change | no |
| `sourcePartition` | When defined reads only elements from this partition code | no |
| `targetPartition` | Used to change the current partition code. Used alone without sourcePartition replaces the tag partitionCode of all xog elements with the defined value. | no |

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="idea" path="idea.xml" type="objects" />
    <file code="application" path="application.xml" type="objects" partitionModel="new-corp" />
    <file code="systems" path="systems.xml" type="objects" sourcePartition="partition10" targetPartition="NIKU.ROOT" />
    <file code="inv" path="inv.xml" type="objects" targetPartition="NIKU.ROOT" />
</xogdriver>
```

### Sub tag `element`
Used to read only the selected elements from the object

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `type` | Defines what element to read. Available: `attribute`, `action` and `link` | yes | 
| `code` | Code of the attribute that you want to include | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="test_subobj" path="test_subobj.xml" type="objects">
        <element type="attribute" code="attr_auto_number" />
        <element type="attribute" code="novo_atr" />
        <element type="action" code="tst_run_proc_lnk" />
        <element type="link" code="test_subobj.link_test" />
    </file>
</xogdriver>
```

# Type `processes`

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Process code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="PRC_0001" path="PRC_0001.xml" type="processes" />
</xogdriver>
```

# Type `lookups`

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Lookup code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="INV_APPLICATION_CATEGORY_TYPE" path="INV_APPLICATION_CATEGORY_TYPE.xml" type="lookups" />
</xogdriver>
```

# Type `portlets`

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Portlet code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="cop.teamCapacityLinkable" path="cop.teamCapacityLinkable.xml" type="portlets"/>
</xogdriver>
```

# Type `queries`

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Query code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="cop.projectCostsPhaseLinkable" path="cop.projectCostsPhaseLinkable.xml" type="queries"/>
</xogdriver>
```

# Type `pages`

| Atribute | Description | Required |
| ------ | ------ | ------ |
| `code` | Page code | yes | 
| `path` | Path to the file to be saved on the file system | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="pma.ideaFrame" path="pma.ideaFrame.xml" type="pages"/>
</xogdriver>
```

# Global Sub Tags
Sub tags that can be used in any type of `file` tag.

### Sub Tag `replace`
Used to do a replace from one string to another in the xog result.

| Tag | Description | Required |
| ------ | ------ | ------ |
| `from` | Defines which string should be searched for to be changed | yes | 
| `to` | String that will replace what was defined in the `from` tag | yes | 

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
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

# XOG Environment example:

```xml
<?xml version="1.0" encoding="utf-8"?>
<xogenvs version="1.0">
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
        <endpoint>http://production.server.com</endpoint>
    </env>
</xogenvs>
```
