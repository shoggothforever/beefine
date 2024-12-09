package imager

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	imaget "github.com/docker/docker/api/types/image"
	"log"
	"os"
	"os/exec"
	exec2 "shoggothforever/beefine/bpf/exec"
	"shoggothforever/beefine/bpf/image_prep"
	"shoggothforever/beefine/bpf/mount"
	"shoggothforever/beefine/internal/cli"
	"sync"
	"unsafe"
)

// MyCustomWidget 是自定义控件，包装了 Select 并添加了额外的字段
type ImageSelect struct {
	widget.BaseWidget                // 嵌入 BaseWidget
	base              *widget.Select // 内嵌 Select
	imageLogs         *widget.TextGrid
	bpfLogs           *widget.TextGrid
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

func (w *ImageSelect) chooseUnionFS(b bool) {
	if b == true {
		fmt.Println("choose watch unionfs")
		req := image_prep.ImagePrepReq{}
		out, cancel := image_prep.Start(&req)
		w.cancelMap["chooseUnionFS"] = cancel
		go func() {
			for event := range out {
				//time.Sleep(time.Second)
				comm := Bytes2String(event.Comm[:])
				w.m.Lock()
				str := fmt.Sprintf("pid:%d,comm:%s,operation:%s", event.Pid, comm, Bytes2String(event.Operation[:]))
				fmt.Println(str)
				w.AppendLogInLock(w.bpfLogs, str)
				w.m.Unlock()
			}
		}()
	} else {
		fmt.Println("cancel watch unionfs")
		if w.cancelMap["chooseUnionFS"] != nil {
			w.cancelMap["chooseUnionFS"]()
		}

	}

}

func (w *ImageSelect) chooseMount(b bool) {
	if b == true {
		fmt.Println("choose watch mount")
		req := mount.MountReq{}
		out, cancel := mount.Start(&req)
		w.cancelMap["chooseMount"] = cancel
		go func() {
			for event := range out {
				w.m.Lock()
				str := fmt.Sprintf("PID: %d, dev_name: %s, dir_name: %s, type: %s\n",
					event.Pid,
					Bytes2String(event.DevName[:]),
					Bytes2String(event.DirName[:]),
					Bytes2String(event.Type[:]))
				w.AppendLogInLock(w.bpfLogs, str)
				w.m.Unlock()
			}
		}()
	} else {
		fmt.Println("cancel watch mount")
		if w.cancelMap["chooseMount"] != nil {
			w.cancelMap["chooseMount"]()
		}

	}

}

func (w *ImageSelect) chooseNetwork(b bool) {
	if b == true {
		fmt.Println("choose watch network")
		ctx, cancel := context.WithCancel(context.TODO())
		w.cancelMap["chooseNetwork"] = cancel
		go w.runBPFTraceScript(ctx, netTracePointScript)
	} else {
		fmt.Println("cancel watch network")
		if w.cancelMap["chooseNetwork"] != nil {
			w.cancelMap["chooseNetwork"]()
		}

	}
}

func (w *ImageSelect) chooseIsolation(b bool) {
	if b == true {
		fmt.Println("choose watch namespace and cgroup")
		ctx, cancel := context.WithCancel(context.TODO())
		w.cancelMap["chooseIsolation"] = cancel
		go w.runBPFTraceScript(ctx, isolationTracePointScript)
	} else {
		fmt.Println("cancel watch namespace and cgroup")
		if w.cancelMap["chooseIsolation"] != nil {
			w.cancelMap["chooseIsolation"]()
		}

	}
}

func (w *ImageSelect) chooseProcess(b bool) {
	if b == true {
		fmt.Println("choose watch process")
		req := &exec2.ExecReq{}
		out, cancel := exec2.Start(req)
		w.cancelMap["chooseProcess"] = cancel
		go func() {
			mp := make(map[string]uint64)
			for e := range out {
				w.m.Lock()
				comm := Bytes2String(e.Comm[:])
				if e.ExitEvent {
					w.AppendLogInLock(w.bpfLogs, fmt.Sprintf("exit duration_ns:%v,prio:%d, pid: %d, comm: %s\n", e.Ts-mp[comm], e.Prio, e.Pid, comm))
				} else {
					mp[comm] = e.Ts
					w.AppendLogInLock(w.bpfLogs, fmt.Sprintf("exec pid: %d, comm: %s\n", e.Pid, comm))
				}
				w.m.Unlock()
			}
		}()
	} else {
		fmt.Println("cancel watch process")
		if w.cancelMap["chooseProcess"] != nil {
			w.cancelMap["chooseProcess"]()
		}
	}
}

func (w *ImageSelect) AppendLogInLock(logs *widget.TextGrid, text string) {
	// 将日志写入到本地，可以与其它框架结合起来使用
	file, err := os.OpenFile("tmplog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	_, err = writer.WriteString(text)
	if err != nil {
		fmt.Println(err)
		return
	}
	logs.SetRow(len(logs.Rows), widget.NewTextGridFromString(text).Row(0))
	if len(logs.Rows) > 500 {
		logs.SetText("")
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
			fmt.Println("select image ", w.images[k])
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
		fmt.Println(err)
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
	mp := make(map[string]struct{})
	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("failed to kill process: %v", err)
			}
			return
		default:
			scanner.Scan()
			// Process each line of output
			w.m.Lock()
			if _, ok := mp[scanner.Text()]; !ok {
				mp[scanner.Text()] = struct{}{}
				w.AppendLogInLock(w.bpfLogs, scanner.Text())
				fmt.Println("bpftrace output:", scanner.Text())
			}
			w.m.Unlock()
			// Check for scanning errors
			if err := scanner.Err(); err != nil {
				log.Printf("Error reading bpftrace output: %v", err)
			}
		}
	}
}

func Bytes2String(b []byte) string {
	trimmedData := bytes.TrimRight(b, "\x00")
	return *(*string)(unsafe.Pointer(&trimmedData))
}
