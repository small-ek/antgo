package queue

import (
	"errors"
	"fmt"
	"log"
	"time"
)

const cicleSectionNum = 100

type TaskFunc func(args ...interface{})

//任务
type Task struct {
	runTime  time.Time //初次运行时间
	cycleNum int       //需要第几圈
	curIndex int       //当前运行到第几格
	//执行的函数
	exec   TaskFunc
	params []interface{}
}

type DelayMessage struct {
	cycleNum  int //当前运行到第几圈了
	curIndex  int //当前运行到第几格
	slots     [cicleSectionNum]map[string]*Task
	closed    chan bool
	taskClose chan bool
	timeClose chan bool
	startTime time.Time
}

func NewDelayMessage() *DelayMessage {
	dm := &DelayMessage{
		cycleNum:  0,
		curIndex:  0,
		closed:    make(chan bool),
		taskClose: make(chan bool),
		timeClose: make(chan bool),
		startTime: time.Now(),
	}
	for i := 0; i < cicleSectionNum; i++ {
		dm.slots[i] = make(map[string]*Task)
	}
	return dm
}

func (dm *DelayMessage) Start() {
	go dm.taskLoop()
	go dm.timeLoop()
	select {
	case <-dm.closed:
		dm.taskClose <- true
		dm.timeClose <- true
		break
	}
}

func (dm *DelayMessage) Stop() {
	dm.closed <- true
}

func (dm *DelayMessage) taskLoop() {
	defer func() {
		log.Println("任务遍历结束！")
	}()
	for {
		select {
		case <-dm.taskClose:
			return
		default:
			{
				tasks := dm.slots[dm.curIndex]
				if len(tasks) > 0 {
					for k, v := range tasks {
						if v.cycleNum == dm.cycleNum {
							go v.exec(v.params...)
							delete(tasks, k)
						}
					}
				}
			}
		}

	}
}

func (dm *DelayMessage) timeLoop() {
	defer func() {
		log.Println("时间遍历结束！")
	}()
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-dm.timeClose:
			return
		case <-tick.C:
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
			dm.curIndex = (dm.curIndex + 1) % cicleSectionNum
			if dm.curIndex == 0 {
				dm.cycleNum += 1
			}
			fmt.Println("当前循环时间", dm.cycleNum, dm.curIndex)
		}
	}

}

//添加任务
func (dm *DelayMessage) AddTask(t time.Time, key string, exec TaskFunc, params []interface{}) error {
	if dm.startTime.After(t) {
		return errors.New("时间错误")
	}
	//当前时间与指定时间相差秒数
	subSecond := t.Unix() - dm.startTime.Unix()
	//计算循环次数
	cycleNum := int(subSecond / cicleSectionNum)
	//计算任务所在的slots的下标
	ix := subSecond % cicleSectionNum
	//把任务加入tasks中
	tasks := dm.slots[ix]
	if _, ok := tasks[key]; ok {
		return errors.New("该slots中已存在key为" + key + "的任务")
	}
	tasks[key] = &Task{
		runTime:  t,
		cycleNum: cycleNum,
		curIndex: int(ix),
		exec:     exec,
		params:   params,
	}
	return nil
}
