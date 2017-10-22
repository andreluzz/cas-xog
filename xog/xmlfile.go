package xog

import (
	"errors"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
	"strings"
)

var docXogReadXML, soapEnvelope *etree.Document

func InitRead() {
	docXogReadXML = etree.NewDocument()
	docXogReadXML.ReadFromFile("xogRead.xml")
	soapEnvelope = etree.NewDocument()
	soapEnvelopeElement := docXogReadXML.FindElement("//xogtype[@type='envelope']/soapenv:Envelope")
	soapEnvelope.SetRoot(soapEnvelopeElement.Copy())
}

func GetLoginXML() *etree.Element {
	return docXogReadXML.FindElement("//xogtype[@type='login']/soapenv:Envelope").Copy()
}

func GetLogoutXML() *etree.Element {
	return docXogReadXML.FindElement("//xogtype[@type='logout']/soapenv:Envelope").Copy()
}

func getStandardXogReadXML(requestType string) (*etree.Element, error) {
	nikuDataBusElement := docXogReadXML.FindElement("//xogtype[@type='" + requestType + "']/NikuDataBus")
	if nikuDataBusElement == nil {
		return nil, errors.New("invalid object type")
	}
	envelope := soapEnvelope.Root().Copy()
	envelope.FindElement("//soapenv:Body").AddChild(nikuDataBusElement.Copy())
	return envelope, nil
}

func read(file *common.DriverFile, env *EnvType) (string, error) {
	readXML, err := getStandardXogReadXML(file.Type)

	if err != nil {
		return "", err
	}

	if file.Code == "" {
		return "", errors.New("no attribute code defined on tag <file>")
	}

	if file.Path == "" {
		return "", errors.New("no attribute path defined on tag <file>")
	}

	req := etree.NewDocument()
	req.SetRoot(readXML)
	req.FindElement("//xog:SessionID").SetText(env.Session)

	switch file.Type {
	case common.LOOKUP:
		req.FindElement("//Filter[@name='code']").SetText(strings.ToUpper(file.Code))
	case common.PORTLET, common.QUERY, common.PROCESS, common.PAGE, common.MENU, common.OBS:
		req.FindElement("//Filter[@name='code']").SetText(file.Code)
	case common.OBJECT:
		req.FindElement("//Filter[@name='object_code']").SetText(file.Code)
	case common.VIEW:
		if file.ObjCode == "" {
			return "", errors.New("no attribute objectCode defined on tag <file>")
		}
		req.FindElement("//Filter[@name='code']").SetText(file.Code)
		req.FindElement("//Filter[@name='object_code']").SetText(file.ObjCode)
		if file.SourcePartition == "" {
			filter := req.FindElement("//Filter[@name='partition_code']")
			filter.Parent().RemoveChild(filter)
		} else {
			req.FindElement("//Filter[@name='partition_code']").SetText(file.SourcePartition)
		}
	case common.CUSTOM_OBJECT_INSTANCE:
		if file.ObjCode == "" {
			return "", errors.New("no attribute objectCode defined on tag <file>")
		}
		req.FindElement("//Filter[@name='instanceCode']").SetText(file.Code)
		req.FindElement("//Filter[@name='objectCode']").SetText(file.ObjCode)
		req.FindElement("//args[@name='documentLocation']").CreateAttr("value", "./" + common.FOLDER_WRITE + "_" + file.Type + "/_document")
	case common.RESOURCE_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='resource_class']").SetText(file.Code)
	case common.WIP_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='wipclass']").SetText(file.Code)
	case common.INVESTMENT_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='investmentclass']").SetText(file.Code)
	case common.TRANSACTION_CLASS_INSTANCE:
		req.FindElement("//Filter[@name='transclass']").SetText(file.Code)
	case common.RESOURCE_INSTANCE:
		req.FindElement("//Filter[@name='resourceID']").SetText(file.Code)
	case common.USER_INSTANCE:
		req.FindElement("//Filter[@name='userName']").SetText(file.Code)
	case common.PROJECT_INSTANCE:
		req.FindElement("//Filter[@name='projectID']").SetText(file.Code)
		req.FindElement("//args[@name='documentLocation']").CreateAttr("value", "./" + common.FOLDER_WRITE + "_" + file.Type + "/_document")
	case common.IDEA_INSTANCE, common.APPLICATION_INSTANCE, common.ASSET_INSTANCE, common.OTHER_INVESTMENT_INSTANCE, common.PRODUCT_INSTANCE, common.SERVICE_INSTANCE:
		req.FindElement("//Filter[@name='objectID']").SetText(file.Code)
		req.FindElement("//args[@name='documentLocation']").CreateAttr("value", "./" + common.FOLDER_WRITE + "_" + file.Type + "/_document")
	}

	nikuDataBusElement := req.FindElement("//NikuDataBus").Copy()
	if nikuDataBusElement != nil {
		readRequest := etree.NewDocument()
		readRequest.SetRoot(nikuDataBusElement)
		readRequest.IndentTabs()
		readRequest.WriteToFile(common.FOLDER_READ + file.Type + "/" + file.Path)
	} else {
		req.IndentTabs()
		req.WriteToFile(common.FOLDER_READ + file.Type + "/" + file.Path)
	}

	return req.WriteToString()
}

func write(file *common.DriverFile, env *EnvType) (string, error) {
	nikuDataBusXML := etree.NewDocument()
	folder := common.FOLDER_WRITE
	if file.Type == common.MIGRATION {
		folder = common.FOLDER_MIGRATION
	}

	err := nikuDataBusXML.ReadFromFile(folder + file.Type + "/" + file.Path)
	if err != nil {
		return "", err
	}

	req := etree.NewDocument()
	req.SetRoot(soapEnvelope.Root().Copy())

	req.FindElement("//soapenv:Body").AddChild(nikuDataBusXML.Root())
	req.FindElement("//xog:SessionID").SetText(env.Session)

	return req.WriteToString()
}

func GetXMLFile(action string, file *common.DriverFile, env *EnvType) (string, error) {
	if action == "r" {
		return read(file, env)
	} else {
		return write(file, env)
	}
}
