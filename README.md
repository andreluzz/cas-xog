# cas-xog
Execute XOG files reading and wrinting in a more easy way

This is a new easy way of creating reading and writing XOG files. Using a XOGDrive.xml file you can define with objects you would like to read and write.


### Description of attributes for tag file:

| TAG | Description |
| ------ | ------ |
| code | Code of the object or instances in case of reading instances |
| path | Path to the file to be written |
| type | Type of file being read. Available: objects, views, portlets, pages, processes, lookups, groups, menus |
| objectCode | Field to set the object code for views, customobjectinstances |
| ignoreReading | Sets whether to ignore the read action for this file. The file must be created manually in the '_extra/type' so it can be writen |
| sourcePartition | Used to define which partition we will read |
| targetPartition | Replace the partition code in the writing file, it is mandatory to use the sourcePartition tag |
| singleView | Remove all other views leaving only the one that has the same code that was filled |
| copyToView | Defines the code where the view will be cloned to. Required to use in conjunction with 'singleView' |
| envTarget | Defines the target environment to get the destination information. Only available to types: 'views' and 'menus' |


### Description of attributes for tag include:

*Obs.: Available only for type "objects"*

| TAG | Description |
| ------ | ------ |
| code | Code of the attribute, link or action for the object |
| type | Defines what are being readed (view: attribute) (object: attribute, link, action) (menu: menuSection, menuLink) |
| sectionCode | Sets the code of the menu section. Required when type is menuLink |
| linkPosition | Used to change the position of the link in the target. Available only if type is menuLink or menuSection |
| targetSectionPosition | Used to change the position of the section in the target. Available only for file type menus |


### Description of attributes for tag menu:

*Obs.: Available only for type "menus"*

| TAG | Description |
| ------ | ------ |
| code | Code of the menu section |
| action | Defines the actions that should be executed: insert, update or replace |
| targetPosition | Defines the position this section should be inserted in the menu |


### Description of attributes for tag attribute:

*Obs.: Available only for types: "menus" inside tag menu*

| TAG | Description |
| ------ | ------ |
| code | Code of the link |
| targetPosition | Defines the position this link should be inserted in the menu section |


### Description of attributes for tag section:

*Obs.: Available only for types: "views"*

| TAG | Description |
| ------ | ------ |
| sourceSectionPosition | Defines the position of the view in the source environment. Optional for remove action |
| targetSectionPosition | Defines the position of the view in the target environment |
| action | Defines the actions that should be executed: replace, update, remove or insert. <br/>OBS.: If executing more than one action to the same view organize the sections in the following order: replace, update, remove and insert |


### Description of attributes for tag attribute:

*Obs.: Available only for types: "views" inside tag section*

| TAG | Description |
| ------ | ------ |
| code | Code of the attribute |
| column | Defines if the attribute should be placed in the left or right column. This attribute is required. |
| insertBefore | Defines the code of the attribute in the target to use as reference for positioning. If not set, the attribute will be inserted at the end of the column |
| remove | Optional attribute, if set as true remove this atribute from target environment. Only available for action update | 


### Description of Errors:

| Code | Description |
| ---- | ---------- |
| <code>ERRO-00</code> | Trying to validate a write file that does not exist |
| <code>ERRO-01</code> | Output file does not have the XOGOutput Status tag |
| <code>ERRO-02</code> | Trying to write view attributes readed from a different target environment |
| <code>ERRO-03</code> | Readed single view attributes from one target environment and trying to write to another target environment |
| <code>ERRO-04</code> | Transform - trying read source file that does not exists |
| <code>ERRO-05</code> | Transform - trying read target file that does not exists |
| <code>ERRO-06</code> | Transform views - section - invalid action at section tag |
| <code>ERRO-07</code> | Transform views - section - invalid TargetSectionPosition |
| <code>ERRO-08</code> | Transform views - section - invalid SourceSectionPosition |
| <code>ERRO-09</code> | Transform views - section - action update without attributes |
| <code>ERRO-10</code> | Transform views - section - column value invalid, only right or left are available |
| <code>ERRO-11</code> | Transform views - section - insertBefore code does not exists in target |
| <code>ERRO-12</code> | Transform views - action - group code does not exist in source environment view |
| <code>ERRO-13</code> | Transform views - action - group code does not exist in target environment view |
| <code>ERRO-14</code> | Transform views - action - cannot remove action because there is no match code in target environment |
| <code>ERRO-15</code> | Transform menus - invalid action at menu tag |
| <code>ERRO-16</code> | Transform menus - cannot update a section that does not exist in target |
| <code>ERRO-17</code> | Transform menus - lacking link tags to update menu |
| <code>ERRO-18</code> | Transform menus - cannot replace a section that does not exist in target |


### XOG Driver example
```sh
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="idea" path="idea.xml" type="objects">
        <include type="attribute" code="tst_def_compat" />
        <include type="action" code="tst_nova_acao_id" />
        <include type="link" code="tst_novo_link_id" />
    </file>
    <file code="ideaProperty" path="idea_edit.xml" type="views" objectCode="idea" sourcePartition="partition10" singleView="true" envTarget="1">
        <section sourceSectionPosition="4" targetSectionPosition="4" action="replace" />
        <section sourceSectionPosition="2" targetSectionPosition="2" action="update" >
            <attribute code="tst_aval_compat" insertBefore="tst_is_compat" column="left" />
            <attribute code="tst_is_compat" column="left" />
            <attribute code="tst_corp_compat" remove="true" />
        </section>
        <section targetSectionPosition="3" action="remove" />
        <section sourceSectionPosition="1" targetSectionPosition="1" action="insert" />
        <action code="tst_nova_acao_id" groupCode="general" insertBefore="odf_copy_srctst_mtz_compat" />
        <action code="odf_copy_srctst_mtz_compat" groupCode="general" remove="true" />
    </file>
    <file code="application" path="cas_menu_app.xml" type="menus" envTarget="1">
        <menu code="cas_tools" action="insert" targetPosition="2"/>
        <menu code="npt.projmgr" action="replace" />
        <menu code="itl.incidentManager" action="update">
            <link code="tst_roadmap_demandas" targetPosition="5" />
        </menu>
    </file>
    <file code="tst_mtz_compat" path="tst_mtz_compat.xml" type="objects" />
    <file code="*" path="tst_mtz_compat_all.xml" type="views" objectCode="tst_mtz_compat" />
    <file code="tst_mtz_compatCreate" singleView="true" path="mtz_compat_create.xml" type="views" objectCode="tst_mtz_compat" sourcePartition="partition10" targetPartition="partition20" copyToView="tst_mtz_compatProperty" />
    <file code="tst_proc_v1" path="tst_proc_v1.xml" type="processes" />
    <file code="cas_running_processes_detail" path="cas_running_processes_detail.xml" type="portlets" />
    <file code="application" path="cas_menu_app.xml" type="menu" ignoreReading="true" />
</xogdriver>
```

### XOG Environment example
```sh
<?xml version="1.0" encoding="utf-8"?>
<xogenvs version="1.0">
    <global>
        <var name="xog_path" value="path_to_xog_bat"/>
    </global>
    <environments>
        <env name="Development">
            <param name="username" value="username"/>
            <param name="password" value="12345"/>
            <param name="servername" value="development.server.com"/>
            <param name="portnumber" value="80"/>
            <param name="sslenabled" value="false"/>
        </env>
        <env name="Quality">
            <param name="username" value="username"/>
            <param name="password" value="12345"/>
            <param name="servername" value="quality.server.com"/>
            <param name="portnumber" value="80"/>
            <param name="sslenabled" value="false"/>
        </env>
        <env name="Production">
            <param name="username" value="username"/>
            <param name="password" value="12345"/>
            <param name="servername" value="production.server.com"/>
            <param name="portnumber" value="80"/>
            <param name="sslenabled" value="false"/>
        </env>
    </environments>
</xogenvs>
```


TODO
----
* Delete views using complete="true" in view propertySet

License
----
MIT
