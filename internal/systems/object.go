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
			if objC, okOC := result.Entity.GetComponentData(myecs.Object); okOC {
				if obj, okO := objC.(*object.Object); okO {
					if !obj.Hidden {
						pos := hcF.Input.World
						if hcF.View != nil {
							pos = hcF.View.ProjectWorld(pos)
							hcF.Hover = obj.PointInside(pos) && hcF.View.PointInside(pos)
						} else {
							hcF.Hover = obj.PointInside(pos)
						}
						if hcF.Func != nil {
							hcF.Func(hcF)
						}
					}
				}
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
		obj, okT := result.Components[myecs.Object].(*object.Object)
		if okT {
			if timer, ok := temp.(*timing.Timer); ok {
				if timer.UpdateDone() {
					obj.Hidden = true
					obj.Killed = true
					myecs.Manager.DisposeEntity(result.Entity)
				}
			} else if check, ok := temp.(myecs.ClearFlag); ok {
				if check {
					obj.Hidden = true
					obj.Killed = true
					myecs.Manager.DisposeEntity(result.Entity)
				}
			}
		}
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
