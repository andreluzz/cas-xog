package transform

import (
	"testing"
	"github.com/andreluzz/cas-xog/common"
	"github.com/beevik/etree"
)

func TestExecuteToReturnMenu(t *testing.T) {
	file := common.DriverFile{
		Code: "application",
		Type: common.MENU,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")

	err := Execute(xog, nil, file)

	if err != nil {
		t.Fatalf("Error transforming OBS XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "menu_result.xml") == false {
		t.Errorf("Error transforming OBS XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnMenuSectionAndLinks(t *testing.T) {
	file := common.DriverFile{
		Code: "application",
		Type: common.MENU,
		Sections: []common.Section {
			{
				Action: common.ACTION_INSERT,
				Code: "menu_sec_cas_xog",
				TargetPosition: "2",
			},
			{
				Action: common.ACTION_UPDATE,
				Code: "npt.personal",
				Links: []common.SectionLink{
					{
						Code: "odf.obj_testeList",
					},
				},
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog.xml")

	err := Execute(xog, aux, file)

	if err != nil {
		t.Fatalf("Error transforming OBS XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "menu_section_link_result.xml") == false {
		t.Errorf("Error transforming OBS XOG file. Invalid result XML.")
	}
}
