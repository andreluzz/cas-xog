# run-xog
Execute XOG files reading and wrinting in a more easy way

This is a way of runnig XOG files in a much more easy way. Using a XOGDrive.xml file you can define with objects you would like to read and right.

### Description of attributes:
| TAG | Description |
| ------ | ------ |
| code | Code of the object or instances in case of reading instances |
| path | Path to the file to be written |
| type | Type of file being read. Available: objects, views, portlets, pages, processes, lookups |
| objectCode | Field to set the object code for views, customobjectinstances |
| ignoreReading | Sets whether to ignore the read action for this file. The file must be created manually in the '_extra / type' so it can be writen |
| sourcePartition | Used to define which partition we will read |
| targetPartition | Replace the partition code in the writing file, it is mandatory to use the sourcePartition tag |
| singleView | Remove the views leaving only the one that has the same code that was filled |

### XOG Driver example
```sh
<?xml version="1.0" encoding="utf-8"?>
<xogdriver version="1.0">
    <files>
        <file code="tst_mtz_compat" path="tst_mtz_compat.xml" type="objects" />
        <file code="*" path="tst_mtz_compat.xml" type="views" objectCode="tst_mtz_compat" />
        <file code="tst_proc_v1" path="tst_proc_v1.xml" type="processes" />
        <file code="cas_running_processes_detail" path="cas_running_processes_detail.xml" type="portlets" />
        <file code="application" path="cas_menu_app.xml" type="menu" ignoreReading="true" />
    </files>
</xogdriver>
```

### Todos

 - Copy one view to another
 - Define what attributes from an object should be readed

License
----

MIT
