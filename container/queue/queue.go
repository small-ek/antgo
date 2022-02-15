package queue

import (
	"errors"
	"time"
)

const circleSectionNum = 100

//TaskFunc ...
type TaskFunc func(args ...interface{})

//Task ...
type Task struct {
	runTime  time.Time //初次运行时间
	cycleNum int       //需要第几圈
	curIndex int       //当前运行到第几格
	exec     TaskFunc  //执行的函数
	params   []interface{}
}

//DelayMessage ...
type DelayMessage struct {
	cycleNum  int //当前运行到第几圈了
	curIndex  int //当前运行到第几格
	slots     [circleSectionNum]map[string]*Task
	closed    chan bool
	taskClose chan bool
	timeClose chan bool
	startTime time.Time
}

//New ...
func New() *DelayMessage {
	dm := &DelayMessage{
		cycleNum:  0,
		curIndex:  0,
		closed:    make(chan bool),
		taskClose: make(chan bool),
		timeClose: make(chan bool),
		startTime: time.Now(),
	}
	for i := 0; i < circleSectionNum; i++ {
		dm.slots[i] = make(map[string]*Task)
	}
	return dm
}

//Start ...
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

//Stop ...
func (dm *DelayMessage) Stop() {
	dm.closed <- true
}

//taskLoop
func (dm *DelayMessage) taskLoop() {
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

//timeLoop
func (dm *DelayMessage) timeLoop() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-dm.timeClose:
			return
		case <-tick.C:
			dm.curIndex = (dm.curIndex + 1) % circleSectionNum
			if dm.curIndex == 0 {
				dm.cycleNum += 1
			}
		}
	}

}

//AddTask ...
func (dm *DelayMessage) AddTask(t time.Time, key string, exec TaskFunc, params []interface{}) error {
	if dm.startTime.After(t) {
		return errors.New("Queue time error")
	}
	//当前时间与指定时间相差秒数
	subSecond := t.Unix() - dm.startTime.Unix()
	//计算循环次数
	cycleNum := int(subSecond / circleSectionNum)
	//计算任务所在的slots的下标
	ix := subSecond % circleSectionNum
	//把任务加入tasks中
	tasks := dm.slots[ix]
	if _, err := tasks[key]; err {
		return errors.New("Queue key name already exists")
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
