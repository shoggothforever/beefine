package imager

import (
	"bufio"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	imaget "github.com/docker/docker/api/types/image"
	"log"
	"os/exec"
	"shoggothforever/beefine/bpf/image_prep"
	"shoggothforever/beefine/bpf/mount"
	"shoggothforever/beefine/internal/cli"
	"shoggothforever/beefine/internal/helper"
	"shoggothforever/beefine/pkg/component"
	"sync"
	"time"
)

// MyCustomWidget 是自定义控件，包装了 Select 并添加了额外的字段
type ImageSelect struct {
	widget.BaseWidget                // 嵌入 BaseWidget
	base              *widget.Select // 内嵌 Select
	imageLogs         *component.LogBoard
	bpfLogs           *component.LogBoard
	cancelMap         map[string]func()
	images            []imaget.Summary
	currentImageIndex int
	watcherID         int //观测的docker 容器id
	m                 sync.Mutex
}

// NewImageSelect 创建自定义控件实例
func NewImageSelect() *ImageSelect {
	// 初始化 Select
	images, err := cli.ListImage()
	if err != nil {
		return nil
	}
	var imagetags []string
	for _, v := range images {
		for _, tag := range v.RepoTags {
			imagetags = append(imagetags, tag)
			break
		}
	}
	selectWidget := widget.NewSelect(imagetags, nil)
	// 创建 MyCustomWidget 实例
	s := &ImageSelect{
		base:              selectWidget,
		cancelMap:         make(map[string]func()),
		m:                 sync.Mutex{},
		images:            images,
		currentImageIndex: 0,
	}
	s.base.OnChanged = s.OnChanged
	s.base.PlaceHolder = "select existed image"
	s.ExtendBaseWidget(s) // 必须扩展 BaseWidget
	return s
}

// CreateRenderer 实现 fyne.WidgetRenderer，用于渲染控件
func (w *ImageSelect) CreateRenderer() fyne.WidgetRenderer {
	// 将 Select 包装为渲染器的一部分
	return widget.NewSimpleRenderer(w.base)
}

func (w *ImageSelect) chooseVFS(b bool) {
	if b == true {
		w.bpfLogs.AppendLogf("choose watch unionfs")
		req := image_prep.ImagePrepReq{}
		out, cancel := image_prep.Start(&req)

		w.cancelMap["chooseVFS"] = cancel
		go func() {
			mp := make(map[string]int)
			st := time.Now()
			for event := range out {
				comm := helper.Bytes2String(event.Comm[:])
				str := fmt.Sprintf("pid:%d,comm:%s,operation:%s", event.Pid, comm, helper.Bytes2String(event.Operation[:]))
				if _, ok := mp[str]; ok {
					mp[str]++
					continue
				}
				mp[str] = 1
				w.bpfLogs.AppendLogf(str)
			}
			for log, count := range mp {
				w.bpfLogs.AppendLogf("%s count:%d during %f s\n ", log, count, time.Since(st).Seconds())
			}
		}()
	} else {
		w.bpfLogs.AppendLogf("cancel watch chooseVFS")
		if w.cancelMap["chooseVFS"] != nil {
			w.cancelMap["chooseVFS"]()
		}

	}

}

func (w *ImageSelect) chooseMount(b bool) {
	if b == true {
		w.bpfLogs.AppendLogf("choose watch mount")
		req := mount.MountReq{}
		out, cancel := mount.Start(&req)
		w.cancelMap["chooseMount"] = cancel
		go func() {
			for event := range out {
				w.m.Lock()
				str := fmt.Sprintf("[mount] pid: %d, dev_name: %s, dir_name: %s, type: %s\n",
					event.Pid,
					helper.Bytes2String(event.DevName[:]),
					helper.Bytes2String(event.DirName[:]),
					helper.Bytes2String(event.Type[:]))
				w.bpfLogs.AppendLogf(str)
				w.m.Unlock()
			}
		}()
	} else {
		w.bpfLogs.AppendLogf("cancel watch mount")
		if w.cancelMap["chooseMount"] != nil {
			w.cancelMap["chooseMount"]()
		}

	}

}

func (w *ImageSelect) chooseNetwork(b bool) {
	if b == true {
		w.bpfLogs.AppendLogf("choose watch network")
		ctx, cancel := context.WithCancel(context.TODO())
		w.cancelMap["chooseNetwork"] = cancel
		go w.runBPFTraceScript(ctx, netTracePointScript)
	} else {
		w.bpfLogs.AppendLogf("cancel watch network")
		if w.cancelMap["chooseNetwork"] != nil {
			w.cancelMap["chooseNetwork"]()
		}

	}
}

func (w *ImageSelect) chooseIsolation(b bool) {
	if b == true {
		w.bpfLogs.AppendLogf("choose watch namespace and cgroup")
		ctx, cancel := context.WithCancel(context.TODO())
		w.cancelMap["chooseIsolation"] = cancel
		go w.runBPFTraceScript(ctx, isolationTracePointScript)
	} else {
		w.bpfLogs.AppendLogf("cancel watch namespace and cgroup")
		if w.cancelMap["chooseIsolation"] != nil {
			w.cancelMap["chooseIsolation"]()
			w.bpfLogs.AppendLogf("cancel watch isolation")
		}

	}
}

func (w *ImageSelect) OnChanged(s string) {
	w.m.Lock()
	defer w.m.Unlock()
	images, err := cli.ListImage()
	if err != nil {
		return
	}
	w.images = images
	var imagetags []string
	for _, v := range images {
		for _, tag := range v.RepoTags {
			imagetags = append(imagetags, tag)
			break
		}
	}
	w.base.SetOptions(imagetags)
	for k, v := range w.images {
		if len(v.RepoTags) > 0 && v.RepoTags[0] == s {
			w.currentImageIndex = k
			break
		}
	}
}

// Function to execute bpftrace script and read its output asynchronously
func (w *ImageSelect) runBPFTraceScript(ctx context.Context, scriptPath string) {
	// Prepare the bpftrace command with the script file as argument
	cmd := exec.Command("bpftrace", scriptPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	// Start the command asynchronously
	err = cmd.Start()
	if err != nil {
		fmt.Printf("failed to start bpftrace: %v\n", err)
		return
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	// Start a goroutine to process the output asynchronously
	scanner := bufio.NewScanner(stdout)
	mp := make(map[string]int)
	t := time.Now()
	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("failed to kill process: %v", err)
			}
			for text, cnt := range mp {
				w.bpfLogs.AppendLogf("%s catch count %d during %d ms \n", text, cnt, time.Now().Sub(t)/1000/1000)
			}
			return
		default:
			scanner.Scan()
			if _, ok := mp[scanner.Text()]; !ok {
				mp[scanner.Text()] = 1
				w.bpfLogs.AppendLogf(scanner.Text())
			} else {
				mp[scanner.Text()]++
			}
			if err := scanner.Err(); err != nil {
				log.Printf("Error reading bpftrace output: %v", err)
			}
		}
	}
}
