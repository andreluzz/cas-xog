package transform

import (
	"github.com/andreluzz/cas-xog/constant"
	"github.com/andreluzz/cas-xog/model"
	"github.com/beevik/etree"
	"testing"
)

func TestExecuteToReturnErrorMenuInvalidSourceSection(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Code: "invalid_code",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog_no_section_code.xml")

	err := Execute(xog, aux, &file)

	if err == nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: not validating if section code exists on source")
	}
}

func TestExecuteToReturnErrorMenuInvalidTargetSection(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action: constant.ActionUpdate,
				Code:   "npt.personal",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog_no_section_code.xml")

	err := Execute(xog, aux, &file)

	if err == nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: not validating if section code exists on target")
	}
}

func TestExecuteToReturnErrorMenuUpdateWithoutLinks(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action: constant.ActionUpdate,
				Code:   "npt.personal",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog.xml")

	err := Execute(xog, aux, &file)

	if err == nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: not validating if section code exists on target")
	}
}

func TestExecuteToReturnErrorMenuLinkInvalidCode(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action: constant.ActionUpdate,
				Code:   "npt.personal",
				Links: []model.SectionLink{
					{
						Code: "invalid_link_code",
					},
				},
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog.xml")

	err := Execute(xog, aux, &file)

	if err == nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: not validating if source section link exists")
	}
}

func TestExecuteToReturnErrorMenuInsertExistentSection(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action: constant.ActionInsert,
				Code:   "npt.personal",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog.xml")

	err := Execute(xog, aux, &file)

	if err == nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: not validating if source section link exists")
	}
}

func TestExecuteToReturnErrorMenuTargetInvalidSectionPosition(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action:         constant.ActionInsert,
				Code:           "menu_sec_cas_xog",
				TargetPosition: "129",
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog.xml")

	err := Execute(xog, aux, &file)

	if err == nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: not validating insert invalid target section position")
	}
}

func TestExecuteToReturnMenu(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")

	err := Execute(xog, nil, &file)

	if err != nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "menu_result.xml") == false {
		t.Errorf("Error transforming Menu XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnMenuSectionAndLinks(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action:         constant.ActionInsert,
				Code:           "menu_sec_cas_xog",
				TargetPosition: "2",
			},
			{
				Action: constant.ActionUpdate,
				Code:   "npt.personal",
				Links: []model.SectionLink{
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

	err := Execute(xog, aux, &file)

	if err != nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "menu_section_link_result.xml") == false {
		t.Errorf("Error transforming Menu XOG file. Invalid result XML.")
	}
}

func TestExecuteToReturnMenuInsertSectionWithLinks(t *testing.T) {
	file := model.DriverFile{
		Code: "application",
		Type: constant.Menu,
		Sections: []model.Section{
			{
				Action:         constant.ActionInsert,
				Code:           "menu_sec_cas_xog",
				TargetPosition: "2",
				Links: []model.SectionLink{
					{
						Code: "cas_proc_running_tab",
					},
				},
			},
		},
	}

	xog := etree.NewDocument()
	xog.ReadFromFile(packageMockFolder + "menu_full_xog.xml")
	aux := etree.NewDocument()
	aux.ReadFromFile(packageMockFolder + "menu_full_aux_xog.xml")

	err := Execute(xog, aux, &file)

	if err != nil {
		t.Fatalf("Error transforming Menu XOG file. Debug: %s", err.Error())
	}

	if readMockResultAndCompare(xog, "menu_insert_section_link_result.xml") == false {
		t.Errorf("Error transforming Menu XOG file. Invalid result XML.")
	}
}
