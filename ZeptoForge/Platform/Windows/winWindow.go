package Windows

import (
	"fmt"
	"unsafe"

	"github.com/MeiKakuTenShi/zeptoforge/ZeptoForge/event/appEvent"
	"github.com/MeiKakuTenShi/zeptoforge/ZeptoForge/event/keyEvent"
	"github.com/MeiKakuTenShi/zeptoforge/ZeptoForge/event/mouseEvent"

	"github.com/MeiKakuTenShi/zeptoforge/ZeptoForge/logsys"
	"github.com/MeiKakuTenShi/zeptoforge/ZeptoForge/window"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type winData struct {
	title         string
	width, height int
	vsync         bool
	callback      window.EventCallBackFn
}

type WinWindow struct {
	win    window.Window
	window *glfw.Window
	data   *winData
}

var (
	glfwInitialized = false
	default_title   = "ZeptoForge Application"
	default_width   = 1024
	default_height  = 720

	// glfwErroCallback =
)

func Create(props *window.WindowProps) WinWindow {
	result := WinWindow{data: &winData{}}

	if props.Title == "" {
		props.Title = default_title
	} else {

	}
	if props.Width == 0 {
		props.Width = default_width
	} else {

	}
	if props.Height == 0 {
		props.Height = default_height
	} else {

	}

	result.init(props)

	return result
}

func (win *WinWindow) Destruct() {
	win.Shutdown()
}
func (win *WinWindow) Shutdown() {
	win.window.Destroy()
}
func (win *WinWindow) OnUpdate() {
	glfw.PollEvents()
	win.window.SwapBuffers()
}
func (win *WinWindow) GetWidth() int {
	return win.data.width
}
func (win *WinWindow) GetHeight() int {
	return win.data.height
}
func (win *WinWindow) SetEventCallback(callback window.EventCallBackFn) {
	win.data.callback = callback
}
func (win *WinWindow) SetVSync(enabled bool) {
	if enabled {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
	win.data.vsync = enabled
}
func (win *WinWindow) IsVSync() bool {
	return win.data.vsync
}

func (win *WinWindow) init(props *window.WindowProps) {
	win.data.title = props.Title
	win.data.width = props.Width
	win.data.height = props.Height

	logsys.ZF_CORE_INFO(fmt.Sprintf("Creating window %s (%v, %v)", props.Title, props.Width, props.Height))

	if !glfwInitialized {
		if err := glfw.Init(); err != nil {
			logsys.ZF_CORE_ERROR(err)
		}
		glfwInitialized = true
		// glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	var err error
	win.window, err = glfw.CreateWindow(props.Width, props.Height, props.Title, nil, nil)
	if err != nil {
		logsys.ZF_CORE_ERROR(err)
	}
	win.window.MakeContextCurrent()
	win.window.SetUserPointer(unsafe.Pointer(win.data))
	win.SetVSync(true)

	// Set GLFW callbacks

	// Window
	win.window.SetSizeCallback(func(w *glfw.Window, width, height int) {
		data := *(*winData)(w.GetUserPointer())
		data.width = width
		data.height = height

		e := appEvent.NewWindowResizeEvent(width, height)
		data.callback.CallbackFn(e)
	})

	win.window.SetCloseCallback(func(w *glfw.Window) {
		data := *(*winData)(w.GetUserPointer())

		e := appEvent.NewWindowCloseEvent()
		data.callback.CallbackFn(e)
	})

	// Keys
	win.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		data := *(*winData)(w.GetUserPointer())

		switch action {
		case glfw.Press:
			{
				data.callback.CallbackFn(keyEvent.NewKeyPressedEvent(int(key), 0))
				break
			}
		case glfw.Release:
			{
				data.callback.CallbackFn(keyEvent.NewKeyReleasedEvent(int(key)))
				break
			}
		case glfw.Repeat:
			{
				data.callback.CallbackFn(keyEvent.NewKeyPressedEvent(int(key), 1))
				break
			}
		}
	})

	// Mouse
	win.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		data := *(*winData)(w.GetUserPointer())

		switch action {
		case glfw.Press:
			{
				data.callback.CallbackFn(mouseEvent.NewMouseButtonPressedEvent(int(button)))
				break
			}
		case glfw.Release:
			{
				data.callback.CallbackFn(mouseEvent.NewMouseButtonReleasedEvent(int(button)))
				break
			}
		}
	})

	win.window.SetScrollCallback(func(w *glfw.Window, xOff, yOff float64) {
		data := *(*winData)(w.GetUserPointer())
		data.callback.CallbackFn(mouseEvent.NewMouseScrolledEvent(xOff, yOff))
	})

	win.window.SetCursorPosCallback(func(w *glfw.Window, xPos, yPos float64) {
		data := *(*winData)(w.GetUserPointer())
		data.callback.CallbackFn(mouseEvent.NewMouseMovedEvent(xPos, yPos))
	})
}