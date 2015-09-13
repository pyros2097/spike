package gdx

// // import (
// // 	log "github.com/Sirupsen/logrus"
// // 	"github.com/veandco/go-sdl2/sdl"
// // 	"time"
// // )

// // type (
// // 	stage struct {
// // 		running   bool
// // 		fpsTicker time.Ticker
// // 	}
// // )

// // var (
// // 	info      sdl.SysWMInfo
// // 	subsystem string
// // 	window    *sdl.Window
// // 	surface   *sdl.Surface
// // 	renderer  *sdl.Renderer
// // 	event     sdl.Event
// // 	err       error
// // 	points    []sdl.Point
// // 	rect      sdl.Rect
// // 	rects     []sdl.Rect
// // 	core      *stage
// // )

// // func Initialize(string title, w, h, fps int) {
// // 	log.Info("Initializing SDL")
// // 	sdl.Init(sdl.INIT_EVERYTHING)
// // 	sdl.VERSION(&info.Version)

// // 	core := &stage{
// // 		fpsTicker: time.NewTicker(time.Tick(1000 / fps * time.Millisecond)),
// // 	}

// // 	window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
// // 		w, h, sdl.WINDOW_SHOWN)
// // 	defer window.Destroy()
// // 	throwError("Failed to create window", err)

// // 	if window.GetWMInfo(&info) {
// // 		switch info.Subsystem {
// // 		case sdl.SYSWM_UNKNOWN:
// // 			subsystem = "An unknown system!"
// // 		case sdl.SYSWM_WINDOWS:
// // 			subsystem = "Microsoft Windows(TM)"
// // 		case sdl.SYSWM_X11:
// // 			subsystem = "X Window System"
// // 		case sdl.SYSWM_DIRECTFB:
// // 			subsystem = "DirectFB"
// // 		case sdl.SYSWM_COCOA:
// // 			subsystem = "Apple OS X"
// // 		case sdl.SYSWM_UIKIT:
// // 			subsystem = "UIKit"
// // 		}
// // 		log.Infof("Running SDL version %d.%d.%d on %s", info.Version.Major,
// // 			info.Version.Minor,
// // 			info.Version.Patch,
// // 			subsystem)
// // 	} else {
// // 		throwError("Couldn't get window information", sdl.GetError())
// // 	}
// // }

// // func Stage() *stage {
// // 	return core
// // }

// // func (self *stage) Start() {
// // 	running = true
// // 	for now := range self.fpsTicker.C {
// // 		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
// // 			switch t := event.(type) {
// // 			case *sdl.QuitEvent:
// // 				self.running = false
// // 				self.fpsTicker.Stop()
// // 				log.Info("Destroying Window")
// // 				sdl.Quit()
// // 			case *sdl.MouseMotionEvent:
// // 				log.Infof("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d",
// // 					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
// // 				rect = sdl.Rect{t.X, t.Y, 200, 200}
// // 				renderer.SetDrawColor(255, 0, 0, 255)
// // 				renderer.FillRect(&rect)
// // 			case *sdl.MouseButtonEvent:
// // 				log.Infof("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d",
// // 					t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
// // 			case *sdl.MouseWheelEvent:
// // 				log.Infof("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d",
// // 					t.Timestamp, t.Type, t.Which, t.X, t.Y)
// // 			case *sdl.KeyUpEvent:
// // 				log.Infof("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d",
// // 					t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
// // 			}
// // 		}
// // 	}
// // }

// // func (self *stage) Pause() {

// // }

// // func (self *stage) Resume() {

// // }

// // func (self *stage) Quit() {
// // 	self.running = false
// // 	self.fpsTicker.Stop()
// // 	log.Info("Destroying Window")
// // 	sdl.Quit()
// // }

// // func throwError(message string, err error) {
// // 	if err != nil {
// // 		log.Fatalf(message+": %s", err)
// // 	}
// // }

// // func main() {
// // 	Initialize("Demo", 800, 480, 1)
// // 	Stage().Start()
// // 	// surface, err = window.GetSurface()
// // 	// if err != nil {
// // 	// 	log.Fatalf("Failed to get window surface: %s", err)
// // 	// }

// // 	// rect := sdl.Rect{0, 0, 200, 200}
// // 	// surface.FillRect(&rect, 0xffff0000)
// // 	// window.UpdateSurface()

// // 	// renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
// // 	// throwError("Failed to create renderer", err)
// // 	// defer renderer.Destroy()
// // }
