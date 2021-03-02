package ds1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

func (p *widget) createFloorShadowRecord() d2ds1.FloorShadowRecord {
	state := p.getState()

	newFloorShadowRecord := d2ds1.FloorShadowRecord{
		Prop1:       byte(state.addFloorShadowState.prop1),
		Sequence:    byte(state.addFloorShadowState.sequence),
		Unknown1:    byte(state.addFloorShadowState.unknown1),
		Style:       byte(state.addFloorShadowState.style),
		Unknown2:    byte(state.addFloorShadowState.unknown2),
		HiddenBytes: state.addFloorShadowState.hidden,
	}

	return newFloorShadowRecord
}

func (p *widget) createWallRecord() d2ds1.WallRecord {
	state := p.getState()

	newWall := d2ds1.WallRecord{
		Type:        d2enum.TileType(state.addWallState.tileType),
		Zero:        byte(state.addWallState.zero),
		Prop1:       byte(state.addWallState.prop1),
		Sequence:    byte(state.addWallState.sequence),
		Unknown1:    byte(state.addWallState.unknown1),
		Style:       byte(state.addWallState.style),
		Unknown2:    byte(state.addWallState.unknown2),
		HiddenBytes: state.addWallState.hidden,
	}

	return newWall
}

func (p *widget) addFloor() {
	state := p.getState()

	state.addFloorShadowState.cb = func() {
		newFloor := p.createFloorShadowRecord()

		for y := range p.ds1.Tiles {
			for x := range p.ds1.Tiles[y] {
				p.ds1.Tiles[y][x].Floors = append(p.ds1.Tiles[y][x].Floors, newFloor)
			}
		}

		p.ds1.NumberOfFloors++

		p.recreateLayerStreamTypes()
	}

	state.mode = ds1EditorModeAddFloorShadow
}

func (p *widget) deleteFloorRecord() {
	state := p.getState()

	for y := range p.ds1.Tiles {
		for x := range p.ds1.Tiles[y] {
			newFloors := make([]d2ds1.FloorShadowRecord, 0)

			for n, floor := range p.ds1.Tiles[y][x].Floors {
				if n != int(state.object) {
					newFloors = append(newFloors, floor)
				}
			}

			p.ds1.Tiles[y][x].Floors = newFloors
		}
	}

	p.ds1.NumberOfFloors--
	p.recreateLayerStreamTypes()
}

func (p *widget) addWall() {
	state := p.getState()

	state.addWallState.cb = func() {
		newWall := p.createWallRecord()

		for y := range p.ds1.Tiles {
			for x := range p.ds1.Tiles[y] {
				p.ds1.Tiles[y][x].Walls = append(p.ds1.Tiles[y][x].Walls, newWall)
			}
		}

		p.ds1.NumberOfWalls++

		p.recreateLayerStreamTypes()
	}

	state.mode = ds1EditorModeAddWall
}

func (p *widget) deleteWall() {
	state := p.getState()

	for y := range p.ds1.Tiles {
		for x := range p.ds1.Tiles[y] {
			newWalls := make([]d2ds1.WallRecord, 0)

			for n, wall := range p.ds1.Tiles[y][x].Walls {
				if n != int(state.object) {
					newWalls = append(newWalls, wall)
				}
			}

			p.ds1.Tiles[y][x].Walls = newWalls
		}
	}

	p.ds1.NumberOfWalls--
	p.recreateLayerStreamTypes()
}

func (p *widget) addShadow() {
	state := p.getState()

	state.addFloorShadowState.cb = func() {
		newShadow := p.createFloorShadowRecord()

		for y := range p.ds1.Tiles {
			for x := range p.ds1.Tiles[y] {
				p.ds1.Tiles[y][x].Shadows = make([]d2ds1.FloorShadowRecord, 1)
				p.ds1.Tiles[y][x].Shadows[0] = newShadow
			}
		}
	}

	p.ds1.NumberOfShadowLayers++
	p.recreateLayerStreamTypes()

	state.mode = ds1EditorModeAddFloorShadow
}

func (p *widget) deleteShadow() {
	state := p.getState()

	yesCB := func() {
		for y := range p.ds1.Tiles {
			for x := range p.ds1.Tiles[y] {
				p.ds1.Tiles[y][x].Shadows = nil
			}
		}

		p.ds1.NumberOfShadowLayers--
		p.recreateLayerStreamTypes()

		state.mode = ds1EditorModeViewer
	}

	state.confirmDialog = hswidget.NewPopUpConfirmDialog(
		"##"+p.id+"removeShadowConfirm",
		"Warning",
		"non-shadow files aren't supported.\n"+
			"If you'll delete shadow, and will not create\n"+
			"a new one, the file will be destroyed and\n"+
			"You will be unable to open it again.\n"+
			"Continue?",
		yesCB,
		func() {
			state.mode = ds1EditorModeViewer
		},
	)
	state.mode = ds1EditorModeConfirm
}
