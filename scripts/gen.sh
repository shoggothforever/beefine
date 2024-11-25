#!/bin/bash
# 生成 Go 文件模板
gen_go_file(){
echo "generating File "$BPF_DIR/$GO_FILE""
cat > "$BPF_DIR/$GO_FILE" <<EOF
package $PACKAGE_NAME

import (
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"
  "github.com/cilium/ebpf/link"
  "github.com/cilium/ebpf/rlimit"
)
// remove -type event if you won't use diy struct in kernel
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -target bpfel -type event bpf $PACKAGE_NAME.c -- -I../headers
type ${PACKAGE_NAME^}Req struct {

}
type ${PACKAGE_NAME^}Res struct {

}
func Start(req ${PACKAGE_NAME^}Req) (<-chan ${PACKAGE_NAME^}Res,func()) {
  stopper := make(chan os.Signal, 1)
  signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
  // Allow the current process to lock memory for eBPF resources.
  if err := rlimit.RemoveMemlock(); err != nil {
    log.Fatal(err)
  }
  // Load pre-compiled programs and maps into the kernel.
  objs := bpfObjects{}
  if err := loadBpfObjects(&objs, nil); err != nil {
    log.Fatalf("loading objects: %v", err)
  }
  defer objs.Close()
  // write your link code here
  return Action(objs, req, stopper),func(){signal.Notify(stopper, os.Interrupt) }

}
func Action(objs bpfObjects , req ${PACKAGE_NAME^}Req , stopper chan os.Signal) <-chan ${PACKAGE_NAME^}Res{
  // add your link logic here
  out := make(chan ${PACKAGE_NAME^}Res)
  go func() {
    for {
      // write your logical code here
      select {
      case <-stopper:
        return
      default:
        time.Sleep(1 * time.Second)
      }

    }
  }()
  return out
}
// add more action function here
EOF
sudo chmod 777 "$BPF_DIR/$GO_FILE"
# 输出结果
echo "Go 文件已生成: $BPF_DIR/$GO_FILE"
}

gen_c_file(){
  # 生成 C 文件模板
  cat > "$BPF_DIR/$C_FILE" <<EOF
  //go:build ignore

  #include "../vmlinux.h"
  #include "../headers/common.h"
  #include "../headers/bpf_endian.h"
  #include "../headers/bpf_tracing.h"

  char LICENSE[] SEC("license") = "Dual BSD/GPL";
  struct event {
  	u8 comm[16];
  	__u16 val;
  };
  struct event *unused __attribute__((unused));

  SEC("XXX")
  int handle_XXX(){
      // write your code here
  	return 0;
  }
EOF
  sudo chmod 777 "$BPF_DIR/$C_FILE"
  echo "C 文件已生成: $BPF_DIR/$C_FILE"
}


# 检查是否提供了包名参数
if [ $# -lt 1 ]; then
    echo "用法: $0 <包名>"
    echo "例如: $0 mypackage"
    exit 1
fi

BPF_DIR=bpf
PACKAGE_NAME=$1
DIR_NAME=$PACKAGE_NAME
GO_FILE="$PACKAGE_NAME/$PACKAGE_NAME.go"
C_FILE="$PACKAGE_NAME/$PACKAGE_NAME.c"

# 创建包名文件夹
sudo mkdir -p "$BPF_DIR/$DIR_NAME"

# 设置文件夹权限
sudo chmod 777 "$BPF_DIR/$DIR_NAME"
echo "$BPF_DIR/$GO_FILE"
# 检查文件是否存在
if [ ! -f "$BPF_DIR/$GO_FILE" ]; then
   echo "File $BPF_DIR/$GO_FILE does not exist. Creating the file..."
   gen_go_file  # 执行生成文件的逻辑
else
 # 文件已存在，询问是否覆盖
    read -p "File "$BPF_DIR/$GO_FILE" already exists. Do you want to overwrite it? [y/N]: " choice
    case "$choice" in
        y|Y )
            echo "Overwriting file "$BPF_DIR/$GO_FILE...""
            gen_go_file
            ;;
        * )
            ;;
    esac
fi

# 检查文件是否存在
if [ ! -f "$BPF_DIR/$C_FILE" ]; then
   echo "File $BPF_DIR/$GO_FILE does not exist. Creating the file..."
   gen_go_file  # 执行生成文件的逻辑
else
 # 文件已存在，询问是否覆盖
    read -p "File "$BPF_DIR/$C_FILE" already exists. Do you want to overwrite it? [y/N]: " choice
    case "$choice" in
        y|Y )
            echo "Overwriting file "$BPF_DIR/$C_FILE...""
            gen_go_file
            ;;
        * )
            ;;
    esac
fi

echo "文件夹 $PACKAGE_NAME 已设置权限为 777"

# bash gen.sh exec

