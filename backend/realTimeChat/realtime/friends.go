package realtime

import (
	"fmt"
	"realTimeChat/servergrpc"
)

func GetUserFriends(userID string) <-chan []string {
	// ch := make(chan []string)
	// go func() {
	// 	defer close(ch)
	// 	switch userID {
	// 	case "1":
	// 		ch <- []string{"2", "3", "4"}
	// 	case "2":
	// 		ch <- []string{"1", "3", "4"}
	// 	case "3":
	// 		ch <- []string{"1", "2", "4"}
	// 	case "4":
	// 		ch <- []string{"1", "2", "3"}
	// 	default:
	// 		ch <- []string{}
	// 	}

	// }()
	// return ch

	ch := make(chan []string)
	go func() {
		defer close(ch)

		// call grp client fun
		userFriends, err := servergrpc.GetFollowingFollowersClient(userID)

		if err != nil {
			fmt.Printf("Error on friends method realtime")
			ch <- []string{}
			return
		}
		var friends []string
		for _, userIDsList := range userFriends {
			friends = append(friends, userIDsList.UseridsList...)
		}

		ch <- friends

	}()
	return ch
}
