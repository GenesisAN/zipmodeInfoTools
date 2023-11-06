package main

import (
	"encoding/json"
	"flag"
	"fmt"
	util "illisonModeInfoTools/utils"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	//是否启用多线程
	MultiThread bool
	//指定扫描路径
	ScanPath string
	Threads  int
)

func main() {
	Args()
	fs := util.GetAllFiles(ScanPath, ".zipmod")
	all := len(fs)
	// 统计运行时间
	start := time.Now()
	if all == 0 {
		fmt.Println("未找到任何zipmod文件")
		return
	}

	var wg sync.WaitGroup    // 创建一个新的 WaitGroup
	fileCh := make(chan int) // 创建一个用于传递文件索引的通道
	// 启动 numCPU 个 goroutines 来处理文件
	numCPU := runtime.NumCPU()
	if Threads > 0 {
		numCPU = Threads
	}
	for i := 0; i < numCPU; i++ {
		// 每个 goroutine 都从通道中获取文件索引，然后处理该文件
		// wg.Add(1) 是什么意思？它告诉 WaitGroup，我们有一个新的任务要处理
		wg.Add(1)
		go func(icc int) {
			defer wg.Done()
			for fi := range fileCh {
				processFile(all, fi, icc, fs[fi], &wg)
			}
		}(i)
	}
	// 将文件索引发送到通道
	for i := range fs {
		fileCh <- i
	}

	close(fileCh) // 关闭通道，通知所有 goroutines 没有更多的文件
	wg.Wait()     // 等待所有 goroutine 完成
	elapsed := time.Since(start)
	for _, v := range failmods {
		fmt.Printf("读取失败:%s\n失败原因:%s", v.Path, v.Error)
	}
	fmt.Printf("读取完成%d个MOD，耗时:%s，使用线程数:%d\n", all, elapsed, numCPU)
	eyt, err := json.Marshal(failmods)
	byt, err := json.Marshal(mods)
	if err != nil {
		return
	}
	os.WriteFile("./ModsInfo.json", byt, 0644)
	os.WriteFile("./ModsInfoFail.json", eyt, 0644)
}

var (
	modsMutex     sync.Mutex
	failmodsMutex sync.Mutex

	mods     []util.ModXml
	failmods []util.ModXml
)

func processFile(all, i, ip int, v string, wg *sync.WaitGroup) {
	//defer wg.Done() // 在函数退出时，通知 WaitGroup 一个任务已完成

	mod, err := util.ReadZip(v)
	if err != nil {
		failmodsMutex.Lock()
		failmods = append(failmods, mod)
		failmodsMutex.Unlock()
		fmt.Printf("[%d/%d] -%2d:Read Fail:%s,Error:%s\n", i+1, all, ip, mod.Path, err.Error())
	} else {
		modsMutex.Lock()
		mods = append(mods, mod)
		modsMutex.Unlock()
		fmt.Printf("[%d/%d] -%2d:%s --by:%s\n", i+1, all, ip, mod.Name, mod.Author)
	}
}

func Args() {
	flag.BoolVar(&MultiThread, "pt", false, "是否启用多线程")
	flag.StringVar(&ScanPath, "p", "./", "指定扫描路径")
	flag.IntVar(&Threads, "t", 0, "指定线程数")
	flag.Parse()
}
