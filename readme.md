
# 目的
为了能够对系统的**goroutine并发**进行管理限制，解决以下问题:

*   对所有的goroutine可能的panic情况进行统一管理
*   对所有的goroutine数量设置上限，达到上限时可以进行处理
*   对所有的goroutine的运行情况进行管理


## 设计


### 添加任务

1.  检查是否有空闲，如果有，执行2，否则执行4
2.  执行任务
3.  进行清理
4.  添加到队列
