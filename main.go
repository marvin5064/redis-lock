package main

import (
	"fmt"
	"time"

	"github.com/marvin5064/redis-lock/redishelper"
)

func main() {
	name := "lock-for-test-1"
	redisPool := redishelper.NewRedisPool("127.0.0.1:6379")
	// lockTime := 50 * time.Second // time will enter the buffer time, and error out before expire the lock silently
	lockTime := 45 * time.Second // time before the buffer time starts, will not error out
	for {
		for i := 0; i < 1; i++ {
			go func(i int) {
				fmt.Println("try to lock...")
				now := time.Now()
				lock, err := redishelper.NewLock(redisPool, name)
				if err != nil {
					fmt.Println("err", err)
					return
				}
				fmt.Println(name, "locked")
				time.Sleep(lockTime / 2)
				fmt.Println("working for i:", i)
				time.Sleep(lockTime / 2)
				lock.Unlock()
				fmt.Println("from lock to unlock, time cost is:", time.Since(now))
				fmt.Println(name, "unlocked")
			}(i)
		}

		time.Sleep(lockTime * 2)
	}
}
