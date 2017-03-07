package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/marvin5064/redis-lock/redishelper"
)

func main() {
	prefix := "lock-for-test-"
	// lockTime := 50 * time.Second // time will enter the buffer time, and error out before expire the lock silently
	lockTime := 45 * time.Second // time before the buffer time starts, will not error out
	for j := 0; j < 7; j++ {
		for i := 0; i < 15; i++ {
			go func(i int) {

				redisPool := redishelper.NewRedisPool("127.0.0.1:6379")
				name := prefix + strconv.Itoa(i)
				fmt.Println("try to lock...")
				now := time.Now()
				lock, err := redishelper.NewLock(redisPool, name)
				if err != nil {
					fmt.Println("err", err)
					return
				}
				fmt.Println(name, "locked")
				time.Sleep(3 * time.Second)
				lock.Unlock()
				fmt.Println("from lock to unlock, time cost is:", time.Since(now))
				fmt.Println(name, "unlocked")
				redisPool.Close()
			}(i)
		}
	}
	time.Sleep(lockTime * 2)
}
