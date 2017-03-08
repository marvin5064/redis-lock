package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/marvin5064/redis-lock/redishelper"
)

func main() {
	prefix := "locking-wtf-"
	// lockTime := 50 * time.Second // time will enter the buffer time, and error out before expire the lock silently
	lockTime := 45 * time.Second // time before the buffer time starts, will not error out
	redisPool := redishelper.NewRedisPool("127.0.0.1:6379")
	for j := 0; j < 7; j++ {
		for i := 0; i < 400/7; i++ {
			go func(i int) {
				name := prefix + strconv.Itoa(i)
				fmt.Println("try to lock...", name)
				now := time.Now()
				redisClient := redisPool.GetDB(0)
				lock, err := redishelper.NewLock(redisClient, name)
				if err != nil {
					fmt.Println("err", err)
					return
				}
				fmt.Println(name, "locked")
				time.Sleep(3 * time.Second)
				boo := lock.Unlock()
				if boo != true {
					fmt.Println("WTF, unable to unlock", name)
					return
				}
				fmt.Println("from lock to unlock, time cost is:", time.Since(now))
				fmt.Println(name, "unlocked")
				redisClient.Close()
			}(i)
		}
	}
	time.Sleep(lockTime)
}
