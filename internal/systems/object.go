package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
)

func ObjectSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsObject) {
		if obj, ok := result.Components[myecs.Object].(*object.Object); ok {
			if obj.Killed {
				myecs.Manager.DisposeEntity(result)
			} else {
				obj.Update()
			}
		}
	}
}

func ParentSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasParent) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		parent, okP := result.Components[myecs.Parent].(*object.Object)
		if okO && okP {
			if parent.Killed {
				myecs.Manager.DisposeEntity(result)
			} else {
				obj.Pos = parent.Pos.Add(parent.Offset)
				obj.Mask = parent.Mask
				if parent.HideChildren {
					obj.Hidden = parent.Hidden
				}
			}
		}
	}
}

func FunctionSystem() {
	for _, result := range myecs.Manager.Query(myecs.HasUpdate) {
		fnA := result.Components[myecs.Update]
		if fnT, ok := fnA.(*data.TimerFunc); ok {
			if fnT.Timer.UpdateDone() {
				if fnT.Func() {
					result.Entity.RemoveComponent(myecs.Update)
				} else {
					fnT.Timer.Reset()
				}
			}
		} else if hcF, ok := fnA.(*data.HoverClick); ok {
			pos := hcF.Input.World
			if hcF.View != nil {
				hcF.Pos = hcF.View.ProjectWorld(pos)
				hcF.ViewHover = hcF.View.PointInside(hcF.Pos)
			} else {
				hcF.Pos = pos
				hcF.ViewHover = true
			}
			hcF.Hover = false
			if objC, okOC := result.Entity.GetComponentData(myecs.Object); okOC {
				if obj, okO := objC.(*object.Object); okO {
					if !obj.Hidden {
						if hcF.View != nil {
							hcF.Hover = obj.PointInside(hcF.Pos) && hcF.ViewHover
						} else {
							hcF.Hover = obj.PointInside(hcF.Pos)
						}
					}
				}
			}
			if hcF.Func != nil {
				hcF.Func(hcF)
			}
		} else if fnF, ok := fnA.(*data.FrameFunc); ok {
			if fnF.Func() {
				result.Entity.RemoveComponent(myecs.Update)
			}
		} else if fnU, ok := fnA.(*data.Funky); ok {
			fnU.Fn()
		}
	}
}

func TemporarySystem() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		temp := result.Components[myecs.Temp]
		del := false
		if timer, ok := temp.(*timing.Timer); ok {
			if timer.UpdateDone() {
				del = true
			}
		} else if check, ok := temp.(myecs.ClearFlag); ok {
			if check {
				del = true
			}
		}
		if del {
			if objC, ok1 := result.Entity.GetComponentData(myecs.Object); ok1 {
				if obj, ok2 := objC.(*object.Object); ok2 {
					obj.Hidden = true
					obj.Killed = true
				}
			}
			myecs.Manager.DisposeEntity(result.Entity)
		}
	}
}

func ClearTemp() {
	for _, result := range myecs.Manager.Query(myecs.IsTemp) {
		if objC, ok1 := result.Entity.GetComponentData(myecs.Object); ok1 {
			if obj, ok2 := objC.(*object.Object); ok2 {
				obj.Hidden = true
				obj.Killed = true
			}
		}
		myecs.Manager.DisposeEntity(result.Entity)
	}
}

func ClearSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsObject) {
		obj, ok := result.Components[myecs.Object].(*object.Object)
		if ok {
			obj.Hidden = true
			obj.Killed = true
		}
		myecs.Manager.DisposeEntity(result.Entity)
	}
}
