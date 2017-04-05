# cas-xog
Execute XOG files reading and wrinting in a more easy way

This is a way of runnig XOG files in a much more easy way. Using a XOGDrive.xml file you can define with objects you would like to read and write.

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
| viewEnvTarget | Defines the target environment to get the destination view attributes. Required to use in conjunction with 'singleView' |

### Description of attributes for tag include:

*Obs.: Available only for types: "objects", "views" and "menus"*

| TAG | Description |
| ------ | ------ |
| code | Code of the attribute, link or action for the object |
| type | Defines what are being readed (object: attribute, link, action) (menu: menuSection, menuLink) |
| sectionCode | Sets the code of the menu section. Required when type is menuLink |
| linkPosition | Used to change the position of the link in the target. Available only if type is menuLink or menuSection |
| sectionPosition | Used to change the position of the section in the target. Available only if type is menuLink or menuSection |
| insertAfter | Used to define the position os the attribute in the view. Available if include type is attribute and file type is views |
| insertBefore | Used to define the position os the attribute in the view. Available if include type is attribute and file type is views |

| Description of Errors |
| ---------- |
| ERROR-0: Trying to validate a write file that does not exist |
| ERROR-1: Output file does not have the XOGOutput Status tag |
| ERROR-2: Readed single view attributes from one target environment and trying to write to another target environment |


### XOG Driver example
```sh
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <file code="idea" path="idea.xml" type="objects">
        <include type="attribute" code="tst_def_compat" />
        <include type="action" code="tst_nova_acao_id" />
        <include type="link" code="tst_novo_link_id" />
    </file>
    <file code="ideaProperty" path="idea_edit.xml" type="views" objectCode="idea" sourcePartition="partition10" singleView="true" viewEnvTarget="1">
        <include type="attribute" code="tst_mtz_compat" insertBefore="name" />
        <include type="attribute" code="tst_is_compat" insertAfter="name" />
    </file>
    <file code="application" path="cas_menu_app.xml" type="menus">
        <include type="menuSection" code="cas_tools" sectionPosition="1000" />
        <include type="menuLink" code="cas_processes_tabs_page" sectionCode ="cas_tools" linkPosition="1000" sectionPosition="1000" />
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
* Remove attribute from view using <remove type="attribute" code="tst_code">
* Change view sections order


License
----

MIT
